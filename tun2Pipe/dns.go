package tun2Pipe

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"io"
	"math"
	"net"
	"sync"
	"time"
)

type QueryState struct {
	QueryTime     time.Time
	ClientQueryID uint16
	ClientIP      net.IP
	ClientPort    layers.UDPPort
	RemoteIp      net.IP
	RemotePort    layers.UDPPort
}

type DnsProxy struct {
	sync.RWMutex
	poxyConn     *net.UDPConn
	VpnWriteBack io.WriteCloser
	cache        map[uint16]*QueryState
}

func NewDnsCache() (*DnsProxy, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: InnerPivotPort,
	})
	if err != nil {
		return nil, err
	}

	if err := ProtectConn(conn); err != nil {
		return nil, err
	}

	proxy := &DnsProxy{
		poxyConn: conn,
		cache:    make(map[uint16]*QueryState),
	}

	go proxy.DnsWaitResponse()

	return proxy, nil
}

func (c *DnsProxy) Put(qs *QueryState) {
	c.Lock()
	defer c.Unlock()
	c.cache[qs.ClientQueryID] = qs
}

func (c *DnsProxy) Get(id uint16) *QueryState {
	c.RLock()
	defer c.RUnlock()
	return c.cache[id]
}

func (c *DnsProxy) Pop(id uint16) *QueryState {
	c.Lock()
	defer c.Unlock()
	qs := c.cache[id]
	delete(c.cache, id)
	return qs
}

func (c *DnsProxy) sendOut(dns *layers.DNS, ip4 *layers.IPv4, udp *layers.UDP) {

	qs := &QueryState{
		QueryTime:     time.Now(),
		ClientQueryID: dns.ID,
		ClientIP:      ip4.SrcIP,
		ClientPort:    udp.SrcPort,
		RemotePort:    udp.DstPort,
		RemoteIp:      ip4.DstIP,
	}
	c.Put(qs)
	if _, e := c.poxyConn.WriteTo(dns.LayerContents(), &net.UDPAddr{
		IP:   ip4.DstIP,
		Port: int(udp.DstPort),
	}); e != nil {
		return
	}
}

func (c *DnsProxy) DnsWaitResponse() {
	buff := make([]byte, math.MaxInt16)
	defer c.poxyConn.Close()

	for {
		n, _, e := c.poxyConn.ReadFromUDP(buff)
		if e != nil {
			return
		}
		var decoded []gopacket.LayerType
		dns := &layers.DNS{}
		p := gopacket.NewDecodingLayerParser(layers.LayerTypeDNS, dns)
		if err := p.DecodeLayers(buff[:n], &decoded); err != nil {
			continue
		}

		qs := c.Get(dns.ID)
		if qs == nil {
			continue
		}

		data := WrapIPPacketForUdp(qs.ClientIP, qs.RemoteIp, int(qs.ClientPort), int(qs.RemotePort), buff[:n])
		if data == nil {
			continue
		}

		_, err := c.VpnWriteBack.Write(data)
		if err != nil {
			continue
		}
	}
}
