package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Fantom-foundation/go-opera/ftmclient"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	client, err := ftmclient.Dial("/Users/hadv/Library/Lachesis/fakenet-1/opera.ipc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ks := keystore.NewKeyStore("/Users/hadv/Library/Lachesis/fakenet-1/keystore/", keystore.StandardScryptN, keystore.StandardScryptP)

	to := common.HexToAddress("0xc94e731C761A0985c5b1212b9b6208362465b328")
	from := common.HexToAddress("0x15060c1ad0484d426fBB1170B528705ce9450C35")
	data := []byte("Fantom")
	value := big.NewInt(10000000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("gasPrice: %v", gasPrice)
	fmt.Println()
	gasTip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("gasTip: %v", gasTip)
	fmt.Println()
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
	fmt.Printf("gasLimit: %v", gasLimit)
	fmt.Println()
	// gasLimit = gasLimit + 21000
	// fmt.Printf("gasLimit: %v", gasLimit)
	// fmt.Println()
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
	signer := types.NewLondonSigner(networkID)
	sender, err := signer.Sender(signedTx)
	if err != nil {
		fmt.Println("sender: " + err.Error())
		return
	}
	fmt.Printf("sender: %v", sender.Hex())
	fmt.Println()

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println("send: " + err.Error())
		return
	}
	fmt.Println("Send tnx succesfully!")
}
