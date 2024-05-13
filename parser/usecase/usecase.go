package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"parser/domain"
	"strings"
	"sync"
	"time"
)

// parser is a struct that holds the current block, subscribed addresses, and transactions
type parser struct {
	currentBlock    int
	subscribedAddrs map[string]chan domain.Transaction
	transactions    map[string][]domain.Transaction
	mu              sync.Mutex // Change sync.Mutex to sync.RWMutex
}

// NewParser represents a new instance for parser Interface
func NewParser() domain.Parser {
	return &parser{}
}

// GetCurrentBlock returns the current block number
func (p *parser) GetCurrentBlock() (currentBlockId int, err error) {
	currentBlockHex, err := getLatestBlockNumber()
	if err != nil {
		return
	}
	if currentBlockHex == nil {
		return 0, errors.New("invalid response format")
	}
	currentBlockId, err = parseHexInt(*currentBlockHex)
	if err != nil {
		return
	}
	return
}

// Subscribe adds an address to the list of subscribed addresses
func (p *parser) Subscribe(address string) (res bool, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.subscribedAddrs == nil {
		p.subscribedAddrs = make(map[string]chan domain.Transaction)
	}

	if _, exists := p.subscribedAddrs[address]; exists {
		return false, nil // Address already subscribed
	}

	ch := make(chan domain.Transaction, 100)
	p.subscribedAddrs[address] = ch

	// Start a goroutine to listen for transactions for this address
	go p.listenForTransactions(address)

	return true, nil
}

// GetTransactions returns a list of transactions for a given address
func (p *parser) GetTransactions(address string) (transactions domain.Transactions, err error) {

	if _, exists := p.subscribedAddrs[address]; !exists {
		return transactions, fmt.Errorf("address %s not subscribed", address)
	}
	if p.transactions == nil {
		p.transactions = make(map[string][]domain.Transaction)
	}
	transactions.Transactions = p.transactions[address]
	return transactions, nil
}

func (p *parser) listenForTransactions(address string) {
	for {
		blockNumHex, err := getLatestBlockNumber()
		if err != nil {
			log.Printf("Failed to get latest block number: %v", err)
			continue
		}
		if blockNumHex == nil {
			log.Printf("Invalid response format")
			continue
		}
		err = p.parseTransactions(*blockNumHex, address)
		if err != nil {
			log.Printf("Failed to parse transactions for block %v: %v", *blockNumHex, err)
		}

		// Wait for a new block before checking again
		time.Sleep(5 * time.Second)
	}
}

func (p *parser) parseTransactions(blockIdHex string, address string) error {
	rpcReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockIdHex, true},
		"id":      1,
	}

	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return err
	}
	resp, err := http.Post("https://cloudflare-eth.com", "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rpcResp map[string]interface{}
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return err
	}

	result, ok := rpcResp["result"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format")
	}

	txsData, ok := result["transactions"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid transactions data")
	}
	blockNumber, err := parseHexInt(blockIdHex)
	if err != nil {
		return err
	}
	for _, txData := range txsData {
		txDataMap, ok := txData.(map[string]interface{})
		if !ok {
			continue
		}
		from, fromOk := txDataMap["from"].(string)
		to, toOk := txDataMap["to"].(string)
		value, valueOk := txDataMap["value"].(string)

		if !fromOk || !toOk || !valueOk {
			continue
		}
		tx := domain.Transaction{
			From:        from,
			To:          to,
			Value:       weiToEther(value),
			BlockNumber: blockNumber,
		}
		p.mu.Lock()
		if strings.ToLower(from) == strings.ToLower(address) || strings.ToLower(to) == strings.ToLower(address) {
			if p.transactions == nil {
				p.transactions = make(map[string][]domain.Transaction)
			}
			p.transactions[address] = append(p.transactions[address], tx)
			p.subscribedAddrs[address] <- tx
		}
		p.mu.Unlock()
	}

	return nil
}
func weiToEther(wei string) string {
	value, _ := new(big.Int).SetString(wei[2:], 16)
	ether := new(big.Float).SetInt(value)
	ether.Quo(ether, big.NewFloat(1e18))
	return fmt.Sprintf("%.18f", ether)
}
