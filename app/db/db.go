package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/scherbakovx/wishlist_api/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Init() *gorm.DB {
	pg_username := os.Getenv("POSTGRES_USER")
	pg_password := os.Getenv("POSTGRES_PASSWORD")
	pg_db := os.Getenv("POSTGRES_DB")

	dbURL := fmt.Sprintf("postgres://%s:%s@db:5432/%s", pg_username, pg_password, pg_db)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	// db.AutoMigrate(&models.LocalWish{})
	// db.AutoMigrate(&models.AirTableConnection{})
	// db.AutoMigrate(&models.User{})

	return db
}

func GetOrCreateUserInDB(db *gorm.DB, chatId string) (*models.User, error) {
	intChatId, _ := strconv.Atoi(chatId)
	user := models.User{ChatId: int64(intChatId)}
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		db.First(&user, "chat_id = ?", user.ChatId)
	}

	return &user, nil
}
