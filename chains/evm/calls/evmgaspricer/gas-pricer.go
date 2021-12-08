package evmgaspricer

import (
	"context"
	"math/big"
)

type LondonGasClient interface {
	GasPriceClient
	BaseFee() (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

type GasPriceClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

// GasPricerOpts is the structure that holds parameters that could be used to configure different gasPRicer implementation
type GasPricerOpts struct {
	UpperLimitFeePerGas *big.Int      // UpperLimitFeePerGas in Static and London gasPricer limits the maximum gas price that could be used. In London gasPricer if BaseFee > UpperLimitFeePerGas, then maxFeeCap will be BaseFee + 2.5 Gwei for MaxTipCap. If nil - not applied
	GasPriceFactor      *big.Float    // GasPriceFactor In static gasPricer multiplies final gasPrice. Could be for example 0.75 or 5.
	Args                []interface{} // Args is the array of dynamic typed args that could be used for other custom GasPricer implementations
}

func multiplyGasPrice(gasEstimate *big.Int, gasMultiplier *big.Float) *big.Int {
	gasEstimateFloat := new(big.Float).SetInt(gasEstimate)
	result := gasEstimateFloat.Mul(gasEstimateFloat, gasMultiplier)
	gasPrice := new(big.Int)
	result.Int(gasPrice)
	return gasPrice
}
