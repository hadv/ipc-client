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
	client, err := ethclient.Dial("/root/.opera/fakenet-7/opera.ipc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ks := keystore.NewKeyStore("/root/.opera/fakenet-7/keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	to := common.HexToAddress("0xc94e731C761A0985c5b1212b9b6208362465b328")
	from := common.HexToAddress("0x802d1c2560B9b5884DA3CB85E8f426e97C354101")
	data := []byte("Fantom")
	value := big.NewInt(1000000000000000000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	gasTip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	msg := ethereum.CallMsg{
		From:      from,
		To:        &to,
		GasFeeCap: gasPrice,
		GasTipCap: gasTip,
		Value:     value,
		Data:      data,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	gasLimit = gasLimit + 21000
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	newTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   networkID,
		Nonce:     nonce,
		To:        &to,
		Value:     value,
		Gas:       gasLimit,
		GasTipCap: gasTip,
		GasFeeCap: gasPrice,
		Data:      data,
	})
	fmt.Printf("tx type: %v", newTx.Type())
	fmt.Println()

	signedTx, err := ks.SignTxWithPassphrase(accounts.Account{Address: from}, "i3nxx1rk", newTx, networkID)
	if err != nil {
		fmt.Println("sign: " + err.Error())
		return
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println("send: " + err.Error())
		return
	}
	fmt.Println("Send tnx succesfully!")
}
