package initializer

import (
	"kubequntumblock/models"
     "fmt"
)
func SyncDatabase(){
	DB.AutoMigrate(&models.User{})
	fmt.Println("Database synchronized successfully!")
}