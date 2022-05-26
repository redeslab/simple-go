package main

import (
	"errors"
	"github.com/redeslab/go-simple/account"
	"math/big"
)

//export NewWallet
func NewWallet(auth string) (bool, error) {
	w, e := account.NewWallet(auth, false)
	if e != nil {
		return false, e
	}

	if e := w.SaveToPath(_appInst.walletPath); e != nil {
		return false, e
	}
	return true, nil
}

type WalletBalance struct {
	Eth      *big.Int `json:"eth"`
	Token    *big.Int `json:"token"`
	Approved *big.Int `json:"approved"`
}

//export ImportWalletFrom
func ImportWalletFrom(path, auth string) error {
	w, e := account.LoadWallet(path)
	if e != nil {
		return e
	}

	if e := w.Open(auth); e != nil {
		return e
	}

	if e := w.SaveToPath(_appInst.walletPath); e != nil {
		return e
	}

	return nil
}

//export ExportWalletTo
func ExportWalletTo(path string) error {
	w, e := account.LoadWallet(_appInst.walletPath)
	if e != nil {
		return e
	}
	str := w.String()
	if str == "" {
		return errors.New("Invalid wallet data")
	}
	if e := w.SaveToPath(path); e != nil {
		return e
	}
	return nil
}

//export isWalletOpen
func isWalletOpen() bool {
	return false
}

//export openWallet
func openWallet(auth string) error {

	return nil
}

//export closeWallet
func closeWallet() {
}
