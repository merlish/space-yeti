package mnet

import "fmt"
import "net"

// writes server->client packets to given connections.

type SConn net.TCPConn

func (s *SConn) write(buf Buffer) {
	if _, err := (*net.IPConn)(s).Write([]byte(buf)); err != nil {
		panic(fmt.Sprintf("failed to write to connection: %v", err))
	}
}

func (s *SConn) wr(pid byte, gens ...BytesGenerator) {
	s.write(New(pid, gens...))
}

func (s *SConn) KeepAlive(id int32) {
	s.write(New(0x00).Add(Int(id)))
}

func (s *SConn) LoginRequest(eid int32, levelType string, serverMode int32, dimension int32, difficulty int8, maxPlayers byte) {
	s.wr(0x01, Int(eid), String(""), String(levelType), Int(serverMode), Int(dimension), SByte(difficulty), UByte(0), UByte(maxPlayers))
}

func (s *SConn) Handshake(challenge string) {
	s.wr(0x02, String(challenge))
}

func (s *SConn) ChatMessage(msg string) {
	s.wr(0x03, String(msg))
}

func (s *SConn) TimeUpdate(time int64) {
	s.wr(0x04, Long(time))
}

func (s *SConn) EntityEquipment(eid int32, slot int16, itemId int16, damage int16) {
	s.wr(0x05, Int(eid), Short(slot), Short(itemId), Short(damage))
}

func (s *SConn) SpawnPosition(x int32, y int32, z int32) {
	s.wr(0x06, Int(x), Int(y), Int(z))
}

func (s *SConn) UpdateHealth(health int16, food int16, foodSaturation float32) {
	s.wr(0x08, Short(health), Short(food), Single(foodSaturation))
}

func (s *SConn) Respawn(dimension int32, difficulty int8, creativeMode int8, worldHeight int16, levelType string) {
	s.wr(0x09, Int(dimension), SByte(difficulty), SByte(creativeMode), Short(worldHeight), String(levelType))
}

func (s *SConn) PlayerPositionAndLook(x float64, stance float64, y float64, z float64, yaw float32, pitch float32, onGround bool) {
	s.wr(0x0d, Double(x), Double(stance), Double(y), Double(z), Single(yaw), Single(pitch), Bool(onGround))
}

func (s *SConn) UseBed(eid int32, x int32, y byte, z int32) {
	s.wr(0x11, Int(eid), SByte(0), Int(x), UByte(y), Int(z))
}

func (s *SConn) Animation(eid int32, anim int8) {
	s.wr(0x12, Int(eid), SByte(anim))
}

func (s *SConn) SpawnNamedEntity(eid int32, name string, x int32, y int32, z int32, packedYaw int8, packedPitch int8, currItem int16) {
	s.wr(0x14, Int(eid), String(name), Int(x), Int(y), Int(z), SByte(packedYaw), SByte(packedPitch), Short(currItem))
}

func (s *SConn) SpawnDroppedItem(eid int32, itemId int16, count int8, damageMetadata int16, x int32, y int32, z int32, packedRot int8, packedPitch int8, packedRoll int8) {
	s.wr(0x15, Int(eid), Short(itemId), SByte(count), Short(damageMetadata), Int(x), Int(y), Int(z), SByte(packedRot), SByte(packedPitch), SByte(packedRoll))
}

func (s *SConn) CollectItem(collectedEid int32, collectorEid int32) {
	s.wr(0x16, Int(collectedEid), Int(collectorEid))
}

func (s *SConn) SpawnObjectVehicleNotFireball(eid int32, objectType int8, x int32, y int32, z int32) {
	s.wr(0x17, Int(eid), SByte(objectType), Int(x), Int(y), Int(z), Int(0))
}

func (s *SConn) SpawnObjectVehicleFireball(eid int32, objectType int8, x int32, y int32, z int32, throwersEid int32, speedX int16, speedY int16, speedZ int16) {
	if throwersEid == 0 {
		panic("should use SpawnObjectVehicleNotFireball if throwersEid is 0!")
	}
	s.wr(0x17, Int(eid), SByte(objectType), Int(x), Int(y), Int(z), Int(throwersEid), Short(speedX), Short(speedY), Short(speedZ))
}

func (s *SConn) SpawnMob(eid int32, mobType int8, x int32, y int32, z int32, packedYaw int8, packedPitch int8, packedHeadYaw int8, metadata Metadata) {
	s.wr(0x18, Int(eid), SByte(mobType), Int(x), Int(y), Int(z), SByte(packedYaw), SByte(packedPitch), SByte(packedHeadYaw), Buffer(metadata), UByte(127)) // 127 finishes metadata
}

func (s *SConn) SpawnPainting(eid int32, title string, x int32, y int32, z int32, direction int32) {
	s.wr(0x19, Int(eid), String(title), Int(x), Int(y), Int(z), Int(direction))
}

func (s *SConn) SpawnExperienceOrb(eid int32, x int32, y int32, z int32, count int16) {
	s.wr(0x1a, Int(eid), Int(x), Int(y), Int(z), Short(count))
}

func (s *SConn) EntityVelocity(eid int32, velX int16, velY int16, velZ int16) {
	s.wr(0x1c, Int(eid), Short(velX), Short(velY), Short(velZ))
}

func (s *SConn) DestroyEntity(eid int32) {
	s.wr(0x1d, Int(eid))
}

func (s *SConn) Entity(eid int32) {
	s.wr(0x1e, Int(eid))
}

func (s *SConn) EntityRelativeMove(eid int32, relX int8, relY int8, relZ int8) {
	s.wr(0x1f, Int(eid), SByte(relX), SByte(relY), SByte(relZ))
}

func (s *SConn) EntityLook(eid int32, yaw int8, pitch int8) {
	s.wr(0x20, Int(eid), SByte(yaw), SByte(pitch))
}

func (s *SConn) EntityLookAndRelativeMove(eid int32, relX int8, relY int8, relZ int8, yaw int8, pitch int8) {
	s.wr(0x21, Int(eid), SByte(relX), SByte(relY), SByte(relZ), SByte(yaw), SByte(pitch))
}

func (s *SConn) EntityTeleport(eid int32, x int32, y int32, z int32, yaw int8, pitch int8) {
	s.wr(0x22, Int(eid), Int(x), Int(y), Int(z), SByte(yaw), SByte(pitch))
}

func (s *SConn) EntityHeadLook(eid int32, headYaw int8) {
	s.wr(0x23, Int(eid), SByte(headYaw))
}

func (s *SConn) EntityStatus(eid int32, status int8) {
	s.wr(0x26, Int(eid), SByte(status))
}

func (s *SConn) AttachEntity(eid int32, vehicleEid int32) {
	s.wr(0x27, Int(eid), Int(vehicleEid))
}

func (s *SConn) EntityMetadata(eid int32, metadata Metadata) {
	s.wr(0x28, Int(eid), Buffer(metadata), UByte(127)) // 127 finishes metadata
}

func (s *SConn) EntityEffect(eid int32, effectId int8, amplifier int8, duration int16) {
	s.wr(0x29, Int(eid), SByte(effectId), SByte(amplifier), Short(duration))
}

func (s *SConn) RemoveEntityEffect(eid int32, effectId int8) {
	s.wr(0x2a, Int(eid), SByte(effectId))
}

func (s *SConn) SetExperience(barFullness float32, level int16, totalExp int16) {
	s.wr(0x2b, Single(barFullness), Short(level), Short(totalExp))
}

func (s *SConn) MapColumnAllocation(columnX int32, columnZ int32, create bool) {
	s.wr(0x32, Int(columnX), Int(columnZ), Bool(create))
}

func (s *SConn) MapChunk(columnX int32, columnZ int32, groundUpContinuous bool, priBitmap uint16, addBitmap uint16, compressedData []byte) {
	s.wr(0x33, Int(columnX), Int(columnZ), Bool(groundUpContinuous), Short(int16(priBitmap)), Short(int16(addBitmap)), Int(len(compressedData)), Int(0), Buffer(compressedData))
}

func (s *SConn) MultiBlockChange(chunkX int32, chunkZ int32, count int16, data []byte) {
	s.wr(0x34, Int(chunkX), Int(chunkZ), Short(count), Int(len(data)), Buffer(data))
}

func (s *SConn) BlockChange(x int32, y byte, z int32, blockType SByte, blockMetadata SByte) {
	s.wr(0x35, Int(x), UByte(y), Int(z), SByte(blockType), SByte(blockMetadata))
}

func (s *SConn) BlockAction(x int32, y int16, z int32, d1 SByte, d2 SByte) {
	s.wr(0x36, Int(x), Short(y), Int(z), SByte(d1), SByte(d2))
}

func (s *SConn) Explosion(x float64, y float64, z float64, d1 float32, rcount int32, recs []byte) {
	s.wr(0x3c, Double(x), Double(y), Double(z), Single(d1), Int(rcount), Buffer(recs))
}

func (s *SConn) SoundParticleEffect(eid int32, x int32, y byte, z int32, data int32) {
	s.wr(0x3d, Int(eid), Int(x), UByte(y), Int(z), Int(data))
}

func (s *SConn) ChangeGameState(reason int8, gameMode int8) {
	s.wr(0x46, SByte(reason), SByte(gameMode))
}

func (s *SConn) Thunderbolt(eid int32, x int32, y int32, z int32) {
	s.wr(0x47, Int(eid), Bool(true), Int(x), Int(y), Int(z))
}

func (s *SConn) OpenWindow(windowId int8, invType int8, windowTitle string, numSlots int8) {
	s.wr(0x64, SByte(windowId), SByte(invType), String(windowTitle), SByte(numSlots))
}

func (s *SConn) CloseWindow(windowId int8) {
	s.wr(0x65, SByte(windowId))
}

func (s *SConn) SetSlot(windowId int8, slot int16, slotData []byte) {
	s.wr(0x67, SByte(windowId), Short(slot), Buffer(slotData))
}

func (s *SConn) SetWindowItems(windowId int8, count int16, slotData []byte) {
	s.wr(0x68, SByte(windowId), Short(count), Buffer(slotData))
}

func (s *SConn) UpdateWindowProperty(windowId int8, property int16, value int16) {
	s.wr(0x69, SByte(windowId), Short(property), Short(value))
}

func (s *SConn) ConfirmTransaction(windowId int8, actionId int16, accepted bool) {
	s.wr(0x6a, SByte(windowId), Short(actionId), Bool(accepted))
}

func (s *SConn) CreativeInventoryAction(slot int16, slotData []byte) {
	s.wr(0x6b, Short(slot), Buffer(slotData))
}

func (s *SConn) UpdateSign(x int32, y int16, z int32, line1 string, line2 string, line3 string, line4 string) {
	s.wr(0x82, Int(x), Short(y), Int(z), String(line1), String(line2), String(line3), String(line4))
}

func (s *SConn) ItemData(itemType int16, itemId int16, data []byte) {
	if len(data) > 255 {
		panic(fmt.Sprintf("extended data is too long. expected 0..255, got %v", len(data)))
	}
	s.wr(0x83, Short(itemType), Short(itemId), UByte(byte(len(data))), Buffer(data))
}

func (s *SConn) UpdateTileEntity(x int32, y int16, z int32, action int8, d1 int32, d2 int32, d3 int32) {
	s.wr(0x84, Int(x), Short(y), Int(z), SByte(action), Int(d1), Int(d2), Int(d3))
}

func (s *SConn) IncrementStatistic(sid int32, amount int8) {
	s.wr(0xc8, Int(sid), SByte(amount))
}

func (s *SConn) PlayerListItem(playerName string, online bool, ping int16) {
	s.wr(0xc9, String(playerName), Bool(online), Short(ping))
}

func (s *SConn) PlayerAbilities(invulnerable bool, flying bool, canFly bool, instantDestroy bool) {
	s.wr(0xca, Bool(invulnerable), Bool(flying), Bool(canFly), Bool(instantDestroy))
}

func (s *SConn) PluginMessage(channel string, data []byte) {
	s.wr(0xfa, String(channel), Short(int16(len(data))), Buffer(data))
}

func (s *SConn) Kick(reason string) {
	s.wr(0xff, String(reason))
}
