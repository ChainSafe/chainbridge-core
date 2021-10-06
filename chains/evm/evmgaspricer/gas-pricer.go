package evmgaspricer

import (
	"context"
	"fmt"
	"math/big"
)

// TODO: how to pass confif there?????
type LondonGasClient interface {
	GasPriceClient
	BaseFee() (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

type GasPriceClient interface {
	EstimateGas() (*big.Int, error)
}

// DefaultGasPrice for when you want to always use generic `GasPrice()` method from an EVM client.
//
// Client should implement `GasPrice()` to query first for a gas price field that is set on client construction
// This way a developer can use a specific gas price for transactions, such as in the CLI
//
// Currently, if the client being used is created by the `EVMClientFromParams` constructor a constant gas price is then set
// and will be returned by this gas pricer
type StaticGasPriceDeterminant struct {
	client GasPriceClient
}

func NewStaticGasPriceDeterminant(client GasPriceClient) *StaticGasPriceDeterminant {
	return &StaticGasPriceDeterminant{client: client}
}

func (gasPricer *StaticGasPriceDeterminant) GasPrice() ([]*big.Int, error) {
	gp, err := gasPricer.client.EstimateGas()
	if err != nil {
		return nil, err
	}

	var gasPrices []*big.Int
	gasPrices[0] = gp

	return gasPrices, nil
}

type LondonGasPriceDeterminant struct {
	client LondonGasClient
	BaseFee *big.Int
	MaxPriorityFee *big.Int
	MaxFeePerGas *big.Int
	UpperLimitFeePerGas *big.Int
}

func NewLondonGasPriceDeterminant(client LondonGasClient) *LondonGasPriceDeterminant {
	return &LondonGasPriceDeterminant{client: client}
}

func (gasPricer *LondonGasPriceDeterminant) GasPrice() ([]*big.Int, error) {
	baseFee, err := gasPricer.client.BaseFee()
	if err != nil {
		return nil, err
	}
	var gasPrices []*big.Int
	// BaseFee could be nil if eip1559 is not implemented or did not started working on the current chain
	if baseFee == nil {
		// we are using staticGasPriceDeterminant because it counts configs in its gasPrice calculations
		// and seem to be the most favorable option
		staticGasPricer := NewStaticGasPriceDeterminant(gasPricer.client)
		return staticGasPricer.GasPrice()
	}
	gasTipCap, gasFeeCap, err := gasPricer.estimateGasLondon(baseFee)
	if err != nil {
		return nil, err
	}
	gasPrices[0] = gasTipCap
	gasPrices[1] = gasFeeCap
	return gasPrices, nil
}

const TwoAndTheHalfGwei = 2500000000

func (gasPricer *LondonGasPriceDeterminant) estimateGasLondon(baseFee *big.Int) (*big.Int, *big.Int, error) {
	var maxPriorityFeePerGas *big.Int
	var maxFeePerGas *big.Int


	// if gasPriceLimit is set and lower than networks baseFee then
	// maxPriorityFee is set to 3 GWEI because that was practically and theoretically defined as optimum
	// and Max Fee set to baseFee + maxPriorityFeePerGas
	if gasPricer.UpperLimitFeePerGas != nil && gasPricer.UpperLimitFeePerGas.Cmp(baseFee) < 0 {
		maxPriorityFeePerGas = big.NewInt(TwoAndTheHalfGwei)
		maxFeePerGas = new(big.Int).Add(baseFee, maxPriorityFeePerGas)
		return maxPriorityFeePerGas, maxFeePerGas, nil
	}

	maxPriorityFeePerGas, err := gasPricer.client.SuggestGasTipCap(context.TODO())
	if err != nil {
		return nil, nil, err
	}
	maxFeePerGas = new(big.Int).Add(
		maxPriorityFeePerGas,
		new(big.Int).Mul(baseFee, big.NewInt(2)),
	)

	if maxFeePerGas.Cmp(maxPriorityFeePerGas) < 0 {
		return nil, nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", maxFeePerGas, maxPriorityFeePerGas)
	}
	// Check we aren't exceeding our limit if gasPriceLimit set
	if gasPricer.UpperLimitFeePerGas != nil && maxFeePerGas.Cmp(gasPricer.UpperLimitFeePerGas) == 1 {
		maxPriorityFeePerGas.Sub(gasPricer.UpperLimitFeePerGas, baseFee)
		maxFeePerGas = gasPricer.UpperLimitFeePerGas
	}
	return maxPriorityFeePerGas, maxFeePerGas, nil
}
