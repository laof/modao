module github.com/nadoo/glider

go 1.17

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da
	github.com/dgryski/go-camellia v0.0.0-20191119043421-69a8a13fb23d
	github.com/dgryski/go-idea v0.0.0-20170306091226-d2fb45a411fb
	github.com/dgryski/go-rc2 v0.0.0-20150621095337-8a9021637152
	github.com/insomniacslk/dhcp v0.0.0-20210827173440-b95caade3eac
	github.com/laof/request v0.0.0-20211205064616-8f55a41107f5
	github.com/nadoo/conflag v0.2.3
	github.com/nadoo/ipset v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/xtaci/kcp-go/v5 v5.6.1
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	golang.org/x/sys v0.0.0-20211204120058-94396e421777
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/chromedp/cdproto v0.0.0-20211130100048-e5358edf9020 // indirect
	github.com/chromedp/chromedp v0.7.6 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ebfe/rc2 v0.0.0-20131011165748-24b9757f5521 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/klauspost/reedsolomon v1.9.13 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7 // indirect
	github.com/mdlayher/raw v0.0.0-20210412142147-51b895745faf // indirect
	github.com/templexxx/cpu v0.0.7 // indirect
	github.com/templexxx/xorsimd v0.4.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/u-root/uio v0.0.0-20210528151154-e40b768296a7 // indirect
	golang.org/x/net v0.0.0-20211205041911-012df41ee64c // indirect
)

// Replace dependency modules with local developing copy
// use `go list -m all` to confirm the final module used
// replace (
//	github.com/nadoo/conflag => ../conflag
// )
