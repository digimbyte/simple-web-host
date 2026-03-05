package main

import (
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

var lastPing atomic.Int64

const heartbeat = `<script>(()=>{setInterval(()=>fetch("/__ping"),3000)})();</script>`

func main() {
	exe, _ := os.Executable()
	exe, _ = filepath.EvalSymlinks(exe)
	root := http.Dir(filepath.Dir(exe))
	fs := http.FileServer(root)
	port := 8080

	lastPing.Store(time.Now().UnixMilli())

	mux := http.NewServeMux()

	mux.HandleFunc("/__ping", func(w http.ResponseWriter, r *http.Request) {
		lastPing.Store(time.Now().UnixMilli())
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path
		if strings.HasSuffix(name, "/") {
			name += "index.html"
		}
		if strings.HasSuffix(name, ".html") {
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
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	for err != nil && port < 8100 {
		port++
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Println("Simple Web Host")
	fmt.Printf("Serving: %s\n", filepath.Dir(exe))
	fmt.Printf("    on: %s\n", url)
	fmt.Println("\nThis window will close automatically when the browser tab is closed.")

	go watchdog()
	openBrowser(url)
	http.Serve(ln, mux)
}

func watchdog() {
	time.Sleep(15 * time.Second) // grace period for initial browser load
	for {
		time.Sleep(3 * time.Second)
		if time.Since(time.UnixMilli(lastPing.Load())) > 10*time.Second {
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
