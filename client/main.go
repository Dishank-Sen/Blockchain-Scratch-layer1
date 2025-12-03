package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"time"

	quic "github.com/quic-go/quic-go"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,        // trust every certificate
		ServerName: "localhost",        // must match SAN DNS entry in server cert
		NextProtos: []string{"quic-example-v1"},
	}

	quicConf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	addr := "127.0.0.1:4242"

	// Dial the server (performs UDP + QUIC + TLS handshake)
	session, err := quic.DialAddr(context.Background(), addr, tlsConf, quicConf)
	if err != nil {
		log.Fatalf("DialAddr failed: %v", err)
	}
	defer func() {
		_ = session.CloseWithError(0, "client closing")
	}()

	// Open a bidirectional stream
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalf("OpenStreamSync failed: %v", err)
	}
	defer stream.Close()

	// Send a simple message (application-level)
	msg := "hello from client"
	_, err = stream.Write([]byte(msg))
	if err != nil {
		log.Fatalf("stream write failed: %v", err)
	}

	// Read reply (simple echo example)
	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			log.Printf("server reply: %s", string(buf[:n]))
		}
		if err != nil {
			if err == io.EOF { break }          // normal close by server
			log.Fatalf("stream read failed: %v", err)
		}
	}

	// Optionally open additional streams for parallel work
}
