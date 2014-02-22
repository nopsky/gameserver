package gameserver

import (
	"hash/crc32"
	"lib/packet"
)

func PacketData(seqId uint64, data []byte) []byte {
	writer := packet.Writer()
	//size uint16
	writer.WriteU16(uint16(len(data)))
	//crc32 uint32
	crc32 := crc32.Checksum(data, crc32.IEEETable)

	writer.WriteU32(crc32)
	//seqid uint64
	writer.WriteU64(seqId)
	//data (uid + msgid + msgpack)
	writer.WriteRawBytes(data)

	return writer.Data()
}
