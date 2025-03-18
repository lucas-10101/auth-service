package repositories

import (
	"context"

	"github.com/lucas-10101/auth-service/data-layer/db"
	"github.com/lucas-10101/auth-service/data-layer/db/dao/pgsqldao"
	"github.com/lucas-10101/auth-service/data-layer/entities"
)

type RealmRepository interface {
	Open(connectionWrapper db.ConnectionWrapper, context context.Context) error
	FindById(realmId string) (*entities.Realm, error)
	ListAll() ([]*entities.Realm, error)
	Update(realm entities.Realm) error
}

func NewRealmRepository(connectionWrapper db.ConnectionWrapper, context context.Context) RealmRepository {
	var resource RealmRepository = &pgsqldao.RealmDao{}
	resource.Open(connectionWrapper, context)
	return resource
}
