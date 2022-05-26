package iosLib

import (
	"encoding/json"
	"fmt"
	"github.com/redeslab/go-simple/contract/ethapi"
	"strings"
)

func SyncServerList() []byte {
	items := ethapi.SyncServerList()
	for _, it := range items {
		appCaller.ipCache[strings.ToLower(it.Addr)] = it.Host
	}
	bs, err := json.Marshal(items)
	if err != nil {
		return nil
	}
	return bs
}
func RefreshHostByAddr(addr string) string {
	newHost := ethapi.RefreshHostByAddr(addr)
	appCaller.ipCache[addr] = newHost
	return newHost
}

func AdvertiseList() []byte {
	items := ethapi.AdvertiseList("")
	if items == nil {
		return nil
	}

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
	fmt.Println("======>>>", string(bs))
	return bs
}
