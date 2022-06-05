package androidLib

import (
	"encoding/json"
	"fmt"
	"github.com/redeslab/go-simple/account"
	"github.com/redeslab/go-simple/network"
	"github.com/redeslab/go-simple/node"
	"github.com/redeslab/go-simple/proxy"
	"github.com/redeslab/go-simple/tun2Pipe"
	"github.com/redeslab/go-simple/util"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

type VpnDelegate interface {
	ByPass(fd int32) bool
	io.Writer
	VpnClosed()
}

const (
	WalletFile = "wallet.json"
)

type AndroidAPP struct {
	IpCache     map[string]string
	wallet      account.Wallet
	WPath       string
	vpnDelegate VpnDelegate
}

var _appInst = &AndroidAPP{
	IpCache: make(map[string]string),
}

var VpnInst *proxy.VpnProxy = nil
var TunInst *tun2Pipe.Tun2Pipe = nil

func InitSystem(bypassIPs, baseDir string) error {
	if err := util.TouchDir(baseDir); err != nil {
		return err
	}
	walletPath := filepath.Join(baseDir, string(filepath.Separator), WalletFile)
	_appInst.WPath = walletPath
	tun2Pipe.ByPassInst().Load(bypassIPs)
	return nil
}
func waitVpnStatus(signal chan struct{}) {
	<-signal
	if _appInst.vpnDelegate != nil {
		_appInst.vpnDelegate.VpnClosed()
		_appInst.vpnDelegate = nil
	}

	if VpnInst != nil {
		VpnInst.Stop()
	}
	if TunInst != nil {
		TunInst.Finish()
	}
}

func StartVPN(srvAddr, minerAddr string, d VpnDelegate) error {
	fmt.Println("-----------start vpn service---------->\n", srvAddr)
	fmt.Println(minerAddr, "\n--------------------------------------<")
	_appInst.vpnDelegate = d

	minerHost := fmt.Sprintf("%s:%d", minerIP(minerAddr), MinerPort(minerAddr))
	if len(minerHost) == 0 {
		return fmt.Errorf("no miner host for addr:%s", minerAddr)
	}
	if !_appInst.wallet.IsOpen() {
		return fmt.Errorf("open the wallet first please")
	}

	subAddr := _appInst.wallet.SubAddress()
	var aesKey account.PipeCryptKey
	minderPub := account.ID(minerAddr).ToPubKey()
	if err := account.GenerateAesKey(&aesKey, minderPub, _appInst.wallet.CryptKey()); err != nil {
		return err
	}
	f := func(fd uintptr) {
		d.ByPass(int32(fd))
	}

	t2s, err := tun2Pipe.New(srvAddr, f, _appInst)
	if err != nil {
		return err
	}
	TunInst = t2s

	conf := &proxy.VpnConfig{
		Saver:        f,
		LocalApi:     t2s,
		MinerHost:    minerHost,
		SelfAddr:     subAddr,
		LocalSrvAddr: srvAddr,
		AesKey:       aesKey[:],
	}
	srv, err := proxy.NewProxyService(conf)
	if err != nil {
		return err
	}
	VpnInst = srv
	signal := make(chan struct{}, 2)
	go srv.ServingThread(signal)
	go t2s.Proxying(signal)
	go waitVpnStatus(signal)
	return nil
}

func StopVpn() {

	if VpnInst != nil {
		fmt.Println("User Stop the VPN service")
		VpnInst.Stop()
		VpnInst = nil
	}

	if TunInst != nil {
		fmt.Println("User Stop Tun2Proxy Service")
		TunInst.Finish()
		TunInst = nil
	}

	_appInst.vpnDelegate = nil
}

func InputPacket(data []byte) error {
	if TunInst == nil {
		return fmt.Errorf("Tun2Proxy has stopped")
	}
	TunInst.InputPacket(data)
	return nil
}

func SetGlobalModel(g bool) {
	tun2Pipe.ByPassInst().SetGlobal(g)
}

func IsGlobalMode() bool {
	return tun2Pipe.ByPassInst().IsGlobal()
}

type PingResult struct {
	IP   string  `json:"ip"`
	Ping float32 `json:"ping"`
}

func minerIP(mid string) string {
	minerIP := _appInst.IpCache[mid]
	if len(minerIP) == 0 {
		minerIP = RefreshHostByAddr(mid)
		fmt.Println("======>>>find miner ip:=>", mid, minerIP)
	}
	return minerIP
}

func TestPing(mid string) []byte {

	minerIP := minerIP(mid)
	if minerIP == "" {
		return nil
	}

	mAddr := &net.UDPAddr{
		IP:   net.ParseIP(minerIP),
		Port: int(account.ID(mid).ToServerPort()),
	}
	timeOut := time.Second * 5

	fmt.Println("=====>start to ping:", minerIP)
	conn, err := net.DialTimeout("udp4", mAddr.String(), timeOut)
	if err != nil {
		fmt.Println("=====>dial miner err:", err)
		return nil
	}
	now := time.Now()
	defer conn.Close()
	testConn := network.JsonConn{Conn: conn}
	_ = testConn.SetDeadline(now.Add(timeOut))
	err = testConn.WriteJsonMsg(node.CtrlMsg{Typ: node.MsgPingTest, PT: &node.PingTest{
		PayLoad: mid,
	}})
	if err != nil {
		fmt.Println("=====>WriteJsonMsg err:", err)
		return nil
	}
	err = testConn.ReadJsonMsg(&node.MsgAck{})
	if err != nil {
		fmt.Println("=====>ReadJsonMsg err:", err)
		return nil
	}

	result := PingResult{
		IP:   minerIP,
		Ping: float32(time.Now().Sub(now)) / float32(time.Millisecond),
	}
	fmt.Println("=====>finish to ping:", minerIP)
	data, _ := json.Marshal(result)
	return data
}

func MinerPort(addr string) int32 {
	mid := account.ID(addr)
	return int32(mid.ToServerPort())
}

func AndroidApkVersion() (ver string, err error) {
	ver = ""
	err = fmt.Errorf("not found")
	resp, err := http.Get("https://redeslab.github.io/version.js")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code is[%d]", resp.StatusCode)
		return
	}
	ver = string(body)
	return ver, nil
}
