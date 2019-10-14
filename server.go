package eio

import (
	event "github.com/swift9/ares-event"
	"net"
)

type Server struct {
	event.Emitter
	tcpListener *net.TCPListener
	Sockets     map[string]*Socket
	Addr        string
	Protocol    Protocol
	Log         ILog
}

func NewServer(addr string, protocol Protocol) *Server {
	server := &Server{
		Addr:     addr,
		Protocol: protocol,
		Sockets:  make(map[string]*Socket),
		Log:      &SysLog{},
	}
	return server
}

func (server *Server) SetLog(log ILog) {
	server.Log = log
}

func (server *Server) Listen(onConnect func(socket *Socket)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", server.Addr)
	if err != nil {
		server.Log.Error(err)
		return err
	}

	server.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		server.Log.Error(err)
		return err
	}

	for {
		conn, err := server.tcpListener.AcceptTCP()
		if err != nil {
			server.Log.Error(err)
			server.Emit("error", err)
			continue
		}
		socket := NewSocket(conn, server.Protocol)
		socket.SetLog(server.Log)
		server.Sockets[socket.Id] = socket
		go server.onConnect(socket, onConnect)
	}
	return nil
}

func (server *Server) onConnect(socket *Socket, callback func(socket *Socket)) {
	callback(socket)
}
