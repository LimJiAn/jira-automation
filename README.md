# jira-automation

Jira status change automation.
## Requirements

* Go >= 1.14
* Jira v6.3.4 & v7.1.2.

## Installation


```bash
git clone https://github.com/LimJiAn/jira-automation

go get github.com/andygrunwald/go-jira
go get github.com/joho/godotenv
go get github.com/go-resty/resty/v2

or

go get -u
```

##### Make env file

```bash
vi .env

JIRA_USERNAME = 'username'
JIRA_PASSWORD = 'password'
JIRA_URL = 'jira_url'
WEBHOOK_URL = 'webhook_url'

```

##### Commands

```bash
# due_date check subtask
go run jira.go -com=due_date

# status change preset
go run jira.go -com=preset

# status change release
go run jira.go -com=release
```


