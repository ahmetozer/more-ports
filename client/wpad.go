package client

//serveProxyConfig Return wpad config for clients
func serveProxyConfig(addr string) string {
	// if (url.split(":")[0] == "http") {
	//    return "HTTP "+host
	// }

	return `function FindProxyForURL(url, host) {
		if (shExpMatch(url, "*://" + host + ":*")) {
			return "PROXY ` + addr + `";
		}
		return "DIRECT";
	}`
}
