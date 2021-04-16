//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	celo "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/consensus/istanbul"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/rs/zerolog/log"
)

const (
	timeToWaitUntilNextBlockAppear = 5
)

type HeaderByNumberGetter interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

func SyncBlockValidators(stopChn <-chan struct{}, errChn chan error, c HeaderByNumberGetter, db *ValidatorsStore, chainID uint8, epochSize uint64) {
	var prevValidators []*istanbul.ValidatorData
	// If DB is empty will return 0 (first epoch by itself)
	block, err := db.GetLatestKnownEpochLastBlock(chainID)
	if err != nil {
		errChn <- fmt.Errorf("error on get latest known block from db: %w", err)
		return
	}
	if block.Cmp(big.NewInt(0)) == 0 {
		// If block is zero, initial validators should be empty array
		log.Info().Msg("Syncing validators from zero block")
		prevValidators = make([]*istanbul.ValidatorData, 0)
	} else {
		prevValidators, err = db.GetValidatorsForBlock(block, chainID)
		if err != nil {
			errChn <- fmt.Errorf("error on get latest known validators from db: %w", err)
			return
		}
		// We already know validators for that block so moving to next one
		block.Add(block, big.NewInt(0).SetUint64(epochSize))
		log.Info().Msg(fmt.Sprintf("Syncing validators from %s block", block.String()))
	}
	for {
		select {
		case <-stopChn:
			return
		default:
			header, err := c.HeaderByNumber(context.Background(), block)
			if err != nil {
				if errors.Is(err, celo.NotFound) {
					// Block not yet mined, waiting
					time.Sleep(timeToWaitUntilNextBlockAppear * time.Second)
					continue
				}
				errChn <- fmt.Errorf("gettings header by number err: %w", err)
				return
			}
			extra, err := types.ExtractIstanbulExtra(header)
			if err != nil {
				errChn <- fmt.Errorf("error on extracting istanbul extra: %w", err)
				return
			}
			b := bytes.NewBuffer(extra.RemovedValidators.Bytes())

			if len(extra.AddedValidators) != 0 || b.Len() > 0 {
				log.Debug().Str("block", block.String()).Msg("New validators data")
				prevValidators, err = applyValidatorsDiff(extra, prevValidators)
				if err != nil {
					errChn <- fmt.Errorf("error applying validators diff: %w", err)
					return
				}
			}
			err = db.SetValidatorsForBlock(block, prevValidators, chainID)
			if err != nil {
				errChn <- fmt.Errorf("error on set validators to db: %w", err)
				return
			}
			// Current validators for next epoch, will be set for next last epoch block and applied with its diff
			block.Add(block, big.NewInt(0).SetUint64(epochSize))
		}
	}
}
