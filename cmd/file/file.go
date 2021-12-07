package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/laof/request"
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

	str := request.Request()

	if str == "" {
		fmt.Println("please check your network")
		return false
	}

	nodes := strings.Split(str, "\n")

	port := "1080"

	var conf = []string{
		"# " + (time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")),
		"verbose=false",
		"listen=:" + port,
		"strategy=ha",
	}
	var forwards []string
	var ssr, ss int
	for _, v := range nodes {

		var curn node.Node

		br := strings.HasPrefix(v, "ssr://")
		bs := strings.HasPrefix(v, "ss://")

		if br {
			curn = node.SSRNode{}
		} else if bs {
			curn = node.SSNode{}
		} else {
			continue
		}

		nw := strings.TrimSpace(v)

		ser, err := curn.Deconde(nw)
		if err != nil {
			continue
		}

		if br {
			ssr++
		} else {
			ss++
		}

		link := curn.Create(ser)
		forwards = append(forwards, "forward="+link)
	}

	if len(forwards) == 0 {
		fmt.Println("empty node")
		return false
	}

	fmt.Printf("configure [%d] nodes (ssr:%d ss:%d)\n", len(forwards), ssr, ss)

	conf = append(conf, forwards...)

	txt := strings.Join(conf, "\n")

	return write(txt)

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
