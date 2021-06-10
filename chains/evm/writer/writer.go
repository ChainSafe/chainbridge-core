// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"time"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/rs/zerolog/log"
)

var BlockRetryInterval = time.Second * 5

type Proposer interface {
	Status() (relayer.ProposalStatus, error)
	Execute() error
	Vote() error
	VotedBy() bool
}

type MessageHandler interface {
	HandleMessage(m *relayer.Message) (Proposer, error)
}

type EVMVoter struct {
	stop <-chan struct{}
	me   MessageHandler
}

func NewWriter(me MessageHandler) *EVMVoter {
	return &EVMVoter{
		me: me,
	}
}

func (w *EVMVoter) VoteProposal(m *relayer.Message, bridgeAddress string) error {
	prop, err := w.me.HandleMessage(m)
	if err != nil {
		return err
	}
	ps, err := prop.Status()
	if err != nil {
		log.Error().Err(err).Msgf("error getting proposal status %+v", prop)
	}

	votedByCurrentExecutor := prop.VotedBy()

	if votedByCurrentExecutor || ps == relayer.ProposalStatusPassed || ps == relayer.ProposalStatusCanceled || ps == relayer.ProposalStatusExecuted {
		if ps == relayer.ProposalStatusPassed {
			// We should not vote for this proposal but it is ready to be executed
			err = prop.Execute()
			if err != nil {
				return err
			}
			return nil
		} else {
			return nil
		}
	}
	err = prop.Vote()
	if err != nil {
		return err
	}
	// Checking every 5 seconds does proposal is ready to be executed
	// TODO: somehow update infinity loop to break after some period of time
	for {
		select {
		case <-time.After(BlockRetryInterval):
			ps, err := prop.Status()
			if err != nil {
				log.Error().Err(err).Msgf("error getting proposal status %+v", prop)
				return err
			}
			if ps == relayer.ProposalStatusPassed {
				err = prop.Execute()
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case <-w.stop:
			return nil

		}
	}
}
