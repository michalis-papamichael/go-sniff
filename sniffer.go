package gosniff

import (
	"time"

	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Sniffer struct {
	InterfaceName  *string
	BpfFilterExpr  *string
	SnapshotLength int32
	Duration       time.Duration
	Promiscuous    bool
	handle         *pcap.Handle
}

func (s *Sniffer) StartSniff() (chan gopacket.Packet, error) {
	iname, err := s.getInterfaceName()
	if err != nil {
		return nil, err
	}
	handle, err := pcap.OpenLive(*iname, s.SnapshotLength, s.Promiscuous, s.Duration)
	if err != nil {
		return nil, err
	}
	s.handle = handle
	if s.BpfFilterExpr != nil {
		if err := handle.SetBPFFilter(*s.BpfFilterExpr); err != nil {
			return nil, err
		}
	}
	pktChan := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	return pktChan, nil
}

func (s *Sniffer) SetNewBpfFilter(expr string) error {
	if err := s.handle.SetBPFFilter(*s.BpfFilterExpr); err != nil {
		return err
	}
	return nil
}

func (s *Sniffer) CloseAndGetStats(getStats bool) (*pcap.Stats, error) {
	defer s.handle.Close()
	if getStats {
		stat, err := s.handle.Stats()
		if err != nil {
			return stat, err
		}
		return stat, nil
	}
	return nil, nil
}

func (s *Sniffer) getInterfaceName() (*string, error) {
	if s.InterfaceName == nil {
		Interface, err := GetPhysicalInterface()
		if err != nil {
			return nil, err
		}
		return &Interface.Name, nil
	} else {
		return s.InterfaceName, nil
	}
}
