// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/netcloth/netcloth-chain/app/v0/crisis/types
package crisis

import (
	"github.com/netcloth/netcloth-chain/app/v0/crisis/internal/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/crisis/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	RegisterCodec         = types.RegisterCodec
	ErrNoSender           = types.ErrNoSender
	ErrUnknownInvariant   = types.ErrUnknownInvariant
	NewGenesisState       = types.NewGenesisState
	DefaultGenesisState   = types.DefaultGenesisState
	NewMsgVerifyInvariant = types.NewMsgVerifyInvariant
	ParamKeyTable         = types.ParamKeyTable
	NewInvarRoute         = types.NewInvarRoute
	NewKeeper             = keeper.NewKeeper

	// variable aliases
	ModuleCdc                = types.ModuleCdc
	ParamStoreKeyConstantFee = types.ParamStoreKeyConstantFee
)

type (
	GenesisState       = types.GenesisState
	MsgVerifyInvariant = types.MsgVerifyInvariant
	InvarRoute         = types.InvarRoute
	Keeper             = keeper.Keeper
)