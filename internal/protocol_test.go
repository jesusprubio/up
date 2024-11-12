package internal

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestHTTPProbe(t *testing.T) {
	tout := 1 * time.Second
	server := newTestHTTPServer(t)
	defer server.Close()
	t.Run(
		"returns the URL if the request is successful",
		func(t *testing.T) {
			u := url.URL{Scheme: "http", Host: server.Addr}
			proto := HTTP{Timeout: tout}
			got, extra, err := proto.Probe(u.String())
			if err != nil {
				t.Fatal(err)
			}
			want := "http://127.0.0.1:8080"
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
			if extra != "200 OK" {
				t.Fatalf("got %q, want %q", extra, "200 OK")
			}
		},
	)
	t.Run("returns an error if the request fails", func(t *testing.T) {
		u := url.URL{Scheme: "http", Host: "localhost"}
		proto := HTTP{Timeout: 1}
		got, extra, err := proto.Probe(u.String())
		if got != "" {
			t.Fatalf("got %q should be zero", got)
		}
		got = err.Error()
		want := `Get "http://localhost": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if extra != "" {
			t.Fatalf("got %q should be zero", extra)
		}
	})
}

// Creates an HTTP server for testing.
func newTestHTTPServer(t *testing.T) *http.Server {
	hostPort := net.JoinHostPort("127.0.0.1", "8080")
	server := &http.Server{Addr: hostPort}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong\n")
	})
	l, err := net.Listen("tcp", hostPort)
	if err != nil {
		t.Fatalf("create listener %v", err)
	}
	go func() {
		err := server.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("starting http server: %v", err)
		}
	}()
	return server
}

func TestTCPProbe(t *testing.T) {
	tout := 1 * time.Second
	listen := newTestTCPServer(t)
	defer listen.Close()
	go func() {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		go func(conn net.Conn) {
			conn.Write([]byte("pong"))
			conn.Close()
		}(connection)
	}()
	hostPort := listen.Addr().String()
	t.Run(
		"returns the remote host/port if the request is successful",
		func(t *testing.T) {
			proto := &TCP{Timeout: tout}
			got, extra, err := proto.Probe(hostPort)
			if err != nil {
				t.Fatal(err)
			}
			if got != hostPort {
				t.Fatalf("got %q, want %q", got, hostPort)
			}
			host, port, err := net.SplitHostPort(extra)
			if err != nil {
				t.Fatal(err)
			}
			if host != "127.0.0.1" {
				t.Fatalf("got %q, want %q", host, "127.0.0.1")
			}
			if port == "" {
				t.Fatalf("got %q, want a valid port", port)
			}
		},
	)
	t.Run("returns an error if the request fails", func(t *testing.T) {
		proto := &TCP{Timeout: 1}
		got, extra, err := proto.Probe("localhost:80")
		if err == nil {
			t.Fatal("got nil, want an error")
		}
		if got != "" {
			t.Fatalf("got %q should be zero", got)
		}
		got = err.Error()
		want := "dial tcp: lookup localhost: i/o timeout"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if extra != "" {
			t.Fatalf("got %q should be zero", extra)
		}
	})
}

// Creates a TCP server for testing.
func newTestTCPServer(t *testing.T) net.Listener {
	hostPort := net.JoinHostPort("127.0.0.1", "8081")
	listen, err := net.Listen("tcp", hostPort)
	if err != nil {
		t.Fatalf("starting tcp server: %v", err)
	}
	return listen
}

func TestDNSProbe(t *testing.T) {
	tout := 1 * time.Second
	// TODO(#31): Implement a simple DNS server to test this.
	// We need to support custom resolvers first.
	t.Run(
		"returns the domain if the request is successful",
		func(t *testing.T) {
			proto := &DNS{Timeout: tout}
			domain := "google.com"
			got, extra, err := proto.Probe(domain)
			if err != nil {
				t.Fatal(err)
			}
			if got != domain {
				t.Fatalf("got %q, want %q", got, domain)
			}
			if !net.ParseIP(extra).IsGlobalUnicast() {
				t.Fatalf("got %q, want a valid IP address", extra)
			}
		},
	)
	t.Run("returns an error if the request fails", func(t *testing.T) {
		proto := &DNS{Timeout: 1}
		got, extra, err := proto.Probe("invalid.aa")
		if err == nil {
			t.Fatal("got nil, want an error")
		}
		if got != "" {
			t.Fatalf("got %q should be zero", got)
		}
		got = err.Error()
		want := "lookup invalid.aa: no such host"
		if os.Getenv("CI") == "true" {
			want = "lookup invalid.aa on 127.0.0.53:53: no such host"
		}
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if extra != "" {
			t.Fatalf("got %q should be zero", extra)
		}
	})
}
