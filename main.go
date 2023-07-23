package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"tailscale.com/client/tailscale"
	"tailscale.com/tsnet"
)

func setAuthHeader(lc *tailscale.LocalClient, req *http.Request) error {
	identity, err := lc.WhoIs(req.Context(), req.RemoteAddr)
	if err != nil {
		return err
	}

	username := identity.UserProfile.LoginName
	username = regexp.MustCompile(`[^a-zA-Z0-9 _-]+`).ReplaceAllString(username, "_")

	log.Println("user connected:", username)
	req.SetBasicAuth(username, "tailscale-authenticated")

	return nil
}

func main() {
	var (
		hostname = flag.String("hostname", "calibre", "hostname for this service on tailnet")
		origin   = flag.String("origin", "http://localhost:8080", "calibre server root URL")
	)
	flag.Parse()

	// check arguments
	originURL, err := url.Parse(*origin)
	if err != nil {
		log.Fatal(err)
	}

	// set up tailscale node
	srv := new(tsnet.Server)
	srv.Hostname = *hostname

	// set up tailscale client
	lc, err := srv.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	// start listening
	defer srv.Close()
	ln, err := srv.ListenTLS("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	// serve
	proxy := httputil.NewSingleHostReverseProxy(originURL)
	log.Fatal(http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := setAuthHeader(lc, r)
		if err != nil {
			http.Error(w, "tailscale authorization check failed", http.StatusNetworkAuthenticationRequired)
			return
		}

		proxy.ServeHTTP(w, r)
	})))
}
