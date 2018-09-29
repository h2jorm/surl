package surl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// 0    1    2    3    4    5    6    7    0    1    2    3    4    5    6    7
// +----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+
// +     partition     +                     length                           +
// +                           position                                       +
// +                           position                                       +

type coordinate struct {
	partition uint8
	len       uint16
	pos       uint32
}

func (c *coordinate) Bytes() (ret []byte, err error) {
	buf := new(bytes.Buffer)
	head := uint16(c.partition&0x0f)<<12 + c.len&0x0fff
	if err = binary.Write(buf, binary.BigEndian, head); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, c.pos); err != nil {
		return
	}
	ret = buf.Bytes()
	return
}

func (c coordinate) offset() int64 {
	return int64(c.partition)*int64(math.Pow(2, 32)) + int64(c.pos)
}

func (c coordinate) String() string {
	return fmt.Sprintf("partition: %d, pos: %d, len: %d", c.partition, c.pos, c.len)
}

func newCoordinate(bs [6]byte) (coor coordinate, err error) {
	head := binary.BigEndian.Uint16(bs[0:2])
	partition := uint8(head & 0xf000 >> 12)
	len := head & 0x0fff
	pos := binary.BigEndian.Uint32(bs[2:6])
	coor = coordinate{partition: partition, len: len, pos: pos}
	return
}
