package main

/*
#include "callback.h"
void bridge_func(UserInterfaceAPI f , int t, int t2, char* v){
	f(t, t2, v);
}
*/
import "C"
import (
	"bytes"
	"fmt"
)

func (app *MacApp) Log(a ...interface{}) {
	buf := bytes.NewBufferString("")
	_, _ = fmt.Fprint(buf, a...)
	C.bridge_func(app.uiAPI, C.ProtocolLog, 0, C.CString(buf.String()))
}
func (app *MacApp) ActionNotify(typ int, a ...interface{}) {
	buf := bytes.NewBufferString("")
	_, _ = fmt.Fprint(buf, a...)
	C.bridge_func(app.uiAPI, C.ProtocolNotification, C.int(typ), C.CString(buf.String()))
}
func (app *MacApp) SysExit(err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	stopService()
	C.bridge_func(app.uiAPI, C.ProtocolExit, 0, C.CString(errStr))
}

func (app *MacApp) ServiceClosed() {
	C.bridge_func(app.uiAPI, C.ServiceClosed, 0, nil)
}
