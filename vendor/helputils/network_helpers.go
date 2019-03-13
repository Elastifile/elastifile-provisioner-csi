package helputils

import "net"

// Get IP address used to access the outside network
func MyExternalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return []byte{0, 0, 0, 0}
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
