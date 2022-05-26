package iosLib

import (
	"flag"
	"fmt"
	"testing"
)

var (
	userAddr = ""
)

func init() {
	flag.StringVar(&userAddr, "uid", "", "--uid")
}

//go test -run  TestServerList
func TestServerList(t *testing.T) {
	ret := SyncServerList()
	fmt.Println("list:", string(ret))
}

//go test -run  TestQueryByAddr --uid
func TestQueryByAddr(t *testing.T) {
	ret := RefreshHostByAddr(userAddr)
	fmt.Println("ip:", string(ret))
}
