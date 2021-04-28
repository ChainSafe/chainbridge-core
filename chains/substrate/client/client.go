package client

import (
	"fmt"
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

const BridgePalletName = "ChainBridge"
const BridgeStoragePrefix = "ChainBridge"

type VoteState struct {
	VotesFor     []types.AccountID
	VotesAgainst []types.AccountID
	Status       struct {
		IsActive   bool
		IsApproved bool
		IsRejected bool
	}
}

func NewSubstrateClient(url string, key *signature.KeyringPair, stop <-chan struct{}) (*SubstrateClient, error) {
	api, err := gsrpc.NewSubstrateAPI(url)
	if err != nil {
		return nil, err
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}
	return &SubstrateClient{url: url, key: key, genesisHash: genesisHash, api: api, meta: meta, stop: stop}, nil
}

type SubstrateClient struct {
	api         *gsrpc.SubstrateAPI
	url         string                 // API endpoint
	meta        *types.Metadata        // Latest chain metadata
	metaLock    sync.RWMutex           // Lock metadata for updates, allows concurrent reads
	genesisHash types.Hash             // Chain genesis hash
	key         *signature.KeyringPair // Keyring used for signing
	nonce       types.U32              // Latest account nonce
	nonceLock   sync.Mutex             // Locks nonce for updates
	stop        <-chan struct{}        // Signals system shutdown, should be observed in all selects and loops

}

func (c *SubstrateClient) GetMetadata() (meta types.Metadata) {
	c.metaLock.RLock()
	meta = *c.meta
	c.metaLock.RUnlock()
	return meta
}

func (c *SubstrateClient) UpdateMetatdata() error {
	c.metaLock.Lock()
	meta, err := c.api.RPC.State.GetMetadataLatest()
	if err != nil {
		c.metaLock.Unlock()
		return err
	}
	c.meta = meta
	c.metaLock.Unlock()
	return nil
}

func (c *SubstrateClient) GetBlockEvents(hash types.Hash, target interface{}) error {
	meta := c.GetMetadata()
	key, err := types.CreateStorageKey(&meta, "System", "Events", nil, nil)
	if err != nil {
		return err
	}
	var records types.EventRecordsRaw
	_, err = c.api.RPC.State.GetStorage(key, &records, hash)
	if err != nil {
		return err
	}
	err = records.DecodeEventRecords(&meta, target)
	if err != nil {
		return err
	}
	return nil
}

// SubmitTx constructs and submits an extrinsic to call the method with the given arguments.
// All args are passed directly into GSRPC. GSRPC types are recommended to avoid serialization inconsistencies.
func (c *SubstrateClient) SubmitTx(method string, args ...interface{}) error {
	meta := c.GetMetadata()

	// Create call and extrinsic
	call, err := types.NewCall(
		&meta,
		string(method),
		args...,
	)
	if err != nil {
		return fmt.Errorf("failed to construct call: %w", err)
	}
	ext := types.NewExtrinsic(call)
	// Get latest runtime version
	rv, err := c.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return err
	}

	c.nonceLock.Lock()
	latestNonce, err := c.getLatestNonce()
	if err != nil {
		c.nonceLock.Unlock()
		return err
	}
	if latestNonce > c.nonce {
		c.nonce = latestNonce
	}

	// Sign the extrinsic
	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(c.nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(*c.key, o)
	if err != nil {
		c.nonceLock.Unlock()
		return err
	}

	// Submit and watch the extrinsic
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	c.nonce++
	c.nonceLock.Unlock()
	if err != nil {
		return fmt.Errorf("submission of extrinsic failed: %w", err)
	}
	defer sub.Unsubscribe()

	return c.watchSubmission(sub)
}

func (c *SubstrateClient) watchSubmission(sub *author.ExtrinsicStatusSubscription) error {
	for {
		select {
		case <-c.stop:
			return fmt.Errorf("terminated")
		case status := <-sub.Chan():
			switch {
			case status.IsInBlock:
				return nil
			case status.IsRetracted:
				return fmt.Errorf("extrinsic retracted: %s", status.AsRetracted.Hex())
			case status.IsDropped:
				return fmt.Errorf("extrinsic dropped from network")
			case status.IsInvalid:
				return fmt.Errorf("extrinsic invalid")
			}
		case err := <-sub.Err():
			return err
		}
	}
}

func (c *SubstrateClient) getLatestNonce() (types.U32, error) {
	acct := &types.AccountInfo{}
	exists, err := c.QueryStorage("System", "Account", c.key.PublicKey, nil, acct)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, nil
	}

	return acct.Nonce, nil
}

func (c *SubstrateClient) ResolveResourceId(id [32]byte) (string, error) {
	var res []byte
	exists, err := c.QueryStorage(BridgeStoragePrefix, "Resources", id[:], nil, &res)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("resource %x not found on chain", id)
	}
	return string(res), nil
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (c *SubstrateClient) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	// Fetch account nonce
	data := c.GetMetadata()
	key, err := types.CreateStorageKey(&data, prefix, method, arg1, arg2)
	if err != nil {
		return false, err
	}
	return c.api.RPC.State.GetStorageLatest(key, result)
}

func (c *SubstrateClient) VoterAccountID() types.AccountID {
	return types.NewAccountID(c.key.PublicKey)
}

func (c *SubstrateClient) GetHeaderLatest() (*types.Header, error) {
	return c.api.RPC.Chain.GetHeaderLatest()
}

func (c *SubstrateClient) GetBlockHash(blockNumber uint64) (types.Hash, error) {
	return c.api.RPC.Chain.GetBlockHash(blockNumber)
}

func (c *SubstrateClient) GetProposalStatus(sourceID []byte, proposalBytes []byte) (bool, *VoteState, error) {
	voteRes := &VoteState{}
	exists, err := c.QueryStorage(BridgeStoragePrefix, "Votes", sourceID, proposalBytes, voteRes)
	return exists, voteRes, err
}
