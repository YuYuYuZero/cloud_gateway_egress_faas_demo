package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "NowTime: %s", time.Now().Format("2006-01-02 15:04:05"))
}

func host(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Host: %v\n\r", req.Host)
	fmt.Fprintf(w, "NowTime: %v", time.Now().Format("2006-01-02 15:04:05"))
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "TZ:%s", os.Getenv("TZ"))
	fmt.Fprintf(w, "NowTime: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func gatewayTest(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(401)
	fmt.Fprintf(w, "gateway_test, req_url:%v, ts:%s", req.URL, time.Now().Format("2006-01-02 15:04:05"))
}

func gatewayNetworkTest(w http.ResponseWriter, req *http.Request) {
	protocol := "http"
	protocols := req.URL.Query()["protocol"]
	if len(protocols) > 0 && protocols[0] == "https" {
		protocol = "https"
	}
	urls := req.URL.Query()["url"]
	if len(urls) == 0 {
		fmt.Fprintf(w, "non url param")
		return
	}

	fmt.Fprintf(w, AccessOpenApiUrl(protocol, urls[0]))
}

func gatewaySleepTest(w http.ResponseWriter, req *http.Request) {
	sleepTimes := req.URL.Query()["sleep_time"]
	if len(sleepTimes) == 0 {
		fmt.Fprintf(w, "invalid param")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sleepTime := sleepTimes[0]
	w.Header().Set("sleep_time", sleepTime)

	if len(sleepTime) > 0 {
		sleepTimeInt, err := strconv.Atoi(sleepTime)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		time.Sleep(time.Duration(sleepTimeInt) * time.Second)
	}

	fmt.Fprintf(w, "ok")
	w.WriteHeader(http.StatusOK)
	return
}

func main() {

	//go RunFaasCliLoop(map[string]string{
	//	"http_plb":      "http://dev.douyincloud.gateway.egress.ivolces.com",
	//	"https_plb":     "https://dev.douyincloud.gateway.egress.ivolces.com",
	//	"http_openapi":  "http://developer.toutiao.com",
	//	"https_openapi": "https://developer.toutiao.com",
	//	"http_openapi2":  "http://open.douyin.com",
	//	"https_openapi2": "https://open.douyin.com",
	//})

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/host", host)
	http.HandleFunc("/v1/ping", ping)
	http.HandleFunc("/gateway_test", gatewayTest)
	http.HandleFunc("/gateway_network_test", gatewayNetworkTest)
	http.HandleFunc("/gateway_dns_test", gatewayDnsTest)
	http.HandleFunc("/gateway_sleep_test", gatewaySleepTest)
	http.HandleFunc("/gateway_ws_push", gatewayWsPush)
	http.HandleFunc("/gateway_ws_handle", gatewayWsHandle)

	http.ListenAndServe(":8000", nil)
}

//func counterParamFormatDycFullLinkCallBackend(httpStatusCode *int, trafficSource, clientOs, serverPath, availableZone *string) {
//	if trafficSource != nil && len(*trafficSource) == 0 {
//		*trafficSource = "-"
//	}
//	if clientOs != nil && len(*clientOs) == 0 {
//		*clientOs = "-"
//	}
//	if serverPath != nil && len(*serverPath) == 0 {
//		*serverPath = "-"
//	}
//	if availableZone != nil && len(*availableZone) == 0 {
//		*availableZone = "-"
//	}
//}
