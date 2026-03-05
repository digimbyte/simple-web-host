package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var activeClients atomic.Int32

const heartbeat = `<script>(()=>{new EventSource("/__ping")})();</script>`

func main() {
	portFlag := flag.Int("port", 0, "port to listen on (default: auto-select from 8080-8099)")
	stable := flag.Bool("stable", false, "persistent server: no auto-close, no heartbeat injection")
	flag.Parse()

	exe, _ := os.Executable()
	exe, _ = filepath.EvalSymlinks(exe)
	root := http.Dir(filepath.Dir(exe))
	fs := http.FileServer(root)
	port := *portFlag

	mux := http.NewServeMux()

	if !*stable {
		mux.HandleFunc("/__ping", func(w http.ResponseWriter, r *http.Request) {
			flusher, ok := w.(http.Flusher)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			activeClients.Add(1)
			defer activeClients.Add(-1)

			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-r.Context().Done():
					return
				case <-ticker.C:
					fmt.Fprint(w, ": keepalive\n\n")
					flusher.Flush()
				}
			}
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path
		if strings.HasSuffix(name, "/") {
			name += "index.html"
		}
		if !*stable && strings.HasSuffix(name, ".html") {
			if f, err := root.Open(name); err == nil {
				defer f.Close()
				if data, err := io.ReadAll(f); err == nil {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					s := string(data)
					if i := strings.LastIndex(s, "</body>"); i >= 0 {
						fmt.Fprint(w, s[:i], heartbeat, s[i:])
					} else {
						fmt.Fprint(w, s, heartbeat)
					}
					return
				}
			}
		}
		fs.ServeHTTP(w, r)
	})

	// Find an open port
	var ln net.Listener
	var err error
	if port > 0 {
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Port %d unavailable: %v\n", port, err)
			os.Exit(1)
		}
	} else {
		port = 8080
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		for err != nil && port < 8100 {
			port++
			ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "No available port in range 8080-8099")
			os.Exit(1)
		}
	}

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Println("Simple Web Host")
	if *stable {
		fmt.Println("Mode:   stable (persistent)")
	}
	fmt.Printf("Serving: %s\n", filepath.Dir(exe))
	fmt.Printf("    on: %s\n", url)

	if !*stable {
		fmt.Println("\nThis window will close automatically when the browser tab is closed.")
		go watchdog()
	} else {
		fmt.Println("\nPress Ctrl+C to stop the server.")
	}

	openBrowser(url)
	http.Serve(ln, mux)
}

func watchdog() {
	time.Sleep(15 * time.Second) // grace period for initial browser load
	for {
		time.Sleep(3 * time.Second)
		if activeClients.Load() == 0 {
			os.Exit(0)
		}
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
