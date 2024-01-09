package model

type Book struct {
	Isbn string `json:"isbn,omitempty" gorm:"primaryKey;autoIncrement:true"`
	Name string `json:"name"`
}
