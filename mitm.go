package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"crypto/sha256"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/elazarl/goproxy"
)

var (
	counter float64
	Zero    float64
	Once    float64
	Twice   float64
	Thrice  float64
	Four    float64
	Five    float64
)

func tojson(resp *http.Response) *http.Response {
	DrawnJson, _ := ioutil.ReadAll(resp.Body)
	drawJson, _ := gabs.ParseJSON(DrawnJson)

	star := 0
	children, _ := drawJson.S("gachaResults").Children()
	for _, child := range children {
		counter++
		CharDrawn := child.Search("characterId").Data().(float64)
		if CharDrawn == 10002000 {
			star++
			fmt.Println("ゆの")
		} else if CharDrawn == 11002000 {
			star++
			fmt.Println("野々原ゆずこ")
		} else if CharDrawn == 12002000 {
			star++
			fmt.Println("丈槍由紀")
		} else if CharDrawn == 13002000 {
			star++
			fmt.Println("一井透")
		} else if CharDrawn == 14002000 {
			star++
			fmt.Println("九条カレン")
		} else if CharDrawn == 15002000 {
			star++
			fmt.Println("涼風青葉")
		} else if CharDrawn == 16002000 {
			star++
			fmt.Println("本田珠輝")
		} else if CharDrawn == 17002000 {
			star++
			fmt.Println("千矢")
		} else if CharDrawn == 17002010 {
			star++
			fmt.Println("千矢【クリスマス】")
		} else if CharDrawn == 14012000 {
			star++
			fmt.Println("アリス【クリスマス】")
		} else if CharDrawn == 13022000 {
			star++
			fmt.Println("ユー子【クリスマス】")
		} else if CharDrawn == 15022000 {
			star++
			fmt.Println("はじめ【クリスマス】")
		}
	}

	if star == 0 {
		Zero++
	} else if star == 1 {
		Once++
	} else if star == 2 {
		Twice++
	} else if star == 3 {
		Thrice++
	} else if star == 4 {
		Four++
	} else {
		Five++
	}

	percent1 := fmt.Sprintf("%.6f", Once/counter*100)
	percent2 := fmt.Sprintf("%.6f", Twice/counter*100)
	percent3 := fmt.Sprintf("%.6f", Thrice/counter*100)
	percent4 := fmt.Sprintf("%.6f", Four/counter*100)
	percent5 := fmt.Sprintf("%.6f", Five/counter*100)

	fmt.Println("1 SSR(5*) Probability: "+percent1+"%　　Draw num:", counter)
	fmt.Println("2 SSR(5*) Probability: "+percent2+"%　　Draw num:", counter)
	fmt.Println("3 SSR(5*) Probability: "+percent3+"%　　Draw num:", counter)
	fmt.Println("4 SSR(5*) Probability: "+percent4+"%　　Draw num:", counter)
	fmt.Println(">5 SSR(5*) Probability: "+percent5+"%　　Draw num:", counter)

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(DrawnJson))
	return resp
}

func FakeDraw(time string) string {
	drawJson, _ := gabs.ParseJSON([]byte(`{"resultCode":0,"resultMessage":"","serverTime":"2017-12-25T05:32:28","gachaSteps":[],"serverVersion":1712241,"gachas":[{"gacha":{"id":1,"name":"チュートリアルガチャ","bannerId":"0","type":1,"unlimitedGem":-1,"gem1":-1,"gem10":0,"first10":-1,"itemId":-1,"itemAmount":-1,"startAt":"2017-01-01T00:00:00","endAt":"2099-01-01T00:00:00","sun":1,"mon":1,"tue":1,"wed":1,"thu":1,"fri":1,"sat":1,"reDraw":1,"drawLimit":-1,"allLimit":1,"box1Id":0,"box2Id":-1,"box2Limit":-1,"box2Type":-1,"webViewUrl":"https://krr-prd-web.star-api.com/gacha/132/","pick1":45,"pick2":45,"pick3":45,"pick4":45,"pick5":45},"uGem1Total":0,"uGem1Daily":0,"gem1Total":0,"gem1Daily":0,"gem10Total":0,"gem10Daily":0,"itemTotal":0,"itemDaily":0},{"gacha":{"id":1000000,"name":"☆4確定ガチャ","bannerId":"1","type":1,"unlimitedGem":-1,"gem1":-1,"gem10":-1,"first10":-1,"itemId":10001,"itemAmount":1,"startAt":"2017-12-11T00:00:00","endAt":"2099-01-01T00:00:00","sun":1,"mon":1,"tue":1,"wed":1,"thu":1,"fri":1,"sat":1,"reDraw":-1,"drawLimit":-1,"allLimit":1,"box1Id":1,"box2Id":-1,"box2Limit":-1,"box2Type":-1,"webViewUrl":"https://krr-prd-web.star-api.com/gacha/910/","pick1":45,"pick2":45,"pick3":45,"pick4":45,"pick5":45},"uGem1Total":0,"uGem1Daily":0,"gem1Total":0,"gem1Daily":0,"gem10Total":0,"gem10Daily":0,"itemTotal":0,"itemDaily":0},{"gacha":{"id":4,"name":"クリスマスキャラピックアップ　その２","bannerId":"2017Christmas2","type":1,"unlimitedGem":-1,"gem1":40,"gem10":400,"first10":300,"itemId":10000,"itemAmount":1,"startAt":"2017-12-22T15:00:00","endAt":"2017-12-27T13:59:59","sun":1,"mon":1,"tue":1,"wed":1,"thu":1,"fri":1,"sat":1,"reDraw":-1,"drawLimit":-1,"allLimit":-1,"box1Id":4,"box2Id":-1,"box2Limit":-1,"box2Type":-1,"webViewUrl":"https://krr-prd-web.star-api.com/gacha/48/","pick1":45,"pick2":45,"pick3":45,"pick4":20,"pick5":45},"uGem1Total":0,"uGem1Daily":0,"gem1Total":0,"gem1Daily":0,"gem10Total":0,"gem10Daily":0,"itemTotal":0,"itemDaily":0}]}`))
	drawJson.Set("serverTime", time)
	return drawJson.String()
}

func CRC(data string) string {
	return fmt.Sprintf("%08x\n", crc32.ChecksumIEEE([]byte(data)))
}

func SHA256withSid(SessionId, apiEndpoint, json string) string {
	hash := sha256.New()
	if json != "" {
		hash.Write([]byte(SessionId + " " + apiEndpoint + " " + json + " " + "85af4a94ce7a280f69844743212a8b867206ab28946e1e30e6c1a10196609a11"))
	} else {
		hash.Write([]byte(SessionId + " " + apiEndpoint + " " + "85af4a94ce7a280f69844743212a8b867206ab28946e1e30e6c1a10196609a11"))
	}
	sha256hash := hash.Sum(nil)
	return hex.EncodeToString(sha256hash)
}

func CharSave(sessionId string) {
	url := "https://krr-prd.star-api.com/api/player/tutorial/party/set"

	charJson := gabs.New()
	charJson.Set(-1, "stepCode")

	req, _ := http.NewRequest("POST", url, strings.NewReader(charJson.String()))

	hash := SHA256withSid(sessionId, "/api/player/tutorial/party/set", charJson.String())

	req.Header.Add("unity-user-agent", "app/0.0.0; iOS 11.0.3; iPhone7Plus")
	req.Header.Add("x-star-requesthash", hash)
	req.Header.Add("x-unity-version", "5.5.4f1")
	req.Header.Add("X-STAR-AB", "3")
	req.Header.Add("X-STAR-SESSION-ID", sessionId)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("user-agent", "kirarafantasia/17 CFNetwork/887 Darwin/17.0.0")
	req.Header.Add("Host", "krr-prd.star-api.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	jsonParsed, _ := gabs.ParseJSON(body)
	if jsonParsed.S("resultCode").Data().(float64) == 0 {
	    log.Println("成功儲存首抽角色...")
	}
}

func main() {
	log.Println("MITM started...")
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*baidu.com$"))).
		HandleConnect(goproxy.AlwaysReject)
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse(goproxy.UrlMatches(regexp.MustCompile("/api/player/gacha/draw"))).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp = tojson(resp)
        sessionId := ctx.Req.Header.Get("X-STAR-SESSION-ID")
        go CharSave(sessionId)
		return resp
	})
	proxy.OnResponse(goproxy.UrlMatches(regexp.MustCompile("/api/player/gacha/get_all"))).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	    response := `{"serverVersion":1712241,"resultMessage":"","gachas":[{"gacha":{"id":1,"name":"チュートリアルガチャ","bannerId":"0","type":1,"unlimitedGem":-1,"gem1":-1,"gem10":0,"first10":-1,"itemId":-1,"itemAmount":-1,"startAt":"2017-01-01T00:00:00","endAt":"2099-01-01T00:00:00","sun":1,"mon":1,"tue":1,"wed":1,"thu":1,"fri":1,"sat":1,"reDraw":1,"drawLimit":-1,"allLimit":1,"box1Id":0,"box2Id":-1,"box2Limit":-1,"box2Type":-1,"webViewUrl":"https://krr-prd-web.star-api.com/gacha/132/","pick1":45,"pick2":45,"pick3":45,"pick4":45,"pick5":45},"uGem1Total":0,"uGem1Daily":0,"gem1Total":0,"gem1Daily":0,"gem10Total":0,"gem10Daily":0,"itemTotal":0,"itemDaily":0},{"gacha":{"id":1000000,"name":"☆4確定ガチャ","bannerId":"1","type":1,"unlimitedGem":-1,"gem1":-1,"gem10":-1,"first10":-1,"itemId":10001,"itemAmount":1,"startAt":"2017-12-11T00:00:00","endAt":"2099-01-01T00:00:00","sun":1,"mon":1,"tue":1,"wed":1,"thu":1,"fri":1,"sat":1,"reDraw":-1,"drawLimit":-1,"allLimit":1,"box1Id":1,"box2Id":-1,"box2Limit":-1,"box2Type":-1,"webViewUrl":"https://krr-prd-web.star-api.com/gacha/910/","pick1":45,"pick2":45,"pick3":45,"pick4":45,"pick5":45},"uGem1Total":0,"uGem1Daily":0,"gem1Total":0,"gem1Daily":0,"gem10Total":0,"gem10Daily":0,"itemTotal":0,"itemDaily":0},{"gacha":{"id":4,"name":"クリスマスキャラピックアップ　その２","bannerId":"2017Christmas2","type":1,"unlimitedGem":-1,"gem1":40,"gem10":400,"first10":300,"itemId":10000,"itemAmount":1,"startAt":"2017-12-22T15:00:00","endAt":"2017-12-27T13:59:59","sun":1,"mon":1,"tue":1,"wed":1,"thu":1,"fri":1,"sat":1,"reDraw":-1,"drawLimit":-1,"allLimit":-1,"box1Id":4,"box2Id":-1,"box2Limit":-1,"box2Type":-1,"webViewUrl":"https://krr-prd-web.star-api.com/gacha/48/","pick1":45,"pick2":45,"pick3":45,"pick4":20,"pick5":45},"uGem1Total":0,"uGem1Daily":0,"gem1Total":0,"gem1Daily":0,"gem10Total":0,"gem10Daily":0,"itemTotal":0,"itemDaily":0}],"gachaSteps":[],"resultCode":0,"serverTime":"2017-12-25T07:14:50"}`
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(response)))
		resp.Header.Set("X-Star-CC", "cdc8ee92")
		return resp
	})
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8888", "proxy listen address")
	flag.Parse()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
