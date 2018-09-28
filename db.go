package main

import (
	"math"
	"os"
)

type DB struct {
	file *os.File
}

func (db *DB) CurrentCoordinate() (coordinate Coordinate, err error) {
	var stat os.FileInfo
	if stat, err = db.file.Stat(); err != nil {
		return
	}
	partition := stat.Size() / int64(math.Pow(2, 32))
	pos := stat.Size() % int64(math.Pow(2, 32))
	coordinate = Coordinate{partition: uint8(partition), pos: uint32(pos)}
	return
}

func (db *DB) WriteAt(url string, coordinate Coordinate) (err error) {
	_, err = db.file.WriteAt([]byte(url), coordinate.Offset())
	return
}

func (db *DB) UrlOfCoordinate(start Coordinate, end Coordinate) (url string, err error) {
	offset1 := start.Offset()
	offset2 := end.Offset()
	buf := make([]byte, offset2-offset1)
	if _, err = db.file.ReadAt(buf, offset1); err != nil {
		return
	}
	url = string(buf)
	return
}

func NewDB() (db *DB, err error) {
	var file *os.File
	if file, err = os.OpenFile("./data/db", os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	db = &DB{file: file}
	return
}
