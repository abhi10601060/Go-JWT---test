package repo

import (
	"apiwithmysql/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dsn = "root:10601060@tcp(localhost:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB

func init() {
	db, _ = connectToDB()
}

func connectToDB() (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error in connectDb : ", err)
	}
	autoMigrate(db)
	return db, nil
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Book{})
	db.AutoMigrate(&model.User{})
}

//all CRUD operations

func Addbook(book *model.Book) int {
	res := db.Create(book)
	if res.Error != nil {
		log.Fatal(res.Error)
		return 0
	}
	return 1
}

func GetAllbooks(books *[]model.Book) {
	res := db.Find(books)

	if res.Error != nil {
		log.Fatal("Error in getting all books : ", res.Error)
	}
}

func DeleteBookByIsbn(isbn string) {
	var book model.Book
	book.Isbn = isbn
	res := db.Delete(&book)

	if res.Error != nil {
		log.Fatal("Error in deleting book : ", res.Error)
	}

	fmt.Println("delete book called in repo with isbn : ", isbn)
}

func GetBookByIsbn(isbn string, b *model.Book) (int, error) {
	fmt.Println("get book in repo called")
	res := db.First(b, isbn)
	if res.Error != nil {
		return 0, res.Error
	}
	fmt.Println("one book received")
	return 1, nil
}

func UpdateBook(b *model.Book) {
	res := db.Save(b)

	if res.Error != nil {
		log.Fatal(res.Error)
	}
}
func IsUserPresent(u *model.User) bool {
	var ru = model.User{Username: u.Username}
	res := db.First(&ru)

	if res.Error != nil {
		return false
	}
	return true
}
func SaveUser(u *model.User) bool {
	res := db.Create(u)
	if res.Error != nil {
		return false
	}
	return true
}

func GetUser(u *model.User) (*model.User, bool) {
	var ru = model.User{Username: u.Username}

	res := db.First(&ru)
	if res.Error != nil {
		return &ru, false
	}
	return &ru, true
}
