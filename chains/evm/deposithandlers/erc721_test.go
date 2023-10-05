package deposithandlers_test

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/deposit"
	"github.com/ChainSafe/sygma-core/chains/evm/deposithandlers"
	"github.com/ChainSafe/sygma-core/chains/evm/eventhandlers"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/suite"
)

type Erc721HandlerTestSuite struct {
	suite.Suite
}

func TestRunErc721HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Erc721HandlerTestSuite))
}

func (s *Erc721HandlerTestSuite) SetupSuite()    {}
func (s *Erc721HandlerTestSuite) TearDownSuite() {}
func (s *Erc721HandlerTestSuite) SetupTest()     {}
func (s *Erc721HandlerTestSuite) TearDownTest()  {}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerEmptyMetadata() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")

	calldata := deposit.ConstructErc721DepositData(recipient.Bytes(), big.NewInt(2), []byte{})
	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenId := calldata[:32]
	recipientAddressParsed := calldata[64:84]
	var metadata []byte

	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.NonFungibleTransfer,
		Payload: []interface{}{
			tokenId,
			recipientAddressParsed,
			metadata,
		},
	}
	conf := &Config{}

	erc721DepositHandler := deposithandlers.Erc721DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	m, err := erc721DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerArbitraryFunctionError() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")

	calldata := deposit.ConstructErc721DepositData(recipient.Bytes(), big.NewInt(2), []byte{})
	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	conf := &Config{}

	erc721DepositHandler := deposithandlers.Erc721DepositHandler{
		ArbitraryFunction: testErrFunc,
		Config:            conf,
	}
	_, err := erc721DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)

	s.NotNil(err)
}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerIncorrectDataLen() {
	metadata := []byte("0xdeadbeef")

	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 16)...)
	calldata = append(calldata, metadata...)

	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		SenderAddress:       common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	conf := &Config{}

	erc721DepositHandler := deposithandlers.Erc721DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	m, err := erc721DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(m)
	s.EqualError(err, "invalid calldata length: less than 84 bytes")
}

func (s *Erc721HandlerTestSuite) TestErc721DepositHandler() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")
	metadata := []byte("metadata.url")

	calldata := deposit.ConstructErc721DepositData(recipient.Bytes(), big.NewInt(2), metadata)
	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenId := calldata[:32]
	recipientAddressParsed := calldata[64:84]
	parsedMetadata := calldata[116:128]

	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.NonFungibleTransfer,
		Payload: []interface{}{
			tokenId,
			recipientAddressParsed,
			parsedMetadata,
		},
	}
	conf := &Config{}

	erc721DepositHandler := deposithandlers.Erc721DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	m, err := erc721DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}
func (s *Erc721HandlerTestSuite) TestErc721DepositHandlerWithPriority() {
	recipient := common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b")
	metadata := []byte("metadata.url")

	calldata := deposit.ConstructErc721DepositDataWithPriority(recipient.Bytes(), big.NewInt(2), metadata, uint8(1))
	depositLog := &eventhandlers.Deposit{
		DestinationDomainID: 0,
		ResourceID:          [32]byte{0},
		DepositNonce:        1,
		Data:                calldata,
		HandlerResponse:     []byte{},
	}

	sourceID := uint8(1)
	tokenId := calldata[:32]

	// 32 - 64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddressParsed := calldata[64:(64 + recipientAddressLength.Int64())]

	// (64 + recipient address length) - ((64 + recipient address length) + 32) is metadata length
	metadataLength := big.NewInt(0).SetBytes(
		calldata[(64 + recipientAddressLength.Int64()):((64 + recipientAddressLength.Int64()) + 32)],
	)

	metadataStart := (64 + recipientAddressLength.Int64()) + 32
	parsedMetadata := calldata[metadataStart : metadataStart+metadataLength.Int64()]

	expected := &types.Message{
		Source:       sourceID,
		Destination:  depositLog.DestinationDomainID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         types.NonFungibleTransfer,
		Payload: []interface{}{
			tokenId,
			recipientAddressParsed,
			parsedMetadata,
		},
		Metadata: types.Metadata{
			Data: map[string]interface{}{
				"Priority": uint8(1),
			},
		},
	}
	conf := &Config{}

	erc721DepositHandler := deposithandlers.Erc721DepositHandler{
		ArbitraryFunction: testFunc,
		Config:            conf,
	}
	m, err := erc721DepositHandler.HandleDeposit(sourceID, depositLog.DestinationDomainID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Data, depositLog.HandlerResponse)
	s.Nil(err)
	s.NotNil(m)
	s.Equal(expected, m)
}
