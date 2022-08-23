package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"gopkg.in/yaml.v2"
)

// Config File Format
type Config struct {
	LocalHostPort  string `yaml:"localHostPort,omitempty"`
	ProxyServerURI string `yaml:"proxyServerURI,omitempty"`
	TargetURL      string `yaml:"targetURL,omitempty"`
	//TODO: If local TLS is needed
	//ServeTLSLocal  bool
	//CACertChainFile string
	//CAKeyFile string
}

var proxyConfig Config

// Handlers
func h1(w http.ResponseWriter, req *http.Request) {

	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	var buf bytes.Buffer

	io.Copy(&buf, req.Body)

	req2, err := http.NewRequest(req.Method, proxyConfig.ProxyServerURI, &buf)
	if err != nil {
		log.Fatalf("Failed to create a new http request:%v", err)
	}
	requestString := proxyConfig.TargetURL + req.URL.RequestURI()
	log.Println("RequestString:%s", requestString)
	for k, _ := range req.Header {
		for _, v := range req.Header[k] {
			req2.Header.Add(k, v)
		}
	}
	req2.Header["X-PROXY-DESTINATION"] = []string{requestString}

	resp, err := client.Do(req2)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("%v", err))
		resp.Body.Close()
		return
	}
	for k, _ := range resp.Header {
		for _, v := range resp.Header[k] {
			w.Header().Add(k, v)
		}
	}
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func (config *Config) initDefaults() {
	if config.LocalHostPort == "" {
		config.LocalHostPort = ":8080"
	}
	if _, err := url.Parse(config.ProxyServerURI); err != nil || config.ProxyServerURI == "" {
		log.Fatalf("Missing/invalid proxyServerURI %s in configFile", config.ProxyServerURI)
	}
	if _, err := url.Parse(config.TargetURL); err != nil || config.TargetURL == "" {
		log.Fatalf("Missing/invalid targetURL %s in configFile", config.TargetURL)
	}
}

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "wsp_local_client.cfg", "Config file location")
	flag.Parse()
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading config file:%v", err)
	}
	if err = yaml.Unmarshal(configData, &proxyConfig); err != nil {
		log.Fatalf("Error parsing config file:%v", err)
	}
	proxyConfig.initDefaults()

	log.Println(proxyConfig.LocalHostPort + "<==>" + proxyConfig.ProxyServerURI + "<==>" + proxyConfig.TargetURL)
	srv := http.Server{
		Addr: proxyConfig.LocalHostPort,
	}
	http.HandleFunc("/", h1)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener
		log.Fatalf("HTTP server ListenAndServeTLS: %v", err)
	}

	<-idleConnsClosed
}
