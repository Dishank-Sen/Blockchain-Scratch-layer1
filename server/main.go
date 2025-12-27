package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"path"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/quic-go/quic-go"
)

func loadCert() *tls.Config{
	certFilePath := path.Join("certificate", "server", "server.crt")
	keyFilePath := path.Join("certificate", "server", "server.key")

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)

	if err != nil{
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-example-v1"},
	}
}

func handleSession(sess *quic.Conn) {
	defer func() {
		_ = sess.CloseWithError(0, "closing")
	}()

	remote := sess.RemoteAddr().String()
	log.Printf("New session from %s\n", remote)

	// Accept streams in a loop
	for {
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			// AcceptStream returns non-nil err when session is closed
			log.Printf("AcceptStream error (%s): %v\n", remote, err)
			return
		}
		go handleStream(stream)
	}
}

func handleStream(s *quic.Stream) {
	defer s.Close()
	buf := make([]byte, 4096)

	// Simple single-read echo example; production should loop and frame messages
	n, err := s.Read(buf)
	if err != nil {
		log.Printf("stream read error: %v\n", err)
		return
	}
	msg := string(buf[:n])
	log.Printf("Received on stream: %s\n", msg)

	// Echo back with a timestamp
	_, _ = s.Write([]byte("pong: " + time.Now().Format(time.RFC3339)))
}

func main(){
	tlsConf := loadCert()

	quicConf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	addr := "127.0.0.1:4242" // UDP port
	listener, err := quic.ListenAddr(addr, tlsConf, quicConf)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("QUIC server listening on %s\n", addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Accept sessions forever
	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("server shutdown complete")
				return
			}
			logger.Info(fmt.Sprintf("listener error: %v\n", err))
			return
		}
		go handleSession(sess)
	}
}