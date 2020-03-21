package main

import (
	"fmt"
	"runtime"
	"senddns/runing"
	"sync"
	"time"
)

func main() {
	conf, err := runing.Parse()
	if err != nil {
		println(err.Error())
		return
	}

	var confirm, stop string
	fmt.Printf("Are you sure you want to send DNS messages ? [default: yes] ")
	fmt.Scanln(&confirm)
	if confirm == "no" {
		println("Back out")
		return
	}

	cpuNum, send, wg := runtime.NumCPU(), &runing.SendInfo{Config: conf.Config}, new(sync.WaitGroup)
	send.Limit = make(chan bool, send.Config.Speed)
	runtime.GOMAXPROCS(cpuNum)
	if send.Config.Speed > uint(0) {
		wg.Add(1)
		go send.ControlSend(wg)
		cpuNum--
	} else {
		close(send.Limit)
	}
	fmt.Printf("Will use %d CPUs to send packages !\n", cpuNum)

	if conf.SrcsLen >= cpuNum {
		send.Count = make(chan uint64, conf.SrcsLen)
		for _, srcAddr := range conf.Srcs {
			wg.Add(1)
			go send.SendPacket(wg, srcAddr)
		}
	} else {
		send.Count = make(chan uint64, cpuNum)
		for n := 0; n < cpuNum; {
			for _, srcAddr := range conf.Srcs {
				wg.Add(1)
				go send.SendPacket(wg, srcAddr)
				n++
				if n == cpuNum {
					break
				}
			}
		}
	}

	println("Start sending packets !!!")
	if conf.RunTime == 0 {
	waitStop:
		fmt.Printf("Input <stop> to stop program running: ")
		fmt.Scanln(&stop)
		if stop != "stop" {
			goto waitStop
		}
	} else {
		time.Sleep(time.Second * time.Duration(conf.RunTime))
	}

	send.Exit = true
	wg.Wait()
	close(send.Count)
	var sendNum uint64

	for num := range send.Count {
		sendNum += num
	}

	fmt.Printf("End sending packets !!! Send %d requests in total !\n", sendNum)
}
