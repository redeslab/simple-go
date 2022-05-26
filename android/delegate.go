package androidLib

func (d *AndroidAPP) ServiceClosed() {
	StopVpn()
	d.vpnDelegate.VpnClosed()
}

func (d *AndroidAPP) Write(p []byte) (n int, err error) {
	return d.vpnDelegate.Write(p)
}
