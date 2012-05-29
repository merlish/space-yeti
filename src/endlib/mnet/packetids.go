package mnet

// list of packet IDs
const (
	KeepAliveID             = 0x00
	LoginRequestID          = 0x01
	HandshakeID             = 0x02
	ChatMessageID           = 0x03
	TimeUpdateID            = 0x04
	EntityEquipmentID       = 0x05
	SpawnPositionID         = 0x06
	UseEntityID             = 0x07
	UpdateHealthID          = 0x08
	RespawnID               = 0x09
	PlayerID                = 0x0a
	PlayerPositionID        = 0x0b
	PlayerLookID            = 0x0c
	PlayerPositionAndLookID = 0x0d
	PlayerDiggingID         = 0x0e
	PlayerBlockPlacementID  = 0x0f
	HeldItemChangeID        = 0x10
	UseBedID                = 0x11
	AnimationID             = 0x12
	EntityActionID          = 0x13
	SpawnNamedEntityID      = 0x14
	SpawnDroppedEntityID    = 0x15
	CollectItemID           = 0x16
	SpawnObjectVehicleID    = 0x17
	SpawnMobID              = 0x18
	SpawnPaintingID         = 0x19
	SpawnExperienceOrbID    = 0x1a

	EntityVelocityID            = 0x1c
	DestroyEntityID             = 0x1d
	EntityID                    = 0x1e
	EntityRelativeMoveID        = 0x1f
	EntityLookID                = 0x20
	EntityLookAndRelativeMoveID = 0x21
	EntityTeleportID            = 0x22
	EntityHeadLookID            = 0x23

	EntityStatusID       = 0x26
	AttachEntityID       = 0x27
	EntityMetadataID     = 0x28
	EntityEffectID       = 0x29
	RemoveEntityEffectID = 0x2a
	SetExperienceID      = 0x2b

	MapColumnAllocationID = 0x32
	MapChunksID           = 0x33
	MultiBlockChangeID    = 0x34
	BlockChangeID         = 0x35
	BlockActionID         = 0x36

	ExplosionID           = 0x3c
	SoundParticleEffectID = 0x3d

	ChangeGameStateID = 0x46
	ThunderboltID     = 0x47

	OpenWindowID              = 0x64
	CloseWindowID             = 0x65
	ClickWindowID             = 0x66
	SetSlotID                 = 0x67
	SetWindowItemsID          = 0x68
	UpdateWindowPropertyID    = 0x69
	ConfirmTransactionID      = 0x6a
	CreativeInventoryActionID = 0x6b
	EnchantItemID             = 0x6c

	UpdateSignID       = 0x82
	ItemDataID         = 0x83
	UpdateTileEntityID = 0x84

	IncrementStatisticID = 0xc8
	PlayerListItemID     = 0xc9

	PlayerAbilitiesID = 0xca

	PluginMessageID  = 0xfa
	ServerListPingID = 0xfe
	DisconnectID     = 0xff
)
