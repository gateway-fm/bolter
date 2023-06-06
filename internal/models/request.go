package models

type JsonBody struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  []any       `json:"params,omitempty"`
	Id      interface{} `json:"id,omitempty"`
}
