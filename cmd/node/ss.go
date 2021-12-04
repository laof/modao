package node

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type SSNode struct {
	Node
}

func (ss SSNode) Deconde(str string) (Server, error) {

	if strings.HasPrefix(str, "ss://") {
		str = strings.Replace(str, "ss://", "", 1)
	} else {
		return Server{}, errors.New("invalid SS link")
	}

	sp := strings.Split(str, "#")

	str = sp[0]

	str = strings.Replace(str, "==", "", 1)

	remark := ""

	if len(sp) == 2 {
		remark, _ = url.QueryUnescape(sp[1])
	}

	data, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return Server{}, errors.New("Error decoding SS")
	}

	arr := strings.Split(string(data), ":")
	remote := strings.Split(string(arr[1]), "@")

	return Server{
		Addr:     remote[1],
		Port:     arr[2],
		Password: remote[0],
		Method:   arr[0],
		Remarks:  remark,
	}, nil

}

// ss://method:pass@1.1.1.1:8443

func (ss SSNode) Create(s Server) string {

	arr := []string{
		"ss://",
		s.Method, ":", s.Password,
		"@", s.Addr, ":", s.Port,
	}

	return strings.Join(arr, "")

}
