package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"strconv"
)

// Por começar com uma letra maiúscula a struct 'Book' pode ser exportada, assim como seus atributos.
type Book struct {
	ID int `json: "id"` // Serve para indicarmos como esse campo vai ser nomeado no JSON
	Title string `json: "title"`
	Author string `json: "author"`
	Quantity int `json: "quantity"`
}

var books = []Book{
	{ID: 1, Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: 3, Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func checkoutBook(c * gin.Context){
	idStr, ok := c.GetQuery("id")	
	id, err := strconv.Atoi(idStr)

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter!"})
		return
	}
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return 
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id int) (*Book, error){ // O '*' indica que o ponteiro retornado será do tipo Book
	for i, book := range books{
		if book.ID == id {
			return &books[i], nil // O & serve para obter o endereço de memória de um determinado elemento
			// o ,nil é o segundo retorno da função e tá falando que não houve erros ao encontrar o livro
		}
	}
	return nil, errors.New("Book not found!")
}

// O gin.Context é um ponteiro que representa uma estrutura de uma requisição HTTP, contendo todas as informações
// relacionadas a uma solicitação e métodos de resposta
func getBooks(c * gin.Context) { 
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c * gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	
	c.IndentedJSON(http.StatusOK, book)
}


func createBook(c * gin.Context){
	var newBook Book 

	// O bindJSON tenta associar os dados enviados na requisição ao newBook, se caso os dados não baterem
	// ele retorna um erro
	if err := c.BindJSON(&newBook); err != nil {
		return 
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func returnBook(c * gin.Context){
	idStr, ok := c.GetQuery("id")	
	id, err := strconv.Atoi(idStr)

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter!"})
		return
	}
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return 
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default() 
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}