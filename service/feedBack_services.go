package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/JoaoGeraldoS/pingword-feedbacks/models"
)

func sendMessage(chatID int64, text string) {
	token := os.Getenv("TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	jsonData, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

func ListFeedBacks(update models.Update, apiUrl string) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	url := fmt.Sprintf("%s/api/FeedBack", apiUrl)

	if text == "/list" {
		resp, err := http.Get(url)
		if err != nil {
			sendMessage(chatID, "Erro ao conectar com a API de tarefas.")
			return
		}
		defer resp.Body.Close()

		apiBody, err := io.ReadAll(resp.Body)
		if err != nil || len(apiBody) == 0 {
			sendMessage(chatID, "Não foi possível ler os dados da API.")
			return
		}

		var feedback []models.FeedBack
		err = json.Unmarshal(apiBody, &feedback)
		if err != nil {
			sendMessage(chatID, "Erro ao processar a lista de tarefas.")
			return
		}

		if len(feedback) == 0 {
			sendMessage(chatID, "A lista de tarefas está vazia.")
			return
		}

		for _, t := range feedback {
			user := fmt.Sprintf("Usuario: %s", t.User)
			getId := fmt.Sprintf("ID: %s", t.ID)
			msg := fmt.Sprintf("Mensagem: %s", t.Message)
			sendMessage(chatID, user)
			sendMessage(chatID, getId)
			sendMessage(chatID, msg)
		}
	}
}

func DeleteFeedBack(update models.Update, apiUrl string) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	if strings.HasPrefix(text, "/del ") {

		idParaDeletar := strings.TrimPrefix(text, "/del ")

		apiURL := fmt.Sprintf("%s/api/FeedBack/%s", apiUrl, idParaDeletar)

		req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
		if err != nil {
			sendMessage(chatID, "Erro ao criar pedido de remoção.")
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			sendMessage(chatID, "Erro ao conectar com a API para eliminar.")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
			sendMessage(chatID, fmt.Sprintf("Tarefa %s eliminada com sucesso! ✅", idParaDeletar))
		} else {
			sendMessage(chatID, "Não foi possível eliminar a tarefa. Verifica o ID.")
		}
	}
}
