package main

import (
	"context"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("/Users/hadv/work/github/hadv/go-ethereum/data/sun/geth.ipc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ks := keystore.NewKeyStore("/Users/hadv/work/github/hadv/go-ethereum/data/sun/keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	to := common.HexToAddress("0xd5089c1fdf8cebf58c6bbb50a86c1c55893634b8")
	from := common.HexToAddress("0xe348073d55ade0ef0e5696ba51d5565003233d0c")
	data := []byte("Lorem ipsum")
	value := big.NewInt(10000000000000000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	signedTx, err := ks.SignTxWithPassphrase(accounts.Account{Address: from}, "i3nxx1rk", newTx, networkID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Send tnx succesfully!")
}
