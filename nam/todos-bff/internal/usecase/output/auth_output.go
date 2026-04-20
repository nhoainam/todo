package output

import "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"

type LoginOutput struct {
	User        *entity.User
	AccessToken string
}

type LogoutOutput struct{}

type RegisterOutput struct {
	User *entity.User
}
