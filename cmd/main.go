package main

import (
	"fmt"
	bitcoinWallet "github.com/Amirilidan78/bitcoin-wallet"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
)

func main() {

	priv := generate()
	create(priv)
}

func generate() string {
	wallet := bitcoinWallet.GenerateBitcoinWallet(enums.MAIN_NODE)
	fmt.Println(wallet.Address)
	return wallet.PrivateKey
}

func create(priv string) {
	wallet, err := bitcoinWallet.CreateBitcoinWallet(enums.MAIN_NODE, priv)
	fmt.Println(err)
	fmt.Println(wallet.Address)
}
