package implementation

import (
	google_protobuf1 "github.com/golang/protobuf/ptypes/empty"
	"github.com/joelclouddistrict/gochecklist/services"
	"github.com/joelclouddistrict/gochecklist/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type TodoServiceServer struct {
	Store store.Storer
}

func (s *TodoServiceServer) CreateTodo(ctx context.Context, todo *services.TodoMessage) (*services.TodoMessage, error) {
	todo, err := s.Store.CreateTodo(todo)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	return todo, nil
}

func (s *TodoServiceServer) ListTodos(ctx context.Context, filter *services.TodoFilter) (*services.TodoArray, error) {
	todos, err := s.Store.ListTodos(filter)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	return todos, nil
}

func (s *TodoServiceServer) GetTodo(ctx context.Context, id *services.IdMessage) (*services.TodoMessage, error) {
	todo, err := s.Store.GetTodo(id)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	return todo, nil
}

func (s *TodoServiceServer) SetAsDone(ctx context.Context, id *services.IdMessage) (*services.TodoMessage, error) {
	todo := &services.TodoMessage{
		Id:   id.Id,
		Done: true,
	}
	todo, err := s.Store.PatchTodo(todo)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	return todo, nil
}

func (s *TodoServiceServer) SetAsUndone(ctx context.Context, id *services.IdMessage) (*services.TodoMessage, error) {
	todo := &services.TodoMessage{
		Id:   id.Id,
		Done: false,
	}
	todo, err := s.Store.PatchTodo(todo)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	return todo, nil
}

func (s *TodoServiceServer) DeleteTodo(ctx context.Context, id *services.IdMessage) (*google_protobuf1.Empty, error) {
	err := s.Store.DeleteTodo(id)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	out := new(google_protobuf1.Empty)
	return out, nil
}
