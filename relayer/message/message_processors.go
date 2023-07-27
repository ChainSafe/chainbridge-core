package message

import (
	"context"
	"errors"
	"math/big"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel"

	"github.com/rs/zerolog/log"
)

type MessageProcessor func(ctx context.Context, message *Message) error

// AdjustDecimalsForERC20AmountMessageProcessor is a function, that accepts message and map[domainID uint8]{decimal uint}
// using this  params processor converts amount for one chain to another for provided decimals with floor rounding
func AdjustDecimalsForERC20AmountMessageProcessor(args ...interface{}) MessageProcessor {
	return func(ctx context.Context, m *Message) error {
		tp := otel.GetTracerProvider()
		_, span := tp.Tracer("relayer-route").Start(ctx, "relayer.core.MessageProcessor.AdjustDecimalsForERC20AmountMessageProcessor")
		span.SetAttributes(attribute.String("msg_id", m.ID()), attribute.String("msg_type", string(m.Type)))
		defer span.End()
		if len(args) == 0 {
			span.SetStatus(codes.Error, "processor requires 1 argument")
			return errors.New("processor requires 1 argument")
		}
		decimalsMap, ok := args[0].(map[uint8]uint64)
		if !ok {
			span.SetStatus(codes.Error, "no decimals map found in args")
			return errors.New("no decimals map found in args")
		}
		sourceDecimal, ok := decimalsMap[m.Source]
		if !ok {
			span.SetStatus(codes.Error, "no source decimals found at decimalsMap")
			return errors.New("no source decimals found at decimalsMap")
		}
		destDecimal, ok := decimalsMap[m.Destination]
		if !ok {
			span.SetStatus(codes.Error, "no destination decimals found at decimalsMap")
			return errors.New("no destination decimals found at decimalsMap")
		}
		amountByte, ok := m.Payload[0].([]byte)
		if !ok {
			span.SetStatus(codes.Error, "could not cast interface to byte slice")
			return errors.New("could not cast interface to byte slice")
		}
		amount := new(big.Int).SetBytes(amountByte)
		if sourceDecimal > destDecimal {
			diff := sourceDecimal - destDecimal
			roundedAmount := big.NewInt(0)
			roundedAmount.Div(amount, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(diff), nil))
			log.Info().Str("msg_id", m.ID()).Msgf("amount %s rounded to %s from chain %v to chain %v", amount.String(), roundedAmount.String(), m.Source, m.Destination)
			m.Payload[0] = roundedAmount.Bytes()
			span.SetStatus(codes.Ok, "msg processed")
			return nil
		}
		if sourceDecimal < destDecimal {
			diff := destDecimal - sourceDecimal
			roundedAmount := big.NewInt(0)
			roundedAmount.Mul(amount, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(diff), nil))
			m.Payload[0] = roundedAmount.Bytes()
			log.Info().Str("msg_id", m.ID()).Msgf("amount %s rounded to %s from chain %v to chain %v", amount.String(), roundedAmount.String(), m.Source, m.Destination)
		}
		return nil
	}
}
