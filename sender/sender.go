package sender

import (
	"crypto/ecdsa"
)

type Sender interface {
	PrivateKey() *ecdsa.PrivateKey
	Address() string
}
