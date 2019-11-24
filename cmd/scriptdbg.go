package cmd

import (
	"fmt"

	"github.com/Jeiwan/scriptdbg/debugger"
	"github.com/Jeiwan/scriptdbg/gui"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// New ...
func New() *cobra.Command {
	return newRootCmd()
}

func newRootCmd() *cobra.Command {
	var nodeAddr, rpcUser, rpcPass, txHash string
	var txInput int

	cmd := &cobra.Command{
		Use: "scriptdbg",
		RunE: func(cmd *cobra.Command, args []string) error {
			btcclient, err := rpcclient.New(&rpcclient.ConnConfig{
				HTTPPostMode: true,
				DisableTLS:   true,
				Host:         nodeAddr,
				User:         rpcUser,
				Pass:         rpcPass,
			}, nil)
			if err != nil {
				logrus.Fatal(fmt.Errorf("new Bitcoin client: %+v", err))
			}
			defer btcclient.Shutdown()

			txHash, err := chainhash.NewHashFromStr(txHash)
			if err != nil {
				logrus.Fatal(fmt.Errorf("parse txid: %+v", err))
			}

			txResp, err := btcclient.GetRawTransaction(txHash)
			if err != nil {
				logrus.Fatal(fmt.Errorf("get raw transaction: %+v", err))
			}

			prevOut := &txResp.MsgTx().TxIn[txInput].PreviousOutPoint
			prevTxHash := prevOut.Hash
			prevTxResp, err := btcclient.GetRawTransaction(&prevTxHash)
			if err != nil {
				logrus.Fatal(fmt.Errorf("get prev transaction: %+v", err))
			}

			en, err := newEngine(txResp.MsgTx(), prevTxResp.MsgTx().TxOut[prevOut.Index].PkScript)
			if err != nil {
				logrus.Fatal(fmt.Errorf("new engine: %+v", err))
			}

			d, err := debugger.New(en)
			if err != nil {
				logrus.Fatalln(err)
			}

			gui, err := gui.New(d)
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

	cmd.Flags().StringVar(&nodeAddr, "node-addr", "127.0.0.1:8332", "Bitcoin node address.")
	cmd.Flags().StringVar(&rpcUser, "rpc-user", "", "Bitcoin JSON-RPC username.")
	cmd.Flags().StringVar(&rpcPass, "rpc-pass", "", "Bitcoin JSON-RPC password.")
	cmd.Flags().StringVar(&txHash, "tx", "", "Hash of the transaction to debug a script from.")
	cmd.Flags().IntVar(&txInput, "input", 0, "Index of the input to debug a script from.")

	return cmd
}

func newEngine(tx *wire.MsgTx, output []byte) (*txscript.Engine, error) {
	e, err := txscript.NewEngine(
		output,
		tx,
		0,
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
