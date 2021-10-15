package evmgaspricer

import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"
)

// StaticGasPriceDeterminant for when you want to always use generic `GasPrice()` method from an EVM client.
//
// Client should implement `GasPrice()` to query first for a gas price field that is set on client construction
// This way a developer can use a specific gas price for transactions, such as in the CLI
//
// Currently, if the client being used is created by the `EVMClientFromParams` constructor a constant gas price is then set
// and will be returned by this gas pricer
type StaticGasPriceDeterminant struct {
	client              GasPriceClient
	upperLimitFeePerGas *big.Int
	gasPriceMultiplayer *big.Float
}

func NewStaticGasPriceDeterminant(client GasPriceClient, opts *GasPricerOpts) *StaticGasPriceDeterminant {
	return &StaticGasPriceDeterminant{client: client, upperLimitFeePerGas: opts.UpperLimitFeePerGas, gasPriceMultiplayer: opts.GasPriceMultiplayer}

}

func (gasPricer *StaticGasPriceDeterminant) GasPrice() ([]*big.Int, error) {
	gp, err := gasPricer.client.SuggestGasPrice(context.TODO())
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
