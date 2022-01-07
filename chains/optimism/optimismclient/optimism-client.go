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
		c.configureVerifier(verifierEndpoint)
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
		c.configureVerifier(cfg.VerifierEndpoint)
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

// NOTE: The code below is an exact replica of the evm-client and will be necessary if we go about declaring the same private
// variables within the EVMClient struct inside the OptimismClient

// func (c *OptimismClient) SubscribePendingTransactions(ctx context.Context, ch chan<- common.Hash) (*rpc.ClientSubscription, error) {
// 	return c.gethClient.SubscribePendingTransactions(ctx, ch)
// }

// // LatestBlock returns the latest block from the current chain
// func (c *OptimismClient) LatestBlock() (*big.Int, error) {
// 	var head *headerNumber
// 	err := c.rpClient.CallContext(context.Background(), &head, "eth_getBlockByNumber", toBlockNumArg(nil), false)
// 	if err == nil && head == nil {
// 		err = ethereum.NotFound
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return head.Number, nil
// }

// type headerNumber struct {
// 	Number *big.Int `json:"number"           gencodec:"required"`
// }

// func (h *headerNumber) UnmarshalJSON(input []byte) error {
// 	type headerNumber struct {
// 		Number *hexutil.Big `json:"number" gencodec:"required"`
// 	}
// 	var dec headerNumber
// 	if err := json.Unmarshal(input, &dec); err != nil {
// 		return err
// 	}
// 	if dec.Number == nil {
// 		return errors.New("missing required field 'number' for Header")
// 	}
// 	h.Number = (*big.Int)(dec.Number)
// 	return nil
// }

// func (c *OptimismClient) WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error) {
// 	retry := 50
// 	for retry > 0 {
// 		receipt, err := c.Client.TransactionReceipt(context.Background(), h)
// 		if err != nil {
// 			retry--
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}
// 		if receipt.Status != 1 {
// 			return receipt, fmt.Errorf("transaction failed on chain. Receipt status %v", receipt.Status)
// 		}
// 		return receipt, nil
// 	}
// 	return nil, errors.New("tx did not appear")
// }

// func (c *OptimismClient) GetTransactionByHash(h common.Hash) (tx *types.Transaction, isPending bool, err error) {
// 	return c.Client.TransactionByHash(context.Background(), h)
// }

// func (c *OptimismClient) FetchDepositLogs(ctx context.Context, contractAddress common.Address, startBlock *big.Int, endBlock *big.Int) ([]*evmclient.DepositLogs, error) {
// 	logs, err := c.FilterLogs(ctx, buildQuery(contractAddress, string(util.Deposit), startBlock, endBlock))
// 	if err != nil {
// 		return nil, err
// 	}
// 	depositLogs := make([]*evmclient.DepositLogs, 0)

// 	abi, err := abi.JSON(strings.NewReader(consts.BridgeABI))
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, l := range logs {
// 		dl, err := c.UnpackDepositEventLog(abi, l.Data)
// 		if err != nil {
// 			log.Error().Msgf("failed unpacking deposit event log: %v", err)
// 			continue
// 		}
// 		log.Debug().Msgf("Found deposit log in block: %d, TxHash: %s, contractAddress: %s, sender: %s", l.BlockNumber, l.TxHash, l.Address, dl.SenderAddress)

// 		depositLogs = append(depositLogs, dl)
// 	}

// 	return depositLogs, nil
// }

// func (c *OptimismClient) UnpackDepositEventLog(abi abi.ABI, data []byte) (*evmclient.DepositLogs, error) {
// 	var dl evmclient.DepositLogs

// 	err := abi.UnpackIntoInterface(&dl, "Deposit", data)
// 	if err != nil {
// 		return &evmclient.DepositLogs{}, err
// 	}

// 	return &dl, nil
// }

// func (c *OptimismClient) FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error) {
// 	return c.FilterLogs(ctx, buildQuery(contractAddress, event, startBlock, endBlock))
// }

// // SendRawTransaction accepts rlp-encode of signed transaction and sends it via RPC call
// func (c *OptimismClient) SendRawTransaction(ctx context.Context, tx []byte) error {
// 	return c.rpClient.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(tx))
// }

// func (c *OptimismClient) CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error) {
// 	var hex hexutil.Bytes
// 	err := c.rpClient.CallContext(ctx, &hex, "eth_call", callArgs, toBlockNumArg(blockNumber))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return hex, nil
// }

// func (c *OptimismClient) CallContext(ctx context.Context, target interface{}, rpcMethod string, args ...interface{}) error {
// 	err := c.rpClient.CallContext(ctx, target, rpcMethod, args)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *OptimismClient) PendingCallContract(ctx context.Context, callArgs map[string]interface{}) ([]byte, error) {
// 	var hex hexutil.Bytes
// 	err := c.rpClient.CallContext(ctx, &hex, "eth_call", callArgs, "pending")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return hex, nil
// }

// func (c *OptimismClient) From() common.Address {
// 	return c.kp.CommonAddress()
// }

// func (c *OptimismClient) SignAndSendTransaction(ctx context.Context, tx evmclient.CommonTransaction) (common.Hash, error) {
// 	id, err := c.ChainID(ctx)
// 	if err != nil {
// 		//panic(err)
// 		// Probably chain does not support chainID eg. CELO
// 		id = nil
// 	}
// 	rawTx, err := tx.RawWithSignature(c.kp.PrivateKey(), id)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	err = c.SendRawTransaction(ctx, rawTx)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	return tx.Hash(), nil
// }

// func (c *OptimismClient) RelayerAddress() common.Address {
// 	return c.kp.CommonAddress()
// }

// func (c *OptimismClient) LockNonce() {
// 	c.nonceLock.Lock()
// }

// func (c *OptimismClient) UnlockNonce() {
// 	c.nonceLock.Unlock()
// }

// func (c *OptimismClient) UnsafeNonce() (*big.Int, error) {
// 	var err error
// 	for i := 0; i <= 10; i++ {
// 		if c.nonce == nil {
// 			nonce, err := c.PendingNonceAt(context.Background(), c.kp.CommonAddress())
// 			if err != nil {
// 				time.Sleep(1 * time.Second)
// 				continue
// 			}
// 			c.nonce = big.NewInt(0).SetUint64(nonce)
// 			return c.nonce, nil
// 		}
// 		return c.nonce, nil
// 	}
// 	return nil, err
// }

// func (c *OptimismClient) UnsafeIncreaseNonce() error {
// 	nonce, err := c.UnsafeNonce()
// 	if err != nil {
// 		return err
// 	}
// 	c.nonce = nonce.Add(nonce, big.NewInt(1))
// 	return nil
// }

// func (c *OptimismClient) BaseFee() (*big.Int, error) {
// 	head, err := c.HeaderByNumber(context.TODO(), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return head.BaseFee, nil
// }

// func toBlockNumArg(number *big.Int) string {
// 	if number == nil {
// 		return "latest"
// 	}
// 	return hexutil.EncodeBig(number)
// }

// // buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
// func buildQuery(contract common.Address, sig string, startBlock *big.Int, endBlock *big.Int) ethereum.FilterQuery {
// 	query := ethereum.FilterQuery{
// 		FromBlock: startBlock,
// 		ToBlock:   endBlock,
// 		Addresses: []common.Address{contract},
// 		Topics: [][]common.Hash{
// 			{crypto.Keccak256Hash([]byte(sig))},
// 		},
// 	}
// 	return query
// }
