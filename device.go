package adbutils

import (
	"fmt"
)

type BaseDevice struct {
	// TODO
	Client       adbClient
	Serial       string
	Transport_id int
}

func (device BaseDevice) open_transport(command string, timeout float32) string {
	// TODO connect has it own timeout
	c = device.client._connect()
	if timeout > 0 {
		c.conn.settimeout(timeout)
	}

	if command != "" {
		if self.transport_id != 0 {
			c.send_command(fmt.Sprintf("host:transport-id: %d: %s", device.transport_id, command))
		} else if self._serial != []string{} {
			c.send_command(fmt.Sprintf("host:serial: %s: %s", device.serial, command))
		} else {
			fmt.Error("RuntimeError")
		}
		c.check_okay()
	} else {
		if device.transport_id {
			c.send_command(fmt.Sprintf("host:transport-id: %d", device.transport_id))
		} else if device.serial != "" {
			// host:tport:serial:xxx is also fine, but receive 12 bytes
			// recv: 4f 4b 41 59 14 00 00 00 00 00 00 00              OKAY........
			// so here use host:transport
			c.send_command("host:transport: " + device.serial)
		} else {
			fmt.Error("RuntimeError")
		}
		c.check_okay()
	}
	return c
}

func (device BaseDevice) shell(cmdargs string, stream bool, timeout float32, rstrip bool) (string, string) {
	// TODO
	/**Run shell inside device and get it's content

	Args:
	rstrip (bool): strip the last empty line (Default: True)
	stream (bool): return stream instead of string output (Default: False)
	timeout (float): set shell timeout

	Returns:
	string of output when stream is False
	AdbConnection when stream is True

	Raises:
	AdbTimeout

	Examples:
	shell("ls -l")
	shell(["ls", "-l"])
	shell("ls | grep data")
	**/

	if stream {
		timeout = 0
	}
	c = device.open_transport("", timeout)
	c.send_command("shell:" + cmdargs)
	c.check_okay()
	if stream {
		return c
	}
	output = c.read_until_close()
	if rstrip {
		return output.rstrip()
	} else {
		return output
	}
}

type AdbDevice struct {
	// TODO
	BaseDevice
}
