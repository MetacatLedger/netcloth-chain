package app

import (
	"encoding/json"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"os"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/app/protocol"
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/version"
)

const (
	appName = "nch"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.nchcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.nchd")
)

func MakeLatestCodec() *codec.Codec {
	return v0.MakeCodec()
}

type NCHApp struct {
	*BaseApp
	invCheckPeriod uint
}

func NewNCHApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *NCHApp {
	bApp := NewBaseApp(appName, logger, db, baseAppOptions...)

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	//bApp.SetAnteHandler(ante.NewAnteHandler(bApp.accountKeeper, bApp.supplyKeeper, ante.DefaultSigVerificationGasConsumer))
	//bApp.SetFeeRefundHandler(auth.NewFeeRefundHandler(app.accountKeeper, app.supplyKeeper, app.refundKeeper))

	protocolKeeper := sdk.NewProtocolKeeper(protocol.MainKVStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	bApp.SetProtocolEngine(&engine)

	if !bApp.fauxMerkleMode {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeIAVL)
	} else {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeDB)
	}

	bApp.MountKVStores(protocol.Keys)
	bApp.MountTransientStores(protocol.TKeys)

	if loadLatest {
		err := bApp.LoadLatestVersion(protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, bApp.DeliverTx, nil))
	loaded, current := engine.LoadCurrentProtocol(bApp.cms.GetKVStore(protocol.MainKVStoreKey))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	bApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	var app = &NCHApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
	}

	return app
}

func NewNCHAppForReplay(logger log.Logger, db dbm.DB, traceStore io.Writer, loadInit, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *NCHApp {
	bApp := NewBaseApp(appName, logger, db, baseAppOptions...)

	protocolKeeper := sdk.NewProtocolKeeper(protocol.MainKVStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	bApp.SetProtocolEngine(&engine)

	if !bApp.fauxMerkleMode {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeIAVL)
	} else {
		bApp.MountStore(protocol.MainKVStoreKey, sdk.StoreTypeDB)
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, bApp.DeliverTx, nil))
	loaded, current := engine.LoadCurrentProtocol(bApp.cms.GetKVStore(protocol.Keys[MainStoreKey]))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}

	bApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	var app = &NCHApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
	}

	app.MountKVStores(protocol.Keys)
	app.MountTransientStores(protocol.TKeys)

	//app.SetBeginBlocker(app.BeginBlocker)
	//app.SetAnteHandler(ante.NewAnteHandler(app.accountKeeper, app.supplyKeeper, ante.DefaultSigVerificationGasConsumer))
	//app.SetFeeRefundHandler(auth.NewFeeRefundHandler(app.accountKeeper, app.supplyKeeper, app.refundKeeper))
	//app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	} else if loadInit {
		err := app.LoadVersion(0, protocol.MainKVStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

func SetBech32AddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
}

func (app *NCHApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.MainKVStoreKey)
}

func (app *NCHApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, jailWhiteList)
}
