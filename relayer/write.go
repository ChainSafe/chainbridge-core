package relayer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type ChainWriter interface {
	Write(m XCMessager)
}

func Write(m XCMessager) {
	log.Info().Str("type", string(m.GetType())).Interface("src", m.GetSource()).Interface("dst", m.GetDestination()).Interface("nonce", m.GetDepositNonce()).Str("rId", fmt.Sprintf("%x", m.GetResourceID())).Msg("Attempting to resolve message")
	data, err := m.CreateProposalData()
	if err != nil {
		panic(err)
	}
	//TODO ??
	handlerAddress := m.GetHandlerAddress()
	dataHash := m.CreateProposalDataHash()

	voteProposal()

}
