package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("not RESP")
	}

	n, _ := strconv.Atoi(line[1:])
	args := make([]string, 0, n)

	for i := 0; i < n; i++ {
		// read $len
		if _, err := reader.ReadString('\n'); err != nil {
			return nil, err
		}

		// read data
		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		args = append(args, strings.TrimSpace(data))
	}

	return args, nil
}

func writeSimpleString(w io.Writer, s string) {
	fmt.Fprintf(w, "+%s\r\n", s)
}

func writeError(w io.Writer, s string) {
	fmt.Fprintf(w, "-%s\r\n", s)
}

func writeInteger(w io.Writer, n int64) {
	fmt.Fprintf(w, ":%d\r\n", n)
}

func writeBulkString(w io.Writer, s string) {
	fmt.Fprintf(w, "$%d\r\n%s\r\n", len(s), s)
}

func writeNull(w io.Writer) {
	fmt.Fprint(w, "$-1\r\n")
}
