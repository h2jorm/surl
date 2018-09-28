package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type Coordinate struct {
	partition uint8
	pos       uint32
}

func (c *Coordinate) Bytes() (ret []byte, err error) {
	buf := new(bytes.Buffer)
	if err = binary.Write(buf, binary.BigEndian, c.partition); err != nil {
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

func (c *Coordinate) Offset() int64 {
	return int64(c.partition)*int64(math.Pow(2, 32)) + int64(c.pos)
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("partition: %d, pos: %d", c.partition, c.pos)
}

func (c *Coordinate) Grow(num int) {
	if uint32(math.Pow(2, 32)-float64(num)) < c.pos {
		c.partition++
		c.pos = 0
	} else {
		c.pos += uint32(num)
	}
}

func NewCoordinate(bs []byte) (coordinate Coordinate) {
	partition := bs[0]
	pos := binary.BigEndian.Uint32(bs[1:5])
	coordinate = Coordinate{partition: partition, pos: pos}
	return
}
