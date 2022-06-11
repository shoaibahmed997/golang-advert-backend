package models

type User struct {
	ID        string
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname"`
	Email     string `json:"email" validate:"required,min=5"`
	Password  string `json:"password" validate:"required,min=7"`
	CreatedAt string
}

type Image struct {
	ID      string
	Imgpath string
	Postid  int
}

type Post struct {
	ID        string
	Title     string `json:"title" validate:"required"`
	Desc      string `json:"desc"`
	Price     uint64 `json:"price" validate:"required"`
	Category  string `json:"category" validate:"required"`
	Location  string `json:"location" validate:"required"`
	UserEmail string
	CreatedAt string
	Images    []Image
}
