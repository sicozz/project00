package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/google/uuid"
)

type Host struct {
	id uuid.UUID
	ip string
}

func NewHost(id uuid.UUID, ip string) Host {
	return Host{id, ip}
}

func (h *Host) Id() uuid.UUID {
	return h.id
}

func (h *Host) Ip() string {
	return h.ip
}

func HostDiscovery(hostsFile string) (map[uuid.UUID]Host, error) {
	hostsDb, err := os.Open(hostsFile)
	if err != nil {
		Error(fmt.Sprintf("Failed to open hosts file: %v", hostsFile))
		return nil, err
	}
	scanner := bufio.NewScanner(hostsDb)
	hosts := make(map[uuid.UUID]Host)
	for scanner.Scan() {
		ip := scanner.Text()
		host := NewHost(uuid.New(), ip)
		hosts[host.Id()] = host
	}
	return hosts, nil
}

func GetSelfHostId(hosts map[uuid.UUID]Host) (uuid.UUID, error) {
	// for _, h := range hosts {
	// 	// if h.Ip() ==
	// }
	getSelfIp()
	return uuid.New(), nil
}

func getSelfIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		Error(fmt.Sprintf("Failed to query net interfaces address: %v", err))
		return nil, err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue // Ignore loopback addresses (e.g., 127.0.0.1)
		}
		if ipNet.IP.To4() != nil {
			Debug(fmt.Sprintf("IPv4: %v", ipNet.IP.String()))
		} else {
			Debug(fmt.Sprintf("IPv6: %v", ipNet.IP.String()))
		}
	}
	return nil, nil
}
