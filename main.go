package main
 
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
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

func getbookByID(c *gin.Context){
	id := c.Param("id")
	book, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"the request  book not found "})
	}
    c.IndentedJSON(http.StatusOK, book)

}

func getBookByID(id string)(*book ,error){
   for i, b := range books{
	if b.ID == id {
		return &books[i],nil
	}
   }

   return nil, errors.New("book not Fount")
}

func getBooks( c *gin.Context){
	c.IndentedJSON(http.StatusAccepted, books)
}

func checkout( c *gin.Context){
	id , ok:= c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest , gin.H{"message" : "Missing id query parameter"})
		return
	}

	book , err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": " Books not avaible."})
		return
	}

	if book.Quantity <= 0{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": " Book not avaible for the checkout"})
		return
	}

	book.Quantity -=1
	c.IndentedJSON(http.StatusOK, book)
}

func createBook(c *gin.Context){
	var newBook book

	if err := c.BindJSON(&newBook); err != nil{
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}
func main(){
	router := gin.Default()
	router.GET("/books",getBooks)
	router.POST("/createBook",createBook)
	router.GET("/books/:id", getbookByID)
	router.PATCH("/checkout", checkout)
	router.Run("localhost:8080")
}