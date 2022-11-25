# bitcoin-wallet
bitcoin wallet package for creating and generating wallet, transferring BTC, getting wallet unspent transactions(UTXOs), getting wallet txIs , getting wallet balance and crawling blocks to find wallet transactions
 
### Supported nodes
check `enums/nodes` file alternatively you can create your own node using trezor blockBook nodes
```
node := enums.CreateNode("https://btc1.trezor.io","wss://btc1.trezor.io",false)
```

### Wallet methods

Generating bitcoin wallet
```
w := GenerateBitcoinWallet(node)
w.Address // strnig 
w.PrivateKey // strnig 
w.PublicKey // strnig 
```

Creating bitcoin wallet from private key
```
w := CreateBitcoinnWallet(node,privateKeyHex)
w.Address // strnig 
w.PrivateKey // strnig 
w.PublicKey // strnig 
```

Getting wallet bitcoin balance
```
balanceInSatoshi,err := w.Balance()
balanceInSatoshi // int64
```

Getting wallet UTXOs
```
utxos,err := w.UTXOs()
utxos // []blockBook.Utxo
```

Getting wallet txIds
```
txIds,err := w.TxIds()
txIds // []string
```

Transfer btc
```
txId,err := w.Transfer("tb1q0r23g66m9rhhak8aahsg53wfp5egt2huuc4tnu",10000,1000)
txId // string
```

### BTC Faucet
check this website https://coinfaucet.eu/en/btc-testnet

### Important
I simplified this repository github.com/btcsuite/btcd repository to create this package You can check go it for better examples and functionalities and do not use this package in production, I created this package for education purposes, 
and thanks to [eltNEG](https://github.com/eltNEG) really helped me to build this package