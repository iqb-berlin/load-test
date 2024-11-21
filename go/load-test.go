package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Hostname        string `json:"hostname"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	IncrementUserID bool   `json:"increment_user_id"`
	BookletName     string `json:"booklet_name"`
	Workspace       string `json:"workspace"`
	ResourceDir     string `json:"resource_dir"`
	Retries         int    `json:"retries"`
	Timeout         int    `json:"timeout"`
	FileServiceMode bool   `json:"file_service_mode"`
}

type LoginResult struct {
	Token      string `json:"token"`
	GroupToken string `json:"groupToken"`
}

var config Config
var unitList []string
var resourceList []string

func main() {
	logfileName := "results/" + strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println("Logging to: " + logfileName)
	logFile, err := os.OpenFile(logfileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	check(err)
	log.SetOutput(logFile)

	config = loadConfig()

	unitList = readResourceFile(config.ResourceDir + "units.txt")
	resourceList = readResourceFile(config.ResourceDir + "resources.txt")

	users, _ := strconv.Atoi(os.Args[1])
	ch := make(chan string)

	start := time.Now()

	for i := 0; i < users; i++ {
		go LoadBooklet(ch, i+1)
	}
	for i := 0; i < users; i++ {
		log.Println(<-ch)
	}

	secs := time.Since(start).Seconds()
	log.Printf("Total time: %.2f\n", secs)

	out, err := exec.Command("./summary.sh", os.Args[1], logfileName).Output()
	check(err)
	fmt.Println(string(out))

	err = logFile.Close()
	check(err)
}

func loadConfig() Config {
	configPath := os.Args[2]
	file, err := os.Open(configPath)
	check(err)
	var cfg Config
	byteValue, _ := io.ReadAll(file)
	err = json.Unmarshal(byteValue, &cfg)
	check(err)
	err = file.Close()
	check(err)
	return cfg
}

func readResourceFile(filePath string) []string {
	var fileLines []string
	file, err := os.Open(filePath)
	check(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	err = file.Close()
	check(err)
	return fileLines
}

func LoadBooklet(ch chan<- string, userNum int) {
	start := time.Now()

	token, groupToken, err := login(userNum)
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}
	testID, err := putTest(token, userNum)
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}
	err = getTest(token, testID, userNum)
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}
	err = getResource(token, groupToken, testID, userNum)
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}
	err = getUnits(token, testID, userNum)
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("SUCCESS: (user "+strconv.Itoa(userNum)+") Loaded booklet in %.2f", secs)
}

func login(usernameSuffix int) (string, string, error) {
	loginURL := config.Hostname + "/api/session/login"
	var payload = ""
	if config.IncrementUserID {
		username := config.Username + strconv.Itoa(usernameSuffix)
		payload = fmt.Sprintf("{\"name\": %q, \"password\": %q}", username, config.Password)
	} else {
		payload = fmt.Sprintf("{\"name\": %q, \"password\": %q}", config.Username, config.Password)
	}

	response, err := makeRequest(http.MethodPut, loginURL, payload, "", http.StatusOK, usernameSuffix)
	if err != nil {
		return "", "", err
	}

	body, err := io.ReadAll(response.Body)
	check(err)
	err = response.Body.Close()
	check(err)

	var result LoginResult
	err = json.Unmarshal(body, &result)
	check(err)
	return result.Token, result.GroupToken, nil
}

func putTest(token string, userNum int) (string, error) {
	url := config.Hostname + "/api/test"
	payload := fmt.Sprintf("{\"bookletName\": %q}", config.BookletName)

	response, err := makeRequest(http.MethodPut, url, payload, token, http.StatusCreated, userNum)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	check(err)
	err = response.Body.Close()
	check(err)
	return string(body), nil
}

func getTest(token string, testID string, userNum int) error {
	url := config.Hostname + "/api/test/" + testID

	response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK, userNum)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	check(err)
	return nil
}

func getResource(token string, groupToken string, testID string, userNum int) error {
	for _, resource := range resourceList {
		url := config.Hostname + "/api/test/" + testID + "/resource/" + resource
		if config.FileServiceMode {
			url = config.Hostname + "/fs/file/" + groupToken + "/" + config.Workspace + "/Resource/" + resource
		}
		response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK, userNum)
		if err != nil {
			return err
		}
		err = response.Body.Close()
		check(err)
	}
	return nil
}

func getUnits(token string, testID string, userNum int) error {
	for _, unit := range unitList {
		url := config.Hostname + "/api/test/" + testID + "/unit/" + unit + "/alias/" + unit
		response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK, userNum)
		if err != nil {
			return err
		}
		err = response.Body.Close()
		check(err)
	}
	return nil
}

func makeRequest(method string, url string, payload string, authToken string, expectedStatus int, userNum int) (*http.Response, error) {
	var (
		request   *http.Request
		response  *http.Response
		err       error = nil
		triesLeft int   = config.Retries + 1
		timeout   int   = config.Timeout
	)

	// Ignore self signed cert
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	for triesLeft > 0 {
		request, err = http.NewRequest(method, url, strings.NewReader(payload))
		request.Header.Add("AuthToken", authToken)
		check(err)
		client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
		response, err = client.Do(request)

		if os.IsTimeout(err) {
			triesLeft -= 1
			timeout *= 2
			log.Printf("WARNING: " + "(user " + strconv.Itoa(userNum) + ") " + url + " failed! Timeout; Retries left: " + strconv.Itoa(triesLeft))
			if triesLeft == 0 {
				err = errors.New("ERROR: " + "(user " + strconv.Itoa(userNum) + ") " + url + " failed! Timeout; Retries left: " + strconv.Itoa(triesLeft))
			}
			continue
		}

		check(err)

		if response.StatusCode != expectedStatus {
			triesLeft -= 1
			timeout *= 2
			log.Printf("WARNING: " + "(user " + strconv.Itoa(userNum) + ") " + method + " " + url + " failed! Response: " + response.Status + "; Retries left: " + strconv.Itoa(triesLeft))
			if triesLeft == 0 {
				err = errors.New("ERROR: " + "(user " + strconv.Itoa(userNum) + ") " + method + " " + url + " failed! Response: " + response.Status)
			}
			continue
		}
		break
	}
	return response, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
