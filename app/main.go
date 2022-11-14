package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/scherbakovx/wishlist_api/app/db"
	"github.com/scherbakovx/wishlist_api/app/models"
	"github.com/scherbakovx/wishlist_bot/app/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

const successfulMessage string = "Added!"

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
	e.GET("/get_user_wish", getUserWish)
	e.GET("/add_wish_to_user", addWishToUser)
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

	var user *models.User
	user, err = db.GetOrCreateUserInDB(database, user_id)
	if err != nil {
		panic(err)
	}

	user.Status = models.UserStatus(intUserStatus)
	database.Save(user)

	return c.String(http.StatusOK, answer)
}

func getUserWish(c echo.Context) error {

	user_id := c.QueryParam("user_id")

	var database *gorm.DB = db.Init()

	var wish models.LocalWish
	result := database.Clauses(clause.OnConflict{DoNothing: true}).Model(&models.LocalWish{}).Joins("JOIN users ON local_wishes.user_id = users.id").Where("users.chat_id = ?", user_id).First(&wish)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return c.String(http.StatusOK, fmt.Sprintf("User %v has no wishes :(", user_id))
		} else {
			panic(result.Error)
		}
	} else {
		return c.String(http.StatusOK, wish.String())
	}

	// if uh.User.AirTable.Board != "" {
	// 	randomObjectData, err := airtable.GetDataFromAirTable(client, randomizer)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return c.String(http.StatusOK, randomObjectData)
	// }
}

func addWishToUser(c echo.Context) error {

	chat_id := c.QueryParam("chat_id")

	var database *gorm.DB = db.Init()

	user, err := db.GetOrCreateUserInDB(database, chat_id)
	if err != nil {
		panic(err)
	}

	link := c.QueryParam("link")

	// if uh.User.AirTable.Board != "" {
	// 	err := airtable.InsertDataToAirTable(client, link)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	openGraphData, _ := utils.GetOGTags(client, link)
	wish := models.LocalWish{
		Wish: models.Wish{
			Name: openGraphData.Title,
			Link: openGraphData.URL,
		},
		UserId: user.Id,
	}
	result := database.Create(&wish)

	if result.Error != nil {
		panic(result.Error)
	}

	return c.String(http.StatusOK, successfulMessage)

}
