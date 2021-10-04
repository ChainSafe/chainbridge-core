package calls

import (
	"testing"
)

func TestPrepareSetDepositNonceInput(t *testing.T) {
	domainId := uint8(0)
	depositNonce := uint64(0)

	bytes, err := PrepareSetDepositNonceInput(domainId, depositNonce)
	if err != nil {
		t.Fatalf("could not prepare set deposit nonce input: %v", err)
	}

	if len(bytes) == 0 {
		t.Fatal("byte slice returned is empty")
	}
}
