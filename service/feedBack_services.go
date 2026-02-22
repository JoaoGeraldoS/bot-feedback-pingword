package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JoaoGeraldoS/pingword-feedbacks/models"
)

var memoriaMensagens = make(map[string][]int)

func sendMessage(chatID int64, text string) int {
	token := os.Getenv("TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err == nil {
		defer resp.Body.Close()
		var result struct {
			Result struct {
				MessageID int `json:"message_id"`
			} `json:"result"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		return result.Result.MessageID
	}
	return 0
}

func ListFeedBacks(update models.Update, apiUrl string) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	url := fmt.Sprintf("%s/api/FeedBack", apiUrl)

	if text == "/list" {

		go func() {
			time.Sleep(10 * time.Minute)
			DeleteUserMessage(update)
		}()

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
			id1 := sendMessage(chatID, fmt.Sprintf("Usuario: %s", t.User))
			id2 := sendMessage(chatID, fmt.Sprintf("ID: %s", t.ID))
			id3 := sendMessage(chatID, fmt.Sprintf("Mensagem: %s", t.Message))

			memoriaMensagens[t.ID] = []int{id1, id2, id3}

		}
	}
}

func DeleteFeedBack(update models.Update, apiUrl string) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	if strings.HasPrefix(text, "/del ") {

		go func() {
			time.Sleep(10 * time.Second)
			DeleteUserMessage(update)
		}()

		idParaDeletar := strings.TrimPrefix(text, "/del ")

		apiURL := fmt.Sprintf("%s/api/FeedBack/%s", apiUrl, idParaDeletar)

		if ids, ok := memoriaMensagens[idParaDeletar]; ok {
			for _, mID := range ids {
				deleteSpecificMessage(chatID, mID)
			}

			delete(memoriaMensagens, idParaDeletar)
		}

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

			confirmacaoMsg := fmt.Sprintf("Tarefa %s eliminada com sucesso! ✅", idParaDeletar)
			msgID := sendMessage(chatID, confirmacaoMsg)

			go func(cID int64, mID int) {
				time.Sleep(10 * time.Second)
				deleteSpecificMessage(cID, mID)
			}(chatID, msgID)

		} else {
			sendMessage(chatID, "Não foi possível eliminar a tarefa. Verifica o ID.")
		}
	}
}

func DeleteUserMessage(update models.Update) {
	token := os.Getenv("TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessage", token)

	payload := map[string]interface{}{
		"chat_id":    update.Message.Chat.ID,
		"message_id": update.Message.MessageID,
	}

	jsonData, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

func deleteSpecificMessage(chatID int64, messageID int) {
	if messageID == 0 {
		return
	}
	token := os.Getenv("TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessage", token)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}

	jsonData, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}
