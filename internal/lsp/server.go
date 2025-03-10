package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var serverState struct {
	shutdown         bool
	initialized      bool
	writer           *os.File
	logger           *log.Logger
	loggerMu         sync.Mutex
	idCount          int
	serverRequestIds map[int]bool
	documents        map[string]*DocumentService
}

func initializeServerState() {
	serverState.initialized = false
	serverState.shutdown = false
	serverState.writer = os.Stdout
	serverState.logger = getLogger("log.txt")
	serverState.idCount = 1
	serverState.serverRequestIds = make(map[int]bool)
	serverState.documents = make(map[string]*DocumentService)
}

func StartServer() {
	initializeServerState()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(split)
	serverState.logger.Print("start log" + time.Now().GoString())

	for scanner.Scan() {
		request := scanner.Text()
		serverState.logger.Print(" request >>" + request)

		response, err := handleRequest(request)
		if err != nil {
			serverState.logger.Print(fmt.Sprintf("Error handling request: %v\n", err))
			continue
		}
		if response == nil {
			continue
		}

		serverState.logger.Print(" response << " + string(response))
		if err := writeMessage(response); err != nil {
			serverState.logger.Print(fmt.Sprintf("Error writing response: %v\n", err))
			break
		}
	}
}

func sendNotification(response []byte) {

	if response == nil {
		return
	}

	if err := writeMessage(response); err != nil {
		serverState.logger.Print(fmt.Sprintf("Error writing response: %v\n", err))
	}
	serverState.logger.Print(string(response))

}

func sendRequest(method string) {
	id := serverState.idCount
	serverState.idCount++
	serverState.serverRequestIds[id] = true
	switch method {
	case "client/registerCapability":
		requestObj := register(id)

		request, err := json.Marshal(requestObj)
		if err != nil {
			serverState.logger.Print(fmt.Sprintf("invalid Request: %v\n", err))
			return
		}
		if err := writeMessage(request); err != nil {
			serverState.logger.Print(fmt.Sprintf("Error writing request: %v\n", err))
		}
		serverState.logger.Print(string(request))
	}

}
func getLogger(fileName string) *log.Logger {
	logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("invalid log file loc")
	}
	return log.New(logfile, "\nPdun>> ", log.Ldate)
}

func writeMessage(response []byte) error {
	header := []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response)))
	serverState.loggerMu.Lock()
	defer serverState.loggerMu.Unlock()

	_, err := serverState.writer.Write(append(header, response...))
	return err
}
