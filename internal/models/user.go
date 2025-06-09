package models

type User struct {
	Id       int64  `json:"id,omitempty" bun:",pk,autoincrement"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
