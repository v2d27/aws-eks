package service

import (
	"context"
	"fmt"
	"time"

	"application/api/proto/common"
	"application/api/proto/hello"
)

type HelloService struct {
	hello.UnimplementedHelloServiceServer
}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	greeting := s.getGreeting(req.Name, req.Language)

	return &hello.HelloResponse{
		Message:   greeting,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (s *HelloService) SayHelloStream(req *hello.HelloRequest, stream hello.HelloService_SayHelloStreamServer) error {
	for i := 0; i < 5; i++ {
		greeting := fmt.Sprintf("%s (stream %d)", s.getGreeting(req.Name, req.Language), i+1)

		response := &hello.HelloResponse{
			Message:   greeting,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(response); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (s *HelloService) GetGreetings(ctx context.Context, req *common.PaginationRequest) (*hello.GreetingsResponse, error) {
	// Mock data for demonstration
	allGreetings := []*hello.Greeting{
		{Id: "1", Message: "Hello World", Language: "en", CreatedAt: time.Now().Unix()},
		{Id: "2", Message: "Hola Mundo", Language: "es", CreatedAt: time.Now().Unix()},
		{Id: "3", Message: "Bonjour le monde", Language: "fr", CreatedAt: time.Now().Unix()},
		{Id: "4", Message: "Hallo Welt", Language: "de", CreatedAt: time.Now().Unix()},
		{Id: "5", Message: "Ciao Mondo", Language: "it", CreatedAt: time.Now().Unix()},
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	total := int32(len(allGreetings))
	totalPages := (total + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return &hello.GreetingsResponse{
			Greetings: []*hello.Greeting{},
			Pagination: &common.PaginationResponse{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
			},
		}, nil
	}

	if end > total {
		end = total
	}

	return &hello.GreetingsResponse{
		Greetings: allGreetings[start:end],
		Pagination: &common.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *HelloService) getGreeting(name, language string) string {
	greetings := map[string]string{
		"en": "Hello",
		"es": "Hola",
		"fr": "Bonjour",
		"de": "Hallo",
		"it": "Ciao",
	}

	greeting, exists := greetings[language]
	if !exists {
		greeting = "Hello"
	}

	if name == "" {
		name = "World"
	}

	return fmt.Sprintf("%s, %s!", greeting, name)
}
