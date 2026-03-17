package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "test-task-miras/api/proto/gen"
	grpcdelivery "test-task-miras/internal/delivery/grpc"
	"test-task-miras/internal/infrastructure/logger"
	"test-task-miras/internal/infrastructure/postgres"
	"test-task-miras/internal/usecase"
)

func main() {

	db, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	logger.L.Info("Successfully connected to PostgreSQL")

	// we can also use  this kind of implementation
	// provided two variance, to show that we not depended on DB
	// repo := memory.NewShipmentRepository()
	repo := postgres.NewShipmentRepository(db)
	uc := usecase.NewShipmentUseCase(repo)
	handler := grpcdelivery.NewShipmentHandler(uc)

	server := grpc.NewServer()
	pb.RegisterShipmentServiceServer(server, handler)

	// To test with postman
	reflection.Register(server)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.L.Error("Failed to listen", "port", port, "error", err)
		os.Exit(1)
	}

	go func() {
		logger.L.Info("gRPC server is running", "port", port)
		if err := server.Serve(listener); err != nil {
			logger.L.Error("Failed to serve gRPC server", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// blocking main until we get signal
	<-stop

	logger.L.Info("Shutting down gRPC server gracefully...")
	server.GracefulStop()
	logger.L.Info("Server stopped. Bye!")
}
