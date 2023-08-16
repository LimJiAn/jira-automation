package issue

import (
	"fmt"
	"log"

	"github.com/LimJiAn/jira-automation/utils"
	"github.com/andygrunwald/go-jira"
)

func (c *Client) CheckDeployIssues() {
	// your jql here
	jql := ""
	transitionID := ""

	issues := c.GetIssues(jql)
	for _, v := range issues {
		transitions, _, _ := c.Issue.GetTransitions(v.Key)
		for _, t := range transitions {
			if t.Name == "Release" {
				transitionID = t.ID
				break
			}
		}

		if transitionID != "" {
			c.Issue.DoTransition(v.Key, transitionID)
			issue, _, err := c.Issue.Get(v.Key, nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Issue: %v | %v ==>  %s\n", issue.Key, issue.Fields.Summary, issue.Fields.Status.Name)
		}
	}

	if len(issues) > 0 {
		utils.WebhookMessage(jql)
	}
}

func (c *Client) GetIssues(jql string) []jira.Issue {
	issues, _, err := c.Issue.Search(jql, nil)
	if err != nil {
		log.Fatal(err)
	}

	result := make([]jira.Issue, len(issues))
	if issues != nil {
		for i, v := range issues {
			result[i] = v
		}
	} else {
		result = make([]jira.Issue, 0)
	}
	return result
}
