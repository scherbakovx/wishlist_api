package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/scherbakovx/wishlist_api/app/db"
	"github.com/scherbakovx/wishlist_api/app/models"

	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Panic("failed to open .env file")
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/wishes", wishes)
	e.GET("/update_user_status", updateUserStatus)
	e.Logger.Fatal(e.Start(":3000"))
}

// e.GET("/wishes", wishes)
func wishes(c echo.Context) error {
	// var database *gorm.DB = db.Init()
	// user_tg_id := c.QueryParam("user_tg_id")

	// var wishes []models.LocalWish
	// result := database.Clauses(clause.OnConflict{DoNothing: true}).Model(&models.LocalWish{}).Joins("JOIN users ON local_wishes.user_id = users.id").Where("users.chat_id = ?", user_tg_id).Find(&wishes)
	// if result.Error != nil {
	// 	if result.Error.Error() == "record not found" {
	// 		return c.String(http.StatusNotFound, "404 not found")
	// 	} else {
	// 		return c.String(http.StatusBadRequest, result.Error.Error())
	// 	}
	// } else {
	// 	return c.String(http.StatusOK, fmt.Sprint(wishes[0]))
	// }
	return c.String(http.StatusOK, "hello!")
}

func updateUserStatus(c echo.Context) error {

	user_id := c.QueryParam("user_id")
	status := c.QueryParam("status")

	var database *gorm.DB = db.Init()
	db.GetOrCreateUserInDB(database, user_id)

	intUserStatus, err := strconv.Atoi(status)
	if err != nil {
		panic(err)
	}

	var answer string
	if intUserStatus == int(models.Writer) {
		answer = "Your status is Writer — just send me link and I'll add it to your wishlist!"
	} else {
		answer = "Your status is Reader — just send me contact card or nickname and I'll be ready to give you advice :)"
	}

	var user models.User
	user, err = db.GetOrCreateUserInDB(database, user_id)
	if err != nil {
		panic(err)
	}

	user.Status = models.UserStatus(intUserStatus)
	database.Save(user)

	return c.String(http.StatusOK, answer)
}
