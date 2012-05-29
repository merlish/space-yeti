package server

import "fmt"

const (
	InvLoadPlayer = iota
	InvGetCompassPos
	InvSetCompassPos
	InvGetSpawnPos
	InvSetSpawnPos
	InvUnloadPlayer
)

type InventoryRequest struct {
	Id     int
	Player string
	Data   interface{}
	Ret    chan interface{}
}

func inventoryServer(in chan *InventoryRequest) {
	for {
		req := <-in

		if req.Id == InvGetCompassPos {
			req.Ret <- &BlockPosition{0, 70, 0}
		} else if req.Id == InvGetSpawnPos {
			req.Ret <- &EntityPosition{0.0, 70.0, 0.0}
		} else {
			panic(fmt.Sprintf("that inventoryServer req.id (%v) is not implemented yet!", req.Id))
		}
	}
}

