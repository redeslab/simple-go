package main

/*
#include "callback.h"
*/
import "C"
import (
	"fmt"
	"github.com/redeslab/go-lib/proxy"
)

const (
	Success    = iota
	WalletFile = "wallet.json"
	DataBase   = "data"
)

type appConf struct {
	baseDir    string
	walletPath string
	dbPath     string
	BasIP      string
}

func (ac appConf) String() string {
	str := fmt.Sprintf("\n+++++++++++++++++++++++++++Application Config+++++++++++++++++++++++++"+
		"\n base dir:%s"+
		"\n wallet path:%s"+
		"\n dbPath path:%s"+
		"\n BAS ip:%s"+
		"\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++",
		ac.baseDir, ac.walletPath, ac.dbPath, ac.BasIP)

	return str
}

type MacApp struct {
	appConf
	HopThreadSig chan struct{}
	vpnSrv       *proxy.VpnProxy
	uiAPI        C.UserInterfaceAPI
}

var _appInst = &MacApp{}

//export initConf
func initConf(baseDir, tokenAddr, mpsAddr, apiUrl, dns string, uiApi C.UserInterfaceAPI) *C.char {

	return nil
}

//export startApp
func startApp() *C.char {
	return nil
}

//export stopApp
func stopApp() {

}

//export startServing
func startServing(srvAddr, poolAddr, minerID string) (int, *C.char) {

	return Success, nil
}

//export stopService
func stopService() {
}
