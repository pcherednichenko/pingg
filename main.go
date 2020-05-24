package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/sparrc/go-ping"
)

const (
	filename  = "graph.html"
	stadiaURL = "stadia.google.com"
	windowsOS = "windows"
)

func main() {
	pinger, err := ping.NewPinger(stadiaURL)
	if err != nil {
		panic(err)
	}
	if runtime.GOOS == windowsOS {
		pinger.SetPrivileged(true)
	}

	url := flag.String("routerIP", "", "router IP to ping")
	flag.Parse()

	if *url == "" {
		panic("empty routerIP parameter")
	}
	pinger2, err := ping.NewPinger(*url)
	if err != nil {
		panic(err)
	}
	if runtime.GOOS == windowsOS {
		pinger2.SetPrivileged(true)
	}

	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
			pinger2.Stop()
		}
	}()

	times := make([]int, 0)
	resultG := make(map[int]int64)
	resultR := make(map[int]int64)

	pinger.OnRecv = func(pkt *ping.Packet) {
		times = append(times, int(time.Now().Unix()))
		resultG[int(time.Now().Unix())] = pkt.Rtt.Milliseconds()
	}
	pinger2.OnRecv = func(pkt *ping.Packet) {
		times = append(times, int(time.Now().Unix()))
		resultR[int(time.Now().Unix())] = pkt.Rtt.Milliseconds()
	}
	fmt.Println("Started! To stop press ctrl + c")

	go pinger.Run()
	pinger2.Run()

	timesFormatted := make([]string, 0)
	router := make([]int64, 0)
	stadia := make([]int64, 0)
	sort.Ints(times)
	for _, t := range times {
		timesFormatted = append(timesFormatted, time.Unix(int64(t), 0).Format("15:04:05"))
		rR, ok := resultR[t]
		if !ok {
			continue
		}
		router = append(router, rR)
		rG, ok := resultG[t]
		if !ok {
			continue
		}
		stadia = append(stadia, rG)
	}

	line := charts.NewLine()
	line.Title = "Ping time to stadia and to your router"
	line.PageTitle = "Ping stadia and router"
	line.AddXAxis(timesFormatted).AddYAxis("stadia", stadia).AddYAxis("Router", router)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = line.Render(f)
	if err != nil {
		panic(err)
	}
	err = openFile(filename)
	if err != nil {
		panic(err)
	}
}

func openFile(filename string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", filename).Start()
	case windowsOS:
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", filename).Start()
	case "darwin":
		err = exec.Command("open", filename).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}
