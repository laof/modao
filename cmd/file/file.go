package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	data "github.com/nadoo/glider/cmd/data"
	"github.com/nadoo/glider/cmd/node"
)

func getDirectory() string {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//将\替换成/
	return strings.Replace(dir, "\\", "/", -1)
}

func Create() bool {

	str := data.Get()

	if str == "" {
		fmt.Println("please check your network")
		return false
	}

	var nodes []string
	old := strings.Split(str, "\n\n")

	var ss, ssr int

	for _, v := range old {

		var ty string = ""
		if strings.HasPrefix(v, "ssr://") {
			ssr++
			ty = "SSR"

		} else if strings.HasPrefix(v, "ss://") {
			ss++
			ty = "SS"
		}

		if ty != "" {
			nw := strings.TrimSpace(v)
			nw = strings.Replace(nw, "\n", "", -1)
			nodes = append(nodes, nw)
			fmt.Printf("Node %d : %s\n", len(nodes), ty)

		}
	}

	if ss == ssr && ssr == 0 {
		fmt.Println("empty node")
		return false
	}

	var index = 1
	var input string

	for {

		fmt.Printf("you want connent node number:")
		fmt.Scanln(&input)

		n, e := strconv.Atoi(input)

		if e != nil || n < 0 || n > len(nodes) {

			fmt.Println("please input retry")
			continue
		} else {
			index = n - 1
			break
		}

	}

	var curn node.Node

	if strings.Contains(nodes[index], "ssr://") {
		curn = node.SSRNode{}
	} else {
		curn = node.SSNode{}
	}

	ser, err := curn.Deconde(nodes[index])
	if err != nil {
		fmt.Println("error node link:", index)
		fmt.Println(err)
		return false
	}

	txt := curn.Create(ser)

	port := "1080"
	text := []string{
		"# " + (time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")),
		"listen=socks5://:" + port,
		"forward=" + txt,
	}

	conf := strings.Join(text, "\n")

	if write(conf) {
		fmt.Println("listen socks5:" + port)
		return true
	}

	return false

}

func init() {
	create()
}

func create() string {
	glider := filepath.Join(getDirectory(), "glider.conf")

	finame := filepath.FromSlash(glider)
	_, er := os.Stat(finame)

	if er != nil {
		file, _ := os.Create(glider)

		temp := []string{
			"# Don't operate configuration files",
			"listen=socks5://:1080",
			"forward=ss://none:test.com@password:789",
		}

		file.WriteString(strings.Join(temp, "\n"))
	}

	return finame
}

func write(data string) bool {

	wr := os.WriteFile(create(), []byte(data), os.ModeDevice)

	return wr == nil
}
