package surl

import (
	"math"
	"os"
)

type store struct {
	file *os.File
}

func (s *store) nextCoordinate() (coor coordinate, err error) {
	var stat os.FileInfo
	if stat, err = s.file.Stat(); err != nil {
		return
	}
	partition := stat.Size() / int64(math.Pow(2, 32))
	pos := stat.Size() % int64(math.Pow(2, 32))
	coor = coordinate{partition: uint8(partition), pos: uint32(pos)}
	return
}

func (s *store) writeAt(url string, coor coordinate) (err error) {
	_, err = s.file.WriteAt([]byte(url), coor.offset())
	return
}

func (s *store) urlOfCoordinate(coor coordinate) (url string, err error) {
	len := coor.len
	offset := coor.offset()
	buf := make([]byte, len)
	if _, err = s.file.ReadAt(buf, offset); err != nil {
		return
	}
	url = string(buf)
	return
}
