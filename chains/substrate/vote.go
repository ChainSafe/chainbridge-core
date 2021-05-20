package substrate

import "github.com/centrifuge/go-substrate-rpc-client/types"

type VoteState struct {
	VotesFor     []types.AccountID
	VotesAgainst []types.AccountID
	Status       struct {
		IsActive   bool
		IsApproved bool
		IsRejected bool
	}
}
