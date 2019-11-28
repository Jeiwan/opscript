package cmd

import (
	"fmt"

	"github.com/Jeiwan/opscript/blockchain/blockstream"
	"github.com/Jeiwan/opscript/blockchain/node"
	"github.com/Jeiwan/opscript/debugger"
	"github.com/Jeiwan/opscript/gui"
	"github.com/Jeiwan/opscript/spec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Blockchain ...
type Blockchain interface {
	GetTransaction(txHash string) (*wire.MsgTx, error)
}

// New ...
func New(spec spec.Script) *cobra.Command {
	rootCmd := newRootCmd(spec)
	buildSpecCmd := newBuildSpecCmd()

	rootCmd.AddCommand(buildSpecCmd)

	return rootCmd
}

func newRootCmd(spec spec.Script) *cobra.Command {
	var useNode, useBlockstream bool
	var nodeAddr, rpcUser, rpcPass, txHash string
	var txInput int

	cmd := &cobra.Command{
		Use: "opscript",
		RunE: func(cmd *cobra.Command, args []string) error {
			var bchain Blockchain

			switch {
			case useBlockstream:
				bchain = blockstream.New()
			case useNode:
				bchain = node.New(nodeAddr, rpcUser, rpcPass)
			}

			tx, err := bchain.GetTransaction(txHash)
			if err != nil {
				return err
			}

			prevOut := tx.TxIn[txInput].PreviousOutPoint

			prevTx, err := bchain.GetTransaction(prevOut.Hash.String())
			if err != nil {
				logrus.Fatal(fmt.Errorf("get prev transaction: %+v", err))
			}

			en, err := newEngine(tx, prevTx.TxOut[prevOut.Index].PkScript, txInput)
			if err != nil {
				logrus.Fatal(fmt.Errorf("new engine: %+v", err))
			}

			d, err := debugger.New(en)
			if err != nil {
				logrus.Fatalln(err)
			}

			gui, err := gui.New(d, spec)
			if err != nil {
				logrus.Fatalln(err)
			}
			defer gui.Stop()

			if err := gui.Start(); err != nil {
				logrus.Fatalln(err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&useNode, "node", true, "Use Bitcoin node to get transactions (requires 'txindex=1').")
	cmd.Flags().BoolVar(&useBlockstream, "blockstream", false, "Use blockstream.info API to get transactions.")
	cmd.Flags().StringVar(&nodeAddr, "node-addr", "127.0.0.1:8332", "Bitcoin node address.")
	cmd.Flags().StringVar(&rpcUser, "rpc-user", "", "Bitcoin JSON-RPC username.")
	cmd.Flags().StringVar(&rpcPass, "rpc-pass", "", "Bitcoin JSON-RPC password.")
	cmd.Flags().StringVar(&txHash, "tx", "", "Hash of the transaction to debug a script from.")
	cmd.Flags().IntVar(&txInput, "input", 0, "Index of the input to debug a script from.")

	return cmd
}

func newEngine(tx *wire.MsgTx, output []byte, inputIdx int) (*txscript.Engine, error) {
	e, err := txscript.NewEngine(
		output,
		tx,
		inputIdx,
		txscript.ScriptVerifyWitness+txscript.ScriptBip16+txscript.ScriptVerifyCleanStack+txscript.ScriptVerifyMinimalData,
		nil,
		nil,
		-1,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
