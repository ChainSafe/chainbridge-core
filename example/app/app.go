// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/attribute"

	secp256k1 "github.com/ethereum/go-ethereum/crypto"

	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/events"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/monitored"
	"github.com/ChainSafe/chainbridge-core/chains/evm/executor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	secp256k12 "github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/e2e/dummy"
	"github.com/ChainSafe/chainbridge-core/flags"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/opentelemetry"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ChainSafe/chainbridge-core/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configuration, err := config.GetConfig(viper.GetString(flags.ConfigFlagName))
	if err != nil {
		panic(err)
	}

	db, err := lvldb.NewLvlDB(viper.GetString(flags.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}
	blockstore := store.NewBlockStore(db)

	OTLPResource := opentelemetry.InitResource(fmt.Sprintf("Relayer-%s", configuration.RelayerConfig.Id), configuration.RelayerConfig.Env)

	mp, err := opentelemetry.InitMetricProvider(ctx, OTLPResource, configuration.RelayerConfig.OpenTelemetryCollectorURL)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Error().Msgf("Error shutting down meter provider: %v", err)
		}
	}()

	tp, err := opentelemetry.InitTracesProvider(ctx, OTLPResource, configuration.RelayerConfig.OpenTelemetryCollectorURL)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error().Msgf("Error shutting down tracer provider: %v", err)
		}
	}()

	metrics, err := opentelemetry.NewRelayerMetrics(mp.Meter("relayer-metric-provider"), attribute.String("relayerid", configuration.RelayerConfig.Id), attribute.String("env", configuration.RelayerConfig.Env))
	if err != nil {
		panic(err)
	}

	chains := []relayer.RelayedChain{}
	for _, chainConfig := range configuration.ChainConfigs {
		switch chainConfig["type"] {
		case "evm":
			{
				config, err := chain.NewEVMConfig(chainConfig)
				if err != nil {
					panic(err)
				}

				privateKey, err := secp256k1.HexToECDSA(config.GeneralChainConfig.Key)
				if err != nil {
					panic(err)
				}

				kp := secp256k12.NewKeypair(*privateKey)

				client, err := evmclient.NewEVMClient(config.GeneralChainConfig.Endpoint, kp)
				if err != nil {
					panic(err)
				}

				dummyGasPricer := dummy.NewStaticGasPriceDeterminant(client, nil)
				t := monitored.NewMonitoredTransactor(evmtransaction.NewTransaction, dummyGasPricer, client, config.MaxGasPrice, config.GasPriceIncreaseFactor)
				go t.Monitor(ctx, time.Minute*3, time.Minute*10, time.Minute)
				bridgeContract := bridge.NewBridgeContract(client, common.HexToAddress(config.Bridge), t)

				depositHandler := listener.NewETHDepositHandler(bridgeContract)
				depositHandler.RegisterDepositHandler(config.Erc20Handler, listener.Erc20DepositHandler)
				depositHandler.RegisterDepositHandler(config.Erc721Handler, listener.Erc721DepositHandler)
				depositHandler.RegisterDepositHandler(config.GenericHandler, listener.GenericDepositHandler)
				eventListener := events.NewListener(client)
				eventHandlers := make([]listener.EventHandler, 0)
				eventHandlers = append(eventHandlers, listener.NewDepositEventHandler(eventListener, depositHandler, common.HexToAddress(config.Bridge), *config.GeneralChainConfig.Id))
				evmListener := listener.NewEVMListener(client, eventHandlers, blockstore, metrics, *config.GeneralChainConfig.Id, config.BlockRetryInterval, config.BlockConfirmations, config.BlockInterval)

				mh := executor.NewEVMMessageHandler(bridgeContract)
				mh.RegisterMessageHandler(config.Erc20Handler, executor.ERC20MessageHandler)
				mh.RegisterMessageHandler(config.Erc721Handler, executor.ERC721MessageHandler)
				mh.RegisterMessageHandler(config.GenericHandler, executor.GenericMessageHandler)

				var evmVoter *executor.EVMVoter
				evmVoter, err = executor.NewVoterWithSubscription(mh, client, bridgeContract)
				if err != nil {
					log.Error().Msgf("failed creating voter with subscription: %s. Falling back to default voter.", err.Error())
					evmVoter = executor.NewVoter(mh, client, bridgeContract)
				}

				chain := evm.NewEVMChain(evmListener, evmVoter, blockstore, *config.GeneralChainConfig.Id, config.StartBlock, config.GeneralChainConfig.LatestBlock, config.GeneralChainConfig.FreshStart)

				chains = append(chains, chain)
			}
		default:
			panic(fmt.Errorf("type '%s' not recognized", chainConfig["type"]))
		}
	}

	r := relayer.NewRelayer(
		chains,
		metrics,
	)

	errChn := make(chan error)
	go r.Start(ctx, errChn)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	select {
	case err := <-errChn:
		log.Error().Err(err).Msg("failed to listen and serve")
		return err
	case sig := <-sysErr:
		log.Info().Msgf("terminating got ` [%v] signal", sig)
		return nil
	}
}
