package main

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"time"

	desc "github.com/iosakurov/auth/grpc/pkg/auth_v1"
)

const (
	address = "localhost:50051"
	userID  = 1337
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := desc.NewUserAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Printf(color.RedString("User Get\n"))
	response, err := client.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", response.GetInfo()))

	log.Printf(color.RedString("User Update\n"))
	updateResponse, updateErr := client.Update(ctx, &desc.UpdateRequest{
		Id:   userID,
		Name: wrapperspb.String(gofakeit.Name()),
	})
	if updateErr != nil {
		log.Fatalf("failed to update user by id: %v", updateErr)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", updateResponse))

	log.Printf(color.RedString("User Create\n"))
	password := gofakeit.Password(true, false, false, false, false, 10)
	role := desc.Role_ROLE_USER
	createResponse, createErr := client.Create(ctx, &desc.CreateRequest{
		Name:            gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        password,
		PasswordConfirm: password,
		Role:            role,
	})
	if createErr != nil {
		log.Fatalf("failed to create user by id: %v", createErr)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", createResponse))

	log.Printf(color.RedString("User Delete\n"))
	deleteResponse, deleteErr := client.Delete(ctx, &desc.DeleteRequest{Id: 666})
	if deleteErr != nil {
		log.Fatalf("failed to delete user by id: %v", deleteErr)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", deleteResponse))

}
