package main

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

func createTransaction(chain *chaincfg.Params, privateKey *btcec.PrivateKey, fromAddress string, toAddress string, amount int64, fee int64) (string, error) {

	// signed tx
	tx, err := createTransactionAndSignTransaction(chain, fromAddress, privateKey, toAddress, amount, fee)
	if err != nil {
		return "", err
	}

	// raw
	return getRawTransaction(tx)
}

// ======= create inputs ======= //

func createTransactionAndSignTransaction(chain *chaincfg.Params, fromAddress string, privateKey *btcec.PrivateKey, toAddress string, amount int64, fee int64) (*wire.MsgTx, error) {

	fromAddr, err := btcutil.DecodeAddress(fromAddress, chain)
	if err != nil {
		return nil, errors.New("DecodeAddress fromAddr err " + err.Error())
	}

	fromAddrByte, err := txscript.PayToAddrScript(fromAddr)

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

	utxoList, totalAmount, err := prepareUTXOForTransaction(chain, fromAddress, amount, fee)
	if err != nil {
		return nil, errors.New("vin err " + err.Error())
	}
	if totalAmount < amount || len(utxoList) == 0 {
		return nil, errors.New("insufficient balance")
	}

	t, err := createTransactionInputsAndSign(privateKey, utxoList, fromAddrByte, toAddrByte, totalAmount, amount, fee)
	if err != nil {
		return nil, errors.New("vin err " + err.Error())
	}

	return t, nil
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

func getAddressUTXOFromBlockBook(chain *chaincfg.Params, address string) ([]blockBook.Utxo, error) {

	bb := blockBook.NewHttpBlockBookService(enums.MAIN_NODE)
	if &chaincfg.TestNet3Params == chain {
		bb = blockBook.NewHttpBlockBookService(enums.TEST_NODE)
	}

	return bb.GetAddressUTXO(address)
}

func createTransactionInputsAndSign(privateKey *btcec.PrivateKey, utxos []blockBook.Utxo, fromAddressByte []byte, toAddressByte []byte, totalAmount int64, amount int64, fee int64) (*wire.MsgTx, error) {

	transaction := wire.NewMsgTx(2)

	// vin
	for i, utxo := range utxos {

		hash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}

		outPoint := wire.NewOutPoint(hash, uint32(i+1))
		txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
		txIn.Sequence = txIn.Sequence - 2
		transaction.AddTxIn(txIn)
	}

	// vout
	changeAmount := totalAmount - amount - fee
	redeemTxOut0 := wire.NewTxOut(amount, toAddressByte)
	transaction.AddTxOut(redeemTxOut0)
	if changeAmount > 0 {
		redeemTxOut1 := wire.NewTxOut(changeAmount, fromAddressByte)
		transaction.AddTxOut(redeemTxOut1)
	}

	transaction.LockTime = 2407372

	// sign
	for index, in := range transaction.TxIn {

		a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
			in.PreviousOutPoint: {},
		})
		sigHashes := txscript.NewTxSigHashes(transaction, a)

		signature, err := txscript.WitnessSignature(transaction, sigHashes, 0, totalAmount, fromAddressByte, txscript.SigHashAll, privateKey, true)
		if err != nil {
			return nil, errors.New("WitnessSignature err " + err.Error())
		}

		transaction.TxIn[index].Witness = signature
	}

	return transaction, nil
}

//func createTransactionVIn(tx *wire.MsgTx, utxos []blockBook.Utxo) (*wire.MsgTx, []*wire.OutPoint, error) {
//
//	var op []*wire.OutPoint
//
//	for i, utxo := range utxos {
//
//		hash, err := chainhash.NewHashFromStr(utxo.Txid)
//		if err != nil {
//			return nil, op, errors.New("NewHashFromStr err " + err.Error())
//		}
//
//		outPoint := wire.NewOutPoint(hash, uint32(i+1))
//
//		op = append(op, outPoint)
//
//		txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
//
//		txIn.Sequence = wire.MaxTxInSequenceNum
//
//		tx.AddTxIn(txIn)
//	}
//
//	return tx, op, nil
//}
//
//func createTransactionVOut(tx *wire.MsgTx, fromAddressByte []byte, toAddressByte []byte, totalAmount int64, amount int64, fee int64) (*wire.MsgTx, error) {
//
//	changeAmount := totalAmount - amount - fee
//
//	redeemTxOut0 := wire.NewTxOut(amount, toAddressByte)
//	tx.AddTxOut(redeemTxOut0)
//
//	if changeAmount > 0 {
//		redeemTxOut1 := wire.NewTxOut(changeAmount, fromAddressByte)
//		tx.AddTxOut(redeemTxOut1)
//	}
//
//	return tx, nil
//}
//
//func signTransaction(fromAddrByte []byte, privateKey *btcec.PrivateKey, amount int64, outPoints []*wire.OutPoint, tx *wire.MsgTx) (*wire.MsgTx, error) {
//
//	for index, _ := range tx.TxIn {
//
//		a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
//			*outPoints[0]: {},
//		})
//		sigHashes := txscript.NewTxSigHashes(tx, a)
//
//		signature, err := txscript.WitnessSignature(tx, sigHashes, index, amount, fromAddrByte, txscript.SigHashAll, privateKey, true)
//		if err != nil {
//			return nil, errors.New("WitnessSignature err " + err.Error())
//		}
//
//		tx.TxIn[index].Witness = signature
//	}
//
//	return tx, nil
//}

// ======= raw ======= //

func getRawTransaction(tx *wire.MsgTx) (string, error) {

	var signedTx bytes.Buffer

	err := tx.Serialize(&signedTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signedTx.Bytes()), nil
}
