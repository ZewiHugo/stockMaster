package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gocarina/gocsv"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

type Symbol struct {
	Symbol    string  `csv:"Symbol"`
	Name      string  `csv:"Name"`
	MarketCap float64 `csv:"MarketCap"`
	Sector    string  `csv:"Sector"`
	Industry  string  `csv:"Industry`
	IPOYear   string  `csv:"IPOyear"`
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

	fmt.Printf("%+v\n", symbols[0])
	fmt.Printf("%+v\n", symbols[0])

	cluster := gocql.NewCluster("127.0.0.1:9042")
	keyspaceInfo := KeyspaceInfo{
		"stock_master",
		1,
		"SimpleStrategy",
	}
	err = createAndUseKeyspace(cluster, keyspaceInfo)
	if err != nil {
		fmt.Println(errors.WithStack(err))
		return
	}
}

type KeyspaceInfo struct {
	KeyspaceName      string
	ReplicationFactor int
	KeyspaceClass     string
}

func createAndUseKeyspace(cluster *gocql.ClusterConfig, keyspaceInfo KeyspaceInfo) error {
	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 20 * time.Second
	session, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	box := rice.MustFindBox("cassandra/template/")
	boxString, err := box.String("createKeyspace.tmpl")
	if err != nil {
		return err
	}
	t, err := template.New("createKeyspace").Parse(boxString)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, keyspaceInfo)
	if err != nil {
		return err
	}
	if err := session.Query(buf.String()).Exec(); err != nil {
		return err
	}

	return nil
}
