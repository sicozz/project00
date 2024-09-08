package utils

import (
	"bufio"
	"errors"
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

func DiscoverHosts(hostsFile string) (map[uuid.UUID]Host, error) {
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

func GetLocalHostIp(hosts map[uuid.UUID]Host) (uuid.UUID, error) {
	ip, err := getLocalIp()
	if err != nil {
		return uuid.New(), err
	}
	for _, h := range hosts {
		if h.Ip() == ip {
			Debug(fmt.Sprintf("localhost: %v", h))
			return h.Id(), nil
		}
	}
	return uuid.New(), errors.New("Failed to get localhost ip")
}

func getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip.IsLoopback() || ip.To4() == nil {
			continue
		}

		return ip.String(), nil
	}
	return "", errors.New("Failed to get local ip")
}
