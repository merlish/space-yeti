package client

import "fmt"
import "endlib/server"

func otherPlayersHandle(c *Conn, quit chan int, srv *server.Server) {

    defer func() { recover() }()

    pmpc := make(chan interface{})

    fmt.Println("dbg: client/o_p: subscribing to server/players")
    srv.Location.Subscribe <- pmpc


    // first packet we get should be a ConnectedPlayersSummaryPacket.
    fmt.Println("dbg: client/o_p: reading first packet from server/players")
    cpsp := (<-pmpc).(server.ConnectedPlayersSummaryPacket)

    for _, v := range cpsp.Eids {
        // v is the eid of an already connected player!
        // so.. let's spawn 'em clientside.
        fmt.Printf("dbg: client/o_p: spawning named entity %v (for ref we are %v)\n", v, c.eid)
        c.s.SpawnNamedEntity(v, fmt.Sprintf("[A]Player%x", v), 0, 0, 0, 0, 0, 0)
        c.s.Entity(v)

        // lie about equipped weapon & armor
        for i := 0; i < 5; i++ {
            c.s.EntityEquipment(v, int16(i), int16(-1), int16(0))
        }
    }

    // should rly tell the server we exist...
    go func() {
        fmt.Println("dbg: client/o_p: telling server/players about us")
        srv.Location.Notify <- server.PlayerStatusPacket{c.eid, true}
    }()
    // ^ on a goproc b/c what if player server is blocked trying to tell us about stuff atm

    // finally, start going asynchronously.
    fmt.Println("asyncA")
    go otherPlayersHandleBackground(pmpc, quit, c, srv)
    fmt.Println("asyncB :)")
}

func otherPlayersHandleBackground(pmpc chan interface{}, quit chan int, c *Conn, srv *server.Server) {

    defer func() { recover() }()
    defer func() {
        srv.Location.Unsubscribe <- pmpc

        // drain channel
        for _ = range pmpc {}
    }()

    defer func() {
        fmt.Println("dbg: client/o_p: telling server/players we're leaving")
        srv.Location.Notify <- server.PlayerStatusPacket{c.eid, false}
    }()

    // wait for & handle incoming player move things
    fmt.Println("dbg: client/o_p: main loop")
    for {
        select {
        case <-quit:
            break
        case <-pmpc:
            pmv := <-pmpc

            switch pmv.(type) {
            case server.PlayerStatusPacket:
                sp := pmv.(server.PlayerStatusPacket)

                if sp.Eid != c.eid {
                    if sp.Connected {
                        fmt.Printf("dbg: client/o_p: spawning named entity for player with eid %v (for ref, ours is %v)\n", sp.Eid, c.eid)
                        c.s.SpawnNamedEntity(sp.Eid, fmt.Sprintf("[>]Player%x", sp.Eid), 0, 0, 0, 0, 0, 0)
                        c.s.Entity(sp.Eid)
                        for i := 0; i < 5; i++ {
                            c.s.EntityEquipment(sp.Eid, int16(i), int16(-1), int16(0))
                        }
                    } else {
                        // ooh! they're going?!
                        fmt.Printf("dbg: client/o_p: deleting entity with eid %v (we are %v)\n", sp.Eid, c.eid)
                        c.s.DestroyEntity(sp.Eid)
                    }
                }

            case server.PlayerMovePacket:
                m := pmv.(server.PlayerMovePacket)

                if m.Eid != c.eid {
                    // it is a hack, but i really don't care for this relative move stuff...
                    //  i just plain don't see the point.  if we want to save bandwidth,
                    //  we should just not send the packet (or send them less often).
                    //  (TODO: send less often, the further away the entity is...?)

                    yaw := (256 * (int32(m.Yaw) % 360)) / 360
                    yaw8 := int8(0)
                    if yaw > 127 {
                        yaw8 = int8(yaw - 128)
                    } else {
                        yaw8 = int8(yaw)
                    }

                    //pitch := int8(((256 * (int32(m.Pitch) % 360)) / 360) - 128)
                    pitch := int8((256 * (int32(m.Pitch) % 360)) / 360)

                    //fmt.Printf("sending %v yaw,pitch as %v, %v\n", c.eid, yaw8, pitch)

                    c.s.EntityTeleport(int32(m.Eid), int32(m.X * 32), int32(m.Y * 32), int32(m.Z * 32), yaw8, pitch)
                    c.s.EntityHeadLook(int32(m.Eid), yaw8)

                    //fmt.Printf("propagating move of %d (%d,%d,%d)(%f,%f,%f)\n", m.Eid, int32(m.X), int32(m.Y), int32(m.Z), m.X, m.Y, m.Z)
                }

            case server.UpdateEntityMetadataPacket:
                uem := pmv.(server.UpdateEntityMetadataPacket)

                if uem.Eid != c.eid {
                    c.s.EntityMetadata(uem.Eid, uem.Metadata)
                }
            }

        }


    }

}

