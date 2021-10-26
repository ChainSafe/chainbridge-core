package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type TestClient interface {
	E2EClient
	LatestBlock() (*big.Int, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
}

func PreSetupTestSuite(fabric1, fabric2 calls.TxFabric, endpoint1, endpoint2 string, adminKey *secp256k1.Keypair) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		fabric1:   fabric1,
		fabric2:   fabric2,
		endpoint1: endpoint1,
		endpoint2: endpoint2,
		adminKey:  adminKey,
	}
}

type IntegrationTestSuite struct {
	suite.Suite
	client            TestClient
	client2           TestClient
	bridgeAddr        common.Address
	erc20HandlerAddr  common.Address
	erc20ContractAddr common.Address
	fabric1           calls.TxFabric
	fabric2           calls.TxFabric
	endpoint1         string
	endpoint2         string
	adminKey          *secp256k1.Keypair
}

func (s *IntegrationTestSuite) SetupSuite()    {}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest() {
	ethClient, err := evmclient.NewEVMClientFromParams(s.endpoint1, s.adminKey.PrivateKey())
	if err != nil {
		panic(err)
	}
	ethClient2, err := evmclient.NewEVMClientFromParams(s.endpoint2, s.adminKey.PrivateKey())
	if err != nil {
		panic(err)
	}
	b, err := ethClient.LatestBlock()
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("Latest block %s", b.String())
	bridgeAddr, erc20Addr, erc20HandlerAddr, err := PrepareEVME2EEnv(ethClient, s.fabric1, 1, big.NewInt(1), s.adminKey.CommonAddress())
	if err != nil {
		panic(err)
	}
	s.client = ethClient
	s.client2 = ethClient2
	s.bridgeAddr = bridgeAddr
	s.erc20ContractAddr = erc20Addr
	s.erc20HandlerAddr = erc20HandlerAddr
	//Contract addresses should be the same
	_, _, _, err = PrepareEVME2EEnv(ethClient2, s.fabric2, 2, big.NewInt(1), s.adminKey.CommonAddress())
	if err != nil {
		panic(err)
	}
}

func (s *IntegrationTestSuite) TestErc20Deposit() {
	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()
	senderBalBefore, err := calls.GetERC20Balance(s.client, s.erc20ContractAddr, EveKp.CommonAddress())
	s.Nil(err)
	destBalanceBefore, err := calls.GetERC20Balance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	resourceID := calls.SliceTo32Bytes(append(common.LeftPadBytes(s.erc20ContractAddr.Bytes(), 31), uint8(0)))
	gasPricer := evmgaspricer.NewStaticGasPriceDeterminant(s.client, nil)
	err = calls.Deposit(s.client, s.fabric1, gasPricer, s.bridgeAddr, dstAddr, amountToDeposit, resourceID, 2)
	s.Nil(err)

	//Wait 120 seconds for relayer vote
	time.Sleep(120 * time.Second)

	senderBalAfter, err := calls.GetERC20Balance(s.client, s.erc20ContractAddr, s.adminKey.CommonAddress())
	s.Nil(err)
	s.Equal(-1, senderBalAfter.Cmp(senderBalBefore))

	destBalanceAfter, err := calls.GetERC20Balance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(1, destBalanceAfter.Cmp(destBalanceBefore))
}
