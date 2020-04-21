package utils

import (
	"crypto/rand"
	"encoding/binary"
	"os"
	"strconv"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func UniqueID() (uint64, error) {
	return sf.NextID()
}

func init() {
	var st sonyflake.Settings
	var tq_id uint64
	var err error
	s := os.Getenv("TQ_ID")
	if len(s) > 0 {
		tq_id, err = strconv.ParseUint(os.Getenv("TQ_ID"), 16, 16)
		if err != nil {
			panic("Unable to parse TQ_ID")
		}
	} else {
		rb := make([]byte, 8)
		_, err = rand.Read(rb)
		if err != nil {
			panic("Unable to generate random TQ_ID")
		}
		tq_id = binary.BigEndian.Uint64(rb)
	}

	st.MachineID = func() (uint16, error) { return uint16(tq_id), nil }
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("Unable to create sonyflake")
	}
}
