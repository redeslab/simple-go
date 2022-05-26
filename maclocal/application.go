package main

import "C"
import (
	"fmt"
	"github.com/redeslab/go-lib/proxy"
)

const (
	Success    = iota
	WalletFile = "wallet.json"
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
}

var _appInst = &MacApp{}

//export initConf
func initConf(baseDir, tokenAddr, mpsAddr, apiUrl, dns string) error {

	return nil
}

//export startApp
func startApp() error {

	return nil
}

//export stopApp
func stopApp() {

}

//export startServing
func startServing(srvAddr, poolAddr, minerID string) (int, error) {
	return Success, nil
}

//export stopService
func stopService() {
	if _appInst.vpnSrv != nil {
		_appInst.vpnSrv.Stop()
	}
}

//export testPings
func testPings(mid string) (string, float32) {
	return "ip", 0
}

//export dnsAddr
func dnsAddr() string {
	return _appInst.BasIP
}

//export syncAllPoolsData
func syncAllPoolsData() {
}
