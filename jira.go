package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"main.go/utils"

	jira "github.com/andygrunwald/go-jira"
)

const (
	layout = "2006-01-02"
)

var (
	options Args
)

type Args struct {
	commands *string
	confirm  *bool
}

func main() {
	options.commands = flag.String("com", "", "Change the status")
	options.confirm = flag.Bool("Y", false, "Ask for confirmation. default(N)")
	flag.Parse()

	err := godotenv.Load(".env")
	utils.CheckError(err)
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}
	jiraClient, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	utils.CheckError(err)

	switch *options.commands {
	case "preset":
		PresetIssues(jiraClient)
	case "release":
		ReleaseIssues(jiraClient)
	case "due_date":
		DueDateIssues(jiraClient)
	default:
		fmt.Println("Please enter the exact command.\n" +
			"1. go run jira.go -com=preset  (in-progress -> preset)\n" +
			"2. go run jira.go -com=release  (preset -> release)\n" +
			"3. go run jira.go -com=due_date  (check due_date)")
	}
}

func PresetIssues(jiraClient *jira.Client) {
	var (
		transitionID string
		content      []string
		jql          string
	)
	ask, deployDay := utils.AskForDeploy("When is the deploy date?")
	if !ask {
		os.Exit(0)
	}

	jql = "Write down the Jql that you want to change the preset status."

	issues := GetAllIssues(jiraClient, jql)
	if !*options.confirm {
		for _, v := range issues {
			fmt.Printf("Issue: %v | %v ==>  %s\n", v.Key, v.Fields.Summary, v.Fields.Status.Name)
		}
		c := utils.AskForConfirmation(
			"\ncount : " + strconv.Itoa(len(issues)) + "\nDo you want to change the Jira state to preset?")
		if !c {
			os.Exit(0)
		}
	}
	if len(issues) > 0 {
		for _, v := range issues {
			possibleTransitions, _, _ := jiraClient.Issue.GetTransitions(v.Key)
			for _, v := range possibleTransitions {
				if v.Name == "Preset" {
					transitionID = v.ID
					break
				}
			}
			jiraClient.Issue.DoTransition(v.Key, transitionID)
			issue, _, _ := jiraClient.Issue.Get(v.Key, nil)
			fmt.Printf("Issue: %v | %v ==>  %s\n", issue.Key, issue.Fields.Summary, issue.Fields.Status.Name)
			content = append(content,
				"\n\nIssue : ", "`"+issue.Key+"`", " |  "+issue.Fields.Summary+" |  `"+issue.Fields.Status.Name+"`",
				"\n담당자 : "+issue.Fields.Assignee.DisplayName,
				// "\n보고자 : "+issue.Fields.Reporter.DisplayName,
			)
		}
		result := strings.Join(content, " ")
		utils.PostMessage(result)
	} else {
		result := deployDay + " 에 배포 이슈가 없습니다."
		utils.PostMessage(result)
	}
}

func ReleaseIssues(jiraClient *jira.Client) {
	var (
		deployDay    string
		transitionID string
		content      []string
		jql          string
	)

	deployDay = time.Now().Format(layout)

	jql = "Write down the Jql that you want to change the release status."

	issues := GetAllIssues(jiraClient, jql)
	if !*options.confirm {
		for _, v := range issues {
			fmt.Printf("Issue: %v | %v ==>  %s\n", v.Key, v.Fields.Summary, v.Fields.Status.Name)
		}
		ask := utils.AskForConfirmation(
			"\ncount : " + strconv.Itoa(len(issues)) + "\nDo you want to change the Jira state to release?")
		if !ask {
			os.Exit(0)
		}
	}
	if len(issues) > 0 {
		for _, v := range issues {
			possibleTransitions, _, _ := jiraClient.Issue.GetTransitions(v.Key)
			for _, v := range possibleTransitions {
				if v.Name == "Release" {
					transitionID = v.ID
					break
				}
			}
			jiraClient.Issue.DoTransition(v.Key, transitionID)
			issue, _, _ := jiraClient.Issue.Get(v.Key, nil)
			fmt.Printf("Issue: %v | %v ==>  %s\n", issue.Key, issue.Fields.Summary, issue.Fields.Status.Name)
			content = append(content,
				"\n\nIssue : ", "`"+issue.Key+"`", " |  "+issue.Fields.Summary+" |  `"+issue.Fields.Status.Name+"`",
				"\n담당자 : "+issue.Fields.Assignee.DisplayName,
				// "\n보고자 : "+issue.Fields.Reporter.DisplayName,
			)
		}
		result := strings.Join(content, " ")
		utils.PostMessage(result)
	} else {
		result := deployDay + " 에 배포 이슈가 없습니다."
		utils.PostMessage(result)
	}
}

func DueDateIssues(jiraClient *jira.Client) {
	var (
		content []string
	)
	// only DOING Issue
	jql := "Write down the Jql that you want to check the due_date."

	issues := GetAllIssues(jiraClient, jql)
	if !*options.confirm {
		for _, v := range issues {
			fmt.Printf("Issue: %v | %v ==>  %s\n", v.Key, v.Fields.Summary, v.Fields.Status.Name)
		}
		c := utils.AskForConfirmation(
			"Count : " + strconv.Itoa(len(issues)) + "\nDo you want to check the jira due_date?")
		if !c {
			os.Exit(0)
		}
	}
	if len(issues) > 0 {
		for _, v := range issues {
			marshal, err := v.Fields.Duedate.MarshalJSON()
			utils.CheckError(err)

			dueDate := bytes.NewBuffer(marshal).String()
			content = append(content,
				"\n\nIssue : ", "`"+v.Key+"`", " |  "+v.Fields.Summary+" |  `"+dueDate+"`",
				"\n담당자 : "+v.Fields.Assignee.DisplayName,
				// "\n기한 : "+bytes.NewBuffer(marshal).String(),
				// "\n보고자 : "+v.Fields.Reporter.DisplayName,
			)
			fmt.Printf("Issue: %v | %v ==>  %s\n", v.Key, v.Fields.Summary, dueDate)
		}
		content = append(content, "\n\n이슈 기한 확인해 주시기 바랍니다.")
		result := strings.Join(content, " ")
		utils.PostMessage(result)
	} else {
		result := "모든 이슈 기한 체크 완료"
		utils.PostMessage(result)
	}
}

func GetAllIssues(jiraClient *jira.Client, jql string) []jira.Issue {
	var result []jira.Issue
	issues, resp, err := jiraClient.Issue.Search(jql, nil)
	utils.CheckError(err)

	total := resp.Total
	if issues == nil {
		result = make([]jira.Issue, 0, total)
	}

	result = append(result, issues...)
	return result
}
