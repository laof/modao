package data

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type workExample struct {
	URL string `json:"url"`
}

type jsonTypeData struct {
	workExample `json:"workExample"`
}

func Proxy(cha chan ProxyData, pd ProxyData) {

	host, _ := Decode("aHR0cHM6Ly9jb2Rlc2FuZGJveC5pby9zL2dpdGh1Yi9sYW9mL3Nzc3NhbmRib3g=")
	res, fail := http.Get(string(host))

	if fail != nil {
		cha <- pd
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		cha <- pd
		return
	}

	scripts := doc.Find("script")
	if len(scripts.Nodes) <= 0 {
		cha <- pd
		return
	}

	var jsonStr string

	scripts.Each(func(i int, s *goquery.Selection) {

		val, ok := s.Attr("type")

		if !ok {
			return
		}

		if val == "application/ld+json" {
			jsonStr = s.Text()
		}

	})

	if jsonStr == "" {
		cha <- pd
		return
	}

	jsonData := jsonTypeData{}
	json.Unmarshal([]byte(jsonStr), &jsonData)

	vurl := jsonData.workExample.URL

	if strings.HasPrefix(vurl, "http") {
		pd.online = vurl
	} else {
		cha <- pd
		return
	}

	vr, ve := http.Get(vurl)

	if ve != nil {
		cha <- pd
		return
	}

	defer vr.Body.Close()

	vd, verr := goquery.NewDocumentFromReader(vr.Body)

	if verr != nil {
		cha <- pd
		return
	}

	vnode := vd.Find("html")

	if len(vnode.Nodes) == 0 {
		d, e := ioutil.ReadAll(vr.Body)

		if e == nil {
			pd.code = string(d)
		}

	}

	cha <- pd

}
