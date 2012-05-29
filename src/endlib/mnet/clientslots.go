package mnet

// helps reads Slots from client packets.

type CSlot struct {
	ID       int16
	Count    int8
	Metadata int16
	NBT      []byte
}

func isEnchantable(id int16) bool {
	return (256 <= id && id <= 259) || (267 <= id && id <= 279) || (283 <= id && id <= 286) || (290 <= id && id <= 294) || (298 <= id && id <= 317) || id == 261 || id == 359 || id == 346
}

func readSlot(c *CConn, isItem bool) *CSlot {
	id := c.readShort()

	// no item?
	if id == -1 {
		return &CSlot{-1, 0, 0, []byte{}}
	}

	count, metadata := c.readSByte(), c.readShort()

	// if it's not enchantable... (ugh!)
	if !isItem || !isEnchantable(id) {
		return &CSlot{id, count, metadata, []byte{}}
	}

	nbtLength := c.readShort()

	// if no nbt data follows...
	if nbtLength == -1 {
		return &CSlot{id, count, metadata, []byte{}}
	}

	return &CSlot{id, count, metadata, c.readBytes(int(nbtLength))}
}

func (slot *CSlot) IsEmpty() bool {
	return slot.ID == -1
}
