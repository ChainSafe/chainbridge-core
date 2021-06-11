package config

import (
	"errors"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/chains"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	evmListener "github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	evmWriter "github.com/ChainSafe/chainbridge-core/chains/evm/writer"
	"github.com/ChainSafe/chainbridge-core/chains/substrate"
	"github.com/ChainSafe/chainbridge-core/crypto/sr25519"
	"github.com/ChainSafe/chainbridge-core/keystore"

	subListener "github.com/ChainSafe/chainbridge-core/chains/substrate/listener"
	subWriter "github.com/ChainSafe/chainbridge-core/chains/substrate/writer"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ChainSafe/chainbridge-core/sender"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
)

type EVMClient interface {
	evmListener.ChainReader
	evmWriter.VoterExecutor
	InitializeClient(config *evm.EVMConfig, sender sender.Sender) error
}

type SubstrateClient interface {
	subListener.SubstrateReader
	subWriter.Voter
	InitializeClient(url string, key *signature.KeyringPair, stop <-chan struct{}) error
}

func InitializeRelayer(
	cfg *chains.Config,
	evmClient EVMClient,
	subClient SubstrateClient,
	sender sender.Sender,
	kvdb blockstore.KeyValueReaderWriter,
	handler evmListener.HandlerFabric,
	stopChn <-chan struct{}) (*relayer.Relayer, error) {

	relayedChains := make([]relayer.RelayedChain, len(cfg.Chains))
	for index, chainConfig := range cfg.Chains {

		if chainConfig.Type == "ethereum" {
			relayedChain, err := InitializeEVMChain(&chainConfig, evmClient, sender, kvdb, handler)
			if err != nil {
				panic(err) // TODO: make these more descriptive and don't panic return
			}
			relayedChains[index] = relayedChain
		} else if chainConfig.Type == "substrate" {
			subChain, err := InitializeSubChain(&chainConfig, subClient, kvdb, stopChn)
			if err != nil {
				panic(err)
			}
			relayedChains[index] = subChain
		} else {
			return nil, errors.New("unrecognized Chain Type")
		}

	}

	r := relayer.NewRelayer(relayedChains)
	return r, nil
}

func InitializeEVMChain(
	config *chains.RawChainConfig,
	client EVMClient,
	sender sender.Sender,
	kvdb blockstore.KeyValueReaderWriter,
	handler evmListener.HandlerFabric) (*evm.EVMChain, error) {
	cfg, err := evm.ParseConfig(config)
	if err != nil {
		return nil, err
	}
	err = client.InitializeClient(cfg, sender)
	if err != nil {
		return nil, err
	}

	listener := evmListener.NewEVMListener(client)
	listener.RegisterHandlerFabric(cfg.Erc20Handler, handler)

	writer := evmWriter.NewWriter(client)
	writer.RegisterProposalHandler(cfg.Erc20Handler, evmWriter.ERC20ProposalHandler)
	evmChain := evm.NewEVMChain(listener, writer, kvdb, cfg.Bridge, cfg.GeneralChainConfig.Id)
	return evmChain, nil
}

func InitializeSubChain(
	config *chains.RawChainConfig,
	client SubstrateClient,
	kvdb blockstore.KeyValueReaderWriter,
	stopChn <-chan struct{}) (*substrate.SubstrateChain, error) {
	cfg := substrate.ParseConfig(config)

	kp, err := keystore.KeypairFromAddress(cfg.GeneralChainConfig.From, keystore.SubChain, "alice", true)
	if err != nil {
		panic(err)
	}
	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	err = client.InitializeClient(cfg.GeneralChainConfig.Endpoint, krp, stopChn)
	if err != nil {
		panic(err)
	}
	subL := subListener.NewSubstrateListener(client)
	subW := subWriter.NewSubstrateWriter(1, client)

	// TODO: really not need this dynamic handler assignment
	subL.RegisterSubscription(relayer.FungibleTransfer, subListener.FungibleTransferHandler)
	subL.RegisterSubscription(relayer.GenericTransfer, subListener.GenericTransferHandler)
	subL.RegisterSubscription(relayer.NonFungibleTransfer, subListener.NonFungibleTransferHandler)

	subW.RegisterHandler(relayer.FungibleTransfer, subWriter.CreateFungibleProposal)
	subChain := substrate.NewSubstrateChain(subL, subW, kvdb, cfg.GeneralChainConfig.Id)
	return subChain, nil
}
