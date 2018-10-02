package surl

import (
	"errors"
	"net/url"
	"sync"
)

// SURL is the main entry to store and retrieve `URLUnit`
type SURL struct {
	Mapping Mapping
	Storage Storage
	Root    *url.URL
	lock    *sync.Mutex
}

// Shorten encodes a url to surl
func (g *SURL) Shorten(u *url.URL, ids ...int64) (surl *url.URL, err error) {
	if g.lock == nil {
		g.lock = &sync.Mutex{}
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	var id int64
	if len(ids) != 1 {
		if id, err = g.Storage.NextID(); err != nil {
			return
		}
	} else {
		id = ids[0]
	}
	if err = g.Storage.Insert(id, u); err != nil {
		return
	}
	surl = &url.URL{
		Scheme: g.Root.Scheme,
		Host:   g.Root.Host,
		Path:   g.Mapping.Itoa(id),
	}
	return
}

// Parse decodes a surl to its original url
func (g *SURL) Parse(surl *url.URL) (u *url.URL, err error) {
	if surl.Scheme != g.Root.Scheme || surl.Host != g.Root.Host {
		err = errors.New("invalid surl")
		return
	}
	path := surl.Path
	if path[0] == '/' {
		path = path[1:]
	}
	var id int64
	if id, err = g.Mapping.Atoi(path); err != nil {
		return
	}
	if u, err = g.Storage.Query(id); err != nil {
		return
	}
	return
}
