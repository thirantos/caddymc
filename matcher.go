package caddymc

import (
	"errors"
	"io"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/mholt/caddy-l4/layer4"
)

func init() {
	caddy.RegisterModule(&MatchMinecraft{})
}

type MatchMinecraft struct {
	Hosts []string `json:"hosts,omitempty"`
}

func (m *MatchMinecraft) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.matchers.minecraft",
		New: func() caddy.Module { return new(MatchMinecraft) },
	}
}

func (m *MatchMinecraft) Provision(ctx caddy.Context) error {
	for i := range m.Hosts {
		m.Hosts[i] = strings.ToLower(m.Hosts[i])
	}

	return nil
}

func (m *MatchMinecraft) Match(cx *layer4.Connection) (bool, error) {
	handshake, err := ReadHandshakeFromByte(cx.MatchingBytes())

	switch {
	case err == nil:
	case errors.Is(err, io.EOF):
		return false, layer4.ErrConsumedAllPrefetchedBytes
	default:
		return false, nil
	}

	host := strings.ToLower(handshake.ServerAddress)
	for _, hostname := range m.Hosts {
		if host == hostname {
			return true, nil
		}
	}

	return false, nil
}

func (m *MatchMinecraft) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()

		if len(args) == 0 {
			return d.Err("expected at least one hostname")
		}

		m.Hosts = append(m.Hosts, args...)

		if d.NextBlock(0) {
			return d.Err("malformed minecraft matcher: blocks are not supported")
		}
	}

	return nil
}

var (
	_ caddy.Module          = (*MatchMinecraft)(nil)
	_ caddy.Provisioner     = (*MatchMinecraft)(nil)
	_ layer4.ConnMatcher    = (*MatchMinecraft)(nil)
	_ caddyfile.Unmarshaler = (*MatchMinecraft)(nil)
)
