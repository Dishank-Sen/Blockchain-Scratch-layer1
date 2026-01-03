package cli

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/quic-go/quic-go"
	"github.com/spf13/cobra"
)

func init(){
	Register("bind", Bind)
}

func Bind() *cobra.Command{
	cmd := &cobra.Command{
		Use: "bind",
		Short: "establish a p2p connection with the given address and port",
		RunE: bindRunE,
		Example: "bloc bind -a 127.0.0.1:8080",
	}

	cmd.Flags().StringP("addr", "a", "", "ip address and port for peer")

	return cmd
}

func bindRunE(cmd *cobra.Command, args []string) error{
	addr, err := cmd.Flags().GetString("addr")
	if err != nil{
		return err
	}

	ctx := cmd.Context()

	// make a connection with that peer
	if err := connect(ctx, addr); err != nil{
		return err
	}
	logger.Info("connection established")
	return nil
}

func connect(ctx context.Context, addr string) error{
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,        // trust every certificate
		ServerName: "localhost",        // must match SAN DNS entry in server cert
		NextProtos: []string{"quic-example-v1"},
	}

	quicConf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	session, err := quic.DialAddr(ctx, addr, tlsConf, quicConf)
	if err != nil {
		return fmt.Errorf("error while establishing p2p connection: %v", err)
	}
	defer func() {
		_ = session.CloseWithError(0, "client closing")
	}()

	// Open a bidirectional stream
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("OpenStreamSync failed: %v", err)
	}
	defer stream.Close()

	// Send a simple message (application-level)
	msg := "hello from client"
	_, err = stream.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("stream write failed: %v", err)
	}

	// Read reply (simple echo example)
	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			logger.Info(fmt.Sprintf("server reply: %s", string(buf[:n])))
		}
		if err != nil {
			if err == io.EOF { break }          // normal close by server
			return fmt.Errorf("stream read failed: %v", err)
		}
	}
	return nil
}