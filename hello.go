package main
 
import (
	"net/http"
	"github.com/gin-gonic/gin"
	// "errors"
)

type book struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}


var books = []book{
	{ID : "2", Title: "burning desire", Author: "Jane Austen", Quantity: 5},
	{ID : "3", Title: "frozen light", Author: "Fyodor Dostoevsky", Quantity: 3},
	{ID : "4", Title: "shattered dreams", Author: "Virginia Woolf", Quantity: 4},	
}

func getBooks( c *gin.Context){
	c.IndentedJSON(http.StatusAccepted, books)
}
func main(){
	router := gin.Default()
	router.GET("/books",getBooks)
	router.Run("localhost:8080")
}