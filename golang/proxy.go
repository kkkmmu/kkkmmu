package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	//"regexp"
	"bytes"
	"fmt"
	"regexp"
)

var PermitStatistics = make(map[string]int, 100)
var RejectStatistics = make(map[string]int, 100)

var PermitHttpsSiteList = `.*{START|google|yowindow|dasanzhone|dasannetworks|bitbucket|github|broadcom|okta|youtube|ggpht|END}{1}.*`
var PermitHttpSiteList = `.*{START|dasanzhone|dasannetworks|broadcom|okta|END}{1}.*`

func Permit(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	/*
		if _, ok := PermitStatistics[r.RemoteAddr]; !ok {
			PermitStatistics[r.RemoteAddr] = 1
		} else {
			PermitStatistics[r.RemoteAddr] += 1
		}
	*/
	fmt.Printf("Permit request[%s(%d)] %s:%s\n", r.RemoteAddr, RejectStatistics[r.RemoteAddr], r.Method, r.URL)
	return r, nil
}

func Reject(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	/*
		if _, ok := RejectStatistics[r.RemoteAddr]; !ok {
			RejectStatistics[r.RemoteAddr] = 1
		} else {
			RejectStatistics[r.RemoteAddr] += 1
		}
	*/
	fmt.Printf("Reject request[%s(%d)] %s:%s\n", r.RemoteAddr, RejectStatistics[r.RemoteAddr], r.Method, r.URL)
	//fmt.Printf("Reject request[%s(%d)] %s:%s\n   Cookie:%s\n", r.RemoteAddr, RejectStatistics[r.RemoteAddr], r.Method, r.URL, r.Cookie)
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeHtml, http.StatusForbidden,
		fmt.Sprintf(`<h1 style="font-size:50px;color:red" align="center">Access to %s is not permmitted by proxy!</h1>
                <h1 style="font-size:50px;color:red" align="center">Do not use proxy to access local network!</h1>
                `, r.URL))
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	buf := new(bytes.Buffer)
	proxy.Logger = log.New(buf, "", 0)

	proxy.OnRequest(goproxy.Not(goproxy.ReqHostMatches(regexp.MustCompile(PermitHttpSiteList)))).DoFunc(Reject)
	proxy.OnRequest(goproxy.Not(goproxy.ReqHostMatches(regexp.MustCompile(PermitHttpsSiteList)))).HandleConnect(goproxy.AlwaysReject)
	log.Fatal(http.ListenAndServe(":9080", proxy))
}
