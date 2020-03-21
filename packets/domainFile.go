package packets

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

//ReadDomainFile is a text file that reads the domain name and request type
func ReadDomainFile(protocol, filePath string, rd uint16) (pkgList [][]byte, err error) {
	qtype := map[string]uint16{"A": 1, "NS": 2, "MD": 3, "MF": 4, "CNAME": 5, "SOA": 6, "MB": 7, "MG": 8, "MR": 9, "NULL": 10, "WKS": 11, "PTR": 12, "HINFO": 13, "MINFO": 14, "MX": 15, "TXT": 16, "RP": 17, "AFSDB": 18, "X25": 19, "ISDN": 20, "RT": 21, "NSAP": 22, "NSAP-PTR": 23, "SIG": 24, "KEY": 25, "PX": 26, "GPOS": 27, "AAAA": 28, "LOC": 29, "NXT": 30, "EID": 31, "NIMLOC": 32, "SRV": 33, "ATMA": 34, "NAPTR": 35, "KX": 36, "CERT": 37, "A6": 38, "DNAME": 39, "SINK": 40, "OPT": 41, "APL": 42, "DS": 43, "SSHFP": 44, "IPSECKEY": 45, "RRSIG": 46, "NSEC": 47, "DNSKEY": 48, "DHCID": 49, "NSEC3": 50, "NSEC3PARAM": 51, "TLSA": 52, "SMIMEA": 53, "HIP": 55, "NINFO": 56, "RKEY": 57, "TALINK": 58, "CDS": 59, "CDNSKEY": 60, "OPENPGPKEY": 61, "CSYNC": 62, "ZONEMD": 63, "SPF": 99, "UINFO": 100, "UID": 101, "GID": 102, "UNSPEC": 103, "NID": 104, "L32": 105, "L64": 106, "LP": 107, "EUI48": 108, "EUI64": 109, "TKEY": 249, "TSIG": 250, "IXFR": 251, "AXFR": 252, "MAILB": 253, "MAILA": 254, "ANY": 255, "URI": 256, "CAA": 257, "AVC": 258, "DOA": 259, "AMTRELAY": 260, "TA": 32768, "DLV": 32769}

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

		buffer := new(bytes.Buffer)

		//rand ID
		binary.Write(buffer, binary.BigEndian, uint16(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(65536)))

		// QR: 1bit; opcode: 4bit; AA: 1bit; TC: 1bit; RD: 1bit; RA: 1bit; zero: 3bit; rcode: 4bit
		flag := uint16(0)<<15 + uint16(0)<<11 + uint16(0)<<10 + uint16(0)<<9 + rd<<8 + uint16(0)<<7 + uint16(0)<<4 + uint16(0)
		binary.Write(buffer, binary.BigEndian, flag)

		//QDCOUNT: 16bit; ANCOUNT: 16bit; NSCOUNT: 16bit; ARCOUNT: 16bit
		binary.Write(buffer, binary.BigEndian, uint64(1)<<48)

		//Domain Name
		for _, i := range strings.Split(strings.Fields(str)[0], ".") {
			var domainStr string
			if i == "*" {
				domainStr = randStr(int(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(17)))
			} else {
				domainStr = i
			}
			binary.Write(buffer, binary.BigEndian, byte(len(domainStr)))
			binary.Write(buffer, binary.BigEndian, []byte(domainStr))
		}
		binary.Write(buffer, binary.BigEndian, byte(0x00))

		//Question Type
		binary.Write(buffer, binary.BigEndian, qtype[strings.Fields(str)[1]])

		//Question Class
		binary.Write(buffer, binary.BigEndian, uint16(1))

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
