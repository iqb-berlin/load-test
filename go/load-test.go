package main

import (
    "bufio"
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

var config Config

type Config struct {
    Hostname string `json:"hostname"`
    Username string `json:"username"`
    Password string `json:"password"`
    BookletName string `json:"booklet_name"`
    ResourceDir string `json:"resource_dir"`
    Retries int `json:"retries"`
    Timeout int `json:"timeout"`
}

type LoginResult struct {
	Token  string  `json:"token"`
}

func main() {
    logfileName := "results/" + strconv.FormatInt(time.Now().Unix(), 10)
    fmt.Println("Logging to: " + logfileName)
    logFile, err := os.OpenFile(logfileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    check(err)
    log.SetOutput(logFile)

    config = loadConfig()

    users, _ := strconv.Atoi(os.Args[1])
    ch := make(chan string)

    start := time.Now()

    for i := 0; i < users; i++ {
        go LoadBooklet(ch)
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
    file, err := os.Open("config.json")
    check(err)
    var cfg Config
    byteValue, _ := io.ReadAll(file)
    err = json.Unmarshal(byteValue, &cfg)
    check(err)
    err = file.Close()
    check(err)
    return cfg
}

func LoadBooklet(ch chan<- string) {
    start := time.Now()

    token, err := login()
    if err != nil {
        ch <- fmt.Sprintf(err.Error())
        return
    }
    testID, err := putTest(token)
    if err != nil {
        ch <- fmt.Sprintf(err.Error())
        return
    }
    err = getTest(token, testID)
    if err != nil {
        ch <- fmt.Sprintf(err.Error())
        return
    }
    err = getResource(token, testID)
    if err != nil {
        ch <- fmt.Sprintf(err.Error())
        return
    }
    err = getUnits(token, testID)
    if err != nil {
        ch <- fmt.Sprintf(err.Error())
        return
    }

    secs := time.Since(start).Seconds()
    ch <- fmt.Sprintf("LoadBooklet success: %.2f elapsed", secs)
}

func login() (string, error) {
    loginURL := config.Hostname + "/api/session/login"
    payload := fmt.Sprintf("{\"name\": %q, \"password\": %q}", config.Username, config.Password)

    response, err := makeRequest(http.MethodPut, loginURL, payload, "", http.StatusOK)
    if err != nil {
        return "", err
    }

    body, err := io.ReadAll(response.Body)
    check(err)
    err = response.Body.Close()
    check(err)

    var result LoginResult
    err = json.Unmarshal(body, &result)
    check(err)
    return result.Token, nil
}

func putTest(token string) (string, error) {
    url := config.Hostname + "/api/test"
    payload := fmt.Sprintf("{\"bookletName\": %q}", config.BookletName)

    response, err := makeRequest(http.MethodPut, url, payload, token, http.StatusCreated)
    if err != nil {
        return "", err
    }

    body, err := io.ReadAll(response.Body)
    check(err)
    err = response.Body.Close()
    check(err)
    return string(body), nil
}

func getTest(token string, testID string) error {
    url := config.Hostname + "/api/test/" + testID

    response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK)
    if err != nil {
        return err
    }

    err = response.Body.Close()
    check(err)
    return nil
}

func getResource(token string, testID string) error {
    file, err := os.Open(config.ResourceDir + "resources.txt")
    check(err)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        url := config.Hostname + "/api/test/" + testID + "/resource/" + scanner.Text()
        response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK)
        if err != nil {
            return err
        }
        err = response.Body.Close()
        check(err)
    }
    err = file.Close()
    check(err)
    return nil
}

func getUnits(token string, testID string) error {
    file, err := os.Open(config.ResourceDir + "units.txt")
    check(err)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        url := config.Hostname + "/api/test/" + testID + "/unit/" + scanner.Text() + "/alias/"  + scanner.Text()
        response, err := makeRequest(http.MethodGet, url, "", token, http.StatusOK)
        if err != nil {
            return err
        }
        err = response.Body.Close()
        check(err)
    }
    err = file.Close()
    check(err)
    return nil
}

func makeRequest(method string, url string, payload string, authToken string, expectedStatus int) (*http.Response, error) {
    var (
        request *http.Request
        response *http.Response
        err      error = nil
        retries  int = config.Retries
    )

    for retries > 0 {
        request, err = http.NewRequest(method, url, strings.NewReader(payload))
        request.Header.Add("AuthToken", authToken)
        check(err)
        client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
        response, err = client.Do(request)

        if os.IsTimeout(err)  {
            err = errors.New("error: " + url + " failed! Timeout")
            retries -= 1
        } else {
          if response.StatusCode != expectedStatus {
              err = errors.New("error: " + method + " " + url + " failed! Response: " + response.Status)
              retries -= 1
          } else {
              break
          }
        }
    }
    return response, err
}

func check(e error) {
    if e != nil {
       panic(e)
    }
}
