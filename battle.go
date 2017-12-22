package main

import (
    "bytes"
    "compress/zlib"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "io"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "strconv"
    "strings"
    "time"

    "crypto/tls"
    "fmt"
    urlP "net/url"

    "github.com/Jeffail/gabs"
)

var (
    SessionId = "6efe8970-0b33-4071-a91b-337fbf145d6a"
    rID       = float64(42418238)
    proxyURL  = "http://192.168.0.216:8888"
)

func random(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max-min) + min
}

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

func Decode(json string) string {
    sDec, _ := base64.RawStdEncoding.DecodeString(json)

    var out bytes.Buffer
    r, _ := zlib.NewReader(bytes.NewReader(sDec))
    io.Copy(&out, r)
    r.Close()

    battleJson, _ := gabs.ParseJSON([]byte("{" + out.String() + "}"))
    battleJson.Set(int(rID), "RecvID")

    var in bytes.Buffer
    w := zlib.NewWriter(&in)
    w.Write(battleJson.Bytes()[1: len(battleJson.Bytes())-1])
    w.Close()

    sEec := base64.RawStdEncoding.EncodeToString(in.Bytes())

    return sEec
}

func battleAdd() {
    url := "https://krr-prd.star-api.com/api/player/quest_log/add"

    json := `{"type":1,"questId":1101030,"managedBattlePartyId":37733776,"supportCharacterId":-1,"questData":"eJztV1FPwjAQ/i99HqTdBpM9Amow06CgPBhDqhRY3MaybkRi+O9eO8ZmbGBYXzQuIbm7ft/H9dpLruiBJchFpImbNjLQKHvOA8IZLilnyG0QA92xl/Wgn9u3GeOpcAjBBFsYgOnmOlqIkOU4FvzaBuqtorS3yqIUuYCY0DUbzN6kfc3T0asfBOVqEbkIFsid04AzA10ks/EmZnId7N6SJrTIoPC9NSQhXa9bAO9HZbQHNkfu45P0JnGFP4lL2CSucPjSjwutOxZSPxpnSbRX56xMO/dn+5QvabaAhN/RAw0A0MRbA4272XwuYuHU8znwHoUJxSl3l7vIbQlzCNsKc3wFIAXhjLZb40Q60aObenRLj27r0VsH6aSkk6OVN7H8bKd9Zjmk1TldrVpJsyO/lmmRtnlmm6erHa7MH95a5dBtvX5R02v3i5peu1/U9Nr9oqZrlq52v+AfvVRqNfLNS6VWO1zYCt360a2p1axvbk2tVvvQHb1+UdNr94uaXrtf1PTa/aKma5buv1/+QL88wZB2HrHQZzxHiVGV7+a0/cLe28jBtwODL8by8srgbnwMp/1kJafMjknyxLQ4DfIJ3fgEbxCR+6/IUtR49LJksyxgg5SFeaql1acplXO6UC4cAPTlf99kcH67zZZI6X6VOMg6KgkPhaG3ezEMvauVH8GjRfofYBf9Tg=="}`

    payload := strings.NewReader(json)

    req, _ := http.NewRequest("POST", url, payload)

    hash := SHA256withSid("/api/player/quest_log/add", json)

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
        rID = jsonParsed.S("orderReceiveId").Data().(float64)
        stamina := jsonParsed.S("player", "stamina").Data().(float64)
        log.Println("加載戰鬥資料...orderReceiveId:", strconv.Itoa(int(rID)), "體力剩餘:", strconv.Itoa(int(stamina)))
    } else {
        log.Println("Error:", jsonParsed.S("resultCode").Data().(float64))
    }
}

func battleSave1() {
    url := "https://krr-prd.star-api.com/api/player/quest_log/save"

    waveJson := gabs.New()
    waveJson.Set(int(rID), "orderReceiveId")
    waveJson.Set(Decode(`eJztWFlv4jAQ/i9+DsjOQYC3QtotK9plocdDVa0McSFqDhQnPbTqf18fgdCVu7ikK7WSefIcnz0z9mcyBlckB32A2rDtAAvMyrlUuEyYrDAloA8tMCWLh1EA+o7vdeyOZ1vgZ0lowVUIQQQd5jMpns/SJVd1feR0feY0zNJimJVpISa5xg9kFD6J8RktZvdRHNfWjeYkXoL+HY4pscBJHl48r2UIbDxc4RyLNV2IbAh3tOMH0Pc8IY8HzEGMLmdc7UsvNqagf+Nb/q1QXK+r6N1K5L42rASBRD0h0VW0Bn3pNiUJjtKLMk83QV1SUich5XCbwDdcLln4v8EVjplDG75Y4GJQ3t1xXfJrHFGGu+FDVqo6VymyjPhwwtJLpP+Og5gQ8RlfrHfiUUO83RDvNMS7DfGeNt5vWH81Xr/+arx+/dV4/fqr8fr1V+P31B/VeLin/s4h8Lr87iHwunp223F7/AfdjmN7na73r9mcZrmo4Y5uLmq4e2AuZl8+576gZrmo4bZuLmr4u3KB+vuy7+5R4/XvXjVe/+7ddzAOw+vfvWq8/t1rNztLarj2WVLDDS8qXANeqPH6vFDj9XmhxuvzQo3X54Ua/zcvbtmH+nFKkohQ6cabF1p9q28NW+lZ9BI91glBKEohlKKH4FKQZ6Lb6NlIRtYI00KvvFuv3FuIx/4louQ1ni1WJCxjMipIIkOtRwEuMB+JfDYCcwh4n8aDOS/ZHlb51s5C1JoFvT0F6xEnY+m5PStTQlm/dhyTROhnBV7cbyZuJrDVZgXlfSH9D1MfzdMkrrKmU7IUTWfVoFadrCjZV7dtUg0iiucx2bTZO4WR9qMwnOTZXGU/J0/F7jaQtbg76u5dW8mCuSb4np+WQZaWVLXaKS7IcIVT+USwa2EZklRofzymJD/HCRMA2ObKTmf2KOoBBjjPI/5kw1bdSvWLxCl/unCcrgUCgsOAFDiK5XvLcEoWmJeRzzJMQpnzJSWCE/LRh9llRjzJjzWJc2q49UVshltvcMuBHWS4ZbhluPXx3EJdu2O4ZbhluPXx3GqZfy3DLMMswyzDrM9mM8x683sQeZ+QW+Kt8HsWpaPwiS8OLWTZlmO5lnf7B0TMKT0=`), "questData")

    payload := strings.NewReader(waveJson.String())

    req, _ := http.NewRequest("POST", url, payload)

    hash := SHA256withSid("/api/player/quest_log/save", waveJson.String())

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
        log.Println("儲存Wave1...")
    } else {
        log.Println("Error:", jsonParsed.S("resultCode").Data().(float64))
    }
}

func battleSave2() {
    url := "https://krr-prd.star-api.com/api/player/quest_log/save"

    waveJson := gabs.New()
    waveJson.Set(int(rID), "orderReceiveId")
    waveJson.Set(Decode(`eJztWFtv2jAU/i9+DsiOExJ4K9CuTLRj0MtDVU2GuBA1CShOulZT//t8CYRO7nDJKsEUnnwun33OsT+TY3BDU9ABqAmbGFhgkk+VwuHCaEEYBR1ogTGdPQ36oIM9t2W3XNsC33PKMqFCCCKIuc8oe7lI5kLlewj7HnfqLZOst8yTTE5yS57oIHjmEAtcsGzyGEZRaV1rzqI56DyQiFELnKXB1ctKhcDHvQVJiVzTgciGcEs7fAId15XysCtX4KPriVB7youPGejceZZ3LxW3qyJ6pxCFrw0LQSJRW0psEa5AR7mNaUzC5CpPk3VQ14yWSSg52CTwheRzHv4vcEMi7tDEDo9aRI4x8j2/5b9a4KqbPzwIn/jHMGR8njsx5KUrc1ciz1AMRzzdWPlvOcgFEGzC11frg3hUEW9XxOOKeKci3jXGexXrr8eb11+PN6+/Hm9efz3evP56/I76oxIPd9Qf7wMvy+/sAy+rZ3N6t8UPOi1suy3f/dtsuFouejg2zUUPd/bMpd6Xw9wXVC0XPdw2zUUP/1Au0Hxfdt09erz53avHm9+9uw7Gfnjzu1ePN7977WpnSQ83Pkt6eM2LAleBF3q8OS/0eHNe6PHmvNDjzXmhx//Ji3v+oX6a0DikTLmJZoYV3+obw0Z6kb1Fm3dG/HPfWitlTyGkfrqU3UfbRiqySpgGeuPdeOPeQCL2o4hS1HgyW9Agj+ggo7EKtRz1SUbESOazFrhDX/RtIpjLnO9hkW/pLEWjWdD7U/CecTRUnpuzMqaM92unEY2lfpKR2eN64moCX22SMdEXsk+Y+mSaxFGRNRvTuWw6i4a16GxlyY7dtk61HzIyjei6zd4qjLKfBMEoXU519kv6nG1vA12pk1u081DVVmihTvnGlYdzS8mjOC/dZZIz3XrnJKO9BUnUo8G2hedIE6n99jOh6SWJuQDAJlt+Ppc/ZUVAl6RpKB5x+KobqXyjOBePGRj7FuhTEvRpRsJIvcD0xnRGRCHFLL04UFlfMypZoZ6BuF1lJLL8tyZ5Umt2HYntk9il5dFO5SFxC8MWqrlVc+vwuPU//HMh327V7KrZdXjsOv5/rkb9v1Uzq2ZWzayaWYdmq5n17vcgcg+QW/K98OsyTAbBs1gcWsiyLWw5lnv/G075LSs=`), "questData")

    payload := strings.NewReader(waveJson.String())

    req, _ := http.NewRequest("POST", url, payload)

    hash := SHA256withSid("/api/player/quest_log/save", waveJson.String())

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
        log.Println("儲存Wave2...")
    } else {
        log.Println("Error:", jsonParsed.S("resultCode").Data().(float64))
    }
}

func battleSave3() {
    url := "https://krr-prd.star-api.com/api/player/quest_log/save"

    waveJson := gabs.New()
    waveJson.Set(int(rID), "orderReceiveId")
    waveJson.Set(Decode(`eJztWetP2zAQ/1/8Oav8yKPJt9HuwcRY18L4gNBkiGkjkrSKEwaa+r/Pjz7C8No06VjRAkKyz/c7353v59YH+MYyEADUgR0CLDAqrrXAFpPBhHIGAmiBIbu5P+6DgDjYIa6HLfC1YDxXIgQhwkJnkD9+TsdS1PUQ6Uql3jTNe9MizZWRC3rPjsMHEIiVzzwf3UVxvF5dSt7HYxDc0pgzC7zPwrPHmXZBjHsTmlG5AbLllrAkPbkHgeOo+cmRUFCj85EUu1pLjDkILj3Lu1KCi5myhKC9mEpdDBcThURdNeOTaAYCrTZkCY3SsyJLl06dc7YOQs/DVQAfaDEW7v8E32gsrHeIK5Il9oSe60ObzC1wdlTc3kqN5PtJxIWVSzkUiVtHrqciPjkciGATrV9SWJiH87m1Ixw1g+NmcNIMbjeDOxvhaA2HWzKPOrIaxY/tuV3iIcff3RoqW/PVj4MJcnHXxrtbIyVrGzPjNSsqM7xyUZnhlYvKDK9cVGZ45aIyw//boirByV5DM1sjNUMzW7M3hgZ3ObUacNQMjpvBK98VZnjD1DlVi6rly6vhSwmO9hqa2RquGZrZWuV63n5qNeCVrwIzvPJVsL1IasAbpq7yVYD3WlRma3WLymztt8xciW/+71KWRIxrLfk24osv/6uF1exRPVV8iCDGKnVKqB4synA/m+rXjHqKEe3dU6D4qwlEVXa8WkP/5HptD1xcE4hRFeBm17GvHmx1PYA1j8utmHVZSKObCQuLmB3nLNH+r0d9mlM5UkEuJ0JhYUTGdVqISkXz7RpWSUpKOIV8vuVONsx7b9dAL+LH38OJXwsMTjRudWsMGY94/i5miZKPcnpztzzEZhOx2yjnsufAK5r+CVZXnG57iAtOXWSwQ6C4WX3f820IXWJ7/nwXP95ep0m8SBgfsrHqfiz6JosGi0rba19bhtqPOL2O2bLfU0qMXn8bhoNsem1aP2UPefnM2AwEb9DqHKDOrZRCk/CJqnDngtE7WVxH07Tgpv0+0pz1JjTVvasKleDoD0jXxsTDiChGiCoes1RZ+PIjZdkpTcQEgFVmROFPf6jsgSOaZZHsOwoPV7N1W+3jTLYe7a4F+oyGfZbTKNZNw96Q3VCZdGmll4Q6Q+ecKYbpzqVY127KjCyX0JMl7JTXNsB+X1JlfSi8fXZA63NHtj4hBB2fQIJ3oGrL273ydjtD983b/XMRux7ZIxcr05SQ18tF3HKx5eJf4CKy5X/a/gEXN35kHhQXW3odEL22Cg+JXJL1L0Gt1/qNs2VWy6yWWS2zWmYdDrOQ032hx9lO3FJ9zk/TKD0OH+Tm0EIWtohlW87VL5+F9pY=`), "questData")

    payload := strings.NewReader(waveJson.String())

    req, _ := http.NewRequest("POST", url, payload)

    hash := SHA256withSid("/api/player/quest_log/save", waveJson.String())

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
        log.Println("儲存Wave3...")
    } else {
        log.Println("Error:", jsonParsed.S("resultCode").Data().(float64))
    }
}

func BattleSet() {
    url := "https://krr-prd.star-api.com/api/player/quest_log/set"

    setJson, _ := gabs.ParseJSON([]byte(`{"orderReceiveId":42298325,"state":2,"clearRank":3,"skillExps":"3432518231:0:0:0,3430779650:0:0:0,3431166508:0:0:0","dropItems":"7010:99,7011:99,7012:99,7013:99,8003:99,8004:99,8005:99,8006:99,8007:99","killedEnemies":"19010002:4","weaponSkillExps":"","friendUseNum":0,"masterSkillUseNum":0,"uniqueSkillUseNum":0,"stepCode":0}`))
    setJson.Set(int(rID), "orderReceiveId")

    payload := strings.NewReader(setJson.String())

    req, _ := http.NewRequest("POST", url, payload)

    hash := SHA256withSid("/api/player/quest_log/set", setJson.String())

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
        log.Println("上級已清...")
    } else {
        log.Println("Error:", jsonParsed.S("resultCode").Data().(float64))
    }
}

func main() {
    log.Println("刷上級活動用...")
    battleAdd()
    /*time.Sleep(5 * time.Second)
    battleSave1()
    time.Sleep(time.Duration(random(5, 10)) * time.Second)
    battleSave2()
    time.Sleep(time.Duration(random(5, 10)) * time.Second)*/
    //battleSave3()
    //time.Sleep(time.Duration(random(5, 10)) * time.Second)
    BattleSet()
}

func main2() {
    sDec := `"Ver":"1.0.3","SubVer":"4","Phase":0,"RecvID":33021357,"QuestID":3100120,"PtyMngID":8713871,"ContCount":0,"WaveIdx":2,"MstSkillCount":0,"MstSkillFlg":false,"FrdType":0,"FrdCharaID":14012000,"FrdCharaLv":55,"FrdLB":1,"FrdUSLv":1,"FrdCSLvs":[7,7],"FrdWpID":1104,"FrdWpLv":20,"FrdWpSLv":17,"Frdship":4,"FrdRemainTurn":0,"FrdUseCount":0,"FrdUsed":false,"Gauge":{"Val":1.6150131225585938},"TBuff":{"m_List":[{"m_CondType":0,"m_Cond":5,"m_Param":{"m_Type":0,"m_Val":2.0}},{"m_CondType":0,"m_Cond":5,"m_Param":{"m_Type":1,"m_Val":2.0}},{"m_CondType":0,"m_Cond":5,"m_Param":{"m_Type":2,"m_Val":2.0}},{"m_CondType":0,"m_Cond":5,"m_Param":{"m_Type":3,"m_Val":2.0}},{"m_CondType":0,"m_Cond":5,"m_Param":{"m_Type":4,"m_Val":2.0}},{"m_CondType":0,"m_Cond":5,"m_Param":{"m_Type":5,"m_Val":2.0}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":1,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":3,"m_Val":1.0}},{"m_CondType":0,"m_Cond":7,"m_Param":{"m_Type":0,"m_Val":2.0}},{"m_CondType":0,"m_Cond":7,"m_Param":{"m_Type":1,"m_Val":2.0}},{"m_CondType":0,"m_Cond":7,"m_Param":{"m_Type":2,"m_Val":2.0}},{"m_CondType":0,"m_Cond":7,"m_Param":{"m_Type":3,"m_Val":2.0}},{"m_CondType":0,"m_Cond":7,"m_Param":{"m_Type":4,"m_Val":2.0}},{"m_CondType":0,"m_Cond":7,"m_Param":{"m_Type":5,"m_Val":2.0}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":1,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":3,"m_Val":1.0}},{"m_CondType":1,"m_Cond":3,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":3,"m_Param":{"m_Type":3,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":3,"m_Param":{"m_Type":4,"m_Val":1.0}},{"m_CondType":0,"m_Cond":0,"m_Param":{"m_Type":0,"m_Val":1.0}},{"m_CondType":0,"m_Cond":0,"m_Param":{"m_Type":1,"m_Val":1.0}},{"m_CondType":0,"m_Cond":0,"m_Param":{"m_Type":2,"m_Val":1.0}},{"m_CondType":0,"m_Cond":0,"m_Param":{"m_Type":3,"m_Val":1.0}},{"m_CondType":0,"m_Cond":0,"m_Param":{"m_Type":4,"m_Val":1.0}},{"m_CondType":0,"m_Cond":0,"m_Param":{"m_Type":5,"m_Val":1.0}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":1,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":0,"m_Param":{"m_Type":3,"m_Val":1.0}},{"m_CondType":1,"m_Cond":3,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":3,"m_Param":{"m_Type":3,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":3,"m_Param":{"m_Type":4,"m_Val":1.0}},{"m_CondType":1,"m_Cond":1,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":1,"m_Param":{"m_Type":2,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":1,"m_Param":{"m_Type":4,"m_Val":1.0}},{"m_CondType":0,"m_Cond":3,"m_Param":{"m_Type":0,"m_Val":1.0}},{"m_CondType":0,"m_Cond":3,"m_Param":{"m_Type":1,"m_Val":1.0}},{"m_CondType":0,"m_Cond":3,"m_Param":{"m_Type":2,"m_Val":1.0}},{"m_CondType":0,"m_Cond":3,"m_Param":{"m_Type":3,"m_Val":1.0}},{"m_CondType":0,"m_Cond":3,"m_Param":{"m_Type":4,"m_Val":1.0}},{"m_CondType":0,"m_Cond":3,"m_Param":{"m_Type":5,"m_Val":1.0}},{"m_CondType":1,"m_Cond":2,"m_Param":{"m_Type":0,"m_Val":1.2000000476837159}},{"m_CondType":1,"m_Cond":2,"m_Param":{"m_Type":2,"m_Val":1.2999999523162842}},{"m_CondType":1,"m_Cond":2,"m_Param":{"m_Type":4,"m_Val":1.0}}]},"Enemies":{"m_Waves":[{"m_Enemies":[{"m_EnemyID":19010223,"m_EnemyLv":24,"m_DropID":1001203},{"m_EnemyID":19010901,"m_EnemyLv":24,"m_DropID":1001203},{"m_EnemyID":19010913,"m_EnemyLv":24,"m_DropID":1001203}]},{"m_Enemies":[{"m_EnemyID":19010901,"m_EnemyLv":24,"m_DropID":1001203},{"m_EnemyID":19010621,"m_EnemyLv":24,"m_DropID":1001203},{"m_EnemyID":19010211,"m_EnemyLv":24,"m_DropID":1001203}]},{"m_Enemies":[{"m_EnemyID":29011101,"m_EnemyLv":24,"m_DropID":1001203},{"m_EnemyID":19010203,"m_EnemyLv":24,"m_DropID":1001203},{"m_EnemyID":19010613,"m_EnemyLv":24,"m_DropID":1001203}]}]},"ScheduleItems":[{"Items":[{"Datas":[{"ID":100101,"Num":1},{"ID":100103,"Num":1}]},{"Datas":[{"ID":100101,"Num":1},{"ID":100103,"Num":99}]},{"Datas":[{"ID":100101,"Num":1},{"ID":100103,"Num":1}]}]},{"Items":[{"Datas":[{"ID":100101,"Num":1}]},{"Datas":[{"ID":100101,"Num":1},{"ID":100103,"Num":1}]},{"Datas":[{"ID":100101,"Num":1}]}]},{"Items":[{"Datas":[{"ID":100101,"Num":99}]},{"Datas":[{"ID":100101,"Num":1}]},{"Datas":[{"ID":100101,"Num":1},{"ID":100103,"Num":1}]}]}],"PLs":[{"Param":{"ResistElem":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"StsBuffs":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"Abnmls":[{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0}],"AbnmlDisableBuff":{"Stack":[]},"AbnmlAddProbBuff":{"Stack":[]},"NextBuffs":[{"Step":-1,"Val":0.0},{"Step":0,"Val":0.0},{"Step":-1,"Val":0.0}],"WeakElemBonusBuff":{"Stack":[]},"HateChange":{"Stack":[]},"Regene":{"OwnerName":"","Turn":0,"Pow":0},"Barrier":0.0,"BarrierCount":0,"Hp":1177,"DeadDetail":0,"MCRecast":0},"Cmds":[{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0}]},{"Param":{"ResistElem":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"StsBuffs":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"Abnmls":[{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0}],"AbnmlDisableBuff":{"Stack":[]},"AbnmlAddProbBuff":{"Stack":[]},"NextBuffs":[{"Step":0,"Val":0.0},{"Step":-1,"Val":0.0},{"Step":-1,"Val":0.0}],"WeakElemBonusBuff":{"Stack":[]},"HateChange":{"Stack":[]},"Regene":{"OwnerName":"","Turn":0,"Pow":0},"Barrier":0.0,"BarrierCount":0,"Hp":1138,"DeadDetail":0,"MCRecast":0},"Cmds":[{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0}]},{"Param":{"ResistElem":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"StsBuffs":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"Abnmls":[{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0}],"AbnmlDisableBuff":{"Stack":[]},"AbnmlAddProbBuff":{"Stack":[]},"NextBuffs":[{"Step":0,"Val":0.0},{"Step":-1,"Val":0.0},{"Step":-1,"Val":0.0}],"WeakElemBonusBuff":{"Stack":[]},"HateChange":{"Stack":[]},"Regene":{"OwnerName":"","Turn":0,"Pow":0},"Barrier":0.0,"BarrierCount":0,"Hp":599,"DeadDetail":0,"MCRecast":0},"Cmds":[{"UseNum":1,"RecastVal":35},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0}]},{"Param":{"ResistElem":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"StsBuffs":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"Abnmls":[{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0}],"AbnmlDisableBuff":{"Stack":[]},"AbnmlAddProbBuff":{"Stack":[]},"NextBuffs":[{"Step":0,"Val":0.0},{"Step":0,"Val":0.0},{"Step":0,"Val":0.0}],"WeakElemBonusBuff":{"Stack":[]},"HateChange":{"Stack":[]},"Regene":{"OwnerName":"","Turn":0,"Pow":0},"Barrier":0.0,"BarrierCount":0,"Hp":763,"DeadDetail":0,"MCRecast":0},"Cmds":[{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0}]},{"Param":{"ResistElem":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"StsBuffs":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"Abnmls":[{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0}],"AbnmlDisableBuff":{"Stack":[]},"AbnmlAddProbBuff":{"Stack":[]},"NextBuffs":[{"Step":0,"Val":0.0},{"Step":0,"Val":0.0},{"Step":0,"Val":0.0}],"WeakElemBonusBuff":{"Stack":[]},"HateChange":{"Stack":[]},"Regene":{"OwnerName":"","Turn":0,"Pow":0},"Barrier":0.0,"BarrierCount":0,"Hp":-1,"DeadDetail":0,"MCRecast":0},"Cmds":[{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0}]},{"Param":{"ResistElem":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"StsBuffs":[{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]},{"Stack":[]}],"Abnmls":[{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0},{"IsRegist":false,"Turn":0}],"AbnmlDisableBuff":{"Stack":[]},"AbnmlAddProbBuff":{"Stack":[]},"NextBuffs":[{"Step":0,"Val":0.0},{"Step":0,"Val":0.0},{"Step":0,"Val":0.0}],"WeakElemBonusBuff":{"Stack":[]},"HateChange":{"Stack":[]},"Regene":{"OwnerName":"","Turn":0,"Pow":0},"Barrier":0.0,"BarrierCount":0,"Hp":1583,"DeadDetail":0,"MCRecast":0},"Cmds":[{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0},{"UseNum":0,"RecastVal":0}]}],"PLJoinIdxs":[0,1,2,3,4,5]`

    var out bytes.Buffer
    r := zlib.NewWriter(&out)
    r.Write([]byte(sDec))
    r.Close()

    sEec := base64.RawStdEncoding.EncodeToString(out.Bytes())

    fmt.Println(sEec)
}
