package api

import "encoding/json"

type Request struct {
	Params *json.RawMessage `json:"params"`
	Method string           `json:"method"`
	ID     interface{}      `json:"id"`
	Client string
}

type Response struct {
	JsonRpc string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   string           `json:"error"`
	ID      interface{}      `json:"id"`
}
