package data

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

func Direct(cha chan string) {

	host, _ := Decode("aHR0cHM6Ly9naXRodWIuY29tL0FsdmluOTk5OS9uZXctcGFjL3dpa2kvc3MlRTUlODUlOEQlRTglQjQlQjklRTglQjQlQTYlRTUlOEYlQjc=")
	resp, err := http.Get(string(host))

	if err != nil {
		cha <- ""
		return
	}
	defer resp.Body.Close()
	body, ierr := ioutil.ReadAll(resp.Body)

	if ierr != nil {
		cha <- ""
		return
	}

	buf := string(body)

	reg := regexp.MustCompile(`</?.+?/?>`)

	old := reg.ReplaceAllStringFunc(buf, func(str string) string {
		return "\n" + str + "\n"
	})

	// ssr
	reg = regexp.MustCompile(`ssr://(.*?)(\s|[\r\n]$)`)

	txt := reg.FindAllString(old, -1)
	var all string
	for _, text := range txt {
		all += "\n" + text
	}

	// ss
	reg = regexp.MustCompile(`ss://(.*?)(\s|[\r\n]$)`)

	sstxt := reg.FindAllString(old, -1)
	for _, txt := range sstxt {
		all += "\n" + txt
	}

	cha <- all

}
