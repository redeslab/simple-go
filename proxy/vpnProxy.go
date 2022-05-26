package proxy

import (
	"fmt"
	"github.com/redeslab/go-simple/account"
	"github.com/redeslab/go-simple/network"
	"github.com/redeslab/go-simple/node"
	"github.com/redeslab/go-simple/util"
	"net"
	"time"
)

type LocalProxyApi interface {
	GetTarget(conn net.Conn) string
	ProxyClose(conn net.Conn)
}

const (
	bufferSize   = 1 << 20
	MinerTimeOut = time.Second * 5
)

type VpnProxy struct {
	conn net.Listener
	conf *VpnConfig
}

type VpnConfig struct {
	Saver        util.ConnSaver
	LocalApi     LocalProxyApi
	MinerHost    string
	SelfAddr     account.ID
	LocalSrvAddr string
	AesKey       []byte
}

func NewProxyService(conf *VpnConfig) (*VpnProxy, error) {

	c, e := net.Listen("tcp", conf.LocalSrvAddr)
	if e != nil {
		return nil, e
	}

	vp := &VpnProxy{
		conn: c,
		conf: conf,
	}
	fmt.Println("VPN service listen at:", c.Addr().String())
	return vp, nil
}

func (vp *VpnProxy) ServingThread(sig chan struct{}) {

	defer vp.conn.Close()
	for {
		c, e := vp.conn.Accept()
		if e != nil {
			fmt.Println("\n======>>>sVPN Service Thread exit", e)
			break
		}
		go vp.newPipeTask(c)
	}

	sig <- struct{}{}
}

func (vp *VpnProxy) setupCryptConn(target string) (net.Conn, error) {
	conn, err := util.GetSavedConn(vp.conf.MinerHost, vp.conf.Saver, MinerTimeOut)
	if err != nil {
		return nil, err
	}
	_ = conn.(*net.TCPConn).SetKeepAlive(true)
	lvConn := network.NewLVConn(conn)

	iv := network.NewSalt()
	req := &node.SetupReq{
		IV:      *iv,
		SubAddr: vp.conf.SelfAddr,
	}
	jsonConn := &network.JsonConn{Conn: lvConn}
	if err := jsonConn.Syn(req); err != nil {
		return nil, err
	}
	aesConn, err := network.NewAesConn(lvConn, vp.conf.AesKey, *iv)
	if err != nil {
		return nil, err
	}
	jsonConn = &network.JsonConn{Conn: aesConn}
	if err := jsonConn.Syn(&node.ProbeReq{
		Target: target,
	}); err != nil {
		return nil, err
	}

	return aesConn, nil
}

func (vp *VpnProxy) forwardDataToProxy(lConn, targetConn net.Conn) {
	buffer := make([]byte, bufferSize)
	for {
		no, err := lConn.Read(buffer)
		if no == 0 {
			fmt.Println("\n======>>read from app failed:", err, no,
				lConn.LocalAddr().String(), lConn.RemoteAddr().String())
			return
		}
		_, err = targetConn.Write(buffer[:no])
		if err != nil {
			fmt.Println("\n======>>write to remote miner err:", err, no)
			return
		}
	}
}

func (vp *VpnProxy) newPipeTask(lConn net.Conn) {
	_ = lConn.(*net.TCPConn).SetKeepAlive(true)
	defer lConn.Close()
	defer vp.conf.LocalApi.ProxyClose(lConn)
	defer lConn.SetDeadline(time.Now().Add(MinerTimeOut))

	tgtHost := vp.conf.LocalApi.GetTarget(lConn)
	if len(tgtHost) < 2 {
		fmt.Printf("\n ======>>Invalid target[%s]:", tgtHost)
		return
	}
	fmt.Printf("\n======>>Request for[%s]\n", tgtHost)
	targetConn, err := vp.setupCryptConn(tgtHost)
	if err != nil {
		fmt.Printf("\n======>>Create connection to miner for [%s] err:%s", tgtHost, err)
		return
	}
	defer targetConn.Close()

	go vp.forwardDataToProxy(lConn, targetConn)
	buffer := make([]byte, bufferSize)
	for {
		no, err := targetConn.Read(buffer)
		if no == 0 {
			fmt.Println("======>> remote read finished", err)
			_ = targetConn.SetDeadline(time.Now().Add(time.Second * 5))
			return
		}
		_, err = lConn.Write(buffer[:no])
		if err != nil {
			fmt.Println("======>> write to local proxy err:=>", err, no)
		}
	}
}

func (vp *VpnProxy) Stop() {

	if vp.conn != nil {
		_ = vp.conn.Close()
		vp.conn = nil
	}
}
