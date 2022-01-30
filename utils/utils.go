package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		LogText(2, err.Error())
		panic(err)
	}
}

func LogText(num int, str ...string) string {
	_, file, line, _ := runtime.Caller(num)
	text := []interface{}{file, "line-" + strconv.Itoa(line), str}
	fmt.Println(text)
	return fmt.Sprintf("%v", text)
}

func MakeRequest(method, url string, reqData []byte) *http.Request {
	buff := bytes.NewBuffer(reqData)
	req, err := http.NewRequest(method, url, buff)
	CheckError(err)
	req.Close = true
	return req
}

func DoRequest(req *http.Request, header map[string]string) map[string]interface{} {
	for key, val := range header {
		req.Header.Add(key, val)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	CheckError(err)
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	var respData map[string]interface{}
	err = json.Unmarshal(respBytes, &respData)
	CheckError(err)
	return respData
}

func AskForDeploy(s string) (bool, string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s ex) 2022-01-13 : ", s)
		response, err := reader.ReadString('\n')
		CheckError(err)
		re := regexp.MustCompile(`(19|20)\d{2}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[0-1])`)
		if re.MatchString(response) {
			return true, response
		} else {
			return false, response
		}
	}
}

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [Y/N]: ", s)
		response, err := reader.ReadString('\n')
		CheckError(err)
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else {
			return false
		}
	}
}

func PostMessage(content string) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	// if slack payload key text
	data, _ := json.Marshal(map[string]string{
		"text": content,
	})
	DoRequest(MakeRequest("POST", os.Getenv("POST_MESSAGE_URI"), data), header)
}
