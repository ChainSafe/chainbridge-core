package processors

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/observability"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/attribute"
)

type MessageProcessor func(ctx context.Context, message *message.Message) error

// AdjustDecimalsForERC20AmountMessageProcessor is a function, that accepts message and map[domainID uint8]{decimal uint}
// using this  params processor converts amount for one chain to another for provided decimals with floor rounding
func AdjustDecimalsForERC20AmountMessageProcessor(args ...interface{}) MessageProcessor {
	return func(ctx context.Context, m *message.Message) error {
		_, span, logger := observability.CreateSpanAndLoggerFromContext(
			ctx,
			"relayer-core",
			"relayer.core.MessageProcessor.AdjustDecimalsForERC20AmountMessageProcessor",
			attribute.String("msg.id", m.ID()), attribute.String("msg.type", string(m.Type)))
		defer span.End()
		if len(args) == 0 {
			return observability.LogAndRecordErrorWithStatus(&logger, span, errors.New("processor requires 1 argument"), "failed to processMessage with AdjustDecimalsForERC20AmountMessageProcessor")
		}
		decimalsMap, ok := args[0].(map[uint8]uint64)
		if !ok {
			return observability.LogAndRecordErrorWithStatus(&logger, span, errors.New("no decimals map found in args"), "failed to processMessage with AdjustDecimalsForERC20AmountMessageProcessor")
		}
		sourceDecimal, ok := decimalsMap[m.Source]
		if !ok {
			return observability.LogAndRecordErrorWithStatus(&logger, span, errors.New("no source decimals found at decimalsMap"), "failed to processMessage with AdjustDecimalsForERC20AmountMessageProcessor")
		}
		destDecimal, ok := decimalsMap[m.Destination]
		if !ok {
			return observability.LogAndRecordErrorWithStatus(&logger, span, errors.New("no destination decimals found at decimalsMap"), "failed to processMessage with AdjustDecimalsForERC20AmountMessageProcessor")
		}
		amountByte, ok := m.Payload[0].([]byte)
		if !ok {
			return observability.LogAndRecordErrorWithStatus(&logger, span, errors.New("could not cast interface to byte slice"), "failed to processMessage with AdjustDecimalsForERC20AmountMessageProcessor")
		}
		amount := new(big.Int).SetBytes(amountByte)
		roundedAmount := big.NewInt(0)
		if sourceDecimal > destDecimal {
			diff := sourceDecimal - destDecimal
			roundedAmount.Div(amount, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(diff), nil))
			m.Payload[0] = roundedAmount.Bytes()
		} else if sourceDecimal < destDecimal {
			diff := destDecimal - sourceDecimal
			roundedAmount.Mul(amount, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(0).SetUint64(diff), nil))
			m.Payload[0] = roundedAmount.Bytes()
		}
		observability.LogAndEvent(logger.Info(), span, fmt.Sprintf("amount %s rounded to %s from chain %v to chain %v", amount.String(), roundedAmount.String(), m.Source, m.Destination))
		return nil
	}
}
