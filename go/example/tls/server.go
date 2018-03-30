package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/meilihao/water"
)

func main() {
	router := water.NewRouter()
	router.Get("/", _Home)

	w := router.Handler()

	pool := x509.NewCertPool()
	caCertPath := "rootca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		panic(err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	server := &http.Server{
		Addr:    ":443",
		Handler: w,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatalln(err)
	}
}

func _Home(ctx *water.Context) {
	// get clinet crt
	tls := ctx.Request.TLS
	if tls != nil && len(tls.PeerCertificates) > 0 {
		fmt.Printf("Client Crt: %+v\n", tls.PeerCertificates[0].Subject)
	}

	ctx.WriteString(ctx.Request.URL.String())
}
