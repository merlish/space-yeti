package main

import "endlib"
import "endlib/server"

func main() {
    srv := server.New()
	if !dispatch.Listen(srv) {
		panic("failed to start server: dispatch.listen call failed.")
	}
}
