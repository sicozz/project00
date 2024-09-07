package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Host struct {
	Id uuid.UUID `json:"id"`
	Ip string    `json:"ip"`
}

func NewHost(id uuid.UUID, ip string) Host {
	return Host{id, ip}
}

// func (h *Host) Id() uuid.UUID {
// 	return h.id
// }
//
// func (h *Host) Ip() net.IP {
// 	return h.ip
// }

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
		hosts[host.Id] = host
	}
	return hosts, nil
}
