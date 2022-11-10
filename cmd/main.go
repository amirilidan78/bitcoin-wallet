package main

import (
	"fmt"
	bitcoinWallet "github.com/Amirilidan78/bitcoin-wallet"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
)

const privateKey = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"

func main() {

	wallet, _ := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, privateKey)
	fmt.Println(wallet.Address)
	fmt.Println(wallet.PrivateKey)
	fmt.Println(wallet.PublicKey)
	fmt.Println(wallet.Balance())
	fmt.Println(wallet.UTXOs())
}

func generate() string {
	wallet := bitcoinWallet.GenerateBitcoinWallet(enums.TEST_NODE)
	fmt.Println(wallet.PrivateKey)
	fmt.Println(wallet.PublicKey)
	fmt.Println(wallet.Address)
	return wallet.PrivateKey
}

func create(priv string) {
	wallet, err := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, priv)
	fmt.Println(err)
	fmt.Println(wallet.Address)
}
