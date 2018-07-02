package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitghi/servo"
)

var (
	esig   chan os.Signal   // exit signal (SIGINT, SIGKILL)
	server *servo.Server    // server instance
	config servo.ServerOpts // server config
	err    error
	ok     bool
)

// connHandler is a deleate function which gets invoked
// when a new connection to server is established. It is
// executed in a non-blocking mode by the server.
func connHandler(conn net.Conn) interface{} {
	var (
		reader *bufio.Reader = bufio.NewReader(conn)
		buff   []byte
		err    error
	)
	defer func() {
		if conn != nil {
			conn.Close()
			conn = nil
		}
	}()
	log.Printf("* handling connection(%s) ....\n", conn.RemoteAddr().String())
	for {
		buff, _, err = reader.ReadLine()
		if err == io.EOF {
			log.Println("servo:", "error while reading packet:", err)
			return nil
		} else if err != nil {
			log.Println("servo:", "error while reading packet:", err)
		}
		fmt.Println("Packet(content): ", string(buff))
	}
	return nil
}

func main() {
	esig = make(chan os.Signal, 1)
	signal.Notify(esig, syscall.SIGINT, syscall.SIGKILL)
	// define server configuration
	config = servo.ServerOpts{
		OnNewConnection: connHandler, // connection delegate
		Addr:            ":80",       // server address
		StartDeadline:   2,           // maximum wait for server start
	}
	// create server instance
	// from configurations.
	server, err = servo.NewServerFromConfig(config)
	if err != nil {
		panic(err)
	}
	// start the server.
	err = server.Start()
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("+ [%s][%s] server started.", server.SID, config.Addr))
	// catch SIGINT or SIGKILL.
	<-esig
	// stop the server.
	ok, err = server.Stop()
	if !ok || err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("+ [%s] server quitted.", server.SID))
}
