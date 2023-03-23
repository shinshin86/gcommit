package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const URL string = "https://api.openai.com/v1/chat/completions"

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func chat(messages []message, apiKey string) string {
	reqBody := chatRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}
	reqJSON, err := json.Marshal(reqBody)

	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(reqJSON))

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	var resBody chatResponse
	json.NewDecoder(res.Body).Decode(&resBody)

	message := resBody.Choices[0].Message

	return message.Content
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatalln("API key must be defined as an environment variable.")
	}

	cmd := exec.Command("git", "diff")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running git diff: %v\n", err)
	}

	var messages []message
	systemMsg := message{Role: "system", Content: "You are a specialist in thinking about git commit messages. Please generate a commit message based on the results of the git diff I'm about to send you. Please return only the commit message. Also, be sure to prefix the commit message with a prefix, e.g. feat: add text function."}
	messages = append(messages, systemMsg)

	diff := strings.TrimSpace(string(output))
	diffMsg := message{Role: "user", Content: diff}
	messages = append(messages, diffMsg)

	fmt.Println(chat(messages, apiKey))
}
