// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package secp256k1

import (
	"crypto/sha256"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewKeypairFromSeed(t *testing.T) {
	kp, err := GenerateKeypair()
	if err != nil {
		t.Fatal(err)
	}

	if kp.PublicKey() == "" || kp.Address() == "" {
		t.Fatalf("key is missing data: %#v", kp)
	}
}

func TestEncodeAndDecodeKeypair(t *testing.T) {
	kp, err := GenerateKeypair()
	if err != nil {
		t.Fatal(err)
	}

	enc := kp.Encode()
	res := new(Keypair)
	err = res.Decode(enc)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, kp) {
		t.Fatalf("Fail: got %#v expected %#v", res, kp)
	}
}

func TestSign(t *testing.T) {
	kp, err := GenerateKeypair()
	if err != nil {
		t.Fatal(err)
	}

	digestHash := sha256.Sum256([]byte{0, 0})

	sig, err := kp.Sign(digestHash[:])
	if err != nil {
		t.Fatal(err)
	}

	if len(sig) != crypto.SignatureLength {
		t.Fatalf("Fail: wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength)
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := sig[64]

	valid := crypto.ValidateSignatureValues(v, r, s, true)

	if !valid {
		t.Fatal("Fail: got invalid signature")
	}
}
