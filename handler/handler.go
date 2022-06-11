package handler

import (
	"fmt"
	"go-ecom/config"
	"go-ecom/database"
	"go-ecom/helper"
	"go-ecom/models"
	"log"
	"strconv"

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
	// validate json with user
	validateError := helper.Validator.Struct(user)
	if validateError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validateError.Error()})
	}

	// find in database

	rows, resError := database.DB.Query("SELECT password FROM users WHERE email=?;", user.Email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	var resultPassword string
	for rows.Next() {
		rowerr := rows.Scan(&resultPassword)
		if rowerr != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": rowerr.Error()})
		}

	}
	rows.Close()

	if resultPassword != user.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "password do not match"})
	}

	tok := jwt.New(jwt.SigningMethodHS256)
	claims := tok.Claims.(jwt.MapClaims)
	claims["email"] = user.Email

	token, tokenerr := tok.SignedString([]byte(config.Config("SECRETKEY")))
	if tokenerr != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": tokenerr.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "token": token})
}

func Signup(c *fiber.Ctx) error {
	var user models.User
	parseError := c.BodyParser(&user)
	if parseError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": parseError.Error()})
	}
	// validate json with user
	validateError := helper.Validator.Struct(user)
	if validateError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validateError.Error()})
	}
	id := uuid.Must(uuid.NewRandom())

	stmt, _ := database.DB.Prepare("INSERT INTO users (ID,firstname,lastname,email,password) VALUES (?,?,?,?,?);")
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

	return c.JSON(fiber.Map{"success": true, "token": token})
}

type newpassword struct {
	Password string `json:"password" validate:"required,min=7"`
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

	stmt, stmtError := database.DB.Prepare("UPDATE users SET password=? WHERE email=?")
	if stmtError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtError.Error(), "error type": "dbStatementError"})
	}

	_, resError := stmt.Exec(password.Password, email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error(), "error type": "dbExecutionError"})
	}

	return c.JSON(fiber.Map{"success": true, "Data": "Password Changed Successfully"})

}

func DeleteUser(c *fiber.Ctx) error {

	email := c.Locals("email")
	stmt, stmtError := database.DB.Prepare("DELETE FROM users WHERE email=?")
	if stmtError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtError.Error(), "error type": "dbStatementError"})
	}
	_, resError := stmt.Exec(email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error(), "error type": "dbExecutionError"})
	}

	return c.JSON(fiber.Map{"success": true, "Data": "User Deleted Successfully"})
}

// -------------------------POST API FUNCTIONS------------------------------------------------------->

func GetAllPost(c *fiber.Ctx) error {
	// get everything from database and give it to user
	rows, rowsError := database.DB.Query("SELECT * FROM posts")
	if rowsError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": rowsError.Error()})
	}

	var Allposts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.UserEmail, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		// get images from images table  and add to post then add the post to all post
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=?", post.ID)
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
	rows, rowsError := database.DB.Query("SELECT * FROM posts WHERE Category=?", category)
	if rowsError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": rowsError.Error()})
	}

	var Allposts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.UserEmail, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		// get images from images table  and add to post then add the post to all post
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=?", post.ID)
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

func UpdatePost(c *fiber.Ctx) error {
	// post id and email from token
	postID := c.Params("id")
	email := c.Locals("email")
	var post = new(models.Post)
	parserError := c.BodyParser(post)
	if parserError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": parserError.Error()})
	}

	validationError := helper.Validator.Struct(post)
	if validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validationError.Error()})
	}

	stmt, stmtErr := database.DB.Prepare("UPDATE posts SET Title=?, Desc=?,Price=?,Category=?,Location=?  WHERE ID=? AND userEmail=?")
	if stmtErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtErr.Error()})
	}
	res, resError := stmt.Exec(post.Title, post.Desc, post.Price, post.Category, post.Location, postID, email)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "data": res})
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	email := c.Locals("email")

	stmt, stmtErr := database.DB.Prepare("DELETE FROM posts WHERE ID=? AND userEmail=? ;")

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
		err := rows.Scan(&post.ID, &post.Title, &post.Desc, &post.Price, &post.Category, &post.Location, &post.UserEmail, &post.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		// get images from images table  and add to post then add the post to all post
		imgRows, imgRowsError := database.DB.Query("SELECT imgpath FROM images WHERE PostId=?", post.ID)
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

func CreatePost(c *fiber.Ctx) error {
	form, formError := c.MultipartForm()
	if formError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": formError.Error(), "error type": "form parsing error"})
	}
	files := form.File["images"]

	var post models.Post

	id := uuid.Must(uuid.NewRandom())

	for _, file := range files {
		if file.Size > 2000000 { // 2mb
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "File size Greater than 2mb", "error type": "form parsing error"})
		}
		log.Println(file.Header["Content-Type"][0])
		c.SaveFile(file, fmt.Sprint("./public/", id, file.Filename))
	}

	priceconvert, priceConvertError := strconv.ParseUint(form.Value["price"][0], 10, 64)
	if priceConvertError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": priceConvertError.Error(), "error type": "Error parsing and converting price"})
	}

	post.Title = form.Value["title"][0]
	post.Desc = form.Value["desc"][0]
	post.Price = priceconvert
	post.Category = form.Value["category"][0]
	post.Location = form.Value["location"][0]
	post.UserEmail = fmt.Sprint(c.Locals("email"))

	validationError := helper.Validator.Struct(&post)
	if validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": validationError.Error()})
	}

	stmt, stmtError := database.DB.Prepare("INSERT INTO posts (ID,title,desc,price,category,location,userEmail) VALUES (?,?,?,?,?,?,?);")
	if stmtError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": stmtError.Error()})
	}

	res, resError := stmt.Exec(id, post.Title, post.Desc, post.Price, post.Category, post.Location, post.UserEmail)
	if resError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": resError.Error()})
	}

	log.Println(&res)

	stmt2, stmt2Error := database.DB.Prepare("INSERT INTO images (ID,imgpath,PostID) VALUES (?,?,?);")
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
