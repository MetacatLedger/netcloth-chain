package cli

import (
	"github.com/NetCloth/netcloth-chain/client/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"

	"github.com/NetCloth/netcloth-chain/client"
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/codec"
	"github.com/NetCloth/netcloth-chain/modules/auth"
	"github.com/NetCloth/netcloth-chain/modules/auth/client/utils"
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// GetTxCmd returns the transaction commands for this module
func IPALCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "IPAL transaction subcommands",
	}
	txCmd.AddCommand(
		IPALClaimCmd(cdc),
		ServerNodeClaimCmd(cdc),
	)
	return txCmd
}

func IPALClaimCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim",
		Short:   "Create and sign a IPALClaim tx",
		Example: "nchcli ipal claim  --user=<user key name> --proxy=<proxy key name> --ip=<server ip>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtxUser := context.NewCLIContextWithFrom(viper.GetString(flagUser)).WithCodec(cdc)

			info, err := txBldr.Keybase().Get(cliCtxUser.GetFromName())
			if err != nil {
				return err
			}
			userAddress := info.GetAddress().String()

			// build user request signature
			serverIP := viper.GetString(flagServerIP)
			expiration := time.Now().UTC().AddDate(0, 0, 1)
			adMsg := types.NewADParam(userAddress, serverIP, expiration)

			// build msg
			passphrase, err := keys.GetPassphrase(cliCtxUser.GetFromName())
			if err != nil {
				return err
			}
			// sign
			sigBytes, pubkey, err := txBldr.Keybase().Sign(info.GetName(), passphrase, adMsg.GetSignBytes())
			if err != nil {
				return err
			}
			stdSig := auth.StdSignature{
				PubKey:    pubkey,
				Signature: sigBytes,
			}

			// build and sign the transaction, then broadcast to Tendermint
			cliCtxProxy := context.NewCLIContextWithFrom(viper.GetString(flagProxy)).WithCodec(cdc)
			msg := types.NewMsgIPALClaim(cliCtxProxy.GetFromAddress(), userAddress, serverIP, expiration, stdSig)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtxProxy, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagServerIP, "", "server ip")
	cmd.Flags().String(flagUser, "", "user account")
	cmd.Flags().String(flagProxy, "", "proxy account")
	cmd.MarkFlagRequired(flagServerIP)
	cmd.MarkFlagRequired(flagUser)
	cmd.MarkFlagRequired(flagProxy)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func ServerNodeClaimCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server-node-claim",
		Short:   "Create and sign a ServerNodeClaim tx",
		Example: "nchcli ipal server-node-claim  --from=<user key name> --moniker=<name> --identity=<identity> --website=<website> --server_endpoint=<server_endpoint> --details=<details>",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			moniker := viper.GetString(flagMoniker)
			website := viper.GetString(flagWebsite)
			serverEndPoint := viper.GetString(flagServerEndPoint)
			details := viper.GetString(flagDetails)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgServiceNodeClaim(cliCtx.GetFromAddress(),moniker, website, serverEndPoint, details)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagMoniker, "", "server node moniker")
	cmd.Flags().String(flagWebsite, "", "server node website")
	cmd.Flags().String(flagServerEndPoint, "", "server node endpoint")
	cmd.Flags().String(flagDetails, "", "server node details")

	cmd.MarkFlagRequired(flagMoniker)
	cmd.MarkFlagRequired(flagServerEndPoint)

	cmd = client.PostCommands(cmd)[0]

	return cmd
}