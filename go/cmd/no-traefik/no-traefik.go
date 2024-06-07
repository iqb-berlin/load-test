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

var config Config
var resourceList []string

func main() {
	logfileName := "results/" + strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println("Logging to: " + logfileName)
	logFile, err := os.OpenFile(logfileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	check(err)
	log.SetOutput(logFile)

	config = loadConfig()

	resourceList = readResourceFile(config.ResourceDir + "resources.txt")

	users, _ := strconv.Atoi(os.Args[1])
	ch := make(chan string)

	start := time.Now()

	for i := 0; i < users; i++ {
		go LoadBooklet(ch, i)
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

func LoadBooklet(ch chan<- string, i int) {
	start := time.Now()

	// TODO insert a valid auth token and group token here
	err := getResource("")
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("SUCCESS: Loaded booklet in %.2f", secs)
}

func getResource(token string) error {
	for _, resource := range resourceList {
		url := config.Hostname + "/" + resource
		response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK)
		if err != nil {
			return err
		}
		err = response.Body.Close()
		check(err)
	}
	return nil
}

func makeRequest(method string, url string, payload string, authToken string, expectedStatus int) (*http.Response, error) {
	var (
		request  *http.Request
		response *http.Response
		err      error = nil
		retries  int   = config.Retries + 1
	)

	// Ignore self signed cert
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	for retries > 0 {
		request, err = http.NewRequest(method, url, strings.NewReader(payload))
		request.Header.Add("AuthToken", authToken)
		check(err)
		client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
		response, err = client.Do(request)

		if os.IsTimeout(err) {
			retries -= 1
			log.Printf("WARNING: " + url + " failed! Timeout; Retries left: " + strconv.Itoa(retries))
			if retries == 0 {
				err = errors.New("ERROR: " + url + " failed! Timeout; Retries left: " + strconv.Itoa(retries))
			}
			continue
		}

		check(err)

		if response.StatusCode != expectedStatus {
			retries -= 1
			log.Printf("WARNING: " + method + " " + url + " failed! Response: " + response.Status + "; Retries left: " + strconv.Itoa(retries))
			if retries == 0 {
				err = errors.New("ERROR: " + method + " " + url + " failed! Response: " + response.Status)
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
