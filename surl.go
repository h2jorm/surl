package main

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type SURL struct {
	index *Index
	db    *DB
}

func (surl *SURL) Create(id int64, url string) (err error) {
	var occupied bool
	if occupied, err = surl.index.IDOccupied(id); occupied || err != nil {
		err = errors.New("id " + strconv.FormatInt(id, 10) + " is occupied")
		return
	}
	var coordinate Coordinate
	if coordinate, err = surl.db.CurrentCoordinate(); err != nil {
		return
	}
	coordinate.len = uint16(len(url))
	if err = surl.index.WriteAt(coordinate, id); err != nil {
		return
	}
	if err = surl.db.WriteAt(url, coordinate); err != nil {
		return
	}
	log.Printf("ID: %d, Coordinate: %s\n", id, coordinate.String())
	return
}

func (surl *SURL) Find(id int64) (url string, err error) {
	var coordinate Coordinate
	if coordinate, err = surl.index.CoordinateOfID(id); err != nil {
		return
	}
	log.Printf("Find ID: %d, Coordinate: %s\n", id, coordinate.String())
	if url, err = surl.db.UrlOfCoordinate(coordinate); err != nil {
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
		if len(os.Args) < 3 {
			log.Fatal("usage surl create <id> <url>")
		}
		var id int
		if id, err = strconv.Atoi(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		url := os.Args[3]
		if url == "" {
			log.Fatal("undefined url")
		}
		err = surl.Create(int64(id), url)
		if err != nil {
			log.Fatal(err)
		}
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
