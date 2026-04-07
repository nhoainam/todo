package mapper

import (
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
	pb "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
)

func TodoFromProto(pbResp *pb.GetTodoResponse) *entity.Todo {
	if pbResp == nil || pbResp.Todo == nil {
		return nil
	}
	return &entity.Todo{
		Title:   pbResp.Todo.Title,
		Content: pbResp.Todo.Content,
		Status:  entity.TodoStatus(pbResp.Todo.Status),
	}
}
