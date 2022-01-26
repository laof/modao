package node

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type SSRNode struct {
	Node
}

func param(s string) (obj map[string]string) {

	obj = make(map[string]string)
	obj["obfsparam"] = ""
	obj["remarks"] = ""
	obj["group"] = ""

	// obfsparam=&remarks=5rSb5p2J55-2U1NSK1RMUytDYWRkeQ&group=aHR0cHM6Ly9naXQuaW8vdjk5OTk
	u, err := url.Parse("http://a.org/s?" + s)
	if err != nil {
		return
	}
	m, _ := url.ParseQuery(u.RawQuery)

	obfsparam, remarks, group := m["obfsparam"], m["remarks"], m["group"]

	if len(obfsparam) > 0 {
		otxt, e := base64.RawURLEncoding.DecodeString(obfsparam[0])
		if e != nil {
			return
		}

		obj["obfsparam"] = string(otxt)
	}

	if len(remarks) > 0 {
		rtxt, e := base64.RawURLEncoding.DecodeString(remarks[0])
		if e != nil {
			return
		}
		obj["remarks"] = string(rtxt)
	}

	if len(group) > 0 {
		gtxt, e := base64.RawURLEncoding.DecodeString(group[0])
		if e != nil {
			return
		}
		obj["group"] = string(gtxt)
	}

	return

}

func (ssr SSRNode) Deconde(str string) (Server, error) {

	if strings.HasPrefix(str, "ssr://") {
		str = strings.Replace(str, "ssr://", "", 1)
	} else {
		return Server{}, errors.New("Error: invalid SSR link")
	}

	data, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return Server{}, errors.New("Error: decoding SSR")
	}

	arr := strings.Split(string(data), ":")
	pwremark := strings.Split(string(arr[5]), "/?")
	pw := pwremark[0]
	remarks := ""
	obfsparam := ""
	group := ""

	if len(pwremark) == 2 {
		obj := param(pwremark[1])
		remarks = obj["remarks"]
		obfsparam = obj["obfsparam"]
		group = obj["group"]
	}

	password, e := base64.RawURLEncoding.DecodeString(pw)
	if e != nil {
		return Server{}, errors.New("Error decoding password")
	}

	return Server{
		Addr:      arr[0],
		Port:      arr[1],
		Password:  string(password),
		Method:    arr[3],
		Protocol:  arr[2],
		Obfs:      arr[4],
		Obfsparam: obfsparam,
		Remarks:   remarks,
		Group:     group,
	}, nil

}

// ssr://method:pass@host:port?protocol=xxx&protocol_param=yyy&obfs=zzz&obfs_param=xyz

func (ssr SSRNode) Create(s Server) string {

	arr := []string{
		"ssr://",
		s.Method, ":", s.Password,
		"@", s.Addr, ":", s.Port,
		"?protocol=", s.Protocol,
		"&protocol_param=&", "obfs=", s.Obfs,
		"&obfs_param=", s.Obfsparam,
	}

	return strings.Join(arr, "")

}
