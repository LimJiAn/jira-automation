# jira-automation`

Jira status change automation.
## Requirements

* Go >= 1.14
* Jira v6.3.4 & v7.1.2.

## Installation


```bash
go get github.com/andygrunwald/go-jira
go get github.com/joho/godotenv
```

##### make env file

```bash
vi .env

JIRA_USERNAME = 'username'
JIRA_PASSWORD = 'password'
JIRA_URL = 'jira_url'
POST_MESSAGE_URI = 'webhook post message url'

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


