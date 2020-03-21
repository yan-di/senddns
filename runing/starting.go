package runing

import (
	"net"
	"sync"
	"time"
)

const port = ":53"

//Basic configuration item
type BasicConfig struct {
	Speed             uint
	PkgList           [][]byte
	DstAddr, Protocol string
}

//Build connection information to send data
type SendInfo struct {
	Exit   bool
	Limit  chan bool
	Count  chan uint64
	Config BasicConfig
}

//Request system resources to start sending packets
func (s *SendInfo) SendPacket(wg *sync.WaitGroup, srcAddr string) {
	defer wg.Done()
	var num uint64
	var dial *net.Dialer

	switch s.Config.Protocol {
	case "tcp":
		dial = &net.Dialer{LocalAddr: &net.TCPAddr{IP: net.ParseIP(srcAddr), Port: 0}}
	case "udp":
		dial = &net.Dialer{LocalAddr: &net.UDPAddr{IP: net.ParseIP(srcAddr), Port: 0}}
	}

	conn, err := dial.Dial(s.Config.Protocol, s.Config.DstAddr+port)
	if err != nil {
		println(err.Error())
		return
	}
	defer conn.Close()

	for {
		for n := range s.Config.PkgList {
			<-s.Limit
			if _, err = conn.Write(s.Config.PkgList[n]); err != nil || s.Exit {
				goto EXIT
			}
			num++
		}
	}

EXIT:
	s.Count <- num
}

//Limit the number of packets sent by the program per second
func (s *SendInfo) ControlSend(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(s.Limit)
	for !s.Exit {
		for i := uint(1); i <= s.Config.Speed; i++ {
			s.Limit <- true
		}
		time.Sleep(time.Second)
	}
}
