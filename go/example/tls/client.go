package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	pool := x509.NewCertPool()
	caCertPath := "rootca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		panic(err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	clientCrt, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{clientCrt},
		},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost")
	if err != nil {
		panic(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
