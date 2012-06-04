package main

import "fmt"

import "endlib"
import "endlib/server"

func main() {
    srv := server.New()

    fmt.Printf("\nyetisrv: the endsrv fork for space-yeti!\n")

	if !dispatch.Listen(srv) {
		panic("failed to start server: dispatch.listen call failed.")
	}
}
