package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vky5/mailcat/internal/db/models"
)

type MailClient struct {
	BaseURL string
}

func NewMailClient(baseURL string) *MailClient {
	return &MailClient{BaseURL: baseURL}
}

func (m *MailClient) Stream(email, mailbox string, pageSize, pageNumber int,
	onBatch func([]models.Email),
	onNew func(models.Email)) error {

	url := fmt.Sprintf("%s/mail/stream", m.BaseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var currentEvent string
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			switch currentEvent {
			case "batch":
				var emails []models.Email
				if err := json.Unmarshal([]byte(data), &emails); err != nil {
					fmt.Println("Error decoding batch:", err)
					continue
				}
				onBatch(emails)

			case "email":
				var mail models.Email
				if err := json.Unmarshal([]byte(data), &mail); err != nil {
					fmt.Println("Error decoding email:", err)
					continue
				}
				onNew(mail)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream read error: %w", err)
	}
	return nil
}
