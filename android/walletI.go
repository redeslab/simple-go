package androidLib

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/redeslab/go-simple/account"
	"github.com/redeslab/go-simple/util"
	"os"
	"strconv"
)

const (
	OpenWalletSuccess = iota
	PasswordError
	WalletError
	WalletSaveError
	NoWallet
)

func OpenWallet(auth string) error {
	w, err := account.LoadWallet(_appInst.WPath)
	if err != nil {
		return err
	}
	err = w.Open(auth)
	if err != nil {
		return err
	}
	_appInst.wallet = w
	return nil
}

func CloseWallet() {
	if _appInst.wallet == nil {
		return
	}
	_appInst.wallet.Close()
}

func IsEmptyWallet() bool {
	if _appInst.wallet == nil {
		return true
	}
	return false
}

func IsWalletOpen() bool {
	if IsEmptyWallet() {
		return false
	}
	w := _appInst.wallet
	return w != nil && w.IsOpen()
}
func WalletJson() string {
	if IsEmptyWallet() {
		return ""
	}
	return _appInst.wallet.String()

}
func NewWallet(auth string) (string, error) {

	w, e := account.NewWallet(auth, true)
	if e != nil {
		return "", e
	}

	if e := w.SaveToPath(_appInst.WPath); e != nil {
		return "", e
	}
	_appInst.wallet = w
	return w.String(), nil
}

func ImportWalletPrivate(data, auth string) (string, error) {
	hresult := data[4:68]
	w, err := account.NewWalletFromPrivateBytes(auth, hresult)
	if err != nil {
		return "", err
	}

	if ok := common.FileExist(_appInst.WPath); ok {
		t := strconv.FormatInt(util.GetNowMsTime(), 10)
		newfile := _appInst.WPath + t

		if err := os.Rename(_appInst.WPath, newfile); err != nil {
			return "", err
		}
	}

	if e := w.SaveToPath(_appInst.WPath); e != nil {
		return "", e
	}
	_appInst.wallet = w
	return w.String(), nil
}

func ImportWallet(data, auth string) int32 {
	w, e := account.LoadWalletByData(data)
	if e != nil {
		return WalletError
	}

	if e := w.Open(auth); e != nil {
		return PasswordError
	}

	if e := w.SaveToPath(_appInst.WPath); e != nil {
		return WalletSaveError
	}
	_appInst.wallet = w
	return OpenWalletSuccess
}
