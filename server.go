package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Common content types
const (
	ContentJSON = "application/json"
	ContentText = "text/plain"
)

// Server is the server
type Server struct {
}

// NewServer creates a new dispatch server
func NewServer() *Server {
	s := new(Server)

	http.HandleFunc("/ip", serveIPInfo)
	http.HandleFunc("/agent", serveUserAgentInfo)
	http.HandleFunc("/headers", serveHeaderInfo)
	http.HandleFunc("/", serveDefault)

	return s
}

// Run the server
func (s Server) Run(address string) {
	log.Infof("starting webserver on %s", address)
	log.Fatal(http.ListenAndServe(address, WriteLogHandler(http.DefaultServeMux)))
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
}

func (w *statusWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

// WriteLogHandler returns a server log handler
func WriteLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := statusWriter{w, 0, 0}

		// calculate the latency
		t := time.Now()
		handler.ServeHTTP(&writer, r)
		latency := time.Since(t)

		clientIP, _ := getClientIP(r)
		statusCode := writer.statusCode
		path := r.URL.Path
		method := r.Method
		log.Printf("%s - %s %s %d %v", clientIP, method, path, statusCode, latency)
	})
}

func serveIPInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, 404, "page not found")
		return
	}

	clientIP, err := getClientIP(r)
	if err != nil {
		respondError(w, r, 400, "could not parse client ip address")
		return
	}
	proxyList, _ := getClientProxyList(r)
	if err != nil {
		respondError(w, r, 400, "could not parse client proxy list")
		return
	}

	var msg string
	if r.Header.Get("Content-Type") == ContentJSON {
		w.Header().Add("Content-Type", ContentJSON)
		msg = fmt.Sprintf("{\"ip\": \"%s\", \"ip_proxy\": \"%v\"}", clientIP, proxyList)
	} else { // default is text response
		w.Header().Add("Content-Type", ContentText)
		msg = fmt.Sprintf("%s", clientIP)
	}

	w.WriteHeader(200)
	w.Write([]byte(msg))
}

func serveUserAgentInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, 404, "page not found")
		return
	}

	userAgent := r.UserAgent()

	var msg string
	if r.Header.Get("Content-Type") == ContentJSON {
		w.Header().Add("Content-Type", ContentJSON)
		msg = fmt.Sprintf("{\"user-agent\": \"%s\"}", userAgent)
	} else { // default is text response
		w.Header().Add("Content-Type", ContentText)
		msg = fmt.Sprintf("%s", userAgent)
	}

	w.WriteHeader(200)
	w.Write([]byte(msg))
}

func serveHeaderInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, 404, "page not found")
		return
	}

	headers, err := getRequestHeaders(r)
	if err != nil {
		respondError(w, r, 400, "could not parse request headers")
		return
	}

	var msg string
	if r.Header.Get("Content-Type") == ContentJSON {
		w.Header().Add("Content-Type", ContentJSON)
		jsonMsg, err := json.Marshal(headers)
		if err != nil {
			respondError(w, r, 400, "could not parse request headers")
			return
		}
		msg = string(jsonMsg)
	} else { // default is text response
		w.Header().Add("Content-Type", ContentText)
		var names []string
		for n := range headers {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, name := range names {
			msg = fmt.Sprintf("%s%s: %s\n", msg, name, headers[name])
		}
	}

	w.WriteHeader(200)
	w.Write([]byte(msg))
}

func serveDefault(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveIPInfo(w, r)
		return
	}
	respondError(w, r, 404, "page not found")
}

func respondError(w http.ResponseWriter, r *http.Request, code int, message string, a ...interface{}) {
	var msg string
	if r.Header.Get("Content-Type") == ContentJSON {
		w.Header().Add("Content-Type", ContentJSON)
		m := fmt.Sprintf(message, a...)
		msg = fmt.Sprintf("{\"status\": \"error\", \"message\": \"%s\"}", m)
	} else { // default is text response
		w.Header().Add("Content-Type", ContentText)
		m := fmt.Sprintf(message, a...)
		msg = fmt.Sprintf("%d %s", code, m)
	}

	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func splitIPList(ipList string) []string {
	ips := strings.Split(ipList, ", ")
	var list []string
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if len(ip) > 0 {
			list = append(list, ip)
		}
	}
	return list
}

func getClientProxyList(r *http.Request) ([]string, error) {
	rawProxyList := splitIPList(r.Header.Get("X-Forwarded-For"))
	var proxyList []string
	for _, hostport := range rawProxyList {
		ip, _, err := net.SplitHostPort(hostport)
		if err != nil {
			proxyList = append(proxyList, "*")
			continue
		}
		proxyList = append(proxyList, ip)
	}
	return proxyList, nil
}

func getClientIP(r *http.Request) (string, error) {
	// first, figure out the correct IP to use
	clientHostPort := r.RemoteAddr
	proxyList := splitIPList(r.Header.Get("X-Forwarded-For"))
	if len(proxyList) > 0 {
		clientHostPort = proxyList[0]
	}

	// clean it up
	clientIP, _, err := net.SplitHostPort(clientHostPort)
	if err != nil {
		clientIP = clientHostPort
	}

	return clientIP, nil
}

func getRequestHeaders(r *http.Request) (map[string]string, error) {
	var headers map[string]string
	headers = make(map[string]string)
	for name, values := range r.Header {
		name = strings.ToLower(name)
		var value string
		for i, v := range values {
			if i > 0 {
				value = fmt.Sprintf("%s,%s", value, v)
			} else {
				value = v
			}
		}
		headers[name] = value
	}
	return headers, nil
}
