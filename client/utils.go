package client

import "bufio"

func readUntilDelimiter(r *bufio.Reader, delim []byte) ([]byte, error) {
	var buf []byte
	match := 0

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		buf = append(buf, b)

		if b == delim[match] {
			match++
			if match == len(delim) {
				return buf, nil
			}
		} else if b == delim[0] {
			match = 1
		} else {
			match = 0
		}
	}
}
