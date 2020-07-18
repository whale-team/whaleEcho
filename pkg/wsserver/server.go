package wsserver

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/mailru/easygo/netpoll"
	"github.com/rs/zerolog/log"
)

func New() *SocketServer {
	poller, _ := netpoll.New(nil)
	return &SocketServer{
		poller:     poller,
		wg:         &sync.WaitGroup{},
		quitSingal: make(chan struct{}),
	}
}

type SocketServer struct {
	poller              netpoll.Poller
	Addr                string
	Port                string
	ln                  net.Listener
	Handler             HandleFunc
	ErrHandler          ErrHandleFunc
	ConnCloseHandler    ConnCloseHandleFunc
	ConnBuildHandleFunc ConnBuildHandleFunc
	handleConn          connHandleFunc
	wg                  *sync.WaitGroup
	quitSingal          chan struct{}
}

func (serv *SocketServer) resolvedAddr() string {
	return serv.Addr + ":" + serv.Port
}

func (serv *SocketServer) ListenAndServe(addr, port string) error {
	if addr == "" {
		addr = "127.0.0.1"
	}
	if port == "" {
		port = "10333"
	}
	serv.Addr = addr
	serv.Port = port
	ln, err := net.Listen("tcp", serv.resolvedAddr())
	if err != nil {
		return err
	}
	serv.ln = &onceCloseListener{Listener: ln}
	log.Info().Msgf("ws: Listening and serving Websocket Server on %s\n", addr+":"+port)
	serv.Serve()
	return nil
}

func (serv *SocketServer) Serve() {
	for {
		conn, err := serv.ln.Accept()
		if err != nil {
			select {
			case <-serv.quitSingal:
				return
			default:
				log.Error().Err(err).Msg("ws: get client connection failed")
				continue
			}
		}
		serv.upgradeToWS(conn)
	}
}

func (serv *SocketServer) upgradeToWS(conn net.Conn) error {
	if _, err := ws.Upgrade(conn); err != nil {
		return err
	}
	onceConn := &onceCloseConn{Conn: conn}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Context{
		ID:        uuid.New().String(),
		Conn:      onceConn,
		Context:   ctx,
		ctxCancel: cancel,
	}
	serv.ConnBuildHandleFunc(c)
	serv.wg.Add(1)
	serv.registerNetpoll(c, onceConn.Conn)
	return nil
}

func (serv *SocketServer) registerNetpoll(c *Context, conn net.Conn) {
	desc := netpoll.Must(netpoll.HandleRead(conn))
	serv.poller.Start(desc, func(event netpoll.Event) {
		if event&netpoll.EventHup != 0 || event&netpoll.EventReadHup != 0 {
			serv.wg.Done()
			serv.poller.Stop(desc)
			serv.ConnCloseHandler(c)
			return
		}
		go func() {
			if err := c.read(); err != nil {
				if err := serv.handleCloseErr(err, c); err != nil {
					serv.ErrHandler(c, err)
				}
				return
			}
			if err := serv.Handler(c); err != nil {
				serv.ErrHandler(c, err)
			}
		}()
	})
}

func (serv *SocketServer) handleCloseErr(err error, c *Context) error {
	if _, ok := err.(wsutil.ClosedError); ok {
		return serv.ConnCloseHandler(c)
	}
	if err.Error() == io.EOF.Error() {
		return serv.ConnCloseHandler(c)
	}
	return err
}

func (serv *SocketServer) Shutdown() error {
	serv.wg.Wait()
	close(serv.quitSingal)
	serv.ln.Close()
	return nil
}

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (l *onceCloseListener) OnceClose() error {
	l.once.Do(l.close)
	return l.closeErr
}

func (l *onceCloseListener) close() {
	l.closeErr = l.Listener.Close()
}

type onceCloseConn struct {
	net.Conn
	once     sync.Once
	closeErr error
}

func (l *onceCloseConn) OnceClose() error {
	l.once.Do(l.close)
	return l.closeErr
}

func (l *onceCloseConn) close() {
	l.closeErr = l.Conn.Close()
}
