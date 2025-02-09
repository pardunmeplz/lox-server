package lsp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func StartServer() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(split)
	writer := os.Stdout
	logger := getLogger("/log.txt")
	logger.Print("start log" + time.Now().GoString())

	for scanner.Scan() {
		request := scanner.Text()
		logger.Print(request)

		response, err := handleRequest(request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling request: %v\n", err)
			continue
		}
		if response == nil {
			continue
		}

		logger.Print(string(response))
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

	bodyStart := len(header) + 4
	totalLength := len(header) + 4 + contentLength

	return totalLength, data[bodyStart:totalLength], nil
}

func getLogger(fileName string) *log.Logger {
	logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("invalid log file loc")
	}
	return log.New(logfile, "Pdun>> ", log.Ldate)
}

func writeMessage(writer io.Writer, response []byte) error {
	_, err := writer.Write(response)
	return err
}
