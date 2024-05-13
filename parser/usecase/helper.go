package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// getLatestBlockNumber returns the latest block number
func getLatestBlockNumber() (*string, error) {
	rpcReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}

	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://cloudflare-eth.com", "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp map[string]interface{}
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, err
	}

	resultStr, ok := rpcResp["result"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return &resultStr, nil
}

// parseHexInt parses a hex-encoded string into an int
func parseHexInt(hexStr string) (int, error) {
	if len(hexStr) < 2 || hexStr[:2] != "0x" {
		return 0, fmt.Errorf("invalid hex string")
	}

	num, err := strconv.ParseInt(hexStr[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

// parseHexFloat parses a hex-encoded string into a float
func parseHexFloat(hexStr string) (float64, error) {
	if len(hexStr) < 2 || hexStr[:2] != "0x" {
		return 0, fmt.Errorf("invalid hex string")
	}

	num, err := strconv.ParseUint(hexStr[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return float64(num), nil
}
