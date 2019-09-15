package main

import (
	"encoding/json"
	"fmt"
	"go-grpc/api"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Authentication struct {
	Login    string
	Password string
}

// Get current Metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil

}

// Require Transport layer security

func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

// var OneTodo = make(map[string]string)

type OneTodo struct {
	Id          string
	Title       string
	IsCompleted bool
}

var Alltodos []OneTodo

// var TodoFormats []OneTodo

func main() {
	var conn *grpc.ClientConn

	// Setup connection Credentials

	cred, err := credentials.NewClientTLSFromFile("cert/server.crt", "")

	if err != nil {
		log.Fatalf("Could not load Credentials : %s", err)
	}

	auth := Authentication{
		Login:    "john",
		Password: "doe",
	}

	// Connect to Remote GRPC Server
	conn, err = grpc.Dial("localhost:7777", grpc.WithTransportCredentials(cred), grpc.WithPerRPCCredentials(&auth))

	if err != nil {
		log.Fatalf("Could not connect to server %s", err)
	}

	defer conn.Close()

	c := api.NewPingClient(conn)

	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "Foo"})

	if err != nil {
		log.Fatalf("Could get response from remote Server: %s", err)
	}

	log.Printf("Response from server: %s", response.Greeting)

	res, err := c.GetTodos(context.Background(), &api.Empty{})

	if err != nil {
		log.Fatalf("Could not get response for Todos from remote Server: %s", err)
	}

	for _, one := range res.Todos {
		Alltodos = append(Alltodos, OneTodo{
			Id:          one.Id,
			Title:       one.Title,
			IsCompleted: one.IsCompleted,
		})
	}

	resp, err := json.MarshalIndent(Alltodos, "", "\t")

	if err != nil {
		fmt.Printf("Error occured :%s", err)
	}
	log.Print("Data Retrieved", string(resp))

	OneTodoResponse, err := c.GetTodo(context.Background(), &api.TodoId{Id: "3"})

	if err != nil {
		fmt.Printf("Error :%s", err)
		return
	}
	respo, err := json.MarshalIndent(OneTodoResponse, "", "\t")

	log.Print("Data Retrieved", string(respo))
}
