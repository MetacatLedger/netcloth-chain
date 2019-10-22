package ipal

import (
	"github.com/NetCloth/netcloth-chain/modules/ipal/keeper"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
)

const (
	DefaultCodespace = types.DefaultCodespace

	StoreKey  = types.StoreKey
	RouterKey = types.RouterKey
)

var (
	// functions aliases
	RegisterCodec                  = types.RegisterCodec
	ErrEmptyInputs                 = types.ErrEmptyInputs
	ErrStringTooLong               = types.ErrStringTooLong
	ErrInvalidSignature            = types.ErrInvalidSignature
	ErrIPALClaimUserRequestExpired = types.ErrIPALClaimUserRequestExpired

	NewIPALObject       = types.NewIPALObject
	NewServerNodeObject = types.NewServerNodeObject

	NewMsgIPALClaim        = types.NewMsgIPALClaim
	NewMsgServiceNodeClaim = types.NewMsgServiceNodeClaim

	// variable aliases
	ModuleCdc = types.ModuleCdc

	NewKeeper = keeper.NewKeeper

	AttributeValueCategory = types.AttributeValueCategory
)

type (
	Keeper = keeper.Keeper

	MsgIPALClaim        = types.MsgIPALClaim
	MsgServiceNodeClaim = types.MsgServiceNodeClaim
)