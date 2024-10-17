package main

import (
	"log"
	"net"

	userpb "github.com/centodiechi/unary_streams/protos/user/v1"
	server "github.com/centodiechi/unary_streams/unaryStreaming/server"
	storage "github.com/centodiechi/unary_streams/unaryStreaming/storage_provider"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

var logger = server.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	storage.DataBase, err = storage.NewStorageProvider("DataBase")
	if err != nil {
		logger.With(zap.Error(err)).Error("error in initializing storage instance")
	}
}

func main() {
	defer storage.DataBase.DB.Close()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to listen on port 8080")
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterRegisterServiceServer(grpcServer, &server.RegisterService{})
	userpb.RegisterAuthServiceServer(grpcServer, &server.LoginService{})

	logger.Info("Starting gRPC server on :8080")

	if err := grpcServer.Serve(lis); err != nil {
		logger.With(zap.Error(err)).Error("Failed to serve gRPC server")
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
