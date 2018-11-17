package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Symbol struct {
}

func main() {
	url := "http://www.nasdaq.com/screening/companies-by-industry.aspx?exchange=NASDAQ&render=download"

	timeout := time.Duration(5) * time.Second
	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	records, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(records) < 1 {
		fmt.Println("can't find any records in response body")
		return
	}

	symbolIdx := -1
	industryIdx := -1
	for i, str := range records[0] {
		if str == "Symbol" {
			symbolIdx = i
		} else if str == "Industry" {
			industryIdx = i
		}
	}
	if symbolIdx == -1 {
		fmt.Println("cant's find Symbol header")
		return
	}

	if industryIdx == -1 {
		fmt.Println("can't find Industry header")
		return
	}

	for _, strings := range records {

		for _, string := range strings {
			fmt.Printf("%v, ", string)
		}
		fmt.Println("")
	}
}
