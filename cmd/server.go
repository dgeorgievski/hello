/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgeorgievski/hello/pkg/version"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func appVersion(c *gin.Context) {
	c.JSON(http.StatusOK,
		version.GetVersion())
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func hey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "joe",
	})
}

func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func readyz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}


// requestLogger returns a Gin middleware that logs each request with method,
// path, status code, and latency using the standard library log package.
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("%s %s %d %s", c.Request.Method, c.Request.RequestURI, c.Writer.Status(), latency)
	}
}

// newRouter constructs the application router with all registered routes.
// Exported for testability so HTTP handlers can be exercised without opening a network port.
func newRouter() *gin.Engine {
	r := gin.Default()
	r.Use(requestLogger())
	r.GET("/hello", hello)
	r.GET("/version", appVersion)
	r.GET("/ping", ping)
	r.GET("/hey", hey)
	r.GET("/healthz", healthz)
	r.GET("/readyz", readyz)
	r.GET("/greet", greet)
	return r
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Start the Hello HTTP server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := &http.Server{
			Addr:              ":8080",
			Handler:           newRouter(),
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			return fmt.Errorf("serve HTTP on %s: %w", server.Addr, err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when
	// this action is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func greet(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = "World"
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Hello, %s!", name),
	})
}
