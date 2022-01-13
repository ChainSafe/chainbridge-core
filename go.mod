module github.com/ChainSafe/chainbridge-core

go 1.17

require (
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.12
	github.com/golang/mock v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	github.com/status-im/keycard-go v0.0.0-20211004132608-c32310e39b86
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v0.24.0
	go.opentelemetry.io/otel/metric v0.24.0
	go.opentelemetry.io/otel/sdk/export/metric v0.24.0
	go.opentelemetry.io/otel/sdk/metric v0.24.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
)

require (
	github.com/imdario/mergo v0.3.12
	github.com/mitchellh/mapstructure v1.4.2
	github.com/pierrec/xxHash v0.1.5 // indirect
)
