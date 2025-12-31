package client

import (
	"fmt"
	"net"
)

func writeRequest(conn net.Conn, method, endpoint string, body []byte) error {
	// Request line
	if _, err := fmt.Fprintf(conn, "%s %s\r\n", method, endpoint); err != nil {
		return err
	}

	// Body framing
	if len(body) > 0 {
		if _, err := fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body)); err != nil {
			return err
		}
	}

	// Header/body delimiter
	if _, err := fmt.Fprint(conn, "\r\n"); err != nil {
		return err
	}

	// Body
	if len(body) > 0 {
		if _, err := conn.Write(body); err != nil {
			return err
		}
	}

	return nil
}
