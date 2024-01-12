package storage

type Storage interface {
	DoesEmailExist(email string) error
	CreateUser(email, password string) (string, error)
	UserByEmail(email string) (interface{}, error)
}
