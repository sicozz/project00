package datatype

import (
	"net"

	"github.com/google/uuid"
)

type ProgramInfo struct {
	Version string
	Banner  string
}

type Host struct {
	id uuid.UUID
	ip net.IP
}

func NewHost(id uuid.UUID, ip net.IP) Host {
	return Host{id, ip}
}

func (h *Host) Id() uuid.UUID {
	return h.id
}

func (h *Host) Ip() net.IP {
	return h.ip
}
