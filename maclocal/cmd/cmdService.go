package command

import (
	"context"
	"fmt"
	"github.com/redeslab/go-lib/maclocal/pbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type cmdService struct{}

var StartServing func(pooladdr, minerid string) error
var StopService func()

func StartCmdService() {
	var address = "127.0.0.1:10022"
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("cmd listen error")
		return
	}

	cmdServer := grpc.NewServer()

	pbs.RegisterCmdServiceServer(cmdServer, &cmdService{})

	reflection.Register(cmdServer)
	if err := cmdServer.Serve(l); err != nil {
		fmt.Println("failed to serve command :%v", err)
		return
	}
}

func DialToCmdService() pbs.CmdServiceClient {

	var address = "127.0.0.1:10022"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("cmd dial to server err:->", err)
		return nil
	}

	client := pbs.NewCmdServiceClient(conn)

	return client
}

func (cs *cmdService) ShowWallet(context.Context, *pbs.EmptyRequest) (*pbs.CommonResponse, error) {
	return &pbs.CommonResponse{
		Msg: "",
	}, nil

}

func (cs *cmdService) StartVpn(ctx context.Context, r *pbs.VpnSrvParam) (*pbs.CommonResponse, error) {
	msg := "success"
	err := StartServing(r.Pool, r.Miner)
	if err != nil {
		msg = err.Error()
	}

	return &pbs.CommonResponse{
		Msg: msg,
	}, nil
}

func (cs *cmdService) StopVpn(context.Context, *pbs.EmptyRequest) (*pbs.CommonResponse, error) {
	StopService()

	return &pbs.CommonResponse{Msg: "success"}, nil
}
