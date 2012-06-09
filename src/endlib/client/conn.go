package client

import "bytes"
import "compress/zlib"
import "endlib/mnet"
import "fmt"
import "net"
import "endlib/server"
import "strings"
import "time"

// handles incoming client connections.

type Conn struct {
	c     *mnet.CConn
	s     *mnet.SConn
	raddr net.Addr
	name  string
    eid int32
}

func Handle(conn *net.TCPConn, srv *server.Server) {
	// casting to get access to sets of mnet functions
	c := &Conn{(*mnet.CConn)(conn), (*mnet.SConn)(conn), conn.RemoteAddr(), "", 0}
	defer conn.Close()               // then, close the underlying connection properly.
	defer tryTellClientAboutPanic(c) // before anything else, try and kick client.

	// ok. at this point in time, we've read nothing from the conn;
	//  so our first step should be to establish the intentions of the client.

	// TODO: set up read timeouts

	// are they querying us, or are they actually trying to join the game?
	handleFirstPacket(c)

	// ok. if we're here, they must be trying to join the game.

	// TODO: check no-one on the server has the same name already. (should only happen w/auth off)
	// TODO: check server isn't full
	// ^^^ must happen before we send login request... otherwise, mc client will hang on kick. :(

	// handle & respond to login request
    c.eid = <-srv.Eids
    defer func() {
        srv.Eids <- c.eid // return free eid
        //c.eid = 0  // actually, don't do this... other things need our eid for cleanup!
    }()

    fmt.Printf("dbg: in case you care, my eid is %d\n", c.eid)

	handleLogin(c, c.eid)

	fmt.Println("dbg: handled login, getting inv data")

	// figure out spawn position
	iinnvv := make(chan interface{})
	srv.Inventory <- &server.InventoryRequest{server.InvGetCompassPos, c.name, nil, iinnvv}
	cp := (<-iinnvv).(*server.BlockPosition)
	srv.Inventory <- &server.InventoryRequest{server.InvGetSpawnPos, c.name, nil, iinnvv}
	sp := (<-iinnvv).(*server.EntityPosition)

	fmt.Println("got inv data; getting map data")

	// s->c: send prechunks & some chunks
	// we'll do it for a 3x3 grid of the columns surrounding tyhe player
	prec := make(chan interface{}, 1)
	//srv.Map <- &server.MapRequest{server.MapGetColumns, -1, -1, 1, 1, nil, prec}
	srv.Map <- &server.MapGetColumnsRequest{0, 0, 0, 0, nil, prec}

	fmt.Println("got map data; doing prechunks")
	/*for z := -1; z < 2; z++ {
	    for x := -1; x < 2; x++ {
	        c.s.MapColumnAllocation(int32(x), int32(z), true)
	    }
	}*/
	c.s.MapColumnAllocation(int32(0), int32(0), true)

	fmt.Println("doing prechunks; doing column sending")
	//for i := 0; i < 9; i++ {
	i := 0
	fmt.Printf("%v", i)
	column := (<-prec).(*server.MapColumn)
	bs := make([]byte, 0)
	for r := range column.Chunks {
		chunk := column.Chunks[r]
		bs = append(bs, chunk.BlockTypes...)
	}
	for r := range column.Chunks {
		chunk := column.Chunks[r]
		bs = append(bs, chunk.BlockMetadata...)
	}
	for r := range column.Chunks {
		chunk := column.Chunks[r]
		bs = append(bs, chunk.BlockLight...)
	}
	for r := range column.Chunks {
		chunk := column.Chunks[r]
		bs = append(bs, chunk.SkyLight...)
	}
	bs = append(bs, make([]byte, 256)...)

	var b bytes.Buffer
	zlwr := zlib.NewWriter(&b)
	zlwr.Write(bs)
	zlwr.Close()

	fbs := b.Bytes()

	c.s.MapChunk(int32(column.X), int32(column.Z), true, 0xf, 0, fbs)
	//}

	fmt.Println("OK")

	fmt.Println("sending spawn pos")

	// s -> c: send spawn position
	c.s.SpawnPosition(int32(cp.X), int32(cp.Y), int32(cp.Z))

	// TODO: send inventory

	fmt.Println("doing PPAL")

	// s->c: let client spawn by sending position & look packet
	c.s.PlayerPositionAndLook(sp.X, sp.Y + 1.6, sp.Y, sp.Z, 0, 0, true)

	// TODO: subscribe to player movement changes

	pmp := server.PlayerMovePacket{c.eid, sp.X, sp.Y, sp.Z, 0, 0, true}

    // keep sending 'keep alive' packets to client
    go func() {
        defer func() { recover() }()

        for {
            c.s.KeepAlive(10231023)
            time.Sleep(1 * time.Second)
        }
    }()

    // otherPlayersHandle does its setup synchronously, then forks for its main loop.
    // this is a pretty cool pattern.
    opQuit := make(chan int)
    otherPlayersHandle(c, opQuit, srv)

    defer func() {
        fmt.Println("dbg: client/conn: telling client/o_p to quit")
        opQuit <- 0
    }()

    // tell the entity manager about our ent
    srv.Entities <- server.EntitiesCreatePlayer{c.eid, c.name}

    defer func() {
        fmt.Printf("in defer func() srv.ents <- delete{eid}\n")
        srv.Entities <- server.EntitiesDelete{c.eid}
    }()

    // terrible, terrible, terrible hack
    /*for i := 1; i < 16; i++ {
        if int32(i) != c.eid {
            c.s.SpawnNamedEntity(int32(i), fmt.Sprintf("Player%x", i), 0, 0, 0, 0, 0, 0)
            c.s.Entity(int32(i))

            // equipped weapon & armour
            for i2 := 0; i2 < 5; i2++ {
                c.s.EntityEquipment(int32(i), int16(i2), int16(-1), int16(0))
            }
        }
    }*/

    fmt.Println("client/conn: dispatching digging/map func")

    // subscribe to digging/map events
    go func() {
        defer func() { recover() }()

        digch := make(chan server.DigputRequest, 1)

        // subscribe to digput events!
        fmt.Println("client/digput: subscribing...")
        srv.Digput <- &server.DigputSubscribe{digch}

        defer func() {
            srv.Digput <- &server.DigputUnsubscribe{digch}

            // drain channel
            for _ = range digch { }
        }()

        // indicate we don't want to quit yet
        digch <- true

        for {
            select {
            case d, _ := <-digch:
                // TODO: implement more messages
                switch d.(type) {
                case *server.DigputFinishDig:
                    dfd := d.(*server.DigputFinishDig)
                    c.s.BlockChange(int32(dfd.X), byte(dfd.Y), int32(dfd.Z), 0, 0)
                }
            }
        }
    }()

    fmt.Println("client/conn: main loop")

	// the main loop :D
	for {
		// TODO: check for flooding. (could starve resources of other threads...)

        //fmt.Println("client/conn: read packet id")
		id := c.c.ReadID()

		switch id {
        case mnet.KeepAliveID:
            c.c.ReadKeepAlive() // cool. whatever

        case mnet.AnimationID:
            // TODO. Ignored for now.  (notchian clients send 1 Swing arm, apparently)
            c.c.ReadAnimation()

		case mnet.PlayerPositionAndLookID:
			ppal := c.c.ReadPlayerPositionAndLook()
			pmp.X = ppal.X
			pmp.Y = ppal.Y
			pmp.Z = ppal.Z
			pmp.Yaw = ppal.Yaw
			pmp.Pitch = ppal.Pitch
			pmp.OnGround = ppal.OnGround

            //fmt.Printf("yaw f: %v, pitch f: %v\n", ppal.Yaw, ppal.Pitch)

            //fmt.Printf("Debug: PPAL: %f,%f,%f\n", pmp.X, pmp.Y, pmp.Z)

			srv.Location.Notify <- pmp

		case mnet.PlayerPositionID:
			pp := c.c.ReadPlayerPosition()
			pmp.X = pp.X
			pmp.Y = pp.Y
			pmp.Z = pp.Z
			pmp.OnGround = pp.OnGround

            //fmt.Printf("Debug: PP: %f,%f,%f\n", pmp.X, pmp.Y, pmp.Z)

			srv.Location.Notify <- pmp

		case mnet.PlayerLookID:
			pl := c.c.ReadPlayerLook()
			pmp.Yaw = pl.Yaw
			pmp.Pitch = pl.Pitch
			pmp.OnGround = pl.OnGround

            //fmt.Printf("Debug: PL: %f,%f,%f\n", pmp.X, pmp.Y, pmp.Z)

			srv.Location.Notify <- pmp

        case mnet.PlayerID:
            pmp.OnGround = c.c.ReadPlayer()

            srv.Location.Notify <- pmp

        case mnet.PlayerDiggingID:
            //fmt.Println("reading player digging")
            pd := c.c.ReadPlayerDigging()

            ok := make(chan bool)
            bp := &server.BlockPosition{int(pd.X), int(pd.Y), int(pd.Z)}

            if pd.Status == 0 { // started digging
                //srv.Digput <- &server.DigputStartDig{bp, ok}
                //fmt.Println("srv.digput <- ...")
                srv.Digput <- &server.DigputFinishDig{bp, ok}
                //fmt.Println("<-ok")
                <-ok // tmp
                //fmt.Println(":D")
            } else if pd.Status == 2 { // finished digging
                //srv.Digput <- &server.DigputFinishDig{bp, ok}
                //<-ok // tmp
            } else if pd.Status == 4 { // drop item
                panic("digging::drop item not implemented yet!")
            } else if pd.Status == 5 { // shoot arrow/finish eating
                panic("digging::shoot arrow/finish eating not implemented yet!")
            } else {
                panic(fmt.Sprintf("digging::(%d) not implemented yet (what is it?)", pd.Status))
            }

            //fmt.Println("processed digput")

        case mnet.PluginMessageID:
            _ = c.c.ReadPluginMessage()

        case mnet.EntityActionID:
            ea := c.c.ReadEntityAction()

            if ea.EID != c.eid {
                panic("client: Entity Action (0x13) sent c->s _not_ referring to the player... weird!")
            }

            switch ea.ActionID {
            case 1: // crouch
                srv.Entities <- server.EntitiesCrouchUpdate{c.eid, true}
            case 2: // uncrouch
                srv.Entities <- server.EntitiesCrouchUpdate{c.eid, false}
            case 3: // leave bed
                panic("TODO: 0x13 leave bed")
            case 4: // start sprinting
                srv.Entities <- server.EntitiesSprintUpdate { c.eid, true }
            case 5: // stop sprinting
                srv.Entities <- server.EntitiesSprintUpdate { c.eid, false }
            default:
                panic(fmt.Sprintf("client/conn: i have no idea what 0x13 entity action '%v' is!", ea.ActionID))
            }

		default:
			panic(fmt.Sprintf("packet with id 0x%x not implemented", id))
		}
	}

	//panic(fmt.Sprintf("so far, so good! (player name is %v)", c.name))
}

func handleLogin(c *Conn, eid int32) {
	// start reading next packet
	id := c.c.ReadID()

	if id != mnet.LoginRequestID {
		panic(fmt.Sprintf("expected 0x01 login request, not %v", id))
	}

	lr := c.c.ReadLoginRequest()

	if lr.ProtocolVersion != 29 {
		// note: we might allow some limited backward compatibility eventually
		panic(fmt.Sprintf("wrong protocol version; expected 29 (mc 1.2.5) (got %v)", lr.ProtocolVersion))
	}

	// name is the same we got from the handshake, right?
	if lr.Username != c.name {
		panic("name in login request wasn't same as in handshake") // wtf
	}

	// login request was fine. let us reply in the affirmative!
	// TODO: correct everything below (spec. player eid!)
	c.s.LoginRequest(eid, "default", 0, 0, 1, 16)
}

func tryTellClientAboutPanic(c *Conn) {
	// tries to kick the client correctly, even though we might be panicking.

	defer func() {
		recover() // make sure we don't propagate any panics we might cause here!
	}()

	// ...are we panicking?
	if r := recover(); r != nil {
		fmt.Printf("closing connection from %v due to panic: %v\n", c.raddr, r)
		c.s.Kick(fmt.Sprintf("Panic: %v", r))
	}
}

func handleFirstPacket(c *Conn) {
	// are they querying us, or are they actually trying to join the game?
	id := c.c.ReadID()

	if id == mnet.ServerListPingID {
		// it's a query!
		// (packet is just the id, so there is no mnet.Read function for it)
		// TODO: actually ask the server for details
		// format: desc§playersConnected§maxPlayers
		c.s.Kick("merlyn's magic box of pain gain§12§34")
		panic("just a server details query")
	} else if id != mnet.HandshakeID {
		// otherwise: not a handshake?
		panic(fmt.Sprintf("that is not a good first packet id (expected %v Handshake, got %v)", mnet.HandshakeID, id))
	}

	// if we're here, then it must be an incoming c->s handshake packet.

	uavh := c.c.ReadHandshake()
	userAndVhostBits := strings.Split(uavh, ";")

	if len(userAndVhostBits) < 2 {
		// hmm. early client protocol/custom client?
		c.name = userAndVhostBits[0]
	} else if len(userAndVhostBits) == 2 {
		// that's more like it!  format is user;vhost:port.
		c.name = userAndVhostBits[0]
		// TODO: vhost!  (code probably not here but elsewhere, but still.)
	} else {
		// :s
		panic(fmt.Sprintf("expected user;vhost:port, got %v", uavh))
	}

	// ok.  if we're here, we were ok with the c->s part of the handshake.
	// we should now respond with an auth challenge.
	// TODO: (optional) authentication against minecraft.net!
	// TODO: ensure name is somewhat sane

	c.s.Handshake("-") // TEMP: no auth challenge. come on in.
}
