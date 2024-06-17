package yeelight

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type colorMode string

const (
	ColorModeRGB         = "1"
	ColorModeTemperature = "2"
	ColorModeHSV         = "3"
)

type rgb string

func (r rgb) Int() (int, error) {
	v, err := strconv.Atoi(string(r))
	if err != nil {
		return 0, fmt.Errorf("parse %q color: %w", r, err)
	}

	return v, nil
}

type BulbInfo struct {
	Power            string
	Bright           string
	ColorMode        colorMode
	ColorTemperature string
	RGB              rgb
	HUE              string
	Saturation       string

	BackgroundPower            string
	BackgroundBright           string
	BackgroundColorMode        colorMode
	BackgroundColorTemperature string
	BackgroundRGB              rgb
	BackgroundHUE              string
	BackgroundSaturation       string
}

type TCPConn interface {
	Write(b []byte) (int, error)
	Read(b []byte) (int, error)
}

type Controller struct {
	conn          TCPConn
	nextCommandID int
}

func NewController(conn TCPConn) *Controller {
	return &Controller{conn: conn, nextCommandID: 1}
}

func (c *Controller) Info() (BulbInfo, error) {
	result, err := c.sendCommand(command{
		Method: "get_prop",
		Params: []any{
			"power",
			"bright",
			"color_mode",
			"ct",
			"rgb",
			"hue",
			"sat",

			"bg_power",
			"bg_bright",
			"bg_lmode",
			"bg_ct",
			"bg_rgb",
			"bg_hue",
			"bg_sat",
		},
	})
	if err != nil {
		return BulbInfo{}, err
	}

	info := BulbInfo{
		Power:            result[0],
		Bright:           result[1],
		ColorMode:        colorMode(result[2]),
		ColorTemperature: result[3],
		RGB:              rgb(result[4]),
		HUE:              result[5],
		Saturation:       result[6],

		BackgroundPower:            result[7],
		BackgroundBright:           result[8],
		BackgroundColorMode:        colorMode(result[9]),
		BackgroundColorTemperature: result[10],
		BackgroundRGB:              rgb(result[11]),
		BackgroundHUE:              result[12],
		BackgroundSaturation:       result[13],
	}

	return info, nil
}

func (c *Controller) PowerToggle() error {
	_, err := c.sendCommand(command{Method: "dev_toggle", Params: []any{}})

	return err
}

func (c *Controller) BackgroundToggle() error {
	_, err := c.sendCommand(command{Method: "bg_toggle", Params: []any{}})

	return err
}

func (c *Controller) Bright(value int, effect Effect, duration int) error {
	_, err := c.sendCommand(command{Method: "set_bright", Params: []any{value, effect, duration}})

	return err
}

func (c *Controller) BackgroundBright(value int, effect Effect, duration int) error {
	_, err := c.sendCommand(command{Method: "bg_set_bright", Params: []any{value, effect, duration}})

	return err
}

type Effect string

var (
	EffectSudden Effect = "sudden"
	EffectSmooth Effect = "smooth"
)

func (c *Controller) ColorTemperature(value int, effect Effect, duration int) error {
	_, err := c.sendCommand(command{Method: "set_ct_abx", Params: []any{value, effect, duration}})

	return err
}

func (c *Controller) BackgroundColorTemperature(value int, effect Effect, duration int) error {
	_, err := c.sendCommand(command{Method: "bg_set_ct_abx", Params: []any{value, effect, duration}})

	return err
}

func (c *Controller) RGB(value int, effect Effect, duration int) error {
	_, err := c.sendCommand(command{Method: "set_rgb", Params: []any{value, effect, duration}})

	return err
}

func (c *Controller) BackgroundRGB(value int, effect Effect, duration int) error {
	_, err := c.sendCommand(command{Method: "bg_set_rgb", Params: []any{value, effect, duration}})

	return err
}

type command struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

type result struct {
	ID     int      `json:"id"`
	Result []string `json:"result"`
	Error  struct {
		Message string `json:"message"`
	} `json:"error"`
}

var (
	ErrResponseTooLong = errors.New("response is too long")
	ErrBulbRepsponse   = errors.New("bulb error")
)

func (c *Controller) sendCommand(command command) ([]string, error) {
	command.ID = c.nextCommandID
	c.nextCommandID++

	data, err := json.Marshal(command)
	if err != nil {
		return nil, fmt.Errorf("prepare command: %w", err)
	}

	data = append(data, '\r', '\n')
	if _, err := c.conn.Write(data); err != nil {
		return nil, fmt.Errorf("send command: %w", err)
	}

	reader := bufio.NewReader(c.conn)
	for {
		line, prefix, err := reader.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}

		if prefix {
			return nil, ErrResponseTooLong
		}

		var result result
		if err := json.Unmarshal(line, &result); err != nil {
			return nil, fmt.Errorf("parse response %q: %w", string(line), err)
		}

		if result.ID != command.ID {
			continue
		}

		if result.Error.Message == "" {
			return result.Result, nil
		}

		return nil, fmt.Errorf("%w: %v", ErrBulbRepsponse, result.Error.Message)
	}
}
