package tun2Pipe

import (
	"fmt"
	"io"
	"net"
	"syscall"
)

const bufferSize = 1 << 20

func (t2s *Tun2Pipe) simpleForward(conn net.Conn) {

	leftConn := conn.(*net.TCPConn)
	_ = leftConn.SetKeepAlive(true)

	port := leftConn.RemoteAddr().(*net.TCPAddr).Port
	s := t2s.getSession(port)
	if s == nil {
		fmt.Println("------>>>can't proxy this one:", leftConn.RemoteAddr())
		return
	}

	defer t2s.removeSession(port)

	fmt.Println("------>>>tun New conn for tcp session:", s.ToString())

	tgtAddr := fmt.Sprintf("%s:%d", s.RemoteIP, s.RemotePort)
	buff := make([]byte, bufferSize)

	for {
		n, e := leftConn.Read(buff)
		if e != nil {
			fmt.Println("------>>>Tun Read from left conn err:", e, n)
			if e != io.EOF {
				leftConn.Close()
			}
			return
		}
		//fmt.Println("------>>>Tun Read from left ", n)
		if s.Pipe == nil {
			d := &net.Dialer{
				Timeout: SysDialTimeOut,
				Control: func(network, address string, c syscall.RawConn) error {
					return c.Control(_config.protector)
				},
			}
			c, e := d.Dial("tcp", tgtAddr)
			if e != nil {
				fmt.Println("------>>>Dial remote err:", e)
				return
			}
			rightConn := c.(*net.TCPConn)
			_ = rightConn.SetKeepAlive(true)
			fmt.Printf("------>>>TCP:Tun Pipe dial success: %s->%s->%s:\n", rightConn.LocalAddr(), tgtAddr, s.ToString())

			s.Pipe = &directPipe{
				Left:  leftConn,
				Right: rightConn,
			}
			go s.Pipe.readingIn()
		}

		if err := s.Pipe.writeOut(buff[:n]); err != nil {
			return
		}
	}
}
