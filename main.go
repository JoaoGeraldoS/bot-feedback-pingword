package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/JoaoGeraldoS/pingword-feedbacks/models"
	"github.com/JoaoGeraldoS/pingword-feedbacks/service"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func webhook(w http.ResponseWriter, r *http.Request) {

	apiUrl := os.Getenv("API_URL")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	var update models.Update
	json.Unmarshal(body, &update)

	service.ListFeedBacks(update, apiUrl)
	service.DeleteFeedBack(update, apiUrl)

	w.WriteHeader(http.StatusOK)
}

func main() {

	http.HandleFunc("POST /webhook", webhook)
	fmt.Println("Servidor iniciado na porta 8080...")

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Println("Acordado")
		fmt.Fprint(w, "Acordado")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor a rodar na porta %s...\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
