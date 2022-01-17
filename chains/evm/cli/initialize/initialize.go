package initialize

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	evmgaspricer "github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/prepare"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"
)

func InitializeClient(
	url string,
	senderKeyPair *secp256k1.Keypair,
) (*evmclient.EVMClient, error) {
	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client initialization error: %v", err))
		return nil, err
	}
	return ethClient, nil
}

// Initialize transactor which is used for contract calls
// if --prepare flag value is set as true (from CLI) call data is outputted to stdout
// which can be used for multisig contract calls
func InitializeTransactor(
	gasPrice *big.Int,
	txFabric calls.TxFabric,
	client *evmclient.EVMClient,
	prepareFlag bool,
) (transactor.Transactor, error) {
	var trans transactor.Transactor
	if prepareFlag {
		trans = prepare.NewPrepareTransactor()
	} else {
		gasPricer := evmgaspricer.NewLondonGasPriceClient(
			client,
			&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice},
		)
		trans = signAndSend.NewSignAndSendTransactor(txFabric, gasPricer, client)
	}

	return trans, nil
}
