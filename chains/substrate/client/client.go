package client

import (
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func NewSubstrateClient(url string, name string, key *signature.KeyringPair) (*SubstrateClient, error) {
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
	return &SubstrateClient{url: url, name: name, key: key, genesisHash: genesisHash, api: api, meta: meta}, nil
}

type SubstrateClient struct {
	api         *gsrpc.SubstrateAPI
	url         string                 // API endpoint
	name        string                 // Chain name
	meta        *types.Metadata        // Latest chain metadata
	metaLock    sync.RWMutex           // Lock metadata for updates, allows concurrent reads
	genesisHash types.Hash             // Chain genesis hash
	key         *signature.KeyringPair // Keyring used for signing
	nonce       types.U32              // Latest account nonce
	nonceLock   sync.Mutex             // Locks nonce for updates
}

func (c *SubstrateClient) getMetadata() (meta types.Metadata) {
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
	meta := c.getMetadata()
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
