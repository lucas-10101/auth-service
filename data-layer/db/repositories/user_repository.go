package repositories

import (
	"context"

	"github.com/lucas-10101/auth-service/data-layer/db"
	"github.com/lucas-10101/auth-service/data-layer/db/dao/pgsqldao"
	"github.com/lucas-10101/auth-service/data-layer/entities"
)

type UserRepository interface {
	Open(connectionWrapper db.ConnectionWrapper, context context.Context) error
	FindById(userId string) (*entities.User, error)
	ListAll() ([]*entities.User, error)
	Update(user entities.User) error
}

func NewUserRepository(connectionWrapper db.ConnectionWrapper, context context.Context) UserRepository {
	var resource UserRepository = &pgsqldao.UserDao{}
	resource.Open(connectionWrapper, context)
	return resource
}
