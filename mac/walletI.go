package main

import "C"
import (
	"github.com/redeslab/go-simple/account"
)

//export NewWallet
func NewWallet(auth string) (bool, *C.char) {
	w, e := account.NewWallet(auth, true)
	if e != nil {
		return false, C.CString(e.Error())
	}

	if e := w.SaveToPath(_appInst.walletPath); e != nil {
		return false, C.CString(e.Error())
	}
	return true, nil
}

//export ImportWalletFrom
func ImportWalletFrom(path, auth string) *C.char {
	w, e := account.LoadWallet(path)
	if e != nil {
		return C.CString(e.Error())
	}

	if e := w.Open(auth); e != nil {
		return C.CString(e.Error())
	}

	if e := w.SaveToPath(_appInst.walletPath); e != nil {
		return C.CString(e.Error())
	}

	return nil
}

//export ExportWalletTo
func ExportWalletTo(path string) *C.char {
	w, e := account.LoadWallet(_appInst.walletPath)
	if e != nil {
		return C.CString(e.Error())
	}
	str := w.String()
	if str == "" {
		return C.CString("Invalid wallet data")
	}
	if e := w.SaveToPath(path); e != nil {
		return C.CString(e.Error())
	}
	return nil
}
