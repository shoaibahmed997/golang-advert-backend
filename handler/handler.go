package handler

import (
	"errors"
	"fmt"
	"go-ecom/config"
	"go-ecom/database"
	"go-ecom/helper"
	"go-ecom/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// -------------------------USER API FUNCTIONS------------------------------------------------------->
type loginuser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c *fiber.Ctx) error {
	var user = new(loginuser)
	parserErr := c.BodyParser(user)
	if parserErr != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": parserErr.Error()})
	}
	user.Email = strings.ToLower(user.Email)

	// validate json with user
	validateError := helper.Validator.Struct(user)
	if validateError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validateError.Error()})
	}

	// find in database

	rows, resError := database.DB.Query("SELECT ID,firstname,password FROM users WHERE email=$1;", user.Email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	var resultUser struct {
		ID        string
		FirstName string
		Password  string
	}
	for rows.Next() {
		rowerr := rows.Scan(&resultUser.ID, &resultUser.FirstName, &resultUser.Password)
		if rowerr != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": rowerr.Error()})
		}

	}
	rows.Close()

	if resultUser.Password != user.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Email or Password is Invalid"})
	}

	tok := jwt.New(jwt.SigningMethodHS256)
	claims := tok.Claims.(jwt.MapClaims)
	claims["email"] = user.Email

	token, tokenerr := tok.SignedString([]byte(config.Config("SECRETKEY")))
	if tokenerr != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": tokenerr.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "token": token, "user": fiber.Map{"id": resultUser.ID, "firstname": resultUser.FirstName}})
}

func Signup(c *fiber.Ctx) error {
	var user models.User
	parseError := c.BodyParser(&user)
	if parseError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": parseError.Error()})
	}
	user.Email = strings.ToLower(user.Email)
	// validate json with user
	validateError := helper.Validator.Struct(user)
	if validateError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validateError.Error()})
	}
	id := uuid.Must(uuid.NewRandom())

	stmt, _ := database.DB.Prepare("INSERT INTO users (ID,firstname,lastname,email,password) VALUES ($1,$2,$3,$4,$5);")
	_, resError := stmt.Exec(id, user.FirstName, user.LastName, user.Email, user.Password)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	tok := jwt.New(jwt.SigningMethodHS256)
	claims := tok.Claims.(jwt.MapClaims)
	claims["email"] = user.Email

	token, tokenerr := tok.SignedString([]byte(config.Config("SECRETKEY")))
	if tokenerr != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": tokenerr.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "token": token, "user": fiber.Map{"id": user.ID, "firstname": user.FirstName}})
}

type newpassword struct {
	Newpassword     string `json:"new_password" validate:"required,min=7"`
	Currentpassword string `json:"current_password" validate:"required"`
}

func ChangePassword(c *fiber.Ctx) error {
	email := c.Locals("email")

	var password newpassword

	parserError := c.BodyParser(&password)
	if parserError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": parserError.Error(), "error type": "parserError"})
	}

	validationError := helper.Validator.Struct(&password)
	if validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validationError.Error(), "error type": "validationError"})
	}

	//  check for previous password

	prevPassRow := database.DB.QueryRow("SELECT password from users WHERE email=$1", email)
	var previousPass string
	prevError := prevPassRow.Scan(&previousPass)
	if prevError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": prevError.Error(), "error type": "dbStatementError"})
	}

	if previousPass != password.Currentpassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": errors.New("password do not match"), "error type": "Password do not Match"})
	}

	stmt, stmtError := database.DB.Prepare("UPDATE users SET password=$1 WHERE email=$2;")
	if stmtError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtError.Error(), "error type": "dbStatementError"})
	}

	_, resError := stmt.Exec(password.Newpassword, email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error(), "error type": "dbExecutionError"})
	}

	return c.JSON(fiber.Map{"success": true, "data": "Password Changed Successfully"})

}

func DeleteUser(c *fiber.Ctx) error {

	email := c.Locals("email")
	stmt, stmtError := database.DB.Prepare("DELETE FROM users WHERE email=$1;")
	if stmtError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtError.Error(), "error type": "dbStatementError"})
	}
	_, resError := stmt.Exec(email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error(), "error type": "dbExecutionError"})
	}

	return c.JSON(fiber.Map{"success": true, "data": "User Deleted Successfully"})
}

// -------------------------POST API FUNCTIONS------------------------------------------------------->

func GetAllPost(c *fiber.Ctx) error {
	// get everything from database and give it to user
	rows, rowsError := database.DB.Query("SELECT * FROM posts ORDER BY CreatedAt DESC;")
	if rowsError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": rowsError.Error()})
	}

	var Allposts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.Lattitude, &post.Longitude, &post.UserEmail, &post.By, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		// get images from images table  and add to post then add the post to all post
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=$1", post.ID)
		if imgRowsError != nil {
			return c.Status(503).JSON(fiber.Map{"success": false, "error": imgRowsError.Error()})
		}
		for imgRows.Next() {
			var img models.Image
			err := imgRows.Scan(&img.Imgpath)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
			}
			post.Images = append(post.Images, img)
		}
		imgRows.Close()
		Allposts = append(Allposts, post)
	}

	rows.Close()
	return c.Status(200).JSON(fiber.Map{"success": true, "data": Allposts})
}
func GetPostByCategory(c *fiber.Ctx) error {
	// get everything from database and give it to user
	category := c.Params("category")
	rows, rowsError := database.DB.Query("SELECT * FROM posts WHERE Category=$1 ORDER BY CreatedAt DESC;", category)
	if rowsError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": rowsError.Error()})
	}

	var Allposts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.Lattitude, &post.Longitude, &post.UserEmail, &post.By, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		// get images from images table  and add to post then add the post to all post
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=$1", post.ID)
		if imgRowsError != nil {
			return c.Status(503).JSON(fiber.Map{"success": false, "error": imgRowsError.Error()})
		}
		for imgRows.Next() {
			var img models.Image
			err := imgRows.Scan(&img.Imgpath)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
			}
			post.Images = append(post.Images, img)
		}
		imgRows.Close()
		Allposts = append(Allposts, post)
	}
	rows.Close()

	return c.Status(200).JSON(fiber.Map{"success": true, "data": Allposts})
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	email := c.Locals("email")

	stmt, stmtErr := database.DB.Prepare("DELETE FROM posts WHERE ID=$1 AND userEmail=$2 ;")

	if stmtErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtErr.Error()})
	}

	_, resError := stmt.Exec(postID, email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "data": "Deleted Successfully"})
}

func GetFirst20post(c *fiber.Ctx) error {

	var Allposts []models.Post

	rows, rowsError := database.DB.Query("SELECT * FROM posts LIMIT 20")

	if rowsError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": rowsError.Error()})
	}

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.UserEmail, &post.Lattitude, &post.Longitude, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		// get images from images table  and add to post then add the post to all post
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=$1", post.ID)
		if imgRowsError != nil {
			return c.Status(503).JSON(fiber.Map{"success": false, "error": imgRowsError.Error()})
		}
		for imgRows.Next() {
			var img models.Image
			err := imgRows.Scan(&img.Imgpath)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
			}
			post.Images = append(post.Images, img)
		}
		imgRows.Close()

		Allposts = append(Allposts, post)
	}

	rows.Close()

	return c.Status(200).JSON(fiber.Map{"success": true, "data": Allposts})
}

type userName struct {
	Firstname string
	Lastname  string
}

func CreatePost(c *fiber.Ctx) error {
	form, formError := c.MultipartForm()
	if formError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": formError.Error(), "error type": "form parsing error"})
	}

	var post models.Post

	id := uuid.Must(uuid.NewRandom())

	priceconvert, priceConvertError := strconv.ParseFloat(form.Value["price"][0], 64)
	if priceConvertError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": priceConvertError.Error(), "error type": "Error parsing and converting price"})
	}
	lattConvert, _ := strconv.ParseFloat(form.Value["lattitude"][0], 64)
	longConvert, _ := strconv.ParseFloat(form.Value["longitude"][0], 64)

	post.Title = form.Value["title"][0]
	post.Desc = form.Value["desc"][0]
	post.Price = priceconvert
	post.Category = form.Value["category"][0]
	post.Location = form.Value["location"][0]
	post.Lattitude = lattConvert
	post.Longitude = longConvert
	post.UserEmail = fmt.Sprint(c.Locals("email"))

	// get the firstname and last name of the user by email
	var username userName
	NameRow := database.DB.QueryRow("SELECT firstname,lastname FROM users WHERE email=$1", post.UserEmail)
	NameRowError := NameRow.Scan(&username.Firstname, &username.Lastname)
	if NameRowError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": NameRowError.Error(), "error type": "No Data found error, Cannot find your name in database"})
	}
	post.By = fmt.Sprint(username.Firstname, " ", username.Lastname)

	validationError := helper.Validator.Struct(&post)
	if validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validationError.Error(), "error type": "validation error"})
	}
	files := form.File["images"]

	for _, file := range files {
		if file.Size > 2000000 { // 2mb
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "File size Greater than 2mb", "error type": "form parsing error"})
		}
		if file.Header["Content-Type"][0] != "image/png" && file.Header["Content-Type"][0] != "image/jpeg" && file.Header["Content-Type"][0] != "image/jpg" && file.Header["Content-Type"][0] != "image/webp" && file.Header["Content-Type"][0] != "image/gif" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Not an image file", "error type": "image parse error"})
		}
		c.SaveFile(file, fmt.Sprint("./public/", id, file.Filename))
	}

	stmt, stmtError := database.DB.Prepare("INSERT INTO posts (ID,title,description,price,category,location,userEmail,lattitude,longitude,By) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);")
	if stmtError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtError.Error()})
	}

	_, resError := stmt.Exec(id, post.Title, post.Desc, post.Price, post.Category, post.Location, post.UserEmail, post.Lattitude, post.Longitude, post.By)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	stmt2, stmt2Error := database.DB.Prepare("INSERT INTO images (ID,imgpath,PostID) VALUES ($1,$2,$3);")
	if stmt2Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmt2Error.Error()})
	}
	tx, txError := database.DB.Begin()

	if txError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": txError.Error()})
	}

	for _, file := range files {
		imageid := uuid.Must(uuid.NewRandom())
		var imgpath string = fmt.Sprint("/", id, file.Filename)
		_, resError2 := tx.Stmt(stmt2).Exec(imageid, imgpath, id)

		if resError2 != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError2.Error(), "error type": "Database error rolling back transaction"})
		}
	}

	tx.Commit()

	return c.Status(200).JSON(fiber.Map{"success": true, "data": "post created successfully"})
}

func SearchPost(c *fiber.Ctx) error {

	searchterm := c.Params("searchterm")
	rows, rowsError := database.DB.Query(fmt.Sprint("SELECT * FROM posts WHERE title LIKE '%", searchterm, "%' OR description LIKE '%", searchterm, "%' ;"))
	if rowsError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": rowsError.Error(), "error type": "Error While Querying database"})
	}
	var Allposts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.Lattitude, &post.Longitude, &post.UserEmail, &post.By, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error(), "error type": "Error While Querying database"})
		}
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=$1", post.ID)
		if imgRowsError != nil {
			return c.Status(503).JSON(fiber.Map{"success": false, "error": imgRowsError.Error()})
		}
		for imgRows.Next() {
			var img models.Image
			err := imgRows.Scan(&img.Imgpath)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
			}
			post.Images = append(post.Images, img)
		}
		imgRows.Close()
		Allposts = append(Allposts, post)

	}
	rows.Close()

	return c.Status(200).JSON(fiber.Map{"success": true, "data": Allposts})
}

func GetPostByUser(c *fiber.Ctx) error {
	email := c.Params("email")

	rows, rowsError := database.DB.Query("SELECT * FROM posts WHERE userEmail=$1 ORDER BY CreatedAt DESC;", email)
	if rowsError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": rowsError.Error(), "error type": "Error While Querying database"})
	}

	var Allposts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.Lattitude, &post.Longitude, &post.UserEmail, &post.By, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error(), "error type": "Error While Querying database"})
		}
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=$1", post.ID)
		if imgRowsError != nil {
			return c.Status(503).JSON(fiber.Map{"success": false, "error": imgRowsError.Error()})
		}
		for imgRows.Next() {
			var img models.Image
			imgerr := imgRows.Scan(&img.Imgpath)
			if imgerr != nil {
				return c.Status(503).JSON(fiber.Map{"success": false, "error": imgerr.Error()})
			}
			post.Images = append(post.Images, img)
		}

		imgRows.Close()
		Allposts = append(Allposts, post)
	}
	rows.Close()

	return c.Status(200).JSON(fiber.Map{"success": true, "data": Allposts})

}
