package apns

import (
	"crypto/tls"
	"net"
	"strings"
)

const (
	ProductionGateway = "gateway.push.apple.com:2195"
	SandboxGateway    = "gateway.sandbox.push.apple.com:2195"

	ProductionFeedbackGateway = "feedback.push.apple.com:2196"
	SandboxFeedbackGateway    = "feedback.sandbox.push.apple.com:2196"
)

type Conn struct {
	NetConn   net.Conn
	Conf      *tls.Config
	gateway   string
	connected bool
}

func NewConnWithCert(gw string, cert tls.Certificate) Conn {
	gatewayParts := strings.Split(gw, ":")
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   gatewayParts[0],
	}
	return Conn{gateway: gw, Conf: &conf}
}

func NewConn(gw string, crt string, key string) (Conn, error) {
	cert, err := tls.X509KeyPair([]byte(crt), []byte(key))
	if err != nil {
		return Conn{}, err
	}
	return NewConnWithCert(gw, cert), nil
}

func NewConnWithFiles(gw string, certFile string, keyFile string) (Conn, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return Conn{}, err
	}
	return NewConnWithCert(gw, cert), nil
}

func (c *Conn) Connect() error {
	if c.NetConn != nil {
		c.NetConn.Close()
	}

	conn, err := net.Dial("tcp", c.gateway)

	if err != nil {
		return err
	}

	tlsConn := tls.Client(conn, c.Conf)
	err = tlsConn.Handshake()

	if err != nil {
		return err
	}

	c.NetConn = tlsConn
	return nil
}

func (c *Conn) Close() error {
	if c.NetConn != nil {
		return c.NetConn.Close()
	}
	return nil
}

func (c *Conn) Read(p []byte) (int, error) {
	i, err := c.NetConn.Read(p)
	return i, err
}

func (c *Conn) Write(p []byte) (int, error) {
	return c.NetConn.Write(p)
}
