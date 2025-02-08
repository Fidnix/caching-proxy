package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
)

type Response struct {
	status int
	header string
	body   string
}

var (
	cache  = make(map[string]Response)
	port   int64
	origin string
)

func main() {

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Value:       3000,
				Usage:       "port to deploy the caching proxy",
				Aliases:     []string{"p"},
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "origin",
				Usage:       "url to cache",
				Aliases:     []string{"o"},
				Destination: &origin,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			router := gin.Default()
			router.Any("/*request", ProxyEndpoint)
			router.Run(fmt.Sprintf(":%d", port))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func ProxyEndpoint(c *gin.Context) {
	request := c.Param("request")

	if val, ok := cache[request]; ok {
		c.Header("X-Cache", "HIT")
		lines := strings.Split(val.header, "\r\n")
		for _, line := range lines {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				c.Header(key, value)
			}
		}
		c.String(val.status, val.body)
		return
	}

	u, err := url.Parse(origin)
	if err != nil {
		fmt.Println("Error building url:", err)
		return
	}
	u.Path = request
	u.RawQuery = c.Request.URL.RawQuery
	u.Fragment = c.Request.URL.Fragment

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var header string
	for key, values := range resp.Header {
		for _, value := range values {
			header += fmt.Sprintf("%s: %s\r\n", key, value)
			c.Header(key, value)
		}
	}

	c.Header("X-Cache", "MISS")
	cache[request] = Response{
		header: header,
		status: resp.StatusCode,
		body:   string(body),
	}
	c.String(resp.StatusCode, string(body))
}
