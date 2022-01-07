package optimismclient

import (
	"context"
	"crypto/ecdsa"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

type EthContext struct {
	BlockNumber uint64 `json:"blockNumber"`
	Timestamp   uint64 `json:"timestamp"`
}

// RollupContext represents the height of the rollup.
// Index is the last processed CanonicalTransactionChain index
// QueueIndex is the last processed `enqueue` index
// VerifiedIndex is the last processed CTC index that was batched
type RollupContext struct {
	Index         uint64 `json:"index"`
	QueueIndex    uint64 `json:"queueIndex"`
	VerifiedIndex uint64 `json:"verifiedIndex"`
}

type rollupInfo struct {
	Mode          string        `json:"mode"`
	Syncing       bool          `json:"syncing"`
	EthContext    EthContext    `json:"ethContext"`
	RollupContext RollupContext `json:"rollupContext"`
}

type OptimismClient struct {
	// NOTE: If we wanted or needed to have the same private variables within the EVMClient struct inside the OptimismClient
	// we would essentially need to replicate the entire EVMClient. Currently it seems that this can be avoided.
	// *ethclient.Client
	// kp               *secp256k1.Keypair
	// gethClient       *gethclient.Client
	// rpClient         *rpc.Client
	// nonce            *big.Int
	// nonceLock        sync.Mutex
	*evmclient.EVMClient
	verifyRollup     bool
	verifierRpClient *rpc.Client
}

// NewEVMClientFromParams creates a client for EVMChain with provided
// private key.
func NewOptimismClientFromParams(url string, privateKey *ecdsa.PrivateKey, verifyRollup bool, verifierEndpoint string) (*OptimismClient, error) {
	c := &OptimismClient{}

	sequencerClient, err := evmclient.NewEVMClientFromParams(url, privateKey)
	if err != nil {
		return nil, err
	}
	c.EVMClient = sequencerClient

	c.verifyRollup = verifyRollup
	if c.verifyRollup {
		err := c.configureVerifier(verifierEndpoint)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewOptimismClient creates a client for the Optimism chain configured with specified config.
func NewOptimismClient(cfg *chain.OptimismConfig) (*OptimismClient, error) {
	c := &OptimismClient{}

	sequencerClient, err := evmclient.NewEVMClient(&cfg.EVMConfig)
	if err != nil {
		return nil, err
	}
	c.EVMClient = sequencerClient

	c.verifyRollup = cfg.VerifyRollup
	if c.verifyRollup {
		err := c.configureVerifier(cfg.VerifierEndpoint)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *OptimismClient) configureVerifier(url string) error {
	// The VerifierEndpoint in the config is for the verifier replica and is read-only.
	// This RPC client is only used for getting info from the verifier as to whether the rollup is valid
	verifierRpClient, err := rpc.DialContext(context.TODO(), url)
	if err != nil {
		log.Debug().Msgf("could not connect to verifier endpoint: %v", url)
		log.Debug().Msgf("dial context err: %v", err)
		return err
	}
	c.verifierRpClient = verifierRpClient
	return nil
}

func (c *OptimismClient) RollupInfo() (*rollupInfo, error) {
	var info *rollupInfo

	err := c.verifierRpClient.CallContext(context.TODO(), &info, "rollup_getInfo")
	if err == nil && info == nil {
		err = ethereum.NotFound
	}
	return info, err
}

func (c *OptimismClient) IsRollupVerified(blockNumber uint64) (bool, error) {
	log.Debug().Msg("Just got inside method IsRollupVerified")

	if !c.verifyRollup {
		return true, nil
	}

	info, err := c.RollupInfo()
	if err != nil {
		return false, err
	}

	log.Debug().Msgf("Block number to check against index: %v", blockNumber)
	log.Debug().Msgf("Rollup info: %v", info)
	if blockNumber <= info.RollupContext.VerifiedIndex {
		return true, nil
	} else {
		return false, nil
	}
}
