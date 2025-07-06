package handler

// @title Hello Service API
// @version 1.0
// @description Hello Service provides greeting functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
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

	"application/api/proto/common"
	"application/api/proto/hello"
	"application/internal/hello/service"

	// Import generated docs
	_ "application/docs/hello"
)

type HelloHandler struct {
	helloService *service.HelloService
	grpcServer   *grpc.Server
	httpServer   *fiber.App
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{
		helloService: service.NewHelloService(),
	}
}

func (h *HelloHandler) StartGRPCServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	h.grpcServer = grpc.NewServer()
	// Register service when proto files are generated
	// hello.RegisterHelloServiceServer(h.grpcServer, h.helloService)
	reflection.Register(h.grpcServer)

	log.Printf("Hello gRPC server listening on port %d", port)
	return h.grpcServer.Serve(lis)
}

func (h *HelloHandler) StartHTTPServer(port int) error {
	h.httpServer = fiber.New(fiber.Config{
		AppName: "Hello Service HTTP API",
	})

	h.httpServer.Use(cors.New())
	h.setupRoutes()

	log.Printf("Hello HTTP server listening on port %d", port)
	return h.httpServer.Listen(fmt.Sprintf(":%d", port))
}

func (h *HelloHandler) setupRoutes() {
	// Health check endpoint
	h.httpServer.Get("/healthz", h.healthCheck)

	// Swagger documentation
	h.httpServer.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// Serve the swagger.json file
	h.httpServer.Static("/swagger/doc.json", "./docs/hello/swagger.json")

	// API routes
	api := h.httpServer.Group("/api/v1")
	api.Get("/hello", h.sayHelloHTTP)
	api.Get("/greetings", h.getGreetingsHTTP)
}

// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthz [get]
func (h *HelloHandler) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "hello-service",
		"version": "1.0.0",
	})
}

// @Summary Say hello
// @Description Get a greeting message
// @Tags hello
// @Accept json
// @Produce json
// @Param name query string false "Name to greet"
// @Param language query string false "Language for greeting"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/hello [get]
func (h *HelloHandler) sayHelloHTTP(c *fiber.Ctx) error {
	name := c.Query("name", "World")
	language := c.Query("language", "en")
	// Call gRPC service method directly for HTTP endpoint
	ctx := context.Background()
	response, err := h.helloService.SayHello(ctx, &hello.HelloRequest{
		Name:     name,
		Language: language,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   response.Message,
		"timestamp": response.Timestamp,
	})
}

// @Summary Get greetings
// @Description Get paginated list of greetings
// @Tags hello
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/greetings [get]
func (h *HelloHandler) getGreetingsHTTP(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	ctx := context.Background()
	response, err := h.helloService.GetGreetings(ctx, &common.PaginationRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"greetings":  response.Greetings,
		"pagination": response.Pagination,
	})
}

func (h *HelloHandler) Shutdown() error {
	if h.grpcServer != nil {
		h.grpcServer.GracefulStop()
	}
	if h.httpServer != nil {
		return h.httpServer.Shutdown()
	}
	return nil
}
