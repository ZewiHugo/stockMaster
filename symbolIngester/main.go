package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

type Symbol struct {
	Symbol  string `csv:"Symbol"`
	Sector  string `csv:"Sector"`
	IPOyear string `csv:"IPOyear"`
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
		fmt.Println(errors.WithStack(err))
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	var symbols []Symbol
	err = gocsv.Unmarshal(resp.Body, &symbols)
	if err != nil {
		fmt.Println(errors.WithStack(err))
		return
	}

	for _, symbol := range symbols {
		fmt.Println(symbol.Symbol, symbol.Sector, symbol.IPOyear)
	}
}
