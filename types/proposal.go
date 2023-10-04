package types

func NewProposal(source, destination uint8, depositNonce uint64, resourceId ResourceID, data []byte, metadata Metadata) *Proposal {
	return &Proposal{
		OriginDomainID: source,
		DepositNonce:   depositNonce,
		ResourceID:     resourceId,
		Destination:    destination,
		Data:           data,
		Metadata:       metadata,
	}
}

type Proposal struct {
	OriginDomainID uint8
	DepositNonce   uint64
	ResourceID     ResourceID
	Data           []byte
	Destination    uint8
	Metadata       Metadata
}
