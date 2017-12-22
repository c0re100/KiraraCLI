package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"crypto/tls"
	urlP "net/url"

	"github.com/Jeffail/gabs"
)

var (
	SessionId = "833525b7-245b-4802-a0c6-0f7931a13ab9"
	proxyURL  = "http://192.168.0.216:8888"
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

func LevelUp() {

	url := "https://krr-prd.star-api.com/api/player/town_facility/item_up"

	lvjson := gabs.New()
	lvjson.Set(27290346, "managedTownFacilityId")
	lvjson.Set(130101, "itemNo")
	lvjson.Set(1, "amount")
	lvjson.Set(time.Now().UnixNano() / 1000000, "actionTime")
	log.Println(time.Now().UnixNano() / 1000000)

	payload := strings.NewReader(lvjson.String())

	req, _ := http.NewRequest("POST", url, payload)

	hash := SHA256withSid("/api/player/town_facility/item_up", lvjson.String())

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
	if jsonParsed.S("resultCode").Data().(float64) == 0 {
		log.Println("Levelup!")
	} else {
		log.Println("Error:", jsonParsed.S("resultCode").Data().(float64))
		log.Println(jsonParsed.String())
	}

}

func main() {
	LevelUp()
}
