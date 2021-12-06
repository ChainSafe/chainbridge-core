package secp256sender

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
)

type SecpInMemory256Sender struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

func (s *SecpInMemory256Sender) PrivateKey() *ecdsa.PrivateKey {
	return s.privateKey
}
func (s *SecpInMemory256Sender) Address() string {
	return s.address.Hex()
}
