package app

import (
	"errors"
	"fmt"
	"net"

	"github.com/pugkong/ylc/yeelight"
)

type Control struct {
	store   *BulbFileStore
	printer Printer
}

func NewControl(store *BulbFileStore, printer Printer) *Control {
	return &Control{store: store, printer: printer}
}

func (c *Control) Info(name string) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	info, err := yeelight.NewController(conn).Info()
	if err != nil {
		return fmt.Errorf("query %q bulb info: %w", name, err)
	}

	c.printInfo(info)

	return nil
}

func (c *Control) printInfo(info yeelight.BulbInfo) {
	c.printer.Printf("Power: %s\n", info.Power)
	c.printer.Printf("Bright: %s\n", info.Bright)

	switch info.ColorMode {
	case yeelight.ColorModeRGB:
		c.printer.Println("Color mode: RGB")

		rgb, err := info.RGB.Int()
		if err != nil {
			c.printer.Printf("Invalid RGB value: %s", info.RGB)
		} else {
			c.printer.Printf("RGB: %06x\n", rgb)
		}
	case yeelight.ColorModeTemperature:
		c.printer.Println("Color mode: temperature")
		c.printer.Printf("Color temperature: %s\n", info.ColorTemperature)
	case yeelight.ColorModeHSV:
		c.printer.Println("Color mode: HSV")
		c.printer.Printf("HUE: %s\n", info.HUE)
		c.printer.Printf("Saturation: %s\n", info.Saturation)
	}

	c.printer.Println()

	c.printer.Printf("Background power: %s\n", info.BackgroundPower)
	c.printer.Printf("Background bright: %s\n", info.BackgroundBright)

	switch info.BackgroundColorMode {
	case yeelight.ColorModeRGB:
		c.printer.Println("Background color mode: RGB")

		rgb, err := info.BackgroundRGB.Int()
		if err != nil {
			c.printer.Printf("Invalid background RGB value: %s", info.BackgroundRGB)
		} else {
			c.printer.Printf("Background RGB: %06x\n", rgb)
		}
	case yeelight.ColorModeTemperature:
		c.printer.Println("Background color mode: temperature")
		c.printer.Printf("Background color temperature: %s\n", info.BackgroundColorTemperature)
	case yeelight.ColorModeHSV:
		c.printer.Println("Background color mode: HSV")
		c.printer.Printf("Background HUE: %s\n", info.BackgroundHUE)
		c.printer.Printf("Background saturation: %s\n", info.BackgroundSaturation)
	}
}

func (c *Control) PowerToggle(name string) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).PowerToggle(); err != nil {
		return fmt.Errorf("toggle %q bulb power: %w", name, err)
	}

	return nil
}

func (c *Control) BackgroundToggle(name string) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).BackgroundToggle(); err != nil {
		return fmt.Errorf("toggle %q bulb background power: %w", name, err)
	}

	return nil
}

func (c *Control) SetBright(name string, value int, effect yeelight.Effect, duration int) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).Bright(value, effect, duration); err != nil {
		return fmt.Errorf("set %q bulb bright: %w", name, err)
	}

	return nil
}

func (c *Control) SetBackgroundBright(name string, value int, effect yeelight.Effect, duration int) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).BackgroundBright(value, effect, duration); err != nil {
		return fmt.Errorf("set %q bulb background bright: %w", name, err)
	}

	return nil
}

func (c *Control) SetTemperature(name string, value int, effect yeelight.Effect, duration int) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).ColorTemperature(value, effect, duration); err != nil {
		return fmt.Errorf("set %q bulb temperature: %w", name, err)
	}

	return nil
}

func (c *Control) SetBackgroundTemperature(name string, value int, effect yeelight.Effect, duration int) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).BackgroundColorTemperature(value, effect, duration); err != nil {
		return fmt.Errorf("set %q bulb background temperature: %w", name, err)
	}

	return nil
}

func (c *Control) SetRGB(name string, value int, effect yeelight.Effect, duration int) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).RGB(value, effect, duration); err != nil {
		return fmt.Errorf("set %q bulb rgb color: %w", name, err)
	}

	return nil
}

func (c *Control) SetBackgroundRGB(name string, value int, effect yeelight.Effect, duration int) (err error) {
	conn, connClose, err := c.connectByName(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, connClose()) }()

	if err := yeelight.NewController(conn).BackgroundRGB(value, effect, duration); err != nil {
		return fmt.Errorf("set %q bulb background rgb color: %w", name, err)
	}

	return nil
}

func (c *Control) connectByName(name string) (*net.TCPConn, func() error, error) {
	bulb, err := c.store.FindByName(name)
	if err != nil {
		return nil, nil, fmt.Errorf("find %q bulb: %w", name, err)
	}

	addr, err := net.ResolveTCPAddr("tcp", bulb.Addr)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve addr for %q bulb: %w", name, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to %q bulb: %w", name, err)
	}

	connClose := func() error {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("close connection to %q bulb: %w", name, err)
		}

		return nil
	}

	return conn, connClose, nil
}
