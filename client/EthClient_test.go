package client

import (
	"fmt"
	"testing"
)

func TestBlockNumber(t *testing.T) {

	ec := EthClient{"http://localhost:22000"}

	bn := ec.BlockNumber()

	fmt.Println("Latest Block Number ", bn)

	if bn == "" {
		t.Error("Couldnt get latest Block number")
	}

}

func TestCoinbase(t *testing.T) {

	ec := EthClient{"http://localhost:22000"}

	cb := ec.Coinbase()

	fmt.Println("Coinbase ", cb)

	if cb == "" {
		t.Error("Couldnt get Coinbase")
	}

}
