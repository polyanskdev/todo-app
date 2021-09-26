package todo

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" bindings:"required"`
	Username string `json:"username" bindings:"required"`
	Password string `json:"password" bindings:"required"`
}
