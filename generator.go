package surl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Generator is the main entry to store and retrieve `URLUnit`
type Generator struct {
	index   *index
	store   *store
	mapping Mapping
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
func (gen Generator) Create(id int64, url string) (unit URLUnit, err error) {
	var occupied bool
	if occupied, err = gen.index.idOccupied(id); occupied || err != nil {
		err = errors.New("id " + strconv.FormatInt(id, 10) + " is occupied")
		return
	}
	var coor coordinate
	if coor, err = gen.store.nextCoordinate(); err != nil {
		return
	}
	coor.len = uint16(len(url))
	if err = gen.index.writeAt(coor, id); err != nil {
		return
	}
	if err = gen.store.writeAt(url, coor); err != nil {
		return
	}
	unit = URLUnit{
		id:  id,
		hex: gen.mapping.Itoa(id),
		url: url,
	}
	return
}

// Append records url in the next id position.
func (gen Generator) Append(url string) (URLUnit, error) {
	var id int64
	var err error
	if id, err = gen.index.nextID(); err != nil {
		return URLUnit{}, err
	}
	return gen.Create(id, url)
}

// Find retrives a record by id.
func (gen Generator) Find(hex string) (unit URLUnit, err error) {
	id := gen.mapping.Atoi(hex)
	var coor coordinate
	if coor, err = gen.index.coordinateOfID(id); err != nil {
		return
	}
	var url string
	if url, err = gen.store.urlOfCoordinate(coor); err != nil {
		return
	}
	unit = URLUnit{
		id:  id,
		hex: hex,
		url: url,
	}
	return
}

// NewSURL returns a new Generator struct. `dirpath` is the directory path to store data.
func NewSURL(dirpath string, mapping Mapping) (surl *Generator, err error) {
	var storeFile *os.File
	var indexFile *os.File
	if storeFile, err = os.OpenFile(filepath.Join(dirpath, "_.db"), os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	if indexFile, err = os.OpenFile(filepath.Join(dirpath, "_.index"), os.O_RDWR|os.O_CREATE, 0700); err != nil {
		return
	}
	surl = &Generator{
		index:   &index{file: indexFile},
		store:   &store{file: storeFile},
		mapping: mapping,
	}
	return
}
