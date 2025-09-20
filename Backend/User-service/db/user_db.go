package db

import (
"log"


"gorm.io/driver/postgres"
"gorm.io/gorm"
"user-service/models"
)


var DB *gorm.DB


func InitDB(dsn string) {
var err error
DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
log.Fatalf("failed to connect db: %v", err)
}
DB.AutoMigrate(&models.User{})
}