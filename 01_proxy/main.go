package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {

	log.Printf("[INFO] starting proxy server")
	if err := http.ListenAndServe(":8989", http.HandlerFunc(handleRequest)); err != nil {
		panic(err)
	}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	log.Printf("[INFO] %s target: %v, client: %v", r.Method, r.RequestURI, r.RemoteAddr)

	switch r.Method {
	case http.MethodHead, http.MethodGet,
		http.MethodPost, http.MethodPut,
		http.MethodDelete, http.MethodPatch:
		// expected http requests
	case http.MethodConnect:
		handleConnect(w, r)
	default:
		log.Printf("[WARNING] unexpected method %s", r.Method)
	}
	handleDefault(w, r)

}

func handleConnect(w http.ResponseWriter, r *http.Request) {

	// if !r.URL.IsAbs() {
	// 	r.URL.Scheme = "https"
	// }

	connToTarget, err := net.Dial("tcp", r.RequestURI)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	hj, ok := w.(http.Hijacker)
	if !ok {
		err := fmt.Errorf("http server doesn't support hijacking connection")
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	connToClient, _, err := hj.Hijack()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	go tunnelConn(connToClient, connToTarget)
	go tunnelConn(connToTarget, connToClient)

}

func tunnelConn(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	_, err := io.Copy(dst, src)
	if err != nil {
		// log.Println(err)
	}
}

func handleDefault(w http.ResponseWriter, r *http.Request) {

	proxyReq := copyRequest(r)

	client := http.Client{}
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] %s", proxyResp.Status)

	copyResponseHeader(w.Header(), proxyResp.Header)

	if _, err := io.Copy(w, proxyResp.Body); err != nil {
		log.Println(w, err, http.StatusInternalServerError)
		return
	}
}

func copyRequest(prev *http.Request) *http.Request {
	next := prev.Clone(prev.Context())

	// According to RFC 2612, the following headers are treated as hop-by-hop by default. Therefore, we delete them
	next.Header.Del("Keep-Alive")
	next.Header.Del("Transfer-Encoding")
	next.Header.Del("TE")
	next.Header.Del("Connection")
	next.Header.Del("Trailer")
	next.Header.Del("Upgrade")
	next.Header.Del("Proxy-Authorization")
	next.Header.Del("Proxy-Authenticate1")

	next.Header.Add("X-Forwarded-For", prev.RemoteAddr)

	next.RequestURI = ""

	return next
}

func copyResponseHeader(dst, src http.Header) {
	for key, vals := range src {
		for _, val := range vals {
			dst.Add(key, val)
		}
	}
}

func handleError(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println("[ERROR]", err)
}
