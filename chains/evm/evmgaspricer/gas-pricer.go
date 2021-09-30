package evmgaspricer

import (
	"context"
	"math/big"
)



type LondonGasClient interface {
	GasPriceClient
	BaseFee() (*big.Int, error)
	EstimateGasLondon(ctx context.Context, baseFee *big.Int) (*big.Int, *big.Int, error)
}

type GasPriceClient interface {
	GasPrice() (*big.Int, error)
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
	gp, err := gasPricer.client.GasPrice()
	if err != nil {
		return nil, err
	}

	var gasPrices []*big.Int
	gasPrices[0] = gp

	return gasPrices, nil
}

type LondonGasPriceDeterminant struct {
	client LondonGasClient
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
	if baseFee != nil {
		gasTipCap, gasFeeCap, err := gasPricer.client.EstimateGasLondon(context.TODO(), baseFee)
		if err != nil {
			return nil, err
		}
		gasPrices[0] = gasTipCap
		gasPrices[1] = gasFeeCap
	} else {
		gp, err := gasPricer.client.GasPrice()
		if err != nil {
			return nil, err
		}
		gasPrices[0] = gp
	}
	return gasPrices, nil
}
