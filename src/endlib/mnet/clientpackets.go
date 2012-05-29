package mnet

import "fmt"
import "math"
import "net"
import "unicode/utf16"

// packets from clients to servers.

/*type CConn struct {
    Conn net.Conn
    // read buffer?
}*/

type CConn net.TCPConn

func (c *CConn) readBytes(n int) []byte {
	// throw hissy fit if n is negative
	if n < 0 {
		panic(fmt.Sprintf("asked to read negative (%v) number of bytes", n))
	} else if n == 0 {
		return []byte{}
	}

	// otherwise...
	buf := make([]byte, n)
	// TODO: buffer net input, like a nice application!

	count, err := (*net.TCPConn)(c).Read(buf)

	if err != nil {
		panic(fmt.Sprintf("read failed: %v", err))
	} else if count < n {
		panic(fmt.Sprintf("client provided %v bytes, but %v were expected. (TODO: buffer net input?)", count, n))
	}

	return buf
}

func (c *CConn) readUByte() byte {
	return c.readBytes(1)[0]
}
func (c *CConn) readSByte() int8 {
	return int8(c.readBytes(1)[0])
}
func (c *CConn) readBool() bool {
	b := c.readUByte()

	if b == 1 {
		return true
	} else if b == 0 {
		return false
	}

	panic(fmt.Sprintf("read boolean from client should have been 0 or 1; instead, was %v", b))
}
func (c *CConn) readShort() int16 {
	bs := c.readBytes(2)
	return (int16(bs[0]) << 8) + int16(bs[1])
}
func (c *CConn) readInt() int32 {
	bs := c.readBytes(4)
	return (int32(bs[0]) << 24) + (int32(bs[1]) << 16) + (int32(bs[2]) << 8) + int32(bs[3])
}
func (c *CConn) readLong() int64 {
	bs := c.readBytes(8)
	return (int64(bs[0]) << 56) + (int64(bs[1]) << 48) + (int64(bs[2]) << 40) + (int64(bs[3]) << 32) + (int64(bs[4]) << 24) + (int64(bs[5]) << 16) + (int64(bs[6]) << 8) + int64(bs[7])
}
func (c *CConn) readSingle() float32 {
	return math.Float32frombits(uint32(c.readInt()))
}
func (c *CConn) readDouble() float64 {
    //bs := c.readBytes(8)
    //f := math.Float64frombits((int64(bs[7]) << 64) + (int64(bs[6]) << 48) + (int64(bs[5]) << 40) + (int64(bs[4]) << 32) + (int64(bs[3]) << 24) + (int64(bs[2]) << 16) + (int64(bs[1]) << 8) + int64(bs[0])
    f := math.Float64frombits(uint64(c.readLong()))
    fmt.Printf("DEEBBUUGG: %f\n", f)
    return f

    /*
    //l := c.readLong()
    bs := c.readBytes(8)
    l := int64(bs[7] << 56) + int64(bs[6] << 48) + int64(bs[5] << 40) + int64(bs[4] << 32) + int64(bs[3] << 24) + int64(bs[2] << 16) + int64(bs[1] << 8) + int64(bs[0])
    sw := int64(bs[0] << 56) + int64(bs[1] << 48) + int64(bs[2] << 40) + int64(bs[3] << 32) + int64(bs[4] << 24) + int64(bs[5] << 16) + int64(bs[6] << 8) + int64(bs[7])
    f := *(*float64)(unsafe.Pointer(&l))
    f = 0.0
    fmt.Printf("DEBUUUGGG: unswapped: %d, swapped: %d, after: %f\n", l, sw, f)
	return math.Float64frombits(uint64(l))
    //return f */
}
func (c *CConn) readString() string {
	// MC max string len is 256 chars (=> 512 bytes, right? ucs-2 not utf-16..)...
	length := c.readShort()
	if length > 256 {
		panic(fmt.Sprintf("this string is too dang long! (%v, expected <= 256 characters)", length))
	}

	// *2 b/c ucs-2
	bs := c.readBytes(int(length * 2))
	shs := make([]uint16, length)

	// TODO: if this is too slow, http://groups.google.com/group/golang-nuts/browse_thread/thread/68f8ef7a65781745 ... but i really don't think it will be, ever.
	for i := 0; i < len(bs); i += 2 {
		shs[i/2] = uint16(bs[i]<<8) + uint16(bs[i+1])
	}

	return string(utf16.Decode(shs))
}

func (c *CConn) ReadID() byte {
	return c.readUByte()
}

// 0x00
func (c *CConn) ReadKeepAlive() int32 {
	return c.readInt()
}

type CLoginRequest struct {
	ProtocolVersion int32
	Username        string
}

// 0x01
func (c *CConn) ReadLoginRequest() *CLoginRequest {
	// actually used fields
	pv, u := c.readInt(), c.readString()

	// unused fields
	c.readString()
	c.readInt()
	c.readInt()
	c.readSByte()
	c.readUByte()
	c.readUByte()

	// unlike in C, it's ok to return the address of a local var. 
	return &CLoginRequest{pv, u}
}

// 0x02
func (c *CConn) ReadHandshake() string {
	return c.readString()
}

// c.reader will have to convert from this type to string themselves...
// hopefully will help remind the c.reader they need to sanitise dat input.
type UnsanitizedChat string

// 0x03
func (c *CConn) ReadChatMessage() UnsanitizedChat {
	return UnsanitizedChat(c.readString())
}

type CUseEntity struct {
	User            int32
	Target          int32
	LeftMouseButton bool
}

// 0x07
func (c *CConn) ReadUseEntity() *CUseEntity {
	return &CUseEntity{c.readInt(), c.readInt(), c.readBool()}
}

// 0x09
func (c *CConn) ReadRespawn() {
	// i don't think we should trust anything the client says here. :/
	// dimension, difficulty, creative mode, world height, level type... c'mon!
	c.readInt()
	c.readSByte()
	c.readSByte()
	c.readShort()
	c.readString()
}

// 0x0a
func (c *CConn) ReadPlayer() (onGround bool) {
	return c.readBool()
}

type CPlayerPosition struct {
	X        float64
	Y        float64
	Stance   float64
	Z        float64
	OnGround bool // b/c 0x0b packet inherits from 0x0a packet
}

// 0x0b
func (c *CConn) ReadPlayerPosition() *CPlayerPosition {
	return &CPlayerPosition{c.readDouble(), c.readDouble(), c.readDouble(), c.readDouble(), c.readBool()}
}

type CPlayerLook struct {
	Yaw      float32
	Pitch    float32
	OnGround bool
}

// 0x0c
func (c *CConn) ReadPlayerLook() *CPlayerLook {
	return &CPlayerLook{c.readSingle(), c.readSingle(), c.readBool()}
}

type CPlayerPositionAndLook struct {
	X        float64
	Y        float64
	Stance   float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

// 0x0d
func (c *CConn) ReadPlayerPositionAndLook() *CPlayerPositionAndLook {
	return &CPlayerPositionAndLook{c.readDouble(), c.readDouble(), c.readDouble(), c.readDouble(), c.readSingle(), c.readSingle(), c.readBool()}
}

type CPlayerDigging struct {
	Status int8
	X      int32
	Y      byte
	Z      int32
	Face   int8
}

// 0x0e
func (c *CConn) ReadPlayerDigging() *CPlayerDigging {
	return &CPlayerDigging{c.readSByte(), c.readInt(), c.readUByte(), c.readInt(), c.readSByte()}
}

type CPlayerBlockPlacement struct {
	X         int32
	Y         byte
	Z         int32
	Direction int8
	Slot      *CSlot
}

// 0x0f
func (c *CConn) ReadPlayerBlockPlacement() *CPlayerBlockPlacement {
	return &CPlayerBlockPlacement{c.readInt(), c.readUByte(), c.readInt(), c.readSByte(), readSlot(c, false)}
}

// 0x10
func (c *CConn) ReadHeldItemChange() int16 {
	return c.readShort()
}

type CAnimation struct {
	EID       int32
	Animation int8
}

// 0x12
func (c *CConn) ReadAnimation() *CAnimation {
	return &CAnimation{c.readInt(), c.readSByte()}
}

type CEntityAction struct {
	EID      int32
	ActionID int8
}

// 0x13
func (c *CConn) ReadEntityAction() *CEntityAction {
	return &CEntityAction{c.readInt(), c.readSByte()}
}

// 0x65
func (c *CConn) ReadCloseWindow() int8 {
	return c.readSByte()
}

type CClickWindow struct {
	WindowID     int8
	SlotClicked  int16
	RightClicked bool
	ActionNumber int16
	Shift        bool
	ClickedItem  *CSlot
}

// 0x66
func (c *CConn) ReadClickWindow() *CClickWindow {
	return &CClickWindow{c.readSByte(), c.readShort(), c.readBool(), c.readShort(), c.readBool(), readSlot(c, true)}
}

type CConfirmTransaction struct {
	WindowID     int8
	ActionNumber int16
	Accepted     bool
}

// 0x6a
func (c *CConn) ReadConfirmTransaction() *CConfirmTransaction {
	return &CConfirmTransaction{c.readSByte(), c.readShort(), c.readBool()}
}

type CCreativeInventoryAction struct {
	SlotID      int16
	ClickedItem *CSlot
}

// 0x6b
func (c *CConn) ReadCreativeInventoryAction() *CCreativeInventoryAction {
	return &CCreativeInventoryAction{c.readShort(), readSlot(c, true)}
}

type CEnchantItem struct {
	WindowID    int8
	Enchantment int8
}

// 0x6c
func (c *CConn) ReadEnchantItem() *CEnchantItem {
	return &CEnchantItem{c.readSByte(), c.readSByte()}
}

type CUpdateSign struct {
	X    int32
	Y    int16
	Z    int32
	Text []string
}

// 0x82
func (c *CConn) ReadUpdateSign() *CUpdateSign {
	return &CUpdateSign{c.readInt(), c.readShort(), c.readInt(), []string{c.readString(), c.readString(), c.readString(), c.readString()}}
}

type CPlayerAbilities struct {
	Invulnerable   bool
	Flying         bool
	CanFly         bool
	InstantDestroy bool
}

// 0xca
func (c *CConn) ReadPlayerAbilities() *CPlayerAbilities {
	return &CPlayerAbilities{c.readBool(), c.readBool(), c.readBool(), c.readBool()}
}

type CPluginMessage struct {
	Channel string
	Data    []byte
}

// 0xfa
func (c *CConn) ReadPluginMessage() *CPluginMessage {
	// TODO: set max limit on Length?
	channel, length := c.readString(), c.readShort()

	// empty?
	if length == 0 {
		return &CPluginMessage{channel, []byte{}}
	}

	// otherwise...
	return &CPluginMessage{channel, c.readBytes(int(length))}
}

// 0xff
func (c *CConn) ReadDisconnect() (reason string) {
	reason = c.readString()
	return
}
