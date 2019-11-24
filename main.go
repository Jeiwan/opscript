package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jeiwan/scriptdbg/debugger"
	"github.com/Jeiwan/scriptdbg/gui"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/sirupsen/logrus"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var nodeAddr = flag.String("node-addr", "127.0.0.1:8332", "Bitcoin node address.")
	var nodeRPCUser = flag.String("node-rpc-user", "", "Bitcoin node RPC user.")
	var nodeRPCPass = flag.String("node-rpc-pass", "", "Bitcoin node RPC password.")
	var txid = flag.String("txid", "", "Transaction ID")
	var input = flag.Int("input", -1, "Input number.")

	flag.Parse()

	btcclient, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         *nodeAddr,
		User:         *nodeRPCUser,
		Pass:         *nodeRPCPass,
	}, nil)
	if err != nil {
		logrus.Fatal(fmt.Errorf("new Bitcoin client: %+v", err))
	}
	defer btcclient.Shutdown()

	txHash, err := chainhash.NewHashFromStr(*txid)
	if err != nil {
		logrus.Fatal(fmt.Errorf("parse txid: %+v", err))
	}

	txResp, err := btcclient.GetRawTransaction(txHash)
	if err != nil {
		logrus.Fatal(fmt.Errorf("get raw transaction: %+v", err))
	}

	prevOut := &txResp.MsgTx().TxIn[*input].PreviousOutPoint
	prevTxHash := prevOut.Hash
	prevTxResp, err := btcclient.GetRawTransaction(&prevTxHash)
	if err != nil {
		logrus.Fatal(fmt.Errorf("get prev transaction: %+v", err))
	}

	en, err := newEngine(txResp.MsgTx(), prevTxResp.MsgTx().TxOut[prevOut.Index].PkScript)
	if err != nil {
		logrus.Fatal(fmt.Errorf("new engine: %+v", err))
	}

	d, err := debugger.NewWithEngine(en)
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
}
