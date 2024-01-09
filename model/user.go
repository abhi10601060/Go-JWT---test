package model

type User struct {
	Username string `json:"username,omitempty" gorm:"primarykey;"`
	Password string `json:"password,omitempty"`
}
