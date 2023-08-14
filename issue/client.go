package issue

import (
	"os"

	"github.com/andygrunwald/go-jira"
)

type Client struct {
	*jira.Client
}

func Authenticate() (*Client, error) {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		return nil, err
	}

	return &Client{
		client}, nil
}
