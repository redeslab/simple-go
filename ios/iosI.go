package iosLib

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/redeslab/go-simple/account"
	"github.com/redeslab/go-simple/network"
	"github.com/redeslab/go-simple/node"
	"net"
	"time"
)

const TXCheckTime = time.Second * 2

type UICallBack interface {
	Log(str string)
	Notify(note string, data string)
	SysExit(err error)
}

type AppHelper struct {
	callback UICallBack
	timer    *time.Timer
	Wallet   account.Wallet
	ipCache  map[string]string
}

var appCaller = &AppHelper{
	timer:   time.NewTimer(TXCheckTime),
	ipCache: make(map[string]string),
}

func InitSystem(cb UICallBack) {
	appCaller.callback = cb
}

type PingResult struct {
	IP   string
	Ping float32
}

func minerIP(mid string) string {
	minerIP := appCaller.ipCache[mid]
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
