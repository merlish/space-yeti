package server

type PlayerMovePacket struct {
	Eid      int
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func playersServer(subscribe chan (chan *PlayerMovePacket), unsubscribe chan (chan *PlayerMovePacket), notify chan *PlayerMovePacket) {

	// TODO: verify packets received via notify make sense (=> anticheat code)
	//  (if they don't, we need to remember to tell the cheating player their new pos)
	//  (can't do it via notifying a subscribed channel. they ignore self notifications.)

	subscribers := make([]chan *PlayerMovePacket, 0)

	for {
		select {
		case s, _ := <-subscribe:
			subscribers = append(subscribers, s)
		case u, _ := <-unsubscribe:
			for i, s := range subscribers {
                if s == u {
                    subscribers = append(subscribers[:i], subscribers[i+1:]...)
                    break
                }
			}
		case n, _ := <-notify:
			// TODO: anticheat goes here!

			// tell everyone about the movement.
			// if we're telling someone about their own movements, this is filtered out
			//  by the code on the other side of the channel.
			for _, c := range subscribers {
				c <- n
			}
		}
	}
}
