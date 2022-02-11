package evmgaspricer

import (
	"context"
	"math/big"
)

type LondonGasPriceDeterminant struct {
	client LondonGasClient
	opts   *GasPricerOpts
}

func NewLondonGasPriceClient(client LondonGasClient, opts *GasPricerOpts) *LondonGasPriceDeterminant {
	return &LondonGasPriceDeterminant{client: client, opts: opts}
}

func (gasPricer *LondonGasPriceDeterminant) GasPrice(priority *uint8) ([]*big.Int, error) {
	baseFee, err := gasPricer.client.BaseFee()
	if err != nil {
		return nil, err
	}
	gasPrices := make([]*big.Int, 2)
	// BaseFee could be nil if eip1559 is not implemented or did not started working on the current chain
	if baseFee == nil {
		// we are using staticGasPriceDeterminant because it counts configs in its gasPrice calculations
		// and seem to be the most favorable option
		staticGasPricer := NewStaticGasPriceDeterminant(gasPricer.client, gasPricer.opts)
		return staticGasPricer.GasPrice(nil)
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
	gasPricer.client = client
}
func (gasPricer *LondonGasPriceDeterminant) SetOpts(opts *GasPricerOpts) {
	gasPricer.opts = opts
}

const TwoAndTheHalfGwei = 2500000000 // Lowest MaxPriorityFee. Defined by some researches...

func (gasPricer *LondonGasPriceDeterminant) estimateGasLondon(baseFee *big.Int) (*big.Int, *big.Int, error) {
	var maxPriorityFeePerGas *big.Int
	var maxFeePerGas *big.Int

	// if gasPriceLimit is set and lower than networks baseFee then
	// maxPriorityFee is set to 3 GWEI because that was practically and theoretically defined as optimum
	// and Max Fee set to baseFee + maxPriorityFeePerGas
	if gasPricer.opts != nil && gasPricer.opts.UpperLimitFeePerGas != nil && gasPricer.opts.UpperLimitFeePerGas.Cmp(baseFee) < 0 {
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

	// Check we aren't exceeding our limit if gasPriceLimit set
	if gasPricer.opts != nil && gasPricer.opts.UpperLimitFeePerGas != nil && maxFeePerGas.Cmp(gasPricer.opts.UpperLimitFeePerGas) == 1 {
		maxPriorityFeePerGas.Sub(gasPricer.opts.UpperLimitFeePerGas, baseFee)
		maxFeePerGas = gasPricer.opts.UpperLimitFeePerGas
	}
	return maxPriorityFeePerGas, maxFeePerGas, nil
}
