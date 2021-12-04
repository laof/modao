package ss

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/nadoo/glider/log"
	"github.com/nadoo/glider/pool"
	"github.com/nadoo/glider/proxy"
	"github.com/nadoo/glider/proxy/protocol/socks"
)

var nm sync.Map

// NewSSServer returns a ss proxy server.
func NewSSServer(s string, p proxy.Proxy) (proxy.Server, error) {
	return NewSS(s, nil, p)
}

// ListenAndServe serves ss requests.
func (s *SS) ListenAndServe() {
	go s.ListenAndServeUDP()
	s.ListenAndServeTCP()
}

// ListenAndServeTCP serves tcp ss requests.
func (s *SS) ListenAndServeTCP() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.F("[ss] failed to listen on %s: %v", s.addr, err)
		return
	}

	log.F("[ss] listening TCP on %s", s.addr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.F("[ss] failed to accept: %v", err)
			continue
		}
		go s.Serve(c)
	}

}

// Serve serves a connection.
func (s *SS) Serve(c net.Conn) {
	defer c.Close()

	if c, ok := c.(*net.TCPConn); ok {
		c.SetKeepAlive(true)
	}

	sc := s.StreamConn(c)

	tgt, err := socks.ReadAddr(sc)
	if err != nil {
		log.F("[ss] failed to get target address: %v", err)
		proxy.Copy(io.Discard, c) // https://github.com/nadoo/glider/issues/180
		return
	}

	network := "tcp"
	dialer := s.proxy.NextDialer(tgt.String())

	rc, err := dialer.Dial(network, tgt.String())
	if err != nil {
		log.F("[ss] %s <-> %s via %s, error in dial: %v", c.RemoteAddr(), tgt, dialer.Addr(), err)
		return
	}
	defer rc.Close()

	log.F("[ss] %s <-> %s via %s", c.RemoteAddr(), tgt, dialer.Addr())

	if err = proxy.Relay(sc, rc); err != nil {
		log.F("[ss] %s <-> %s via %s, relay error: %v", c.RemoteAddr(), tgt, dialer.Addr(), err)
		// record remote conn failure only
		if !strings.Contains(err.Error(), s.addr) {
			s.proxy.Record(dialer, false)
		}
	}
}

// ListenAndServeUDP serves udp requests.
func (s *SS) ListenAndServeUDP() {
	lc, err := net.ListenPacket("udp", s.addr)
	if err != nil {
		log.F("[ss] failed to listen on UDP %s: %v", s.addr, err)
		return
	}
	defer lc.Close()

	log.F("[ss] listening UDP on %s", s.addr)

	lc = s.PacketConn(lc)
	for {
		c := NewPktConn(lc, nil, nil, true)
		buf := pool.GetBuffer(proxy.UDPBufSize)

		n, srcAddr, err := c.ReadFrom(buf)
		if err != nil {
			log.F("[ssu] remote read error: %v", err)
			continue
		}

		var session *Session
		sessionKey := srcAddr.String()

		v, ok := nm.Load(sessionKey)
		if !ok || v == nil {
			session = newSession(sessionKey, srcAddr, c)
			nm.Store(sessionKey, session)
			go s.serveSession(session)
		} else {
			session = v.(*Session)
		}

		session.msgCh <- buf[:n]
	}
}

func (s *SS) serveSession(session *Session) {
	dstC, dialer, writeTo, err := s.proxy.DialUDP("udp", session.srcPC.tgtAddr.String())
	if err != nil {
		log.F("[ssu] remote dial error: %v", err)
		return
	}
	dstPC := NewPktConn(dstC, writeTo, nil, false)
	defer dstPC.Close()

	go func() {
		proxy.RelayUDP(session.srcPC, session.src, dstPC, 2*time.Minute)
		nm.Delete(session.key)
		close(session.finCh)
	}()

	fmt.Printf("[ssu] %s <-> %s via %s\n", session.src, session.srcPC.tgtAddr, dialer.Addr())

	for {
		select {
		case p := <-session.msgCh:
			_, err = dstPC.WriteTo(p, writeTo)
			if err != nil {
				log.F("[ssu] writeTo %s error: %v", writeTo, err)
			}
			pool.PutBuffer(p)
		case <-session.finCh:
			return
		}
	}
}

// Session is a udp session
type Session struct {
	key   string
	src   net.Addr
	srcPC *PktConn
	msgCh chan []byte
	finCh chan struct{}
}

func newSession(key string, src net.Addr, srcPC *PktConn) *Session {
	return &Session{key, src, srcPC, make(chan []byte, 32), make(chan struct{})}
}
