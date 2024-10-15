package handlers

import (
	"encoding/json"
	"strconv"
	"fmt"
	"os"
	"os/signal"
	"github.com/go-ping/ping"
	"github.com/gofiber/fiber/v2"
)

func TestPing(context *fiber.Ctx) error {

	ip := context.Params("ip")
	countStr := context.Params("count")
	

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return context.Status(500).SendString(err.Error())
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
    	return context.Status(400).SendString("Invalid count value")
	}
	pinger.Count = count

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {

			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)

	
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		if stats.PacketsRecv == 0 {
			// Ping başarısız oldu
			context.SendString("ping access is fail")
		} else {
			// Ping başarılı oldu
			context.SendString("ping successful")

			result := struct {
				Address     string  `json:"address"`
				PacketsSent int     `json:"packets_sent"`
				PacketsRecv int     `json:"packets_recv"`
				PacketLoss  float64 `json:"packet_loss"`
				AvgRtt      string  `json:"avg_rtt"`
			}{
				Address:     stats.Addr,
				PacketsSent: stats.PacketsSent,
				PacketsRecv: stats.PacketsRecv,
				PacketLoss:  stats.PacketLoss,
				AvgRtt:      stats.AvgRtt.String(),
			}
	

			jsonResult, err := json.Marshal(result)
			if err != nil {
				context.Status(500).SendString(err.Error())
				return
			}
	
			context.Set("Content-Type", "application/json")
			context.SendString(string(jsonResult))
		}
	}


	
	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		return context.Status(500).SendString(err.Error())
	}

	return nil
}
