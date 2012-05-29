package mnet

import "fmt"

// helps compose entity metadata, to be sent to the client.

type Metadata []byte

// note: we don't include the 'stop reading' byte, 127, here.  that's done in the conn wr funcs.

func (m Metadata) addToMetadata(key int, vtype int, gen ...BytesGenerator) Metadata {
	buf := Buffer(m)

	if key > 0x1f || key < 0 {
		panic(fmt.Sprintf("invalid metadata key: expected 0..0x1f, got %v", key))
	}

	// write combined key&valuetype byte
	buf = buf.Add(UByte(byte((key & 0x1f) + (vtype << 5))))

	for _, g := range gen {
		buf = buf.Add(g)
	}

	return Metadata(buf)
}

func (m Metadata) AddToMetadataSByte(key int, data int8) Metadata {
	return m.addToMetadata(key, 0, SByte(data))
}

func (m Metadata) AddToMetadataShort(key int, data int16) Metadata {
	return m.addToMetadata(key, 1, Short(data))
}

func (m Metadata) AddToMetadataInt(key int, data int32) Metadata {
	return m.addToMetadata(key, 2, Int(data))
}

func (m Metadata) AddToMetadataSingle(key int, data float32) Metadata {
	return m.addToMetadata(key, 3, Single(data))
}

func (m Metadata) AddToMetadataString(key int, data string) Metadata {
	return m.addToMetadata(key, 4, String(data))
}

func (m Metadata) AddToMetadataSBS(key int, v1 int16, v2 int8, v3 int16) Metadata {
	return m.addToMetadata(key, 5, Short(v1), SByte(v2), Short(v3))
}

func (m Metadata) AddToMetadataIII(key int, v1 int32, v2 int32, v3 int32) Metadata {
	return m.addToMetadata(key, 6, Int(v1), Int(v2), Int(v3))
}
