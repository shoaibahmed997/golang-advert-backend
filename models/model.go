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
	Title     string  `form:"title" validate:"required,max=20"`
	Desc      string  `form:"desc"`
	Price     float64 `form:"price" validate:"required"`
	Category  string  `form:"category" validate:"required,max=20"`
	Location  string  `form:"location" validate:"required"`
	Lattitude float64 `form:"lattitude" validate:"required"`
	Longitude float64 `form:"longitude" validate:"required"`
	UserEmail string
	By        string
	CreatedAt string
	Images    []Image
}
