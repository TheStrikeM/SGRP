package protocol

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
)

type StrikeRequest struct {
	Body string
}
type StrikeResponse struct {
	Result string
}

type handlerFunction func(StrikeRequest) StrikeResponse
type handlerMap map[string]handlerFunction

type SGRP struct {
	log         *slog.Logger
	handlersMap handlerMap
	Port        int16
	TcpListener *net.Listener
}

func New(log *slog.Logger, port int16) *SGRP {
	return &SGRP{
		log:         log,
		Port:        port,
		handlersMap: make(handlerMap),
		TcpListener: nil,
	}
}

func (protoc *SGRP) AddRoute(addr string, callback handlerFunction) {
	_, ok := protoc.handlersMap[addr]
	if ok {
		return
	}
	protoc.handlersMap[addr] = callback
}

func (protoc *SGRP) MustRun() *sync.WaitGroup {
	var wg sync.WaitGroup

	log := slog.With(
		slog.Int("port", int(protoc.Port)),
	)
	log.Info("Starting SGRP server ...")

	tcpListener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", protoc.Port))
	if err != nil {
		panic(err)
	}
	protoc.TcpListener = &tcpListener

	wg.Add(1)
	go func(tcpListener *net.Listener) {
		defer func() {
			(*tcpListener).Close()
			wg.Done()
		}()
		for {
			conn, err := (*tcpListener).Accept()
			if err != nil {
				log.Error(fmt.Sprintf("Error with message: %s", err.Error()))
				return
			}
			go protoc.requestHandler(&conn)
		}
	}(&tcpListener)

	return &wg
}

func (protoc *SGRP) requestHandler(conn *net.Conn) {
	who := (*conn).RemoteAddr().String()
	protoc.log.Info(fmt.Sprintf("User with %s success connected", who))

	message, err := bufio.NewReader(*conn).ReadString('\n')
	if err != nil {
		protoc.log.Error(fmt.Sprintf("Eror with message: %s", err.Error()))
	}

	_, addr, body := protoc.parse(message)
	var response StrikeResponse
	for key, callback := range protoc.handlersMap {
		if key == addr {
			response = callback(StrikeRequest{Body: body})
		}
	}

	_, err = (*conn).Write([]byte(fmt.Sprintf("%s %s %s", "RES", addr, response.Result)))
	if err != nil {
		protoc.log.Error(fmt.Sprintf("Eror with message: %s", err.Error()))
	}
}

func (protoc *SGRP) parse(message string) (string, string, string) {
	splitMessage := strings.Split(message, " ")
	return splitMessage[0], splitMessage[1], strings.Join(splitMessage[2:], "")
}

//type requestFunction[T any] func(request SGRPRequest[T]) SGRPResponse
//type requestMap map[SGRPRequest[any]]requestFunction[any]
//type SGRP struct {
//	log         *slog.Logger
//	port        int16
//	tcpListener *net.Listener
//	requests    requestMap
//}
//
//type SGRPRequest[T any] struct {
//	Addr string
//	body T
//}
//
//type SGRPResponse string
//
//func New(log *slog.Logger, port int16) *SGRP {
//	return &SGRP{
//		log:         log,
//		port:        port,
//		tcpListener: nil,
//		requests:    make(requestMap),
//	}
//}
//
//func (protocol *SGRP) Run() (err error) {
//	const op = "Protocol.Run"
//	defer func() { err = e.Wrap(op, err) }()
//
//	protocol.log.Info(
//		"Starting SGRP server...",
//		slog.Int("port", int(protocol.port)),
//	)
//	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", protocol.port))
//	if err != nil {
//		return err
//	}
//	protocol.tcpListener = &tcpListener
//
//	defer tcpListener.Close()
//	for {
//		conn, err := tcpListener.Accept()
//		if err != nil {
//			return err
//		}
//		go protocol.requestHandler(conn)
//	}
//}
//
//func (protocol *SGRP) requestHandler(conn net.Conn) {
//	who := conn.RemoteAddr().String()
//	protocol.log.Info(fmt.Sprintf("Welcome, %s", who))
//	message, err := bufio.NewReader(conn).ReadString('\n')
//	if err != nil {
//		protocol.log.Error(err.Error())
//	}
//	_, addr, body := protocol.parse(message)
//	var response SGRPResponse
//	for req, callback := range protocol.requests {
//		if req.Addr == addr {
//			response = callback(req)
//		}
//	}
//	_, err = conn.Write([]byte(fmt.Sprintf("%s %s %s", "RES", addr, response)))
//	if err != nil {
//		protocol.log.Error(err.Error())
//	}
//}
//
//func (protocol *SGRP) parse(message string) (string, string, string) {
//	splitMessage := strings.Split(message, " ")
//	return splitMessage[0], splitMessage[1], strings.Join(splitMessage[2:], "")
//}
//
//func (protocol *SGRP) Route(request SGRPRequest[any], callback requestFunction[any]) {
//	protocol.requests[request] = callback
//}
