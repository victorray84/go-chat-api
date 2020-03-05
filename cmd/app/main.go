package main

import (
	"fmt"
	"github.com/leartgjoni/go-chat-api/http"
	"github.com/leartgjoni/go-chat-api/websocket"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/signal"
)

func main() {
	m := NewMain()

	// Load configuration.
	if err := m.LoadConfig(); err != nil {
		_, _ = fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	// Execute program.
	if err := m.Run(); err != nil {
		_, _ = fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	_, _ = fmt.Fprintln(m.Stdout, "received interrupt, shutting down...")
	_ = m.Close()
}

// Main represents the main program execution.
type Main struct {
	ConfigPath string
	Config     Config

	// Input/output streams
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	closeFn func() error
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		closeFn: func() error { return nil },
	}
}

// Close cleans up the program.
func (m *Main) Close() error { return m.closeFn() }

// LoadConfig parses the configuration file.
func (m *Main) LoadConfig() error {

	if os.Getenv("CONFIG_PATH") != "" {
		m.ConfigPath = os.Getenv("CONFIG_PATH")
	} else {
		m.ConfigPath = ".env"
	}

	viper.SetConfigFile(m.ConfigPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	m.Config = Config{
		// Todo
	}

	return nil
}

func (m *Main) Run() error {
	// Todo: connect db

	clientService := websocket.NewClientService()
	hubService := websocket.NewHubService()

	// Initialize Http server.
	httpServer := http.NewServer()
	httpServer.Addr = ":8080"

	httpServer.ClientService = clientService
	httpServer.HubService = hubService

	// Start HTTP server.
	if err := httpServer.Start(); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(m.Stdout, "Listening on port: %s\n", httpServer.Addr)

	// Assign close function.
	m.closeFn = func() error {
		_ = httpServer.Close()
		// Todo: close db
		return nil
	}

	return nil
}

type Config struct {
	// Todo
}