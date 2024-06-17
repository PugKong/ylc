package app

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"time"

	"github.com/pugkong/ylc/pokemon"
	"github.com/pugkong/ylc/yeelight"
)

type Manager struct {
	store   *BulbFileStore
	names   *pokemon.Names
	printer Printer
}

func NewManager(store *BulbFileStore, names *pokemon.Names, printer Printer) *Manager {
	return &Manager{
		store:   store,
		names:   names,
		printer: printer,
	}
}

func (m *Manager) List() error {
	m.printList(m.store.All())

	return nil
}

func (m *Manager) printList(bulbs []Bulb) {
	const fmt = " %12s %18s %18s\n"

	m.printer.Printf(fmt, "Name", "Address", "ID")
	for _, bulb := range bulbs {
		m.printer.Printf(fmt, bulb.Name, bulb.Addr, bulb.ID)
	}
}

func (m *Manager) Discover(listen string, duration time.Duration) error {
	conn, err := net.ListenPacket("udp", listen)
	if err != nil {
		return fmt.Errorf("listen %q udp: %w", listen, err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(duration)); err != nil {
		return fmt.Errorf("set timeout: %w", err)
	}

	rawBulbs, err := m.discoverBulbs(yeelight.NewDiscoverer(conn))
	if err != nil {
		return err
	}

	for _, name := range m.store.AllNames() {
		m.names.Occupy(name)
	}

	bulbs := make([]Bulb, 0, len(rawBulbs))
	for _, rawBulb := range rawBulbs {
		bulb, err := m.saveBulb(rawBulb)
		if err != nil {
			return err
		}

		bulbs = append(bulbs, bulb)
	}

	if err := m.store.Flush(); err != nil {
		return err
	}

	m.printList(bulbs)

	return nil
}

func (m *Manager) discoverBulbs(discoverer *yeelight.Discoverer) ([]yeelight.Bulb, error) {
	if err := discoverer.SendDiscover(); err != nil {
		return nil, fmt.Errorf("discover: %w", err)
	}

	var bulbs []yeelight.Bulb
	for {
		bulb, err := discoverer.ReadBulb()
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break
			}

			return nil, fmt.Errorf("discover: %w", err)
		}

		if i := slices.IndexFunc(bulbs, func(b yeelight.Bulb) bool { return b.ID == bulb.ID }); i == -1 {
			bulbs = append(bulbs, bulb)
		} else {
			bulbs[i] = bulb
		}
	}

	return bulbs, nil
}

func (m *Manager) saveBulb(rawBulb yeelight.Bulb) (Bulb, error) {
	bulb, err := m.store.FindByID(rawBulb.ID)
	switch {
	case errors.Is(err, ErrBulbNotFound):
		bulb, err = m.makeBulb(rawBulb)
		if err != nil {
			return Bulb{}, err
		}
	case err != nil:
		return Bulb{}, err
	default:
		bulb.Addr = rawBulb.Addr
	}

	m.store.Save(bulb)

	return bulb, nil
}

func (m *Manager) makeBulb(rawBulb yeelight.Bulb) (Bulb, error) {
	name, err := m.names.Generate()
	if err != nil {
		return Bulb{}, fmt.Errorf("generate bulb name: %w", err)
	}

	bulb := Bulb{
		ID:   rawBulb.ID,
		Name: name,
		Addr: rawBulb.Addr,
	}

	return bulb, nil
}

func (m *Manager) Delete(name string) error {
	bulb, err := m.store.FindByName(name)
	if err != nil {
		return err
	}

	m.store.Delete(bulb)

	if err := m.store.Flush(); err != nil {
		return err
	}

	return nil
}
