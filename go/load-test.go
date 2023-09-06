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
    "strconv"
    "strings"
    "time"
)

//const LogFilePath = "log.txt"

var config Config

type Config struct {
    Hostname string `json:"hostname"`
    Username string `json:"username"`
    Password string `json:"password"`
    BookletName string `json:"booklet_name"`
    ResourceDir string `json:"resource_dir"`
}

type LoginResult struct {
	Token  string  `json:"token"`
    //Error error
}

func main() {
    logfileName := "results/" + strconv.FormatInt(time.Now().Unix(), 10)
    logFile, err := os.OpenFile(logfileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    check(err)
    defer logFile.Close()
    log.SetOutput(logFile)

    config = loadConfig()

    users, _ := strconv.Atoi(os.Args[1])
    ch := make(chan string)
    for i := 0; i < users; i++ {
        go LoadBooklet(ch, strconv.Itoa(i))
    }
    for i := 0; i < users; i++ {
        log.Println(<-ch)
    }
}

func loadConfig() Config {
    file, err := os.Open("config.json")
    check(err)
    defer file.Close()

    var cfg Config
    byteValue, _ := io.ReadAll(file)
    json.Unmarshal(byteValue, &cfg)
    return cfg
}

func LoadBooklet(ch chan<- string, index string) {
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
    ch <- fmt.Sprintf("%s - %.2f elapsed", index, secs)
}

func login() (string, error) {
    loginURL := config.Hostname + "/api/session/login"
    payload := fmt.Sprintf("{\"name\": %q, \"password\": %q}", config.Username, config.Password)
    req, err := http.NewRequest(http.MethodPut, loginURL, strings.NewReader(payload))
    check(err)
    client := &http.Client{}
    resp, err := client.Do(req)
    check(err)

    if resp.StatusCode != http.StatusOK {
        return "", errors.New("login failed: " + resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    check(err)
    defer resp.Body.Close()

    var result LoginResult
    err = json.Unmarshal(body, &result)
    check(err)

    return result.Token, nil
}

func putTest(token string) (string, error) {
    url := config.Hostname + "/api/test"
    payload := fmt.Sprintf("{\"bookletName\": %q}", config.BookletName)
    req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(payload))
    check(err)
    req.Header.Add("AuthToken", token)
    client := &http.Client{}
    resp, err := client.Do(req)
    check(err)

    if resp.StatusCode != http.StatusCreated {
        return "", errors.New("putTest Failed: " + resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    defer resp.Body.Close()
    return string(body), nil
}

func getTest(token string, testID string) error {
    url := config.Hostname + "/api/test/" + testID
    req, _:= http.NewRequest(http.MethodGet, url, nil)
    req.Header.Add("AuthToken", token)
    client := &http.Client{}
    resp, _ := client.Do(req)

    if resp.StatusCode != http.StatusOK {
        return errors.New("getTest Failed: " + resp.Status)
    }
    defer resp.Body.Close()
    return nil
}

func getResource(token string, testID string) error {
    file, err := os.Open(config.ResourceDir + "resources.txt")
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        req, _:= http.NewRequest(http.MethodGet, config.Hostname + "/api/test/" + testID + "/resource/" + scanner.Text(), nil)
        req.Header.Add("AuthToken", token)
        client := &http.Client{}
        resp, _ := client.Do(req)
        if resp.StatusCode != http.StatusOK {
            return errors.New("getResource Failed: " + scanner.Text() + " Response: " + resp.Status)
        }
        defer resp.Body.Close()
    }
    return nil
}

func getUnits(token string, testID string) error {
    file, err := os.Open(config.ResourceDir + "units.txt")
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        req, _:= http.NewRequest(http.MethodGet,
            config.Hostname + "/api/test/" + testID + "/unit/" + scanner.Text() + "/alias/"  + scanner.Text(),
            nil)
        req.Header.Add("AuthToken", token)
        client := &http.Client{}
        resp, _ := client.Do(req)
        if resp.StatusCode != http.StatusOK {
            return errors.New("getUnit Failed: " + scanner.Text() + " Response: " + resp.Status)
        }
        defer resp.Body.Close()
    }
    return nil
}


func check(e error) {
    if e != nil {
       panic(e)
    }
    //if e != nil {
    //    log.Fatal(e)
    //}
}
