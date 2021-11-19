package message

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
)

type MessageProcessor func(message *Message) error

// AdjustDecimalsForERC20AmountMessageProcessor is a function, that accepts message and map[domainID uint8]{decimal uint}
// using this  params processor converts amount for one chain to another for provided decimals with floor rounding
func AdjustDecimalsForERC20AmountMessageProcessor(args ...interface{}) MessageProcessor {
	return func(m *Message) error {
		if len(args) == 0 {
			return errors.New("processor requires 1 argument")
		}
		decimalsMap, ok := args[0].(map[uint8]uint64)
		if !ok {
			return errors.New("no decimals map found in args")
		}
		sourceDecimal, ok := decimalsMap[m.Source]
		if !ok {
			return errors.New("no source decimals found at decimalsMap")
		}
		destDecimal, ok := decimalsMap[m.Destination]
		if !ok {
			return errors.New("no destination decimals found at decimalsMap")
		}
		amountByte, ok := m.Payload[0].([]byte)
		if !ok {
			return errors.New("could not cast interface to byte slice")
		}
		amount := new(big.Int).SetBytes(amountByte)
		if sourceDecimal > destDecimal {
			diff := sourceDecimal - destDecimal
			roundedAmount := big.NewInt(0)
			roundedAmount.Div(amount, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(diff), nil))
			log.Info().Msgf("amount %s rounded to %s from chain %v to chain %v", amount.String(), roundedAmount.String(), m.Source, m.Destination)
			m.Payload[0] = roundedAmount.Bytes()
			return nil
		}
		if sourceDecimal < destDecimal {
			diff := destDecimal - sourceDecimal
			roundedAmount := big.NewInt(0)
			roundedAmount.Mul(amount, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(diff), nil))
			m.Payload[0] = roundedAmount.Bytes()
			log.Info().Msgf("amount %s rounded to %s from chain %v to chain %v", amount.String(), roundedAmount.String(), m.Source, m.Destination)
		}
		return nil
	}
}
