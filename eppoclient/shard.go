package eppoclient

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

type ShardRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func getShard(input string, subjectShards int64) int64 {
	hash := md5.Sum([]byte(input))
	hashOutput := hex.EncodeToString(hash[:])

	// get the first 4 bytes of the md5 hex string and parse it using base 16
	// (8 hex characters represent 4 bytes, e.g. 0xffffffff represents the max 4-byte integer)
	intVal, _ := strconv.ParseInt(hashOutput[0:8], 16, 0)

	return intVal % subjectShards
}

func isShardInRange(shard int, shardRange ShardRange) bool {
	return shard >= shardRange.Start && shard < shardRange.End
}
