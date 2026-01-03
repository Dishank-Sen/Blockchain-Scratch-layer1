package peer

import "crypto/tls"

func getTlsConfig() *tls.Config{
	return &tls.Config{
		InsecureSkipVerify: true,        // trust every certificate
		ServerName: "bloc",
		NextProtos: []string{"quic-example-v1"},
	}
}