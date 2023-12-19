package readonlyconn

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

type ReadOnlyConn struct {
	reader io.Reader
}

func (conn ReadOnlyConn) Read(p []byte) (int, error)         { return conn.reader.Read(p) }
func (conn ReadOnlyConn) Write(p []byte) (int, error)        { return 0, io.ErrClosedPipe }
func (conn ReadOnlyConn) Close() error                       { return nil }
func (conn ReadOnlyConn) LocalAddr() net.Addr                { return nil }
func (conn ReadOnlyConn) RemoteAddr() net.Addr               { return nil }
func (conn ReadOnlyConn) SetDeadline(t time.Time) error      { return nil }
func (conn ReadOnlyConn) SetReadDeadline(t time.Time) error  { return nil }
func (conn ReadOnlyConn) SetWriteDeadline(t time.Time) error { return nil }

func ReadClientHello(reader io.Reader) (*tls.ClientHelloInfo, error) {
	var hello *tls.ClientHelloInfo

	err := tls.Server(ReadOnlyConn{reader: reader}, &tls.Config{
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			hello = &tls.ClientHelloInfo{}
			*hello = *argHello
			return nil, nil
		},
	}).Handshake()

	if hello == nil {
		return nil, err
	}

	return hello, nil
}

func PeekClientHello(reader io.Reader) (*tls.ClientHelloInfo, io.Reader, error) {
	peekedBytes := &bytes.Buffer{}
	hello, err := ReadClientHello(io.TeeReader(reader, peekedBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("can't read client hello: %w", err)
	}
	return hello, io.MultiReader(peekedBytes, reader), nil
}
