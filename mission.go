package main

import (
    "bufio"
    "crypto/sha256"
    "crypto/tls"
    "encoding/hex"
    "io/ioutil"
    "log"
    "net/http"
    urlP "net/url"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/Jeffail/gabs"
    "github.com/jbrodriguez/mlog"
    "fmt"
)

var (
    SessionId = "1e33c699-09b6-4422-aa8d-c444095cadd5"
    proxyURL  = "http://192.168.0.216:8888"
    count     int
    BoxID     string
)

func SHA256withSid(apiEndpoint string, json string) string {
    hash := sha256.New()
    if json != "" {
        hash.Write([]byte(SessionId + " " + apiEndpoint + " " + json + " " + "85af4a94ce7a280f69844743212a8b867206ab28946e1e30e6c1a10196609a11"))
    } else {
        hash.Write([]byte(SessionId + " " + apiEndpoint + " " + "85af4a94ce7a280f69844743212a8b867206ab28946e1e30e6c1a10196609a11"))
    }
    sha256hash := hash.Sum(nil)
    return hex.EncodeToString(sha256hash)
}

func missionGet() {
    url := "https://krr-prd.star-api.com/api/player/mission/get_all"

    req, _ := http.NewRequest("GET", url, nil)

    hash := SHA256withSid("/api/player/mission/get_all", "")

    req.Header.Add("unity-user-agent", "app/0.0.0; iOS 11.0.3; iPhone7Plus")
    req.Header.Add("x-star-requesthash", hash)
    req.Header.Add("x-unity-version", "5.5.4f1")
    req.Header.Add("X-STAR-AB", "3")
    req.Header.Add("X-STAR-SESSION-ID", SessionId)
    req.Header.Add("content-type", "application/json; charset=UTF-8")
    req.Header.Add("user-agent", "kirarafantasia/17 CFNetwork/887 Darwin/17.0.0")
    req.Header.Add("Host", "krr-prd.star-api.com")

    proxy, _ := urlP.Parse(proxyURL)
    tr := &http.Transport{
        Proxy:           http.ProxyURL(proxy),
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{
        Transport: tr,
        Timeout:   time.Second * 10,
    }

    res, _ := client.Do(req)
    defer res.Body.Close()

    jsonParsed, _ := gabs.ParseJSONBuffer(res.Body)
    if jsonParsed.S("resultCode").Data().(float64) == 0 {
        log.Println("讀取任務列表...")
        var idList string
        children, _ := jsonParsed.S("missionLogs").Children()
        for _, child := range children {
            CharDrawn := child.Search("managedMissionId").Data().(float64)
            idList += strconv.Itoa(int(CharDrawn)) + "\n"
        }
        ioutil.WriteFile("idlist.txt", []byte(idList), 0644)
    }
}

func finish(id int) {
    url := "https://krr-prd.star-api.com/api/player/mission/complete"

    missionjson := gabs.New()
    missionjson.Set(id, "managedMissionId")

    payload := strings.NewReader(missionjson.String())

    req, _ := http.NewRequest("POST", url, payload)

    hash := SHA256withSid("/api/player/mission/complete", missionjson.String())

    req.Header.Add("unity-user-agent", "app/0.0.0; Android OS 7.1.2 / API-25 N2G48C/4104010; LGE Nexus 5X")
    req.Header.Add("x-star-requesthash", hash)
    req.Header.Add("x-unity-version", "5.5.4f1")
    req.Header.Add("X-STAR-AB", "3")
    req.Header.Add("X-STAR-SESSION-ID", SessionId)
    req.Header.Add("content-type", "application/json; charset=UTF-8")
    req.Header.Add("user-agent", "Dalvik/2.1.0 (Linux; U; Android 7.1.2; Nexus 5X Build/N2G48C)")
    req.Header.Add("Host", "krr-prd.star-api.com")

    proxy, _ := urlP.Parse(proxyURL)
    tr := &http.Transport{
        Proxy:           http.ProxyURL(proxy),
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{
        Transport: tr,
        Timeout:   time.Second * 10,
    }

    res, _ := client.Do(req)
    defer res.Body.Close()
    jsonParsed, _ := gabs.ParseJSONBuffer(res.Body)
    mlog.Start(mlog.LevelInfo, "mission.log")
    if jsonParsed.S("resultCode").Data().(float64) == 0 {
        log.Println(id, "Complete!")
        mlog.Info("Mission Complete! mID: %d", id)
    } else {
        mlog.Warning("Error: %.0f mID: %d", jsonParsed.S("resultCode").Data().(float64), id)
    }
}

func getPresent() {
    url := "https://krr-prd.star-api.com/api/player/present/get"

    managedPresentId := "?managedPresentId=" + BoxID + "&stepCode=0"
    req, _ := http.NewRequest("GET", url+managedPresentId, nil)

    hash := SHA256withSid("/api/player/present/get"+managedPresentId, "")

    req.Header.Add("unity-user-agent", "app/0.0.0; Android OS 7.1.2 / API-25 N2G48C/4104010; LGE Nexus 5X")
    req.Header.Add("x-star-requesthash", hash)
    req.Header.Add("x-unity-version", "5.5.4f1")
    req.Header.Add("X-STAR-AB", "3")
    req.Header.Add("X-STAR-SESSION-ID", SessionId)
    req.Header.Add("content-type", "application/json; charset=UTF-8")
    req.Header.Add("user-agent", "Dalvik/2.1.0 (Linux; U; Android 7.1.2; Nexus 5X Build/N2G48C)")
    req.Header.Add("Host", "krr-prd.star-api.com")

    proxy, _ := urlP.Parse(proxyURL)
    tr := &http.Transport{
        Proxy:           http.ProxyURL(proxy),
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{
        Transport: tr,
        Timeout:   time.Second * 10,
    }

    res, _ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    jsonParsed, _ := gabs.ParseJSON(body)
    if jsonParsed.S("resultCode").Data().(float64) == 0 {
        log.Println("領取禮物盒...")
    }
}

func presentGet() {
    url := "https://krr-prd.star-api.com/api/player/present/get_all"

    req, _ := http.NewRequest("GET", url, nil)

    hash := SHA256withSid("/api/player/present/get_all", "")

    req.Header.Add("unity-user-agent", "app/0.0.0; Android OS 7.1.2 / API-25 N2G48C/4104010; LGE Nexus 5X")
    req.Header.Add("x-star-requesthash", hash)
    req.Header.Add("x-unity-version", "5.5.4f1")
    req.Header.Add("X-STAR-AB", "3")
    req.Header.Add("X-STAR-SESSION-ID", SessionId)
    req.Header.Add("content-type", "application/json; charset=UTF-8")
    req.Header.Add("user-agent", "Dalvik/2.1.0 (Linux; U; Android 7.1.2; Nexus 5X Build/N2G48C)")
    req.Header.Add("Host", "krr-prd.star-api.com")

    proxy, _ := urlP.Parse(proxyURL)
    tr := &http.Transport{
        Proxy:           http.ProxyURL(proxy),
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{
        Transport: tr,
        Timeout:   time.Second * 10,
    }

    res, _ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    jsonParsed, _ := gabs.ParseJSON(body)
    if jsonParsed.S("resultCode").Data().(float64) == 0 {
        log.Println("讀取禮物盒...")

        presents, err := jsonParsed.S("presents").Children()
        if err != nil {
            return
        }
        for _, box := range presents {
            BoxID += strconv.FormatFloat(box.Search("managedPresentId").Data().(float64), 'f', 0, 64) + ","
        }
        BoxID = BoxID[:len(BoxID)-1]
    }
}

func ReadFile() {
    f, _ := os.Open("idlist2.txt")
    defer f.Close()
    s := bufio.NewScanner(f)
    for s.Scan() {
        if count >= 80 {
            var next string
            for next == "" {
                fmt.Println("請輸入帖子ID: ")
                fmt.Scanln(&next)
            }
            if next == "1" {
                count = 0
            } else {
                time.Sleep(60 * time.Second)
            }
        }
        /*if s.Text() == "3221484531" || s.Text() == "3221484605" || s.Text() == "3221484650" || s.Text() == "3221484682" || s.Text() == "3221484747" || s.Text() == "3221484789" || s.Text() == "3221484826" || s.Text() == "3221484856" || s.Text() == "3221484905" || s.Text() == "3221484922" || s.Text() == "3221484967" || s.Text() == "3221484988" {
            continue
        }
        f, err := os.OpenFile("idlist2.txt", os.O_APPEND|os.O_WRONLY, 0600)
        if err != nil {
            panic(err)
        }
        defer f.Close()
        _, err = f.WriteString("\n" + s.Text())*/
        count++
        i, _ := strconv.Atoi(s.Text())
        finish(i)
    }
}

func main() {
    /*for {
        missionGet()
    }*/
    //ReadFile()
    presentGet()
    getPresent()
}
