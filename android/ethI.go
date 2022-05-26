package androidLib

import (
	"encoding/json"
	"fmt"
	"github.com/redeslab/go-simple/contract/ethapi"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

func directNetWork() {
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			Control: func(network, address string, c syscall.RawConn) error {
				if _appInst.vpnDelegate != nil {
					f := func(fd uintptr) {
						_appInst.vpnDelegate.ByPass(int32(fd))
					}
					return c.Control(f)
				}
				return nil

			},
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func SyncServerList() []byte {
	directNetWork()
	items := ethapi.SyncServerList()
	for _, it := range items {
		_appInst.IpCache[strings.ToLower(it.Addr)] = it.Host
	}
	bs, err := json.Marshal(items)
	if err != nil {
		return nil
	}
	return bs
}
func RefreshHostByAddr(addr string) string {
	directNetWork()
	newHost := ethapi.RefreshHostByAddr(addr)
	_appInst.IpCache[addr] = newHost
	return newHost
}

func AdvertiseList() []byte {
	items := ethapi.AdvertiseList("")
	if items == nil {
		return nil
	}
	directNetWork()

	result := make([]*ethapi.AdvertiseConfig, 0)
	for _, item := range items {
		adItem := &ethapi.AdvertiseConfig{}
		if err := json.Unmarshal([]byte(item.ConfigInJson), adItem); err != nil {
			fmt.Println("======>>>ad config json str err:=>", err)
			continue
		}

		result = append(result, adItem)
	}

	bs, _ := json.Marshal(result)
	return bs
}
