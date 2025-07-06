package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"

	"application/internal/hello/handler"
)

type Config struct {
	Server struct {
		GRPCPort int    `yaml:"grpc_port"`
		HTTPPort int    `yaml:"http_port"`
		Host     string `yaml:"host"`
	} `yaml:"server"`
	Service struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"service"`
}

func main() {
	// Load configuration
	config, err := loadConfig("configs/hello.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create handler
	helloHandler := handler.NewHelloHandler()

	// Start servers
	go func() {
		if err := helloHandler.StartGRPCServer(config.Server.GRPCPort); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	go func() {
		if err := helloHandler.StartHTTPServer(config.Server.HTTPPort); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down servers...")
	if err := helloHandler.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
