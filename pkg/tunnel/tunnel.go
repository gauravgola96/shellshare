package tunnel

import (
	"io"
	"sync"
)

var Tunnel Tunnel_

func init() {
	tm := make(map[string]chan ConnectionTunnel)
	Tunnel = Tunnel_{tmap: tm}
}

type Tunnel_ struct {
	sync.RWMutex
	tmap map[string]chan ConnectionTunnel
}

func (t *Tunnel_) Store(key string, val chan ConnectionTunnel) {
	t.Lock()
	defer t.Unlock()
	t.tmap[key] = val
}

func (t *Tunnel_) Delete(key string) {
	t.Lock()
	defer t.Unlock()
	delete(t.tmap, key)
}

func (t *Tunnel_) Get(key string) (chan ConnectionTunnel, bool) {
	t.Lock()
	defer t.Unlock()
	val, ok := t.tmap[key]
	return val, ok
}

func (t *Tunnel_) GetWaitTunnel(key string) chan ConnectionTunnel {
	t.Lock()
	defer t.Unlock()
	return t.tmap[key]
}

// ConnectionTunnel ConnectionTunnel HTTP request and SSH session
type ConnectionTunnel struct {
	W    io.Writer
	Done chan struct{}
}
