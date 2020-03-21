package packets

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

//ReadHexFile supports the hexadecimal byte stream of the Domain Name System in the dns request message exported by wireshark
func ReadHexFile(protocol, filePath string) (pkgList [][]byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	read := bufio.NewReader(file)
	for {
		str, err := read.ReadString('\n')
		if err != nil {
			break
		}
		byteHex, _ := hex.DecodeString(strings.TrimSuffix(str, "\n"))
		if protocol == "udp" {
			pkgList = append(pkgList, byteHex)
		} else {
			pkgList = append(pkgList, append([]byte{byte(len(byteHex) >> 8), byte(len(byteHex))}, byteHex...))
		}
	}
	fmt.Printf("Get %d DNS request messages !\n", len(pkgList))
	return
}
