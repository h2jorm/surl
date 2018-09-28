package main

import (
	"log"
	"os"
	"strconv"
)

type SURL struct {
	index *Index
	db    *DB
}

func (surl *SURL) Create(url string) (err error) {
	var id int64
	var coordinate Coordinate
	if id, err = surl.index.NextID(); err != nil {
		return
	}
	if coordinate, err = surl.db.CurrentCoordinate(); err != nil {
		return
	}
	nextCoordinate := coordinate
	nextCoordinate.Grow(len(url))
	if err = surl.index.WriteAt(nextCoordinate, id); err != nil {
		return
	}
	if err = surl.db.WriteAt(url, coordinate); err != nil {
		return
	}
	log.Printf("ID: %d, Coordinate: [(%s), (%s)]\n", id, coordinate.String(), nextCoordinate.String())
	return
}

func (surl *SURL) Find(id int64) (url string, err error) {
	var start, end Coordinate
	if start, end, err = surl.index.CoordinateOfID(id); err != nil {
		return
	}
	log.Printf("Find ID: %d, Coordinate:[(%s), (%s)]\n", id, start.String(), end.String())
	if url, err = surl.db.UrlOfCoordinate(start, end); err != nil {
		return
	}
	return
}

func NewSURL() (surl *SURL, err error) {
	var db *DB
	var index *Index
	if db, err = NewDB(); err != nil {
		return
	}
	if index, err = NewIndex(); err != nil {
		return
	}
	surl = &SURL{
		index: index,
		db:    db,
	}
	return surl, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	var err error
	var surl *SURL
	surl, err = NewSURL()

	if len(os.Args) < 2 {
		log.Fatal("usage surl <command>")
	}

	switch os.Args[1] {
	case "create":
		url := os.Args[2]
		if url == "" {
			log.Fatal("undefined url")
		}
		surl.Create(url)
	case "find":
		var id int
		var url string
		if id, err = strconv.Atoi(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		if url, err = surl.Find(int64(id)); err != nil {
			log.Fatal(err)
		}
		log.Printf("surl: /%s, url: %s\n", decimalToHexAny(int64(id)), url)
	default:
		log.Fatal("invalid command")
	}
}
