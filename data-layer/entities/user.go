package entities

type User struct {
	Id       string `db:"ID"`
	Username string `db:"USERNAME"`
	Password string `db:"PASSWORD"`
}
