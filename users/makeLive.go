package users

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"path"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/internal/identity"
	"github.com/quic-go/quic-go"
)

func loadCert() (*tls.Config, error){
	certFilePath := path.Join("certificate", "server", "server.crt")
	keyFilePath := path.Join("certificate", "server", "server.key")

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)

	if err != nil{
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-example-v1"},
	}
	return config, nil
}

func handleSession(sess *quic.Conn, addr string) {
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
		go handleStream(stream, addr)
	}
}

func handleStream(s *quic.Stream, addr string) {
	defer removeIdentity(addr)
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

func saveIdentity(addr string) error{
	privatePEM, publicPEM, err := identity.GenerateKeyPairPEM()
	if err != nil{
		return err
	}
	dir := path.Join("internal", "storage", addr)
	privPath := path.Join(dir, "private.key")
    pubPath  := path.Join(dir, "public.key")

    if _, err := os.Stat(privPath); err == nil {
        if _, err := os.Stat(pubPath); err == nil {
            return nil
        }
    }

	return identity.SaveKeyPair(dir, privatePEM, publicPEM)
}

func removeIdentity(addr string) error{
	path := path.Join("internal", "storage", addr)
	return os.RemoveAll(path)
}

func MakeLive(ctx context.Context, addr string) error{
	tlsConf, err := loadCert()
	if err != nil{
		return err
	}

	quicConf := &quic.Config{
		MaxIdleTimeout: 24 * time.Hour,
	}

	listener, err := quic.ListenAddr(addr, tlsConf, quicConf)

	if err != nil {
		return err
	}
	log.Printf("QUIC server listening on %s\n", addr)

	err = saveIdentity(addr)
	if err != nil{
		return err
	}

	// Accept sessions forever
	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			log.Printf("listener accept error: %v\n", err)
			continue
		}
		go handleSession(sess, addr)
	}
}