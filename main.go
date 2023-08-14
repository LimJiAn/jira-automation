package main

import (
	"log"

	"github.com/LimJiAn/jira-automation/issue"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	client, err := issue.Authenticate()
	if err != nil {
		log.Fatal(err)
	}
	client.CheckDeployIssues()
}
