package main

import (
	"context"
	"fmt"
	cmd "github.com/redeslab/go-lib/maclocal/cmd"
	"github.com/redeslab/go-lib/maclocal/pbs"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
)

var param struct {
	version bool
	pool    string
	token   string
	user    string
}

func getbasdir() string {
	home, _ := os.UserHomeDir()

	basdir := path.Join(home, ".testpirate")

	return basdir
}

var rootCmd = &cobra.Command{
	Use: "Pool",

	Short: "minerPool is the miner pool logic service for microPayment system",

	Long: `usage description`,

	Run: mainRun,
}

var walletShowCmd = &cobra.Command{
	Use: "wallet",

	Short: "show wallet",

	Long: `usage description`,

	Run: showWallet,
}

var startvpncmd = &cobra.Command{
	Use: "startvpn",

	Short: "start vpn",

	Long: `usage description`,

	Run: startvpn,
}

var stopvpncmd = &cobra.Command{
	Use: "stopvpn",

	Short: "stop vpn",

	Long: `usage description`,

	Run: stopvpn,
}

func init() {

	rootCmd.Flags().BoolVarP(&param.version, "version",
		"v", false, "Pool -v")

	rootCmd.AddCommand(walletShowCmd)
	rootCmd.AddCommand(startvpncmd)
	rootCmd.AddCommand(stopvpncmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startService(pooladdr, minerid string) error {
	return nil
}

func mainRun(_ *cobra.Command, _ []string) {
	if param.version {
		fmt.Println("1.0.0")
		return
	}

	sigCh := make(chan os.Signal, 1)

	waitSignal(sigCh)
}

func showWallet(_ *cobra.Command, _ []string) {
	c := cmd.DialToCmdService()
	md, e := c.ShowWallet(context.Background(), &pbs.EmptyRequest{})
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(md.Msg)
}

func startvpn(_ *cobra.Command, p []string) {

}

func stopvpn(_ *cobra.Command, _ []string) {
}

func waitSignal(sigCh chan os.Signal) {
	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>miner pool start at pid(%s)<<<<<<<<<<\n", pid)
	pidf := path.Join(getbasdir(), "pidfile")
	if err := ioutil.WriteFile(pidf, []byte(pid), 0644); err != nil {
		fmt.Print("failed to write running pid", err)
	}

	signal.Notify(sigCh,
		//syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sigCh
	//stopService()

	stopApp()

	fmt.Printf("\n>>>>>>>>>>process finished(%s)<<<<<<<<<<\n", sig)
}
