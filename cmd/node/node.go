package node

type Node interface {
	Deconde(string) (Server, error)
	Create(Server) string
}

type Server struct {
	Addr      string
	Port      string
	Password  string
	Method    string
	Protocol  string
	Obfs      string
	Obfsparam string
	Remarks   string
	Group     string
}
