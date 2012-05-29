package mnet

import "math"
import "unicode/utf16"

// Does low-level sending and receiving of basic Minecraft types.
//
// Note that we need to send a minecraft packet as one net packet, as much as we can.
//  Or rather calling code does, and we need to help it.
// (The real Minecraft client & servers don't buffer incoming packets in any way :/)

type BytesGenerator interface {
	Bytes() []byte
}

type Buffer []byte

type UByte byte
type SByte int8
type Bool bool
type Short int16
type Int int32
type Long int64
type Single float32
type Double float64

type String string

func (x UByte) Bytes() []byte {
	return []byte{byte(x)} // wow
}

func (sb SByte) Bytes() []byte {
	return []byte{byte(sb)}
}

func (b Bool) Bytes() []byte {
	if b {
		return []byte{1}
	}

	// otherwise...
	return []byte{0}
}

func (short Short) Bytes() []byte {
	return []byte{byte(short >> 8), byte(short)}
}

func (i Int) Bytes() []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func (f Single) Bytes() []byte {
	return Int(int32(math.Float32bits(float32(f)))).Bytes()
	//return []byte{byte(f >> 24), byte(f >> 16), byte(f >> 8), byte(f)}
}

func (d Double) Bytes() []byte {
	return Long(int64(math.Float64bits(float64(d)))).Bytes()
	//return []byte{byte(d >> 56), byte(d >> 48), byte(d >> 40), byte(d >> 32),
	//byte(d >> 24), byte(d >> 16), byte(d >> 8), byte(d)}
}

func (l Long) Bytes() []byte {
	return []byte{byte(l >> 56), byte(l >> 48), byte(l >> 40), byte(l >> 32),
		byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}
}

func (s String) Bytes() []byte {
	u16str := utf16.Encode([]rune(string(s)))

	// write string length
	outv := Short(len(u16str)).Bytes()

	// write string data
	for i := range u16str {
		outv = append(outv, Short(u16str[i]).Bytes()...) // three dots unpacks received slice
	}

	return outv
}

func (b Buffer) Bytes() []byte {
	return []byte(b) // ahem
}

func New(pid byte, gens ...BytesGenerator) (buf Buffer) {
	buf = Buffer([]byte{pid})

	for _, g := range gens {
		buf = buf.Add(g)
	}

	return
}

func (buf Buffer) Add(gen BytesGenerator) Buffer {
	return append(buf, gen.Bytes()...)
}

func (buf Buffer) AddBytes(bs []byte) Buffer {
	return append(buf, bs...)
}
