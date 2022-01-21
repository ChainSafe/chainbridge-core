package voter_test

import (
	"errors"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"
)

var errIncorrectERC20PayloadLen = errors.New("malformed payload. Len  of payload should be 3")
var errIncorrectERC20Amount = errors.New("wrong payloads amount format")
var errIncorrectERC20Recipient = errors.New("wrong payloads recipient format")
var errIncorrectERC20Priority = errors.New("wrong payloads priority format")

var errIncorrectERC721PayloadLen = errors.New("malformed payload. Len  of payload should be 4")
var errIncorrectERC721TokenID = errors.New("wrong payloads tokenID format")
var errIncorrectERC721Recipient = errors.New("wrong payloads recipient format")
var errIncorrectERC721Metadata = errors.New("wrong payloads metadata format")
var errIncorrectERC721Priority = errors.New("wrong payloads priority format")

var errIncorrectGenericPayloadLen = errors.New("malformed payload. Len  of payload should be 1")
var errIncorrectGenericMetadata = errors.New("unable to convert metadata to []byte")

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
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{1}, // priority
		},
	}

	prop, err := voter.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(err)
	s.NotNil(prop)
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectDataLen() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
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
	}

	prop, err := voter.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC20PayloadLen.Error())
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectAmount() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			"incorrectAmount", // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{1}, // priority
		},
	}

	prop, err := voter.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC20Amount.Error())
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectRecipient() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2},            // amount
			"incorrectRecipient", // recipientAddress
			[]byte{1},            // priority
		},
	}

	prop, err := voter.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC20Recipient.Error())
}

func (s *Erc20HandlerTestSuite) TestErc20HandleMessageIncorrectPriority() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			"incorrectPriority", // priority
		},
	}

	prop, err := voter.ERC20MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC20Priority.Error())
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
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // tokenID
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{},  // metadata
			[]byte{1}, // priority
		},
	}

	prop, err := voter.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(err)
	s.NotNil(prop)
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectDataLen() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
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
	}

	prop, err := voter.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC721PayloadLen.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectAmount() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			"incorrectAmount", // tokenID
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{},  // metadata
			[]byte{1}, // priority
		},
	}

	prop, err := voter.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC721TokenID.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectRecipient() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			"incorrectRecipient",
			[]byte{},  // metadata
			[]byte{1}, // priority
		},
	}

	prop, err := voter.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC721Recipient.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectMetadata() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
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
			[]byte{1},           // priority
		},
	}

	prop, err := voter.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC721Metadata.Error())
}

func (s *Erc721HandlerTestSuite) TestErc721MessageHandlerIncorrectPriority() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{2}, // amount
			[]byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}, // recipientAddress
			[]byte{},            // metadata
			"incorrectPriority", // priority
		},
	}

	prop, err := voter.ERC721MessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectERC721Priority.Error())
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
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			[]byte{}, // metadata
		},
	}

	prop, err := voter.GenericMessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(err)
	s.NotNil(prop)
}

func (s *GenericHandlerTestSuite) TestGenericHandleEventIncorrectDataLen() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload:      []interface{}{},
	}

	prop, err := voter.GenericMessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectGenericPayloadLen.Error())
}

func (s *GenericHandlerTestSuite) TestGenericHandleEventIncorrectMetadata() {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(2), 32)...)
	message := &message.Message{
		Source:       1,
		Destination:  0,
		DepositNonce: 1,
		ResourceId:   [32]byte{0},
		Type:         message.FungibleTransfer,
		Payload: []interface{}{
			"incorrectMetadata", // metadata
		},
	}

	prop, err := voter.GenericMessageHandler(message, common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"), common.HexToAddress("0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b"))

	s.Nil(prop)
	s.NotNil(err)
	s.EqualError(err, errIncorrectGenericMetadata.Error())
}
