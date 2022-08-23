package client

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client connects to one or more Server using HTTP websockets.
// The Server can then send HTTP requests to execute.
type Client struct {
	Config *Config

	client *http.Client
	dialer *websocket.Dialer
	pools  map[string]*Pool
}

// NewClient creates a new Client.
func NewClient(config *Config) (c *Client) {
	c = new(Client)
	c.Config = config
	c.client = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	c.dialer = &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		//func(*http.Request) (*url.URL, error) {
		//	return url.Parse(os.Getenv("HTTP_PROXY"))
		//},
		// TODO Remove: Skip server cert validation for now
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.pools = make(map[string]*Pool)
	return
}

// Start the Proxy
func (c *Client) Start(ctx context.Context) {
	for _, target := range c.Config.Targets {
		pool := NewPool(c, target, c.Config.SecretKey)
		c.pools[target] = pool
		go pool.Start(ctx)
	}
}

// Shutdown the Proxy
func (c *Client) Shutdown() {
	for _, pool := range c.pools {
		pool.Shutdown()
	}
}
