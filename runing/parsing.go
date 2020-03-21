package runing

import (
	"errors"
	"flag"
	"net"
	"senddns/packets"
)

//Structure information of parameter configuration
type Config struct {
	Srcs             []string
	SrcsLen, RunTime int
	Config           BasicConfig
}

//Initialize parameter configuration
func Parse() (conf Config, err error) {
	var number, rdValue int
	var netStr, domain, hexFile, domainFile, pcapFile string
	flag.StringVar(&netStr, "s", "", "配置来源IP地址段，支持IPv6，可用地址集为指定地址段与本地生效IP地址的交集")
	flag.StringVar(&conf.Config.DstAddr, "d", "", "配置目的IP地址，支持IPv6")
	flag.StringVar(&conf.Config.Protocol, "p", "udp", "请求的网络协议，tcp或udp")
	flag.UintVar(&conf.Config.Speed, "speed", 0, "并发请求数，默认为0，即：最大限度的进行并发请求")
	flag.IntVar(&conf.RunTime, "time", 0, "请求时间，单位：秒，0为运行后直至输入stop退出运行")
	flag.StringVar(&domain, "domain", "", "指定请求域，不能与hexfile/domainfile/pcapfile中任一参数共用")
	flag.IntVar(&number, "number", 1000, "构造与指定域相关的请求脏数据，依赖domain参数")
	flag.StringVar(&hexFile, "hexfile", "", "十六进制请求数据的文件路径，不能与domain/domainfile/pcapfile中任一参数共用")
	flag.StringVar(&domainFile, "domainfile", "", "以请求记录与请求类型为一行的文件的路径，不能与domain/hexfile/pcapfile中任一参数共用")
	flag.IntVar(&rdValue, "rdvalue", 0, "RD标志位，0 或 1，依赖domainfile参数")
	flag.StringVar(&pcapFile, "pcapfile", "", "从网络数据包文件中提取DNS请求，不能与domain/domainfile/hexfile中任一参数共用")
	flag.Parse()

	_, ipNet, err := net.ParseCIDR(netStr)
	if err != nil {
		return
	}
	addrs, _ := net.InterfaceAddrs()
	for _, i := range addrs {
		if addr, ok := i.(*net.IPNet); ok && ipNet.Contains(addr.IP) {
			conf.Srcs = append(conf.Srcs, addr.IP.String())
		}
	}

	if conf.SrcsLen = len(conf.Srcs); conf.SrcsLen == 0 {
		err = errors.New("Local IP address is not in the network segment, please re-specify")
		return
	}

	if ipAddr := net.ParseIP(conf.Config.DstAddr); ipAddr == nil {
		err = errors.New("Invalid destination IP address")
		return
	} else if ipAddr.To4() == nil {
		conf.Config.DstAddr = "[" + conf.Config.DstAddr + "]"
	}

	if conf.Config.Protocol != "udp" && conf.Config.Protocol != "tcp" {
		err = errors.New("Protocol setting error")
		return
	}

	if domain != "" && number != 0 {
		conf.Config.PkgList, err = packets.CreateMessage(conf.Config.Protocol, domain, number)
	} else if hexFile != "" {
		conf.Config.PkgList, err = packets.ReadHexFile(conf.Config.Protocol, hexFile)
	} else if domainFile != "" && (rdValue == 0 || rdValue == 1) {
		conf.Config.PkgList, err = packets.ReadDomainFile(conf.Config.Protocol, domainFile, uint16(rdValue))
	} else if pcapFile != "" {
		conf.Config.PkgList, err = packets.ReadPcapFile(conf.Config.Protocol, pcapFile)
	} else {
		err = errors.New("No specified run mode or configuration error")
	}

	return
}
