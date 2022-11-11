package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	bitcoinWallet "github.com/Amirilidan78/bitcoin-wallet"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const privateKeyHex = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"

const utxoHash = "c77fe6125260df1702382ff89716dd5873051ffb95872d6d9407132ef52c4e84"
const utxoAmount = 67000

func createWallet() *bitcoinWallet.BitcoinWallet {
	w, _ := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, privateKeyHex)
	return w
}

func walletPrivateAndPublicKey(wallet *bitcoinWallet.BitcoinWallet) (*btcec.PrivateKey, *btcec.PublicKey) {
	b, _ := wallet.PrivateKeyBytes()
	return btcec.PrivKeyFromBytes(b)
}

func getFromAddress(wallet *bitcoinWallet.BitcoinWallet) btcutil.Address {
	address, _ := btcutil.DecodeAddress(wallet.Address, &chaincfg.TestNet3Params)
	return address
}

func getToAddress() btcutil.Address {
	address, _ := btcutil.DecodeAddress("tb1qtvnf9xcnyw34qrxc0aufqr34el7l4fec4fnknp", &chaincfg.TestNet3Params)
	return address
}

func main() {

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	amount := int64(utxoAmount) - 50

	wallet := createWallet()
	priv, pub := walletPrivateAndPublicKey(wallet)

	fmt.Println(wallet.Address)
	fromAddress := getFromAddress(wallet)
	toAddress := getToAddress()
	toAddressByte, err := txscript.PayToAddrScript(toAddress)

	bldr := txscript.NewScriptBuilder()
	bldr.AddOp(txscript.OP_1)
	bldr.AddData(priv.PubKey().SerializeCompressed())

	// sigScript, err := bldr.Script()

	hash, err := chainhash.NewHashFromStr(utxoHash)

	outPoint := wire.NewOutPoint(hash, 0)
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	redeemTxOut := wire.NewTxOut(amount, toAddressByte)
	redeemTx.AddTxOut(redeemTxOut)

	subscript := fromAddress.ScriptAddress()
	a := txscript.NewMultiPrevOutFetcher(nil)
	txSigHashes := txscript.NewTxSigHashes(redeemTx, a)
	wit, err := txscript.WitnessSignature(redeemTx, txSigHashes, 0, amount, subscript, txscript.SigHashAll, priv, true)

	fmt.Println(wit, err)

	witness := [][]byte{
		pub.SerializeCompressed(),
	}

	redeemTx.TxIn[0].Witness = witness

	finalRawTx, err := GetRawTransaction(redeemTx)
	fmt.Println("==========")
	fmt.Println(finalRawTx, err)

}

func GetRawTransaction(redeemTx *wire.MsgTx) (string, error) {

	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
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
