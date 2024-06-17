package yeelight

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Bulb struct {
	ID   string
	Addr string
}

type Discoverer struct {
	conn net.PacketConn
}

func NewDiscoverer(conn net.PacketConn) *Discoverer {
	return &Discoverer{conn: conn}
}

var (
	discoverAddr = &net.UDPAddr{
		IP:   net.IPv4(239, 255, 255, 250),
		Port: 1982,
	}
	discoverMsg = []byte(strings.Join([]string{
		"M-SEARCH * HTTP/1.1",
		"MAN: \"ssdp:discover\"",
		"ST: wifi_bulb",
	}, "\r\n"))
)

func (d *Discoverer) SendDiscover() error {
	if _, err := d.conn.WriteTo(discoverMsg, discoverAddr); err != nil {
		return fmt.Errorf("send discover message: %w", err)
	}

	return nil
}

var (
	bulbIDPrefix       = []byte("id: ")
	bulbLocationPrefix = []byte("Location: yeelight://")
)

func (d *Discoverer) ReadBulb() (Bulb, error) {
	bulb := Bulb{}

	buffer := make([]byte, 4096)
	n, _, err := d.conn.ReadFrom(buffer)
	if err != nil {
		return bulb, fmt.Errorf("read bulb response: %w", err)
	}

	for _, line := range bytes.Split(buffer[:n], []byte("\r\n")) {
		if len(line) == 0 {
			continue
		}

		switch {
		case bytes.HasPrefix(line, bulbIDPrefix):
			bulb.ID = string(bytes.TrimPrefix(line, bulbIDPrefix))
		case bytes.HasPrefix(line, bulbLocationPrefix):
			bulb.Addr = string(bytes.TrimPrefix(line, bulbLocationPrefix))
		}
	}

	return bulb, nil
}
