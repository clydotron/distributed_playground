package main

import (
	"auth-service/data"
	auth "auth-service/protos"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

var gRpcPort = "50001"

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	UserRepo *data.PostgresUserRepo
}

func (a *AuthServer) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {

	user, err := a.UserRepo.GetByEmail(req.GetEmail())
	if err != nil {
		return &auth.AuthResponse{Result: "Invalid credentials"}, nil
	}

	valid, err := user.PasswordMatches(req.GetPassword())
	if err != nil || !valid {
		return &auth.AuthResponse{Result: "Invalid credentials"}, nil
	}

	return &auth.AuthResponse{Result: "Authenticated"}, nil
}

func (app *App) gRPCListen() {
	listener, err := net.Listen("tcp", ":"+gRpcPort)
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	srv := grpc.NewServer()
	auth.RegisterAuthServiceServer(srv, &AuthServer{UserRepo: app.UserRepo})

	log.Println("gRPC server started on port", gRpcPort)

	if err := srv.Serve(listener); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}
