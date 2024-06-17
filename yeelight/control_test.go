package yeelight

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TCPConnDummy struct {
	output []byte
}

func (t *TCPConnDummy) Write(data []byte) (int, error) {
	return len(data), nil
}

func (t *TCPConnDummy) Read(buffer []byte) (int, error) {
	return copy(buffer, t.output), nil
}

func TestController_sendCommand(t *testing.T) {
	t.Run("it handles multiple messages on single conn.read", func(t *testing.T) {
		output := "{\"method\":\"props\",\"params\":{\"bg_power\":\"off\",\"power\":\"off\"}}\r\n"
		output += "{\"id\":1,\"result\":[\"ok\"]}\r\n"

		conn := &TCPConnDummy{output: []byte(output)}
		controller := NewController(conn)
		result, err := controller.sendCommand(command{})

		require.NoError(t, err)
		require.Equal(t, []string{"ok"}, result)
	})
}
