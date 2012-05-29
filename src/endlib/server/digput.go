package server

// digput: manages digging up & putting down of blocks.

// TODO: get timing right, dependent on tool used and on which block
// TODO: actually check that they're digging up an existing block/putting down in a space
// TODO: ensure player is actually close enough to the block to pick it up/put it down

type DigputRequest interface{}

type DigputStartDig struct {
    *BlockPosition
    Ret chan bool
}

type DigputFinishDig DigputStartDig

type DigputPut struct {
    *BlockPosition
    BlockType byte
    Ret chan bool
}

type DigputSubscribe struct {
    Listener chan DigputRequest
}

type DigputUnsubscribe struct {
    Listener chan DigputRequest
}

func digputServer(in chan DigputRequest) {

    subscribers := make([]chan DigputRequest, 0)

    for {
        r := <-in

        switch r.(type) {
        case *DigputStartDig:
            dpd := r.(*DigputStartDig)

            // confirm & then tell everyone about the dig.
            dpd.Ret <- true

            for _, s := range subscribers {
                s <- dpd
            }

        case *DigputFinishDig:
            dfd := r.(*DigputFinishDig)
            dfd.Ret <- true
            for _, s := range subscribers {
                s <- dfd
            }

        case *DigputSubscribe:
            dps := r.(*DigputSubscribe)
            subscribers = append(subscribers, dps.Listener)

        case *DigputUnsubscribe:
            dpu := r.(*DigputUnsubscribe)
            for i, s := range subscribers {
                if dpu.Listener == s {
                    subscribers = append(subscribers[:i], subscribers[i+1:]...)
                    break
                }
            }
            close(dpu.Listener)

        default:
            // TODO: actually list type here
            panic("request of type ? isn't understood by digputServer() yet")
        }

    }

}
