package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	bitcoinWallet "github.com/Amirilidan78/bitcoin-wallet"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

var privateKeyHex = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
var toAddressHex = "tb1qppv790u4dz48ctnk3p7ss7fmspckagp3wrfyp0"
var chain = &chaincfg.TestNet3Params

func createWallet() *bitcoinWallet.BitcoinWallet {
	w, _ := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, privateKeyHex)
	return w
}

func main() {

	tx()

}

func tx() {

	amount := int64(2789695)
	fee := int64(1000)

	wallet := createWallet()

	priv, _ := wallet.PrivateKeyBTCE()

	tx, _ := createTransaction(chain, priv, wallet.Address, toAddressHex, amount, fee)

	fmt.Println(tx)

}

func test() {

	totalAmount := int64(788000)
	amount := int64(100000)
	fee := int64(1000)

	privKey := "cS9Zef6XdN3jHTFJFSsyJAtmDgCCdnygyVUJsLoyB8neuwhidUNJ"
	spendAddrStr := "tb1qppv790u4dz48ctnk3p7ss7fmspckagp3wrfyp0"
	destAddrStr := toAddressHex
	chain := &chaincfg.TestNet3Params
	txHash := "0d36d447b86f6ec93b4ff6c743fd5a50f9e1f884ec9bdc27f5c53365837cc29e"
	position := 1

	spendAddr, err := btcutil.DecodeAddress(spendAddrStr, chain)
	if err != nil {
		log.Println("DecodeAddress spendAddr err", err)
		return
	}

	destAddr, err := btcutil.DecodeAddress(destAddrStr, chain)
	if err != nil {
		log.Println("DecodeAddress destAddrStr err", err)
		return
	}

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		log.Println("wif err", err)
		return
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return
	}

	destAddrByte, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		log.Println("destAddr PayToAddrScript err", err)
		return
	}

	// == //

	utxoHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Println("NewHashFromStr err", err)
		return
	}

	redeemTx := wire.NewMsgTx(2)

	outPoint := wire.NewOutPoint(utxoHash, uint32(position))
	txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
	txIn.Sequence = txIn.Sequence - 2
	redeemTx.AddTxIn(txIn)

	redeemTxOut0 := wire.NewTxOut(amount, destAddrByte)
	redeemTxOut1 := wire.NewTxOut(totalAmount-amount-fee, spenderAddrByte)

	redeemTx.AddTxOut(redeemTxOut0)
	redeemTx.AddTxOut(redeemTxOut1)
	redeemTx.LockTime = 2407372

	if err != nil {
		log.Println("DecodeString pkScript err", err)
		return
	}

	sigHashes := txscript.NewTxSigHashes(redeemTx, txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: {},
	}))

	fmt.Println("totalAmount")
	fmt.Println(totalAmount)

	signature, err := txscript.WitnessSignature(redeemTx, sigHashes, 0, totalAmount, spenderAddrByte, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		log.Println("WitnessSignature err", err)
		return
	}
	redeemTx.TxIn[0].Witness = signature

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("")
	fmt.Println("hexSignedTx")
	fmt.Println(hexSignedTx)

}
