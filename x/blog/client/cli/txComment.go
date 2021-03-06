package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/example/blog/x/blog/types"
)

// CmdCreateComment() returns create-comment trasaction for this module
func CmdCreateComment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-comment [body] [postID]",
		Short: "Creates a new comment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsBody := string(args[0])
			argsPostID := string(args[1])
			s, err1 := strconv.ParseUint(argsPostID, 10, 64)
			if err1 != nil {
				return err1
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateComment(clientCtx.GetFromAddress().String(), string(argsBody), s)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdateComment() returns update-comment trasaction for this module
func CmdUpdateComment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-comment [id] [body] [postID]",
		Short: "Update a comment",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argsBody := string(args[1])
			argsPostID := string(args[2])
			a, err1 := strconv.ParseUint(argsPostID, 10, 64)
			if err1 != nil {
				return err1
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateComment(clientCtx.GetFromAddress().String(), id, string(argsBody), a)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdDeleteComment() returns delete-comment trasaction for this module
func CmdDeleteComment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-comment [id]",
		Short: "Delete a comment by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteComment(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
