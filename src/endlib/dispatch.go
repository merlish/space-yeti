package dispatch

import "fmt"
import "net"

import "endlib/client"
import "endlib/server"

//import "endlib/mnet"

// Per server, dispatches connections for incoming clients.
// (Not multihome-aware... yet!)

func Listen(srv *server.Server) (ok bool) {
	// TODO: take args for where to listen on

	// create listener & await connections. close at end of func.
	/*addr, err0 := net.ResolveTCPAddr("tcp", ":31337")
	  if err0 != nil {
	      panic("failed to resolve local tcp addr. fuck")
	  }*/

	//ln, err := net.ListenTCP("tcp", addr)
	ln, err := net.Listen("tcp", ":31337")

	if err != nil {
		fmt.Println("error listening: %s", err.Error())
		return false
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %s\n", err.Error())
			continue
		}

		// so, `conn' is the new connection.
		// for now, let's just send a kick 0xFF..
		//kick := mnet.Buffer([]byte{}).Add(mnet.UByte(0xff)).Add(mnet.String("this is a test!"))
		//conn.Write(kick)

		//if closeErr := conn.Close(); closeErr != nil {
		//fmt.Println("failed to close connection in kick test: %s", closeErr.Error())
		//}

		//fmt.Println("Served client with test kick.")

		//ipconn, ok := conn.(net.IPConn)
		//if !ok {
		//panic("type assertion failure")
		//}

		//tcconn := conn.(*net.TCPConn)
		tcconn := conn.(*net.TCPConn)
		fmt.Printf("dispatching goroutine for new connection from %v\n", tcconn.RemoteAddr())

		//go client.Handle(&ipconn)
		go client.Handle(tcconn, srv)

	}

	return true

}
