package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func StartServer() {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		message, err := readMessage(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
			break
		}

		response, err := handleRequest(message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling request: %v\n", err)
			continue
		}

		if err := writeMessage(writer, response); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
			break
		}

	}
}

func writeMessage(writer io.Writer, response any) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	_, err = writer.Write([]byte(header))
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err

}

func readMessage(reader *bufio.Reader) ([]byte, error) {
	var contentLength int

	//read header
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %v", err)
			}
		}
	}
	if contentLength == 0 {
		fmt.Println(3)
		return nil, fmt.Errorf("missing Content-Length header")
	}

	//read body
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		fmt.Println(4)
		return nil, err
	}

	return body, nil
}
