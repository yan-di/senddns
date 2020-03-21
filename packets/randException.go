package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

func randStr(long int) string {
	rand.Seed(time.Now().UnixNano())
	oldStr := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(oldStr)
	newStr := []byte{}
	for i := 0; i < long; i++ {
		newStr = append(newStr, bytes[rand.Intn(len(bytes))])
	}
	return string(newStr)
}

//CreateMessage randomly generate a specified number of DNS exception messages
func CreateMessage(protocol, domain string, numberPackets int) (pkgList [][]byte, err error) {
	var domainStr string
	for i := 0; i < numberPackets; i++ {
		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, rand.New(rand.NewSource(time.Now().UnixNano())).Uint32())
		binary.Write(buffer, binary.BigEndian, rand.New(rand.NewSource(time.Now().UnixNano())).Uint64())
		for j := 0; j < int(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(17)); j++ {
			domainStr = randStr(int(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(17)))
			binary.Write(buffer, binary.BigEndian, byte(len(domainStr)))
			binary.Write(buffer, binary.BigEndian, []byte(domainStr))
		}
		binary.Write(buffer, binary.BigEndian, byte(0x00))
		binary.Write(buffer, binary.BigEndian, rand.New(rand.NewSource(time.Now().UnixNano())).Uint32())
		if protocol == "udp" {
			pkgList = append(pkgList, buffer.Bytes())
		} else {
			pkgList = append(pkgList, append([]byte{byte(buffer.Len() >> 8), byte(buffer.Len())}, buffer.Bytes()...))
		}
		buffer.Reset()
	}
	fmt.Printf("Get %d DNS request messages !\n", len(pkgList))
	return
}
