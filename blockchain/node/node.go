package node

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/sirupsen/logrus"
)

// Node is a Bitcoin node.
type Node struct {
	host    string
	rpcUser string
	rpcPass string
}

// New returns a new Node.
func New(host, rpcUser, rpcPass string) *Node {
	return &Node{
		host:    host,
		rpcUser: rpcUser,
		rpcPass: rpcPass,
	}
}

// GetTransaction returns a transaction by its hash.
func (n Node) GetTransaction(txHash string) (*wire.MsgTx, error) {
	btcclient, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         n.host,
		User:         n.rpcUser,
		Pass:         n.rpcPass,
	}, nil)
	if err != nil {
		logrus.Fatal(fmt.Errorf("new Bitcoin client: %+v", err))
	}
	defer btcclient.Shutdown()

	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, fmt.Errorf("parse txid: %+v", err)
	}

	txResp, err := btcclient.GetRawTransaction(hash)
	if err != nil {
		return nil, fmt.Errorf("get raw transaction: %+v", err)
	}

	return txResp.MsgTx(), nil
}
