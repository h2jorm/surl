package surl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// SURL is a operation collection of surl db and index file.
type SURL struct {
	index   *index
	store   *store
	Mapping Mapping
}

// URLUnit is a base unit of surl
type URLUnit struct {
	id  int64
	hex string
	url string
}

// ID returns id of urlUnit
func (unit URLUnit) ID() int64 {
	return unit.id
}

// Hex returns hex of urlUnit
func (unit URLUnit) Hex() string {
	return unit.hex
}

// URL returns url of urlUnit
func (unit URLUnit) URL() string {
	return unit.url
}

func (unit URLUnit) String() string {
	return fmt.Sprintf("id: %d, hex: %s, url: %s", unit.id, unit.hex, unit.url)
}

// Create creates a new record with id.
func (surl SURL) Create(id int64, url string) (unit URLUnit, err error) {
	var occupied bool
	if occupied, err = surl.index.idOccupied(id); occupied || err != nil {
		err = errors.New("id " + strconv.FormatInt(id, 10) + " is occupied")
		return
	}
	var coor coordinate
	if coor, err = surl.store.nextCoordinate(); err != nil {
		return
	}
	coor.len = uint16(len(url))
	if err = surl.index.writeAt(coor, id); err != nil {
		return
	}
	if err = surl.store.writeAt(url, coor); err != nil {
		return
	}
	unit = URLUnit{
		id:  id,
		hex: surl.Mapping.Itoa(id),
		url: url,
	}
	return
}

// Append records url in the next id position.
func (surl SURL) Append(url string) (URLUnit, error) {
	var id int64
	var err error
	if id, err = surl.index.nextID(); err != nil {
		return URLUnit{}, err
	}
	return surl.Create(id, url)
}

// Find retrives a record by id.
func (surl SURL) Find(hex string) (unit URLUnit, err error) {
	id := surl.Mapping.Atoi(hex)
	var coor coordinate
	if coor, err = surl.index.coordinateOfID(id); err != nil {
		return
	}
	log.Printf("Find ID: %d, coordinate: %s\n", id, coor.String())
	var url string
	if url, err = surl.store.urlOfCoordinate(coor); err != nil {
		return
	}
	unit = URLUnit{
		id:  id,
		hex: hex,
		url: url,
	}
	return
}

// NewSURL returns a new SURL struct. `dirpath` is the directory path to store data.
func NewSURL(dirpath string) (surl *SURL, err error) {
	var storeFile *os.File
	var indexFile *os.File
	if storeFile, err = os.OpenFile(filepath.Join(dirpath, "_.db"), os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	if indexFile, err = os.OpenFile(filepath.Join(dirpath, "_.index"), os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	surl = &SURL{
		index:   &index{file: indexFile},
		store:   &store{file: storeFile},
		Mapping: hex62{},
	}
	return
}
