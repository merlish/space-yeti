package server

const (
	MapGetColumns = iota
	//MapSubscribeColumns
	//MapUnsubscribeColumns
)

type Chunk struct {
	BlockTypes    []byte // byte per block
	BlockMetadata []byte // nibble per block
	BlockLight    []byte // nibble per block
	SkyLight      []byte // nibble per block
	//Add []byte // nibble per block
}

type MapColumn struct {
	X      int
	Z      int
	Chunks []*Chunk
}

type MapRequest interface{}

type MapGetColumnsRequest struct {
	StartX int
	EndX   int
	StartZ int
	EndZ   int
	Data   interface{}
	Ret    chan interface{}
}

type MapHintEntityPosition struct {
    *EntityPosition
}

type ChunkPosition BlockPosition

func mapServer(in chan MapRequest) {
	fakeground := buildFakeGroundChunk()
	fakeair := buildFakeAirChunk()

	for {
		req := (<-in).(*MapGetColumnsRequest) // temp

		//if req.Id == MapGetColumns {

			for z := req.StartZ; z <= req.EndZ; z++ {
				for x := req.StartX; x <= req.EndX; x++ {
					chunks := make([]*Chunk, 16)

					// need 4 ground chunks, 12 air chunks (.'. 1/4 full world)
					for i := 0; i < 4; i++ {
						chunks[i] = fakeground
					}
					for i := 4; i < 16; i++ {
						chunks[i] = fakeair
					}

					col := &MapColumn{x, z, chunks}
					req.Ret <- col
				}
			}

		/*} else {
			panic(fmt.Sprintf("that mapServer req.id (%v) is not implemented yet!", req.Id))
		}*/
	}
}

func buildFakeGroundChunk() *Chunk {
	numblocks := 16 * 16 * 16

	types := make([]byte, numblocks)
	for i := range types {
		types[i] = 87 // dat netherrack
	}
	metadatas := make([]byte, numblocks/2)
	blocklight := make([]byte, numblocks/2)
	for i := range blocklight {
		blocklight[i] = 0xff // both blocks should have light 0xf
	}
	skylight := make([]byte, numblocks/2)
	for i := range skylight {
		skylight[i] = 0xff // both blocks should have skylight 0xf
	}

	return &Chunk{types, metadatas, blocklight, skylight}
}

func buildFakeAirChunk() *Chunk {
	numblocks := 16 * 16 * 16

	types := make([]byte, numblocks)
	metadatas := make([]byte, numblocks/2)
	blocklight := make([]byte, numblocks/2)
	for i := range blocklight {
		blocklight[i] = 0xff // both blocks should have light 0xf
	}
	skylight := make([]byte, numblocks/2)
	for i := range skylight {
		skylight[i] = 0xff // both blocks should have skylight 0xf
	}

	return &Chunk{types, metadatas, blocklight, skylight}
}
