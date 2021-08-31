package evmgaspricer

import (
	"context"
	"math/big"
)

type DefaultGasPricer interface {
	GasPrice() (*big.Int, error)
}

type LondonGasPricer interface {
	DefaultGasPricer
	BaseFee() (*big.Int, error)
	EstimateGasLondon(ctx context.Context, baseFee *big.Int) (*big.Int, *big.Int, error)
}

// DefaultGasPrice for when you want to always use generic `GasPrice()` method from an EVM client.
//
// Client should implement `GasPrice()` to query first for a gas price field that is set on client construction
// This way a developer can use a specific gas price for transactions, such as in the CLI
//
// Currently, if the client being used is created by the `EVMClientFromParams` constructor a constant gas price is then set
// and will be returned by this gas pricer
type StaticGasPriceDeterminant struct {
	client DefaultGasPricer
}

func NewStaticGasPriceDeterminant(client DefaultGasPricer) *StaticGasPriceDeterminant {
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

type LondonGasPricerDeterminant struct {
	client LondonGasPricer
}

func NewLondonGasPricerDeterminant(client LondonGasPricer) *LondonGasPricerDeterminant {
	return &LondonGasPricerDeterminant{client: client}
}

func (gasPricer *LondonGasPricerDeterminant) GasPrice() ([]*big.Int, error) {
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
