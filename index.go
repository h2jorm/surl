package main

import (
	"os"
)

type Index struct {
	file *os.File
}

const (
	indexSize int64 = 5
)

func (idx *Index) NextID() (id int64, err error) {
	var stat os.FileInfo
	if stat, err = idx.file.Stat(); err != nil {
		return
	}
	id = stat.Size() / indexSize
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

func (idx *Index) CoordinateOfID(id int64) (start Coordinate, end Coordinate, err error) {
	buf := make([]byte, 5)
	if _, err = idx.file.ReadAt(buf, id*indexSize); err != nil {
		return
	}
	end = NewCoordinate(buf)
	prevID := id - 1
	if prevID < 0 {
		start = NewCoordinate([]byte{0, 0, 0, 0, 0})
		return
	}
	if _, err = idx.file.ReadAt(buf, prevID*indexSize); err != nil {
		return
	}
	start = NewCoordinate(buf)
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
