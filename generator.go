package surl

import (
	"errors"
	"net/url"
	"sync"
)

// Generator is the main entry to store and retrieve `URLUnit`
type Generator struct {
	Mapping Mapping
	Store   Store
	Root    *url.URL
	lock    *sync.Mutex
}

// Shorten encodes a url to surl
func (g *Generator) Shorten(u *url.URL, ids ...int64) (surl *url.URL, err error) {
	if g.lock == nil {
		g.lock = &sync.Mutex{}
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	var id int64
	if len(ids) != 1 {
		if id, err = g.Store.NextID(); err != nil {
			return
		}
	} else {
		id = ids[0]
	}
	if err = g.Store.Insert(id, u); err != nil {
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
func (g *Generator) Parse(surl *url.URL) (u *url.URL, err error) {
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
	if u, err = g.Store.Query(id); err != nil {
		return
	}
	return
}
