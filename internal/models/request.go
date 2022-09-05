package models

type JsonBody struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Id      interface{} `json:"id,omitempty"`
}
