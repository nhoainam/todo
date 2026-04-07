package input

type GetTodoInput struct {
	Name string
}

type CreateTodoInput struct {
	Title       string
	Description string
}

type UpdateTodoInput struct {
	Name        string
	Title       *string
	Description *string
	Status      *string
}

type ListTodosInput struct {
	ListName  string
	PageSize  int32
	PageToken string
}
