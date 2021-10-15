package evmgaspricer

import (
	"context"
	"fmt"
	"math/big"
)

type LondonGasPriceDeterminant struct {
	client              LondonGasClient
	upperLimitFeePerGas *big.Int
}

func NewLondonGasPriceClient(client LondonGasClient, opts *GasPricerOpts) *LondonGasPriceDeterminant {
	return &LondonGasPriceDeterminant{client: client, upperLimitFeePerGas: opts.UpperLimitFeePerGas}
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
		staticGasPricer := NewStaticGasPriceDeterminant(gasPricer.client, &GasPricerOpts{GasPriceMultiplayer: nil, UpperLimitFeePerGas: gasPricer.upperLimitFeePerGas})
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

func (gasPricer *LondonGasPriceDeterminant) SetClient(client LondonGasClient) {

}
func (gasPricer *LondonGasPriceDeterminant) SetOpts(opts *GasPricerOpts) {

}

const TwoAndTheHalfGwei = 2500000000

func (gasPricer *LondonGasPriceDeterminant) estimateGasLondon(baseFee *big.Int) (*big.Int, *big.Int, error) {
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
	if gasPricer.upperLimitFeePerGas != nil && maxFeePerGas.Cmp(gasPricer.upperLimitFeePerGas) == 1 {
		maxPriorityFeePerGas.Sub(gasPricer.upperLimitFeePerGas, baseFee)
		maxFeePerGas = gasPricer.upperLimitFeePerGas
	}
	return maxPriorityFeePerGas, maxFeePerGas, nil
}
