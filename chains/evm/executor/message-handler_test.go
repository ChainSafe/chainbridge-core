package executor_test

import (
	"errors"

	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/executor"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

var errIncorrectERC20PayloadLen = errors.New("malformed payload. Len  of payload should be 2")
var errIncorrectERC721PayloadLen = errors.New("malformed payload. Len  of payload should be 3")
var errIncorrectGenericPayloadLen = errors.New("malformed payload. Len  of payload should be 1")

var errIncorrectAmount = errors.New("wrong payload amount format")
var errIncorrectRecipient = errors.New("wrong payload recipient format")
var errIncorrectTokenID = errors.New("wrong payload tokenID format")
var errIncorrectMetadata = errors.New("wrong payload metadata format")

//ERC20
type Erc20HandlerTestSuite struct {
	suite.Suite
}

func TestRunErc20HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Erc20HandlerTestSuite))
}

func (s *Erc20HandlerTestSuite) SetupSuite()    {}
func (s *Erc20HandlerTestSuite) TearDownSuite() {}
func (s *Erc20HandlerTestSuite) SetupTest()     {}
func (s *Erc20HandlerTestSuite) TearDownTest()  {}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessage() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(err)
	s.NotNil(prop)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectDataLen() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC20PayloadLen.Error())
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectAmount() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			"incorrectAmount", // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectAmount.Error())
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectRecipient() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2},            // amount
			"incorrectRecipient", // recipientAddress
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectRecipient.Error())
}

// ERC721
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

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerEmptyMetadata() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // tokenID
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{}, // metadata
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(err)
	s.NotNil(prop)
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectDataLen() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // tokenID
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC721PayloadLen.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectAmount() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			"incorrectAmount", // tokenID
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{}, // metadata
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectTokenID.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectRecipient() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			"incorrectRecipient",
			[]byte{}, // metadata
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectRecipient.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectMetadata() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			"incorrectMetadata", // metadata
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectMetadata.Error())
}

// GENERIC
type GenericHandlerTestSuite struct {
	suite.Suite
}

func TestRunGenericHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GenericHandlerTestSuite))
}

func (s *GenericHandlerTestSuite) SetupSuite()    {}
func (s *GenericHandlerTestSuite) TearDownSuite() {}
func (s *GenericHandlerTestSuite) SetupTest()     {}
func (s *GenericHandlerTestSuite) TearDownTest()  {}
func (s *GenericHandlerTestSuite) TestGenericHandleEvent() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{}, // metadata
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.GenericMessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(err)
	s.NotNil(prop)
}

func (s *GenericHandlerTestSuite) TestGenericHandleEventIncorrectDataLen() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload:      []interface{}{},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.GenericMessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectGenericPayloadLen.Error())
}

func (s *GenericHandlerTestSuite) TestGenericHandleEventIncorrectMetadata() {
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			"incorrectMetadata", // metadata
		},
		Metadata: message.Metadata{
			Priority: uint8(1),
		},
	}

	prop, err := executor.GenericMessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectMetadata.Error())
}
