package domain

// Transaction represents an Ethereum transaction
type Transaction struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber int    `json:"blockNumber"`
}
type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

// Parser is an interface for parsing Ethereum transactions
type Parser interface {
	// last parsed block
	GetCurrentBlock() (int, error)
	// add address to observer
	Subscribe(address string) (bool, error)
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) (Transactions, error)
}
