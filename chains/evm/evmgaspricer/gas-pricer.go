package evmgaspricer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
)

type LondonGasClient interface {
	GasPriceClient
	BaseFee() (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

type GasPriceClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

// StaticGasPriceDeterminant for when you want to always use generic `GasPrice()` method from an EVM client.
//
// Client should implement `GasPrice()` to query first for a gas price field that is set on client construction
// This way a developer can use a specific gas price for transactions, such as in the CLI
//
// Currently, if the client being used is created by the `EVMClientFromParams` constructor a constant gas price is then set
// and will be returned by this gas pricer
type StaticGasPriceDeterminant struct {
	upperLimitFeePerGas *big.Int
	gasPriceMultiplayer *big.Float
}
func NewStaticGasPriceDeterminant(upperLimitFeePerGas *big.Int, gasPriceMultiplayer *big.Float) *StaticGasPriceDeterminant {
	return &StaticGasPriceDeterminant{upperLimitFeePerGas, gasPriceMultiplayer}

}

func (gasPricer *StaticGasPriceDeterminant) GasPrice(client GasPriceClient) ([]*big.Int, error) {
	gp, err := client.SuggestGasPrice(context.TODO())
	log.Debug().Msgf("Suggested GP %s", gp.String())
	if err != nil {
		return nil, err
	}
	if gasPricer.gasPriceMultiplayer != nil {
		gp = multiplyGasPrice(gp, gasPricer.gasPriceMultiplayer)
	}
	if gasPricer.upperLimitFeePerGas != nil {
		if gp.Cmp(gasPricer.upperLimitFeePerGas) == 1 {
			gp = gasPricer.upperLimitFeePerGas
		}
	}
	var gasPrices []*big.Int
	gasPrices[0] = gp
	return gasPrices, nil
}

type LondonGasPriceDeterminant struct {
	upperLimitFeePerGas *big.Int
}

func NewLondonGasPriceDeterminant(UpperLimitFeePerGas *big.Int) *LondonGasPriceDeterminant {
	return &LondonGasPriceDeterminant{upperLimitFeePerGas: UpperLimitFeePerGas}
}

func (gasPricer *LondonGasPriceDeterminant) GasPrice(client LondonGasClient) ([]*big.Int, error) {
	baseFee, err := client.BaseFee()
	if err != nil {
		return nil, err
	}
	var gasPrices []*big.Int
	// BaseFee could be nil if eip1559 is not implemented or did not started working on the current chain
	if baseFee == nil {
		// we are using staticGasPriceDeterminant because it counts configs in its gasPrice calculations
		// and seem to be the most favorable option
		staticGasPricer := NewStaticGasPriceDeterminant(gasPricer.upperLimitFeePerGas, big.NewFloat(1))
		return staticGasPricer.GasPrice(client)
	}
	gasTipCap, gasFeeCap, err := gasPricer.estimateGasLondon(client, baseFee)
	if err != nil {
		return nil, err
	}
	gasPrices[0] = gasTipCap
	gasPrices[1] = gasFeeCap
	return gasPrices, nil
}

const TwoAndTheHalfGwei = 2500000000

func (gasPricer *LondonGasPriceDeterminant) estimateGasLondon(client LondonGasClient, baseFee *big.Int) (*big.Int, *big.Int, error) {
	var maxPriorityFeePerGas *big.Int
	var maxFeePerGas *big.Int

	// if gasPriceLimit is set and lower than networks baseFee then
	// maxPriorityFee is set to 3 GWEI because that was practically and theoretically defined as optimum
	// and Max Fee set to baseFee + maxPriorityFeePerGas
	if gasPricer.upperLimitFeePerGas != nil && gasPricer.upperLimitFeePerGas.Cmp(baseFee) < 0 {
		maxPriorityFeePerGas = big.NewInt(TwoAndTheHalfGwei)
		maxFeePerGas = new(big.Int).Add(baseFee, maxPriorityFeePerGas)
		return maxPriorityFeePerGas, maxFeePerGas, nil
	}

	maxPriorityFeePerGas, err := client.SuggestGasTipCap(context.TODO())
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
	if gasPricer.upperLimitFeePerGas != nil && maxFeePerGas.Cmp(gasPricer.upperLimitFeePerGas) == 1 {
		maxPriorityFeePerGas.Sub(gasPricer.upperLimitFeePerGas, baseFee)
		maxFeePerGas = gasPricer.upperLimitFeePerGas
	}
	return maxPriorityFeePerGas, maxFeePerGas, nil
}

func multiplyGasPrice(gasEstimate *big.Int, gasMultiplier *big.Float) *big.Int {
	gasEstimateFloat := new(big.Float).SetInt(gasEstimate)
	result := gasEstimateFloat.Mul(gasEstimateFloat, gasMultiplier)
	gasPrice := new(big.Int)
	result.Int(gasPrice)
	return gasPrice
}

