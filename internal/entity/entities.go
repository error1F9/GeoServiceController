package entity

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
	House  string `json:"house"`
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
