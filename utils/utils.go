package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	WebhookContent     = "[배포 완료]\n" + "배포 목록 리스트 입니다."
	webhookLinkContent = "배포 내용 확인하기"
)

func WebhookMessage(jql string) {
	data, err := json.Marshal(map[string]string{
		"content":      WebhookContent,
		"link_content": webhookLinkContent,
		"link":         os.Getenv("JIRA_URL") + "issues/?jql=" + jql,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	if resp.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("Error: %s", string(body)))
	}
}
