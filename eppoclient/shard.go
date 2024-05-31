package eppoclient

import (
	"crypto/md5"
	"encoding/binary"
)

func getShard(input string, totalShards int64) int64 {
	hash := md5.Sum([]byte(input))
	// Only first 4 bytes of md5 are used for the shard value.
	intVal := int64(binary.BigEndian.Uint32(hash[:4]))
	return intVal % totalShards
}

func isShardInRange(shard int64, shardRange shardRange) bool {
	return shard >= shardRange.Start && shard < shardRange.End
}
