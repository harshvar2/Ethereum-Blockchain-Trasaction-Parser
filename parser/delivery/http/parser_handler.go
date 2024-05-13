package http

import (
	"parser/domain"

	"github.com/labstack/echo/v4"
)

// ParserHandler : http handler for parser
type ParserHandler struct {
	parser domain.Parser
}

// NewParserHandler will initialize the parser http handler
func NewParserHandler(e *echo.Echo, parser domain.Parser) {
	handler := &ParserHandler{
		parser: parser,
	}
	e.POST("/subscribe/:address", handler.Subscribe)
	e.GET("/transactions/:address", handler.GetTransactions)
	e.GET("/currentBlock", handler.GetCurrentBlock)
}

// GetCurrentBlock : get current block
func (p *ParserHandler) GetCurrentBlock(c echo.Context) error {
	currentBlock, err := p.parser.GetCurrentBlock()
	if err != nil {
		return c.JSON(500, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(200, map[string]int{
		"id": currentBlock,
	})
}

// Subscribe : subscribe to an address
func (p *ParserHandler) Subscribe(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		return c.JSON(400, map[string]string{
			"status":  "error",
			"message": "Missing address parameter"})
	}

	res, err := p.parser.Subscribe(address)
	if err != nil {
		return c.JSON(500, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if !res {
		return c.JSON(400, map[string]string{
			"status":  "error",
			"message": "Address: " + address + " already subscribed",
		})
	}
	return c.JSON(200, map[string]string{
		"status":  "success",
		"message": "Subscribed to address: " + address,
	})
}

// GetTransactions : get transactions for an address
func (p *ParserHandler) GetTransactions(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		return c.JSON(400, map[string]string{
			"status":  "error",
			"message": "Missing address parameter"})
	}

	transactions, err := p.parser.GetTransactions(address)
	if err != nil {
		return c.JSON(500, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(200, transactions)
}
