package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/Jud1k/web_server/internal/request"
	"github.com/Jud1k/web_server/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	statusCode response.StatusCode
	message    string
}

func (e *HandlerError)Write(w io.Writer)error{
	_,err:=fmt.Fprint(w,e)
	if err!=nil{
		return err
	}
	return nil
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener, handler: handler}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if s.closed.Load() {
			return
		}
		if err != nil {
			if s.closed.Load(){
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req,err := request.RequestFromReader(conn)
	if err!=nil{
		hErr := &HandlerError{
			statusCode: response.StatusCodeBadRequest,
			message: err.Error(),
		}
		hErr.Write(conn)
	}
	writer := response.NewWriter()
	s.handler(writer,req)
	writer.WriteTo(conn)
}

