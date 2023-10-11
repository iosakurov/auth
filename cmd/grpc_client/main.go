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

	desc "github.com/iosakurov/auth/pkg/auth_v1"
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
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal("Произошла ошибка")
		}
	}()

	client := desc.NewUserAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Print(color.RedString("User Get\n"))
	response, err := client.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", response.GetName()))

	log.Print(color.RedString("User Update\n"))
	updateResponse, err := client.Update(ctx, &desc.UpdateRequest{
		Id:   userID,
		Name: wrapperspb.String(gofakeit.Name()),
		Role: desc.Role_ADMIN,
	})
	if err != nil {
		log.Fatalf("failed to update user by id: %v", err)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", updateResponse))

	log.Print(color.RedString("User Create\n"))
	password := gofakeit.Password(true, false, false, false, false, 10)
	role := desc.Role_USER
	createResponse, err := client.Create(ctx, &desc.CreateRequest{
		Name:            gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        password,
		PasswordConfirm: password,
		Role:            role,
	})
	if err != nil {
		log.Fatalf("failed to create user by id: %v", err)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", createResponse))

	log.Print(color.RedString("User Delete\n"))
	deleteResponse, err := client.Delete(ctx, &desc.DeleteRequest{Id: 666})
	if err != nil {
		log.Fatalf("failed to delete user by id: %v", err)
	}
	log.Printf(color.RedString("User Info:\n"), color.GreenString("%+v", deleteResponse))

}
