package ethdial

// strip the '0x' from a hex string
// TODO I think the go-ethereum has a function to do this - use it if so
func unHex(a string) string {
	if len(a) > 2 {
		if a[0] == '0' && (a[1] == 'x' || a[1] == 'X') {
			a = a[2:]
		}
	}
	return a
}

// add the '0x' to a hex string
// TODO I think the go-ethereum has a function to do this - use it if so
func reHex(a string) string {
	return "0x" + unHex(a)
}

// pad the hex string with leading '0's
// TODO I think the go-ethereum has a function to do this - use it if so
func leftPadHex(a string, b int) string {
	for {
		if len(a) >= b {
			break
		}
		a = "0" + a
	}
	return a
}
