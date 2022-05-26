package main

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {

	rootCmd.Flags().BoolVarP(&param.version, "version",
		"v", false, "PoolTool -v")

	rootCmd.Flags().StringVarP(&param.basIP, "bas",
		"b", "103.192.253.122", "PoolTool -b")

	rootCmd.Flags().StringVarP(&param.token, "token", "t",
		"0x72F391A5fC31b026739C8C26e0c5C01b2783F786", "PoolTool -t [token address]")

	rootCmd.Flags().StringVarP(&param.payment, "payment", "p",
		"0xb7b93d75690C4d1E8110D8D86b09Ff43BcA4335a", "PoolTool -p [payment address]")

}

const (
	ToolVersion = "1.0.2"
	infuraUrl   = "https://rinkeby.infura.io/v3/5ccfb3c5e4dc4117b91d08ceb1b22b67"
)

var (
	param struct {
		version bool
		basIP   string
		token   string
		payment string
	}

	rootCmd = &cobra.Command{
		Use: "Pool",

		Short: "minerPool is the miner pool logic service for microPayment system",

		Long: `usage description`,

		Run: mainRun,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type ReportItem struct {
	Address   string
	IP        string
	Name      string
	LastDay   float64
	LastMonth float64
	TotalUsed float64
}

func (ri *ReportItem) String() string {
	return fmt.Sprintf("[address=>%s,\tip=>%s,\tname=>%s,\tday=%0.3f,\tmonth=%0.3f,\ttotal=%0.3f]",
		ri.Address,
		ri.IP,
		ri.Name,
		ri.LastDay,
		ri.LastMonth,
		ri.TotalUsed)
}

type PoolReport []*ReportItem

func writeExcel(report PoolReport) {
	t := time.Now().Format("2006_01_02_15_04")
	path := fmt.Sprintf("pool_%s.csv", t)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	w.Comma = ','
	w.UseCRLF = true
	if err := w.Write([]string{"no", "address", "ip", "name", "lastDay", "lastMonth", "total"}); err != nil {
		panic(err)
	}

	for i, item := range report {
		no := fmt.Sprintf("%d", i)
		last := fmt.Sprintf("%0.3fM", item.LastDay)
		month := fmt.Sprintf("%0.3fG", item.LastMonth)
		total := fmt.Sprintf("%0.3fG", item.TotalUsed)
		if err := w.Write([]string{no, item.Address, item.IP, item.Name, last, month, total}); err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(item.String())
	}

	w.Flush()
}

func mainRun(_ *cobra.Command, _ []string) {

	if param.version {
		fmt.Println(ToolVersion)
		return
	}

}
