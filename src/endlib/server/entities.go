package server

import "fmt"
import "endlib/mnet"

const (
    player = iota
)

type entityMetadata struct {
    Eid int32
    Type int
    Flags byte
}

var entityData = make(map[int32] *entityMetadata, 0)

type EntitiesCreatePlayer struct {
    Eid int32
    Name string
}

type EntitiesDelete struct {
    Eid int32
}

type EntitiesCrouchUpdate struct {
    Eid int32
    Crouching bool
}

type EntitiesSprintUpdate struct {
    Eid int32
    Sprinting bool
}

func entityManager(in chan interface{}, pin chan interface{}) {

    for {
        im := <-in

        switch im.(type) {
        case EntitiesCreatePlayer:
            ecp := im.(EntitiesCreatePlayer)

            em := &entityMetadata{ecp.Eid, player, byte(0)}
            entityData[ecp.Eid] = em

        case EntitiesCrouchUpdate:
            ecu := im.(EntitiesCrouchUpdate)
            doUpdateFlag(ecu.Eid, pin, 0x02, ecu.Crouching)

        case EntitiesSprintUpdate:
            esu := im.(EntitiesSprintUpdate)
            doUpdateFlag(esu.Eid, pin, 0x08, esu.Sprinting)

        case EntitiesDelete:
            ed := im.(EntitiesDelete)

            delete(entityData, ed.Eid)
        }

    }

}

func doUpdateFlag(eid int32, playersIn chan interface{}, bitmask int, on bool) {
    doUpdate(eid, playersIn, func(m *entityMetadata) {
        if on {
            m.Flags = byte((int(m.Flags) & (0xff - bitmask)) + bitmask)
        } else {
            m.Flags = byte(int(m.Flags) & (0xff - bitmask))
        }
    })
}

func doUpdate(eid int32, playersIn chan interface{}, uf func(*entityMetadata)) {
    m := getMetadata(eid)

    uf(m)

    postUpdate(m, playersIn)
}

func postUpdate(m *entityMetadata, playersIn chan interface{}) {
    // format dat metadata
    fm := mnet.NewMetadata().AddToMetadataSByte(0, int8(m.Flags)).AddToMetadataShort(1, 300)
    fm = fm.AddToMetadataInt(8, 0)

    /*if m.type == animal {
        fm = fm.AddToMetadataInt(12, 0)
    }*/

    playersIn <- UpdateEntityMetadataPacket{m.Eid, fm}
}

func getMetadata(eid int32) *entityMetadata {
    em, ok := entityData[eid]

    if !ok {
        panic(fmt.Sprintf("server/entities: tried to get metadata for non-recorded entity w/eid %v", eid))
    }

    return em
}

