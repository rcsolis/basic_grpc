package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/rcsolis/basic_grpc/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ADDRESS = "localhost:50051"
)

type NewTodoTask struct {
	Name        string
	Description string
}

type TodoTask struct {
	Date        string
	Name        string
	Description string
	Done        bool
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcServer, err := grpc.NewClient(ADDRESS, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer grpcServer.Close()

	todoClient := pb.NewTodoServiceClient(grpcServer)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	todos := []*NewTodoTask{
		&NewTodoTask{Name: "Task 1", Description: "Description 1"},
		&NewTodoTask{Name: "Task 2", Description: "Description 2"},
		&NewTodoTask{Name: "Task 3", Description: "Description 3"},
		&NewTodoTask{Name: "Task 4", Description: "Description 4"},
	}

	for _, todo := range todos {
		createTodoRequest := &pb.CreateTodoRequest{Name: todo.Name, Description: todo.Description}
		createTodoResponse, err := todoClient.CreateTodo(ctx, createTodoRequest)
		if err != nil {
			log.Fatalf("Failed to create todo: %v", err)

		}
		log.Printf("Created todo: %v", createTodoResponse.Todo)
	}

	todosStream, err := todoClient.GetAllTodos(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to get all todos: %v", err)
	}
	for {
		todo, err := todosStream.Recv()
		if err != io.EOF {
			log.Printf("End of Stream: %v", err)
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive todo: %v", err)
		}
		log.Printf("Received todo: %v", todo)
	}

}
