package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	desc "github.com/iosakurov/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserAPIServer
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Get id: %d\n", req.GetId())

	password := gofakeit.Password(true, false, false, false, false, 10)
	role := desc.Role_ROLE_ADMIN

	return &desc.GetResponse{
		Info: &desc.UserInfo{
			Id:              req.GetId(),
			Name:            gofakeit.Name(),
			Email:           gofakeit.Email(),
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
			CreatedAt:       timestamppb.New(gofakeit.Date()),
			UpdatedAt:       timestamppb.New(gofakeit.Date()),
		},
	}, nil

}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Update id: %d, name: %s, email: %s\n", req.GetId(), req.GetName(), req.GetEmail())

	return &emptypb.Empty{}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create: %#v", req)
	return &desc.CreateResponse{Id: 666}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Delete: %#v", req.GetId())
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
