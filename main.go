package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func checkIP(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	now := time.Now().Unix()
	fmt.Fprintf(w, `{"timestamp":%d,"status":"healthy"}`, now)
}

func portHandler(headerName *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		port := r.PathValue("port")
		var host string
		if *headerName != "" {
			hosts := r.Header.Get(*headerName)
			host = strings.Split(hosts, ", ")[0]
		}
		if host == "" {
			host = r.RemoteAddr
		}
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		fmt.Println(host)
		portState := checkIP(host, port)
		if portState {
			w.Write([]byte("1"))
		} else {
			w.Write([]byte("0"))
		}
	}
}

func main() {
	headerName := flag.String("realipheader", "", "name of the ip address header")
	flag.Parse()
	http.HandleFunc("/{port}", portHandler(headerName))
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
