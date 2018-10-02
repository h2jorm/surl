package surl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"path"
	"strconv"
)

// 0    1    2    3    4    5    6    7    0    1    2    3    4    5    6    7
// +----+----+----+----+----+----+----+----+----+----+----+----+----+----+----+
// +          partition          +                   length                   +
// +                           position                                       +
// +                           position                                       +

type coordinate struct {
	partition uint8
	len       uint16
	pos       uint32
}

func (c *coordinate) Bytes() (ret []byte, err error) {
	buf := new(bytes.Buffer)
	head := uint16(c.partition)<<10 + c.len<<6>>6
	if err = binary.Write(buf, binary.BigEndian, head); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, c.pos); err != nil {
		return
	}
	ret = buf.Bytes()
	return
}

func (c coordinate) offset() int64 {
	return int64(c.partition)*int64(math.Pow(2, 32)) + int64(c.pos)
}

func (c coordinate) String() string {
	return fmt.Sprintf("partition: %d, pos: %d, len: %d", c.partition, c.pos, c.len)
}

func newCoordinate(bs [6]byte) (coor coordinate, err error) {
	head := binary.BigEndian.Uint16(bs[0:2])
	partition := uint8(head >> 10)
	len := head << 6 >> 6
	pos := binary.BigEndian.Uint32(bs[2:6])
	coor = coordinate{partition: partition, len: len, pos: pos}
	return
}

const (
	indexSize int64 = 6
)

type BasicStorage struct {
	IndexFile *os.File
	DBFile    *os.File
}

func (s *BasicStorage) NextID() (id int64, err error) {
	var stat os.FileInfo
	if stat, err = s.IndexFile.Stat(); err != nil {
		return
	}
	id = stat.Size() / indexSize
	return
}

func (s *BasicStorage) Insert(id int64, u *url.URL) (err error) {
	var occupied bool
	if occupied, err = s.idOccupied(id); occupied || err != nil {
		err = errors.New("id " + strconv.FormatInt(id, 10) + " is occupied")
		return
	}
	var coor coordinate
	if coor, err = s.nextCoordinate(); err != nil {
		return
	}
	coor.len = uint16(len(u.String()))
	if err = s.writeCoordinateAtID(coor, id); err != nil {
		return
	}
	if err = s.writeUrlAtCoordinate(u, coor); err != nil {
		return
	}
	return
}

func (s *BasicStorage) Query(id int64) (u *url.URL, err error) {
	var coor coordinate
	if coor, err = s.coordinateOfID(id); err != nil {
		return
	}
	u, err = s.urlOfCoordinate(coor)
	return
}

func (store *BasicStorage) nextCoordinate() (coor coordinate, err error) {
	var stat os.FileInfo
	if stat, err = store.DBFile.Stat(); err != nil {
		return
	}
	partition := stat.Size() / int64(math.Pow(2, 32))
	pos := stat.Size() % int64(math.Pow(2, 32))
	coor = coordinate{partition: uint8(partition), pos: uint32(pos)}
	return
}

func (store *BasicStorage) writeUrlAtCoordinate(u *url.URL, coor coordinate) (err error) {
	_, err = store.DBFile.WriteAt([]byte(u.String()), coor.offset())
	return
}

func (store *BasicStorage) urlOfCoordinate(coor coordinate) (u *url.URL, err error) {
	len := coor.len
	offset := coor.offset()
	buf := make([]byte, len)
	if _, err = store.DBFile.ReadAt(buf, offset); err != nil {
		return
	}
	u, err = url.Parse(string(buf))
	return
}

func (store *BasicStorage) idOccupied(id int64) (occupied bool, err error) {
	var nextID int64
	if nextID, err = store.NextID(); err != nil {
		return
	}
	if id >= nextID {
		occupied = false
		return
	}
	buf := make([]byte, 6, 6)
	if _, err = store.IndexFile.ReadAt(buf, id*indexSize); err != nil {
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

func (store *BasicStorage) writeCoordinateAtID(coor coordinate, id int64) (err error) {
	var bytes []byte
	if bytes, err = coor.Bytes(); err != nil {
		return
	}
	offset := id * indexSize
	_, err = store.IndexFile.WriteAt(bytes, offset)
	return
}

func (store *BasicStorage) coordinateOfID(id int64) (coor coordinate, err error) {
	buf := make([]byte, 6, 6)
	if _, err = store.IndexFile.ReadAt(buf, id*indexSize); err != nil {
		return
	}
	var bytes [6]byte
	copy(bytes[:], buf)
	if coor, err = newCoordinate(bytes); err != nil {
		return
	}
	return
}

func NewBasicStore(dirname string) (store *BasicStorage, err error) {
	var indexFile, dbFile *os.File
	if indexFile, err = os.OpenFile(path.Join(dirname, "_.index"), os.O_CREATE|os.O_RDWR, 0700); err != nil {
		log.Fatal(err)
	}
	if dbFile, err = os.OpenFile(path.Join(dirname, "_.db"), os.O_CREATE|os.O_RDWR, 0700); err != nil {
		log.Fatal(err)
	}
	store = &BasicStorage{
		IndexFile: indexFile,
		DBFile:    dbFile,
	}
	return
}
