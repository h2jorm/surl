package surl

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// SURL is a operation collection of surl db and index file.
type SURL struct {
	index *index
	db    *store
}

// Create creates a new record with id.
func (surl *SURL) Create(id int64, url string) (err error) {
	var occupied bool
	if occupied, err = surl.index.idOccupied(id); occupied || err != nil {
		err = errors.New("id " + strconv.FormatInt(id, 10) + " is occupied")
		return
	}
	var coor coordinate
	if coor, err = surl.db.nextCoordinate(); err != nil {
		return
	}
	coor.len = uint16(len(url))
	if err = surl.index.writeAt(coor, id); err != nil {
		return
	}
	if err = surl.db.writeAt(url, coor); err != nil {
		return
	}
	log.Printf("ID: %d, coordinate: %s\n", id, coor.String())
	return
}

// Find retrives a record by id.
func (surl *SURL) Find(id int64) (url string, err error) {
	var coor coordinate
	if coor, err = surl.index.coordinateOfID(id); err != nil {
		return
	}
	log.Printf("Find ID: %d, coordinate: %s\n", id, coor.String())
	if url, err = surl.db.urlOfCoordinate(coor); err != nil {
		return
	}
	return
}

// NewSURL returns a new SURL struct
func NewSURL(datapath string) (surl *SURL, err error) {
	var dbfile *os.File
	var indexfile *os.File
	if dbfile, err = os.OpenFile(filepath.Join(datapath, "_.db"), os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	if indexfile, err = os.OpenFile(filepath.Join(datapath, "_.index"), os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	surl = &SURL{
		index: &index{file: indexfile},
		db:    &store{file: dbfile},
	}
	return surl, nil
}
