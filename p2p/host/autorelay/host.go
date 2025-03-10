package autorelay

import (
	"github.com/zbliujia/go-libp2p/core/host"
)

type AutoRelayHost struct {
	host.Host
	ar *AutoRelay
}

func (h *AutoRelayHost) Close() error {
	_ = h.ar.Close()
	return h.Host.Close()
}

func (h *AutoRelayHost) Start() {
	h.ar.Start()
}

func NewAutoRelayHost(h host.Host, ar *AutoRelay) *AutoRelayHost {
	return &AutoRelayHost{Host: h, ar: ar}
}
