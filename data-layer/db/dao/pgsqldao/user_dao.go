package pgsqldao

import (
	"context"
	"errors"

	"github.com/lucas-10101/auth-service/data-layer/db"
	"github.com/lucas-10101/auth-service/data-layer/entities"
)

type UserDao struct {
	connection db.ConnectionWrapper
	context    context.Context
}

func (dao *UserDao) Open(connectionWrapper db.ConnectionWrapper, context context.Context) error {
	if dao.connection != nil || dao.context != nil {
		return errors.New("already open")
	}

	dao.connection = connectionWrapper
	dao.context = context
	return nil
}

func (dao *UserDao) FindById(userId string) (*entities.User, error) {
	dest := &entities.User{}
	err := dao.connection.QueryRowContext(dao.context, `
		SELECT
			ID,
			USERNAME,
			PASSWORD
		FROM
			USER
		WHERE
			ID = $1
	`, userId).Scan(&dest.Id, &dest.Username, &dest.Password)
	return dest, err
}

func (dao *UserDao) ListAll() ([]*entities.User, error) {
	rows, err := dao.connection.QueryContext(dao.context, `
		SELECT
			ID,
			USERNAME,
			PASSWORD
		FROM
			USER
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := []*entities.User{}
	for rows.Next() {
		dest := &entities.User{}
		if err = rows.Scan(&dest.Id, &dest.Username, &dest.Password); err != nil {
			return nil, err
		}
		result = append(result, dest)
	}

	return result, nil
}

func (dao *UserDao) Update(user entities.User) error {
	_, err := dao.connection.ExecContext(dao.context, `
		UPDATE USER SET
			USERNAME = $2,
			PASSWORD = $3
		WHERE ID = $1
	`, user.Id, user.Username, user.Password)

	return err
}
