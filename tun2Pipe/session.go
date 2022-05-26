package tun2Pipe

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"io"
	"net"
	"time"
)

type Session struct {
	byPass     bool
	Pipe       *directPipe
	UPTime     time.Time
	RemoteIP   net.IP
	RemotePort int
	ServerPort int
	byteSent   int
	packetSent int
	HostName   string
}

func (s *Session) ToString() string {
	return fmt.Sprintf("(srvPort=%d, bypass=%t) %s:%d t=%s", s.ServerPort, s.byPass,
		s.RemoteIP, s.RemotePort,
		s.UPTime.Format("2006-01-02 15:04:05"))
}

func newSession(ip4 *layers.IPv4, tcp *layers.TCP, srvPort int, bp bool) *Session {
	s := &Session{
		UPTime:     time.Now(),
		RemoteIP:   ip4.DstIP,
		RemotePort: int(tcp.DstPort),
		ServerPort: srvPort,
		byPass:     bp,
	}
	return s
}

type directPipe struct {
	Left  *net.TCPConn
	Right *net.TCPConn
}

func (dp *directPipe) readingIn() {
	defer dp.Left.Close()
	defer dp.Right.Close()
	if _, err := io.Copy(dp.Left, dp.Right); err != nil {
		fmt.Println("------>>>Tun Proxy pipe right 2 left finished:", err)
	}
}

func (dp *directPipe) writeOut(buf []byte) error {
	_, e := dp.Right.Write(buf)
	if e != nil {
		fmt.Println("------>>>tun Proxy pipe left 2 right err:", e)
		return e
	}
	//fmt.Println("------>>>tun Proxy pipe left 2 right", n)
	return nil
}
