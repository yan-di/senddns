package packets

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"io"
)

//ReadPcapFile Read the DNS request data in the pcap file and convert it to a DNS UDP request message.
func ReadPcapFile(protocol, filePath string) (pkgList [][]byte, err error) {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		}
		if packet.TransportLayer() != nil && packet.TransportLayer().TransportFlow().Dst().String() == "53" {
			if packet.ApplicationLayer() != nil {
				if packet.TransportLayer().LayerType().String() == "TCP" && protocol == "udp" {
					pkgList = append(pkgList, packet.ApplicationLayer().LayerContents()[2:])
				} else if packet.TransportLayer().LayerType().String() == "UDP" && protocol == "udp" {
					pkgList = append(pkgList, packet.ApplicationLayer().LayerContents())
				} else if packet.TransportLayer().LayerType().String() == "TCP" && protocol == "tcp" {
					pkgList = append(pkgList, packet.ApplicationLayer().LayerContents())
				} else if packet.TransportLayer().LayerType().String() == "UDP" && protocol == "tcp" {
					pkgList = append(pkgList, append([]byte{byte(len(packet.ApplicationLayer().LayerContents())) >> 8, byte(len(packet.ApplicationLayer().LayerContents()))}, packet.ApplicationLayer().LayerContents()...))
				}
			}
		}
	}
	fmt.Printf("Get %d DNS request messages !\n", len(pkgList))
	return
}
