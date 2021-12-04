package data

type ProxyData struct {
	code   string
	online string
}

func Get() string {

	chaDirect := make(chan string)
	chaProxy := make(chan ProxyData)

	proxyData := ProxyData{"", ""}

	go Direct(chaDirect)
	go Proxy(chaProxy, proxyData)

	direct := <-chaDirect
	proxyData = <-chaProxy

	nodes := ""

	if direct != "" {
		nodes = direct
	} else if proxyData.code != "" {
		nodes = proxyData.code
	} else if proxyData.online != "" {
		nodes = Body(proxyData.online)
	}

	return nodes
}
