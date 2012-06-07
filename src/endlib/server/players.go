package server

import "fmt"
import "endlib/mnet"

type PlayerStatusPacket struct {
    Eid int32
    Connected bool
}

type PlayerMovePacket struct {
	Eid      int32
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

type ConnectedPlayersSummaryPacket struct {
    Eids []int32
}

type UpdateEntityMetadataPacket struct {
    Eid int32
    Metadata mnet.Metadata
}

func playersServer(subscribe chan (chan interface{}), unsubscribe chan (chan interface{}), notify chan interface{}) {

	// TODO: verify packets received via notify make sense (=> anticheat code)
	//  (if they don't, we need to remember to tell the cheating player their new pos)
	//  (can't do it via notifying a subscribed channel. they ignore self notifications.)

	subscribers := make([]chan interface{}, 0)

    connected := make([]int32, 0)

	for {
		select {
		case s, _ := <-subscribe:
			subscribers = append(subscribers, s)

            // start by composing & sending a ConnectedPlayersSummaryPacket.
            connectedCopy := make([]int32, len(connected))
            copy(connectedCopy, connected)
            fmt.Println("dbg: server/players.go: got subscriber; shoving CPSP")
            s <- ConnectedPlayersSummaryPacket{connectedCopy}
            fmt.Println("dbg: server/players.go: subscriber took CPSP.  we continue.")

		case u, _ := <-unsubscribe:
			for i, s := range subscribers {
                if s == u {
                    subscribers = append(subscribers[:i], subscribers[i+1:]...)
                    break
                }
			}
		case n, _ := <-notify:
			// TODO: anticheat goes here!

            switch n.(type) {
            case PlayerStatusPacket:
                fmt.Println("dbg: server/players.go: got PSP")
                psp := n.(PlayerStatusPacket)

                // are they already in our list?
                inList := false
                for k, v := range connected {
                    if psp.Eid == v {
                        inList = true

                        if !psp.Connected {
                            connected = append(connected[:k], connected[k+1:]...)
                        }

                        break
                    }
                }

                if (!inList && psp.Connected) {
                    connected = append(connected, psp.Eid)
                }

                if (inList && psp.Connected) || (!inList && !psp.Connected) {
                    fmt.Printf("srv/players: caught player with eid %v repeating status notification connected=%v; dropped\n", psp.Eid, psp.Connected)
                    continue
                }

                fmt.Println("dbg: server/players.go: made it through PSP special case")
            }

			// tell everyone about the movement.
			// if we're telling someone about their own movements, this is filtered out
			//  by the code on the other side of the channel.
			for _, c := range subscribers {
				c <- n
			}

            //fmt.Println("dbg: server/players.go: send notification onwards")
		}
	}
}
