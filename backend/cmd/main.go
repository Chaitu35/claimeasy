package main


import("fmt"
	
	"os"
	"log"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(".env.dev"); err != nil {
        log.Fatal("Error loading .env.dev file")
    }

    // Retrieve the GitHub token from the environment
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        log.Fatal("GITHUB_TOKEN not set in .env.dev")
    }

    fmt.Println("GitHub token loaded:", token)
}