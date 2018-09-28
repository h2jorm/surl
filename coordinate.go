package main

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

type Coordinate struct {
	partition uint8
	len       uint16
	pos       uint32
}

func (c *Coordinate) Bytes() (ret []byte, err error) {
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

func (c *Coordinate) Partition() uint8 {
	return c.partition
}

func (c *Coordinate) Len() uint16 {
	return c.len
}

func (c *Coordinate) Offset() int64 {
	return int64(c.partition)*int64(math.Pow(2, 32)) + int64(c.pos)
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("partition: %d, pos: %d, len: %d", c.partition, c.pos, c.len)
}

func (c *Coordinate) Grow(num int) {
	if uint32(math.Pow(2, 32)-float64(num)) < c.pos {
		c.partition++
		c.pos = 0
	} else {
		c.pos += uint32(num)
	}
}

func NewCoordinate(bs [6]byte) (coordinate Coordinate, err error) {
	head := binary.BigEndian.Uint16(bs[0:2])
	partition := uint8(head & 0xf000 >> 12)
	len := head & 0x0fff
	pos := binary.BigEndian.Uint32(bs[2:6])
	coordinate = Coordinate{partition: partition, len: len, pos: pos}
	return
}
