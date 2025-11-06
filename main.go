package main
import  (
	"fmt"
	"log"
	"os"
	"net/http"

	"github.com/joho/godotenv"
	// "github.com/jackc/pgx/v5"
	// "github.com/gorilla/mux"
)

func httpHandler(writer http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	fmt.Fprintf(writer, "Hello World, Welcome to Go, The requested URL path is %s",request.URL.Path)
}

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables...")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port ="3001"
		fmt.Println("Add port to env file")
	}
	http.HandleFunc("/getUser", getUserHandler)
	http.HandleFunc("/getBook", GetBookHandler)
	http.HandleFunc("/", httpHandler)
	fmt.Println("Server running at localhost ", port)
	if  err:= http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}