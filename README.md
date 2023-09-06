# jira-automation

Jira status change automation.

Code that registers as a system service file and then executes it as a cron to change the status of the deployed issue
## ✓ Requirements

* Go >= 1.14
* Jira v6.3.4 & v7.1.2.

## ⚙️ Installation


```bash
$ git clone https://github.com/LimJiAn/jira-automation
```

## 👀 Usage
##### 1. Make env file

```bash
$ vi .env

JIRA_USERNAME = 'username'
JIRA_PASSWORD = 'password'
JIRA_URL = 'jira_url'
WEBHOOK_URL = 'webhook_url'

```
##### 2. Registers as a system service
