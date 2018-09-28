package main

import (
	"os"
)

type Index struct {
	file *os.File
}

const (
	indexSize int64 = 6
)

func (idx *Index) NextID() (id int64, err error) {
	var stat os.FileInfo
	if stat, err = idx.file.Stat(); err != nil {
		return
	}
	id = stat.Size() / indexSize
	return
}

func (idx *Index) IDOccupied(id int64) (occupied bool, err error) {
	var nextID int64
	if nextID, err = idx.NextID(); err != nil {
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

func (idx *Index) WriteAt(coordinate Coordinate, id int64) (err error) {
	var bytes []byte
	if bytes, err = coordinate.Bytes(); err != nil {
		return
	}
	offset := id * indexSize
	_, err = idx.file.WriteAt(bytes, offset)
	return
}

func (idx *Index) CoordinateOfID(id int64) (coordinate Coordinate, err error) {
	buf := make([]byte, 6, 6)
	if _, err = idx.file.ReadAt(buf, id*indexSize); err != nil {
		return
	}
	var bytes [6]byte
	copy(bytes[:], buf)
	if coordinate, err = NewCoordinate(bytes); err != nil {
		return
	}
	return
}

func NewIndex() (index *Index, err error) {
	var file *os.File
	if file, err = os.OpenFile("./data/index", os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	index = &Index{file: file}
	return
}
