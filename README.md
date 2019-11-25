## OP_SCRIPT

A viewer and debugger of Bitcoin scripts. **Early development.**

![Screenshot](./screenshot.png)


## Features
1. Reads transactions from a Bitcoin node (requires a full node with `txindex=1`).
1. Automatically finds related output.
1. Allows to navigate forward and backward.
1. Shows stack per line of code.
1. Supports witness data (SegWit).
1. Uses [`btcd/txscript`](https://github.com/btcsuite/btcd/tree/master/txscript) under the hood.


## Usage

1. `go get github.com/Jeiwan/opscript`
1. `opscript --help`
    ```shell
    Usage:
    opscript [flags]

    Flags:
    -h, --help               help for opscript
        --input int          Index of the input to debug a script from.
        --node-addr string   Bitcoin node address. (default "127.0.0.1:8332")
        --rpc-pass string    Bitcoin JSON-RPC password.
        --rpc-user string    Bitcoin JSON-RPC username.
        --tx string          Hash of the transaction to debug a script from.
    ```


## Key bindings

* `q` – quit
* `↑`/`↓` – navigate between lines of code


## Example
```shell
opscript --rpc-user=woot --rpc-pass=woot --tx=70fde4687efab8dae09737f87e30042030288fec42fd9e12f34c435cdeb7812c
```