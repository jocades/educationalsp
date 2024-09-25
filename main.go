package main

import (
	"bufio"
	"educationalsp/analysis"
	"educationalsp/lsp"
	"educationalsp/rpc"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	logger.Println("Started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout
	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
		}
		handleMessage(writer, state, method, contents)
	}
}

func handleMessage(w io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("initialize: %s", err)
			return

		}

		logger.Printf("Connected to: %s, %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(w, msg)

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
			return
		}

		logger.Printf("Oppened: %s", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}

		logger.Printf("Change: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.OpenDocument(request.Params.TextDocument.URI, change.Text)
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return
		}
		prettyPrint(request)

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		prettyPrint(response)

		writeResponse(w, response)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return
		}
		prettyPrint(request)

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		prettyPrint(response)

		writeResponse(w, response)
	}
}

func writeResponse(w io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	w.Write([]byte(reply))
}

func prettyPrint(v any) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	logger.Printf("%s", pretty)
}

var logger *log.Logger

func init() {
	filename, _ := os.Executable()
	path := path.Join(path.Dir(filename), "server.log")
	logfile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(logfile, "[educationalsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}

/* func newLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[educationalsp] ", log.Ldate|log.Ltime|log.Lshortfile)
} */
