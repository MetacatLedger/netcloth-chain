package types

import (
    sdk "github.com/NetCloth/netcloth-chain/types"
)

const (
    DefaultCodespace sdk.CodespaceType = ModuleName

    CodeEmptyInputs         sdk.CodeType = 100
    CodeStringTooLong       sdk.CodeType = 101
    CodeEndpointsFormatErr  sdk.CodeType = 102

    CodeBadDenom            sdk.CodeType = 111
    CodeBondInsufficient    sdk.CodeType = 112

    CodeMonikerExist        sdk.CodeType = 113
)

func ErrEmptyInputs(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}

func ErrStringTooLong(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeStringTooLong, msg)
}

func ErrBadDenom(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeBadDenom, msg)
}

func ErrBondInsufficient(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeBondInsufficient, msg)
}

func ErrMonikerExist(msg string) sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeMonikerExist, msg)
}

func ErrEndpointsFormatErr() sdk.Error {
    return sdk.NewError(DefaultCodespace, CodeEndpointsFormatErr, "endpoints format err, should be in format: serviceType|endpoint,serviceType|endpoint, serviceType is a number, endpoint is a string")
}
