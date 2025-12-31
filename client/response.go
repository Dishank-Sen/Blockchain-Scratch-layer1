package client

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)

func readResponse(conn net.Conn) (*types.Response, error) {
	reader := bufio.NewReader(conn)

	rawHeaders, err := readUntilDelimiter(reader, []byte("\r\n\r\n"))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(rawHeaders), "\r\n")
	status := strings.SplitN(lines[0], " ", 2)

	resp := &types.Response{
		StatusCode: atoi(status[0]),
		Message:    status[1],
		Headers:    make(map[string]string),
	}

	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			resp.Headers[strings.TrimSpace(kv[0])] =
				strings.TrimSpace(kv[1])
		}
	}

	if cl, ok := resp.Headers["Content-Length"]; ok {
		n, err := strconv.Atoi(cl)
		if err != nil {
			return nil, err
		}

		resp.Body = make([]byte, n)
		if _, err := io.ReadFull(reader, resp.Body); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
