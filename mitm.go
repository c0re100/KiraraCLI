package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

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
    counter++
	DrawnJson, _ := ioutil.ReadAll(resp.Body)
	drawJson, _ := gabs.ParseJSON(DrawnJson)

	star := 0
	children, _ := drawJson.S("gachaResults").Children()
	for _, child := range children {
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

func main() {
	log.Println("MITM started...")
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*baidu.com$"))).
		HandleConnect(goproxy.AlwaysReject)
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse(goproxy.UrlMatches(regexp.MustCompile("/api/player/gacha/draw"))).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp = tojson(resp)
		return resp
	})
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8888", "proxy listen address")
	flag.Parse()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
