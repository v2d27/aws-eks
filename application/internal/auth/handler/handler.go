package handler

// @title Auth Service API
// @version 1.0
// @description Auth Service provides authentication and authorization functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"application/internal/auth/service"

	// Import generated docs
	_ "application/docs/auth"
)

type AuthHandler struct {
	authService *service.AuthService
	grpcServer  *grpc.Server
	httpServer  *fiber.App
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(jwtSecret),
	}
}

func (h *AuthHandler) StartGRPCServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	h.grpcServer = grpc.NewServer()
	// Register service when proto files are generated
	// auth.RegisterAuthServiceServer(h.grpcServer, h.authService)
	reflection.Register(h.grpcServer)

	log.Printf("Auth gRPC server listening on port %d", port)
	return h.grpcServer.Serve(lis)
}

func (h *AuthHandler) StartHTTPServer(port int) error {
	h.httpServer = fiber.New(fiber.Config{
		AppName: "Auth Service HTTP API",
	})

	h.httpServer.Use(cors.New())
	h.setupRoutes()

	log.Printf("Auth HTTP server listening on port %d", port)
	return h.httpServer.Listen(fmt.Sprintf(":%d", port))
}

func (h *AuthHandler) setupRoutes() {
	// Health check endpoint
	h.httpServer.Get("/healthz", h.healthCheck)

	// Swagger documentation
	h.httpServer.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// Serve the swagger.json file
	h.httpServer.Static("/swagger/doc.json", "./docs/auth/swagger.json")

	// API routes
	api := h.httpServer.Group("/api/v1")
	api.Post("/auth/login", h.login)
	api.Post("/auth/logout", h.logout)
	api.Post("/auth/validate", h.validateToken)
	api.Post("/auth/refresh", h.refreshToken)
}

// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthz [get]
func (h *AuthHandler) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "auth-service",
		"version": "1.0.0",
	})
}

// @Summary Login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body object true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.authService.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"expires_in":    response.ExpiresIn,
		"token_type":    response.TokenType,
	})
}

// @Summary Logout
// @Description Invalidate access token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body object true "Token to invalidate"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) logout(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.authService.Logout(context.Background(), req.Token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// @Summary Validate token
// @Description Validate access token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body object true "Token to validate"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/validate [post]
func (h *AuthHandler) validateToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.authService.ValidateToken(context.Background(), req.Token)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"valid":      response.Valid,
		"user_id":    response.UserID,
		"username":   response.Username,
		"expires_at": response.ExpiresAt,
	})
}

// @Summary Refresh token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body object true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) refreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.authService.RefreshToken(context.Background(), req.RefreshToken)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"expires_in":    response.ExpiresIn,
	})
}

func (h *AuthHandler) Shutdown() error {
	if h.grpcServer != nil {
		h.grpcServer.GracefulStop()
	}
	if h.httpServer != nil {
		return h.httpServer.Shutdown()
	}
	return nil
}
