// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/opentelemetry"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Run() error {
	errChn := make(chan error)
	stopChn := make(chan struct{})

	db, err := lvldb.NewLvlDB(viper.GetString(config.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}

	//EVM1 setup
	evm1Client := evmclient.NewEVMClient()
	err = evm1Client.Configurate(viper.GetString(config.ChainConfigFlagName), "config_evm1.json")
	if err != nil {
		panic(err)
	}
	evm1Cfg := evm1Client.GetConfig()

	eventHandler := listener.NewETHEventHandler(common.HexToAddress(evm1Cfg.SharedEVMConfig.Bridge), evm1Client)
	eventHandler.RegisterEventHandler(evm1Cfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	eventHandler.RegisterEventHandler(evm1Cfg.SharedEVMConfig.Erc721Handler, listener.Erc721EventHandler)
	eventHandler.RegisterEventHandler(evm1Cfg.SharedEVMConfig.GenericHandler, listener.GenericEventHandler)
	evm1Listener := listener.NewEVMListener(evm1Client, eventHandler, common.HexToAddress(evm1Cfg.SharedEVMConfig.Bridge))

	mh := voter.NewEVMMessageHandler(evm1Client, common.HexToAddress(evm1Cfg.SharedEVMConfig.Bridge))
	mh.RegisterMessageHandler(common.HexToAddress(evm1Cfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	mh.RegisterMessageHandler(common.HexToAddress(evm1Cfg.SharedEVMConfig.Erc721Handler), voter.ERC721MessageHandler)
	mh.RegisterMessageHandler(common.HexToAddress(evm1Cfg.SharedEVMConfig.GenericHandler), voter.GenericMessageHandler)

	evmeVoter, err := voter.NewVoterWithSubscription(mh, evm1Client, evmtransaction.NewTransaction, evmgaspricer.NewLondonGasPriceClient(evm1Client, nil))
	if err != nil {
		panic(err)
	}
	evm1Chain := evm.NewEVMChain(evm1Listener, evmeVoter, db, *evm1Cfg.SharedEVMConfig.GeneralChainConfig.Id, &evm1Cfg.SharedEVMConfig)

	////EVM2 setup
	evm2Client := evmclient.NewEVMClient()
	err = evm2Client.Configurate(viper.GetString(config.ChainConfigFlagName), "config_evm2.json")
	if err != nil {
		panic(err)
	}

	evm2Config := evm2Client.GetConfig()

	eventHandlerEVM := listener.NewETHEventHandler(common.HexToAddress(evm2Config.SharedEVMConfig.Bridge), evm2Client)
	eventHandlerEVM.RegisterEventHandler(evm2Config.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	eventHandlerEVM.RegisterEventHandler(evm2Config.SharedEVMConfig.Erc721Handler, listener.Erc721EventHandler)
	eventHandlerEVM.RegisterEventHandler(evm2Config.SharedEVMConfig.GenericHandler, listener.GenericEventHandler)
	evm2Listener := listener.NewEVMListener(evm2Client, eventHandlerEVM, common.HexToAddress(evm2Config.SharedEVMConfig.Bridge))

	mhEVM := voter.NewEVMMessageHandler(evm2Client, common.HexToAddress(evm2Config.SharedEVMConfig.Bridge))
	mhEVM.RegisterMessageHandler(common.HexToAddress(evm2Config.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	mhEVM.RegisterMessageHandler(common.HexToAddress(evm2Config.SharedEVMConfig.Erc721Handler), voter.ERC721MessageHandler)
	mhEVM.RegisterMessageHandler(common.HexToAddress(evm2Config.SharedEVMConfig.GenericHandler), voter.GenericMessageHandler)

	evm2Voter, err := voter.NewVoterWithSubscription(mhEVM, evm2Client, evmtransaction.NewTransaction, evmgaspricer.NewLondonGasPriceClient(evm2Client, nil))
	if err != nil {
		panic(err)
	}
	evm2Chain := evm.NewEVMChain(evm2Listener, evm2Voter, db, *evm2Config.SharedEVMConfig.GeneralChainConfig.Id, &evm2Config.SharedEVMConfig)

	r := relayer.NewRelayer([]relayer.RelayedChain{evm2Chain, evm1Chain}, &opentelemetry.ConsoleTelemetry{})

	go r.Start(stopChn, errChn)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	select {
	case err := <-errChn:
		log.Error().Err(err).Msg("failed to listen and serve")
		close(stopChn)
		return err
	case sig := <-sysErr:
		log.Info().Msgf("terminating got ` [%v] signal", sig)
		return nil
	}
}
