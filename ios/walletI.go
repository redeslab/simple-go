package iosLib

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/redeslab/go-simple/account"
)

func NewWallet(auth string) []byte {
	w, e := account.NewWallet(auth, true)
	if e != nil {
		return nil
	}
	b, e := json.Marshal(w)
	if e != nil {
		return nil
	}
	appCaller.Wallet = w
	return b
}

func ImportWallet(data, auth string) bool {
	w, e := account.LoadWalletByData(data)
	if e != nil {
		return false
	}

	if e := w.Open(auth); e != nil {
		return false
	}
	appCaller.Wallet = w
	return true
}

func OpenWallet(auth string) bool {
	if appCaller.Wallet == nil {
		return false
	}

	if appCaller.Wallet.IsOpen() {
		return true
	}

	if err := appCaller.Wallet.Open(auth); err != nil {
		return false
	}
	return true
}

func Address() string {
	if appCaller.Wallet == nil {
		return ""
	}

	return appCaller.Wallet.MainAddress().String()
}

func PriKeyData() []byte {
	if appCaller.Wallet == nil {
		return nil
	}

	priKey := appCaller.Wallet.SignKey()
	if priKey == nil {
		return nil
	}
	return priKey.D.Bytes()
}

func Verify(data, sig []byte) []byte {

	msg := crypto.Keccak256(data)
	fmt.Println("lib hash=>", hexutil.Encode(msg))

	addr := appCaller.Wallet.MainAddress()
	result := account.VerifyAbiSig(addr, sig, msg)
	fmt.Println("verify result=>", result, addr.String())

	msg2 := crypto.Keccak256([]byte("foo"))
	fmt.Println("lib hash2=>", hexutil.Encode(msg))

	sig2, err := appCaller.Wallet.Sign(msg2)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return sig2
}

func SubAddress() string {
	if appCaller.Wallet == nil {
		return ""
	}

	return appCaller.Wallet.SubAddress().String()
}
func SubPriKeyData() []byte {
	if appCaller.Wallet == nil {
		return nil
	}
	return appCaller.Wallet.CryptKey()
}

func IsOpen() bool {
	return appCaller.Wallet != nil && appCaller.Wallet.IsOpen()
}

func LoadWallet(walletJSOn string) bool {
	w, err := account.LoadWalletByData(walletJSOn)
	if err != nil {
		//TODO::
		fmt.Println("=======>LoadWallet Err:", err)
		return false
	}
	appCaller.Wallet = w
	return true
}

func AesKeyForMiner(minerAddr string) []byte {
	id, err := account.ConvertToID(minerAddr)
	if err != nil {
		fmt.Println("=======>AesKeyForMiner Err:", err)
		return nil
	}

	peerPub := id.ToPubKey()
	var aesKey account.PipeCryptKey
	if err := account.GenerateAesKey(&aesKey, peerPub, appCaller.Wallet.CryptKey()); err != nil {
		fmt.Println("=======>AesKeyForMiner Err:", err)
		return nil
	}
	return aesKey[:]
}
