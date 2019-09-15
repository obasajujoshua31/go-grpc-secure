package api

import (
	"errors"
	fmt "fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var todos = &Todos{
	Todos: []*Todo{
		&Todo{
			Id:          "2",
			Title:       "This is the first Title",
			IsCompleted: false,
		},
		&Todo{
			Id:          "2",
			Title:       "This is the second Title",
			IsCompleted: true,
		},
		&Todo{
			Id:          "3",
			Title:       "This is the third Title",
			IsCompleted: true,
		},
	},
}

func getOneTodo(id string) (*Todo, error) {
	for _, todo := range todos.GetTodos() {
		if todo.Id == id {
			return todo, nil
		}
	}
	return nil, errors.New("Cannot find todo")
}

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Received Message :%s", in.Greeting)
	return &PingMessage{Greeting: in.GetGreeting()}, nil
}

func (s *Server) GetTodos(ctx context.Context, in *Empty) (*Todos, error) {
	log.Print("Request to Get Notes")
	return todos, nil
}

func (s *Server) GetTodo(ctx context.Context, in *TodoId) (*Todo, error) {
	log.Printf("Received request to get one Todo :%s", in.GetId())
	todo, err := getOneTodo(in.GetId())

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Todo of id :%s was not found", in.GetId()))
	}
	return todo, nil

}
