package input

type LoginInput struct {
	Username string
	Password string
}

type LogoutInput struct {
	UserID int64
}

type RegisterInput struct {
	Username string
	Password string
}
