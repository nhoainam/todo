package output

type LoginOutput struct {
	UserID      int64
	Username    string
	AccessToken string
}

type LogoutOutput struct {
}

type RegisterOutput struct {
	UserID   int64
	Username string
}
