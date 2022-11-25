package bitcoinWallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/Amirilidan78/bitcoin-wallet/blockBook"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"strconv"
)

func getAddressUTXOFromBlockBook(chain *chaincfg.Params, address string) ([]blockBook.Utxo, error) {

	bb := blockBook.NewHttpBlockBookService(enums.MAIN_NODE)
	if &chaincfg.TestNet3Params == chain {
		bb = blockBook.NewHttpBlockBookService(enums.TEST_NODE)
	}

	return bb.GetAddressUTXO(address)
}

func prepareUTXOForTransaction(chain *chaincfg.Params, address string, amount int64, fee int64) ([]blockBook.Utxo, int64, error) {

	records, err := getAddressUTXOFromBlockBook(chain, address)
	if err != nil {
		return nil, 0, err
	}

	var final []blockBook.Utxo
	var total int64

	for _, record := range records {

		if total >= (amount + fee) {
			break
		}

		if record.Confirmations > -1 {

			final = append(final, record)

			txAmount, err := strconv.Atoi(record.Value)
			if err != nil {
				continue
			}

			total += int64(txAmount)
		}
	}

	return final, total, nil
}

func createTransactionAndSignTransaction(chain *chaincfg.Params, fromAddress string, privateKey *btcec.PrivateKey, toAddress string, amount int64, fee int64) (*wire.MsgTx, error) {

	fromAddr, err := btcutil.DecodeAddress(fromAddress, chain)
	if err != nil {
		return nil, errors.New("DecodeAddress fromAddr err " + err.Error())
	}

	fromAddrScriptByte, err := txscript.PayToAddrScript(fromAddr)

	if err != nil {
		return nil, errors.New("fromAddr PayToAddrScript err " + err.Error())
	}

	toAddr, err := btcutil.DecodeAddress(toAddress, chain)
	if err != nil {
		return nil, errors.New("DecodeAddress destAddrStr err " + err.Error())
	}

	toAddrByte, err := txscript.PayToAddrScript(toAddr)
	if err != nil {
		return nil, errors.New("toAddr PayToAddrScript err " + err.Error())
	}

	fromAddrByte, err := txscript.PayToAddrScript(fromAddr)
	if err != nil {
		return nil, errors.New("fromAddrByte PayToAddrScript err " + err.Error())
	}

	utxoList, totalAmount, err := prepareUTXOForTransaction(chain, fromAddress, amount, fee)
	if err != nil {
		return nil, errors.New("vin err " + err.Error())
	}
	if totalAmount < amount || len(utxoList) == 0 {
		return nil, errors.New("insufficient balance")
	}

	t, err := createTransactionInputsAndSign(privateKey, utxoList, fromAddrByte, fromAddrScriptByte, toAddrByte, totalAmount, amount, fee)
	if err != nil {
		return nil, errors.New("vin err " + err.Error())
	}

	return t, nil
}

func createTransactionInputsAndSign(privateKey *btcec.PrivateKey, utxos []blockBook.Utxo, fromAddressByte []byte, fromAddressScriptByte []byte, toAddressByte []byte, totalAmount int64, amount int64, fee int64) (*wire.MsgTx, error) {

	transaction := wire.NewMsgTx(2)

	// vin
	for _, utxo := range utxos {

		hash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}

		txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, [][]byte{})
		txIn.Sequence = txIn.Sequence - 2
		transaction.AddTxIn(txIn)
	}

	// vout
	changeAmount := totalAmount - amount - fee
	transaction.AddTxOut(wire.NewTxOut(amount, toAddressByte))
	if changeAmount > 0 {
		transaction.AddTxOut(wire.NewTxOut(changeAmount, fromAddressByte))
	}

	transaction.LockTime = 0

	signerMap := make(map[wire.OutPoint]*wire.TxOut)
	for _, in := range transaction.TxIn {
		signerMap[in.PreviousOutPoint] = &wire.TxOut{}
	}
	sigHashes := txscript.NewTxSigHashes(transaction, txscript.NewMultiPrevOutFetcher(signerMap))

	// sign
	for index, utxo := range utxos {

		amount, err := strconv.ParseInt(utxo.Value, 10, 64)
		if err != nil {
			return nil, errors.New("ParseInt utxo value err " + err.Error())
		}

		signature, err := txscript.WitnessSignature(transaction, sigHashes, index, amount, fromAddressScriptByte, txscript.SigHashAll, privateKey, true)
		if err != nil {
			return nil, errors.New("WitnessSignature err " + err.Error())
		}

		transaction.TxIn[index].Witness = signature
	}

	return transaction, nil
}

func getRawTransaction(tx *wire.MsgTx) (string, error) {

	var signedTx bytes.Buffer

	err := tx.Serialize(&signedTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signedTx.Bytes()), nil
}

func broadcastHex(chain *chaincfg.Params, hex string) (string, error) {

	bb := blockBook.NewHttpBlockBookService(enums.MAIN_NODE)
	if &chaincfg.TestNet3Params == chain {
		bb = blockBook.NewHttpBlockBookService(enums.TEST_NODE)
	}

	res, err := bb.BroadcastTransaction(hex)
	if err != nil {
		return "", err
	}

	return res.TxId, nil
}

func createSignAndBroadcastTransaction(chain *chaincfg.Params, privateKey *btcec.PrivateKey, fromAddress string, toAddress string, amount int64, fee int64) (string, error) {

	// signed tx
	tx, err := createTransactionAndSignTransaction(chain, fromAddress, privateKey, toAddress, amount, fee)
	if err != nil {
		return "", err
	}

	// raw
	raw, err := getRawTransaction(tx)
	if err != nil {
		return "", err
	}

	// broadcast
	return broadcastHex(chain, raw)
}
