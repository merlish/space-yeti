package server

type BlockPosition struct {
	X int
	Y int
	Z int
}

type EntityPosition struct {
	X float64
	Y float64
	Z float64
}

type Server struct {
	Inventory chan *InventoryRequest // inventory system keeps track of general player data
	Map       chan MapRequest       // keeps track of world data, in columns
	Location  *PlayersServer
    Digput    chan DigputRequest
    Eids      chan int32
    Entities  chan interface{}
}

type PlayersServer struct {
	Subscribe   chan (chan interface{})
	Unsubscribe chan (chan interface{})
	Notify      chan interface{}
}

func New() *Server {
	invIn := make(chan *InventoryRequest)
	go inventoryServer(invIn)
	mapIn := make(chan MapRequest)
	go mapServer(mapIn)

	plIn := &PlayersServer{make(chan (chan interface{})), make(chan (chan interface{})), make(chan interface{})}
	go playersServer(plIn.Subscribe, plIn.Unsubscribe, plIn.Notify)

    dpIn := make(chan DigputRequest)
    go digputServer(dpIn)

    eids := make(chan int32)
    go eidServer(eids)

    ents := make(chan interface{})
    go entityManager(ents, plIn.Notify)

	return &Server{invIn, mapIn, plIn, dpIn, eids, ents}
}

