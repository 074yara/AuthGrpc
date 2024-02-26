package entities

type User struct {
	ID       uint
	Email    string
	PassHash []byte
}
