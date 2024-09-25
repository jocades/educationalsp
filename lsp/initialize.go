package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	// ... there's tons more that goes here
	// Capabilities ClientCapabilities
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	/**
	 * Defines how the host (editor) should sync document changes to the language server.
	 * (https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_synchronization)
	 */
	TextDocumentSync   int  `json:"textDocumentSync"`
	HoverProvider      bool `json:"hoverProvider"`
	DefinitionProvider bool `json:"definitionProvider"`
	// CodeActionProvider bool           `json:"codeActionProvider"`
	// CompletionProvider map[string]any `json:"completionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				/**
				 * Documents are synced by always sending the full content of the document.
				 */
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				// CodeActionProvider: true,
				// CompletionProvider: map[string]any{},
			},
			ServerInfo: ServerInfo{
				Name:    "educationalsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}
