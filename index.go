package surl

import (
	"os"
)

type index struct {
	file *os.File
}

const (
	indexSize int64 = 6
)

func (idx *index) nextID() (id int64, err error) {
	var stat os.FileInfo
	if stat, err = idx.file.Stat(); err != nil {
		return
	}
	id = stat.Size() / indexSize
	return
}

func (idx *index) idOccupied(id int64) (occupied bool, err error) {
	var nextID int64
	if nextID, err = idx.nextID(); err != nil {
		return
	}
	if id >= nextID {
		occupied = false
		return
	}
	buf := make([]byte, 6, 6)
	if _, err = idx.file.ReadAt(buf, id*indexSize); err != nil {
		return
	}

	for _, byte := range buf {
		if byte != 0 {
			occupied = true
			return
		}
	}
	occupied = false
	return
}

func (idx *index) writeAt(coor coordinate, id int64) (err error) {
	var bytes []byte
	if bytes, err = coor.Bytes(); err != nil {
		return
	}
	offset := id * indexSize
	_, err = idx.file.WriteAt(bytes, offset)
	return
}

func (idx *index) coordinateOfID(id int64) (coor coordinate, err error) {
	buf := make([]byte, 6, 6)
	if _, err = idx.file.ReadAt(buf, id*indexSize); err != nil {
		return
	}
	var bytes [6]byte
	copy(bytes[:], buf)
	if coor, err = newCoordinate(bytes); err != nil {
		return
	}
	return
}
