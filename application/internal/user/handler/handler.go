package handler

// @title User Service API
// @version 1.0
// @description User Service provides user management functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"application/internal/user/service"

	// Import generated docs
	_ "application/docs/user"
)

type UserHandler struct {
	userService *service.UserService
	grpcServer  *grpc.Server
	httpServer  *fiber.App
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) StartGRPCServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	h.grpcServer = grpc.NewServer()
	// Register service when proto files are generated
	// user.RegisterUserServiceServer(h.grpcServer, h.userService)
	reflection.Register(h.grpcServer)

	log.Printf("User gRPC server listening on port %d", port)
	return h.grpcServer.Serve(lis)
}

func (h *UserHandler) StartHTTPServer(port int) error {
	h.httpServer = fiber.New(fiber.Config{
		AppName: "User Service HTTP API",
	})

	h.httpServer.Use(cors.New())
	h.setupRoutes()

	log.Printf("User HTTP server listening on port %d", port)
	return h.httpServer.Listen(fmt.Sprintf(":%d", port))
}

func (h *UserHandler) setupRoutes() {
	// Health check endpoint
	h.httpServer.Get("/healthz", h.healthCheck)

	// Swagger documentation
	h.httpServer.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// Serve the swagger.json file
	h.httpServer.Static("/swagger/doc.json", "./docs/user/swagger.json")

	// API routes
	api := h.httpServer.Group("/api/v1")
	api.Post("/users", h.createUser)
	api.Get("/users/:id", h.getUser)
	api.Put("/users/:id", h.updateUser)
	api.Delete("/users/:id", h.deleteUser)
	api.Get("/users", h.listUsers)
}

// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthz [get]
func (h *UserHandler) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "user-service",
		"version": "1.0.0",
	})
}

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body object true "User data"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/users [post]
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	var req service.CreateNewUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.CreateUser(context.Background(), &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"user": user,
	})
}

// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	user, err := h.userService.GetUser(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// @Summary Update user
// @Description Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body object true "User data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	var req service.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.UpdateUser(context.Background(), &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.userService.DeleteUser(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users [get]
func (h *UserHandler) listUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	users, total, err := h.userService.ListUsers(context.Background(), int32(page), int32(limit))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := (total + int32(limit) - 1) / int32(limit)

	return c.JSON(fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *UserHandler) Shutdown() error {
	if h.grpcServer != nil {
		h.grpcServer.GracefulStop()
	}
	if h.httpServer != nil {
		return h.httpServer.Shutdown()
	}
	return nil
}
