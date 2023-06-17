package tunnel

import (
	"io"
	"sync"
)

var Tunnel Tunnel_

func init() {
	tm := make(map[string]chan SSHTunnel)
	Tunnel = Tunnel_{tmap: tm}
}

type Tunnel_ struct {
	sync.RWMutex
	tmap map[string]chan SSHTunnel
}

func (t *Tunnel_) Store(key string, val chan SSHTunnel) {
	t.Lock()
	defer t.Unlock()
	t.tmap[key] = val
}

func (t *Tunnel_) Delete(key string) {
	t.Lock()
	defer t.Unlock()
	delete(t.tmap, key)
}

func (t *Tunnel_) Get(key string) (chan SSHTunnel, bool) {
	t.Lock()
	defer t.Unlock()
	val, ok := t.tmap[key]
	return val, ok
}

func (t *Tunnel_) GetWaitTunnel(key string) chan SSHTunnel {
	t.Lock()
	defer t.Unlock()
	return t.tmap[key]
}

// SSHTunnel SSHTunnel HTTP request and SSH session
type SSHTunnel struct {
	W    io.Writer
	Done chan struct{}
}
