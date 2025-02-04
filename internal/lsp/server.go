package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func StartServer() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(split)
	writer := os.Stdout

	for {
		message, err := readMessage(scanner)
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

func split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return 0, nil, nil
	}
	contentLength, err := strconv.Atoi(string(header[len("Content-Length: "):]))
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength

	return totalLength, data[:totalLength], nil

}

func getLogger(fileName string) *log.Logger {
	logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("invalid log file loc")
	}
	return log.New(logfile, "Pdun>> ", log.Ldate)
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

// func readMessage(scanner bufio.Scanner) ([]byte, error) {
// 	var contentLength int

//     for scanner.Scan()

// 	for {
// 		line, err := reader.ReadHeader()

// 		if err != nil {
// 			return nil, err
// 		}

// 		if line == "" {
// 			break
// 		}

// 		// just ignoring non-content-length headers for now
// 		if strings.HasPrefix(line, "Content-Length: ") {
// 			lengthStr := strings.TrimPrefix(line, "Content-Length: ")
// 			contentLength, err = strconv.Atoi(lengthStr)
// 			if err != nil {
// 				return nil, fmt.Errorf("invalid Content-Length: %v", err)
// 			}
// 		}
// 	}

// 	if contentLength == 0 {
// 		return nil, fmt.Errorf("missing Content-Length header")
// 	}

// 	body, err := reader.ReadBody(contentLength)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return body, nil

// }
