package lead

import (
	"net"
	"strings"
	"time"
)

const (
	discoveryPort   = 48899
	discoveryReq    = "HF-A11ASSISTHREAD"
	discoveryProbes = 5
	discoveryIntv   = 125 * time.Millisecond
	serverPort      = "8899"
)

func Discover(network string) ([]*Controller, error) {
	_, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: discoveryPort})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	bc := bcast(ipnet)

	go func() {
		for i := 0; i < discoveryProbes; i++ {
			conn.WriteTo([]byte(discoveryReq), &net.UDPAddr{IP: bc.IP, Port: discoveryPort})
			time.Sleep(discoveryIntv)
		}
	}()

	t0 := time.Now()
	buf := make([]byte, 128)
	res := make(map[string]struct{})
	for time.Since(t0) < discoveryProbes*discoveryIntv {
		conn.SetReadDeadline(time.Now().Add(discoveryIntv))
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}

		ret := string(buf[:n])
		if ret == discoveryReq {
			continue
		}
		res[ret] = struct{}{}
	}

	ctrls := make([]*Controller, 0, len(res))
	for r := range res {
		fields := strings.Split(r, ",")
		if len(fields) != 3 {
			continue
		}
		ctrls = append(ctrls, &Controller{
			Address: net.JoinHostPort(fields[0], serverPort),
			Serial:  fields[1],
			Model:   fields[2],
		})
	}
	return ctrls, nil
}

func bcast(ip *net.IPNet) *net.IPNet {
	var bc = &net.IPNet{}
	bc.IP = make([]byte, len(ip.IP))
	copy(bc.IP, ip.IP)
	bc.Mask = ip.Mask

	offset := len(bc.IP) - len(bc.Mask)
	for i := range bc.IP {
		if i-offset >= 0 {
			bc.IP[i] = ip.IP[i] | ^ip.Mask[i-offset]
		}
	}
	return bc
}
