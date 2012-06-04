package server

import "fmt"

// TODO: make not terrible!

func eidServer(eids chan int32) {

    const numEids = 32768

    // list of available eids
    available := make([]int32, numEids)

    for i := 1; i < numEids + 1; i++ {
        //available = append(available, int32(i))
        available[i-1] = int32(i)
    }

    for {
        // probably will panic unnecessarily when 1 eid remains, but who cares?
        if len(available) == 0 {
            panic("out of eids! :( assumedly, there's a leak somewhere.")
        }

        select {
        case freed := <-eids: // someone returning an eid
            // TODO: in debug mode only: check eid is not already in free list?

            available = append(available, freed)

        case eids <- available[0]:  // someone trying to read an eid
            fmt.Printf("dbg: eid server: just handed out %v... next eid is %v\n", int(available[0]), available[1])
            available = available[1:]
        }
    }

}

