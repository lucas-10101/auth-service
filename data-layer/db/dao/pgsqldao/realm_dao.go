package pgsqldao

import (
	"context"
	"errors"

	"github.com/lucas-10101/auth-service/data-layer/db"
	"github.com/lucas-10101/auth-service/data-layer/entities"
)

type RealmDao struct {
	connection db.ConnectionWrapper
	context    context.Context
}

func (dao *RealmDao) Open(connectionWrapper db.ConnectionWrapper, context context.Context) error {
	if dao.connection != nil || dao.context != nil {
		return errors.New("already open")
	}

	dao.connection = connectionWrapper
	dao.context = context
	return nil
}

func (dao *RealmDao) FindById(realmId string) (*entities.Realm, error) {
	result := &entities.Realm{}
	err := dao.connection.QueryRowContext(dao.context, `
		SELECT
			ID,
			NAME
		FROM
			REALM
		WHERE
			ID = $1
		FETCH NEXT
			1 ROWS ONLY
	`, realmId).Scan(&result.Id, &result.Name)

	return result, err
}

func (dao *RealmDao) ListAll() ([]*entities.Realm, error) {
	rows, err := dao.connection.QueryContext(dao.context, `
		SELECT
			ID,
			NAME
		FROM
			REALM
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := []*entities.Realm{}
	for rows.Next() {
		dest := entities.Realm{}
		if err = rows.Scan(&dest.Id, &dest.Name); err != nil {
			return nil, err
		}
		result = append(result, &dest)
	}

	return result, err
}

func (dao *RealmDao) Update(realm entities.Realm) error {
	_, err := dao.connection.ExecContext(dao.context, `
		UPDATE REALM SET
			NAME = $2
		WHERE ID = $1
	`, &realm.Id, &realm.Name)

	if err != nil {
		return err
	}

	return nil
}
