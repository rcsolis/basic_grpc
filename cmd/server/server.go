package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	pb "github.com/rcsolis/basic_grpc/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PORT = ":50051"
)

// TodoServer is the server that provides the TodoService
type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	todos map[string]*pb.Todo
}

// CreateTodo creates a new todo
func (s *TodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	newDateTime := time.Now().Format(time.RFC3339)
	newId := uuid.New().String()
	newName := req.GetName()
	newDescription := req.GetDescription()
	if newName == "" {
		log.Print("ERROR: Name cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "Name cannot be empty")
	}
	if newDescription == "" {
		log.Println("ERROR: Description cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "Description cannot be empty")
	}
	todo := &pb.Todo{
		Date:        newDateTime,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Done:        false,
	}
	s.todos[newId] = todo
	log.Printf("Created todo with id: %s", newId)
	return &pb.CreateTodoResponse{Todo: todo}, nil
}

// DeleteTodo deletes a todo
func (s *TodoServer) DeleteTodo(ctx context.Context, req *pb.TodoIdRequest) (*pb.Empty, error) {
	id := req.GetId()
	if _, ok := s.todos[id]; !ok {
		log.Println("ERROR: Todo not found")
		return nil, status.Error(codes.NotFound, "Todo not found")
	}
	delete(s.todos, id)
	log.Printf("Deleted todo with id: %s", id)
	return &pb.Empty{}, nil
}

// Get All Todos using a stream
func (s *TodoServer) GetAllTodos(req *pb.Empty, stream pb.TodoService_GetAllTodosServer) error {
	log.Println("Getting all todos")

	for _, todo := range s.todos {
		if err := stream.Send(todo); err != nil {
			log.Println("ERROR: Failed to send todo")
			return err
		}
	}
	return nil
}

// Implement main function
func main() {
	// Get a new listener on the port
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Failed connection on port: %v", err)
	}
	// Crete a grpc server
	server := grpc.NewServer()
	// Register the TodoServer implementation
	pb.RegisterTodoServiceServer(server, &TodoServer{todos: make(map[string]*pb.Todo)})
	// Start the server
	log.Printf("Starting server: %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
