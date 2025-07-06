package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNewUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserRequest struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserService struct {
	users map[string]*User
	mutex sync.RWMutex
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateNewUserRequest) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if username already exists
	for _, user := range s.users {
		if user.Username == req.Username {
			return nil, fmt.Errorf("username already exists")
		}
		if user.Email == req.Email {
			return nil, fmt.Errorf("email already exists")
		}
	}

	user := &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // In production, this should be hashed
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users[user.ID] = user
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[req.ID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Check if username/email already exists for other users
	for id, existingUser := range s.users {
		if id != req.ID {
			if existingUser.Username == req.Username {
				return nil, fmt.Errorf("username already exists")
			}
			if existingUser.Email == req.Email {
				return nil, fmt.Errorf("email already exists")
			}
		}
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	// Only update password if it's provided
	if req.Password != "" {
		user.Password = req.Password // In production, this should be hashed
	}
	user.UpdatedAt = time.Now()

	s.users[req.ID] = user
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(s.users, id)
	return nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int32) ([]*User, int32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var users []*User
	for _, user := range s.users {
		users = append(users, user)
	}

	total := int32(len(users))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []*User{}, total, nil
	}

	if end > total {
		end = total
	}

	return users[start:end], total, nil
}

// UserResponse represents a user without sensitive data like password
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts a User to UserResponse (without password)
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserResponseList converts a slice of Users to UserResponse slice
func ToUserResponseList(users []*User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToUserResponse()
	}
	return responses
}
