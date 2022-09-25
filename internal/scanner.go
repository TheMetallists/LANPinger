package internal

import (
	"LANPinger/internal/ipnetgen"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

func (this *LANPinger) runScanner() {
	log.Println("Running scanner...")
	this.PTNBtn.Disable()
	this.IPEntry.Disable()
	this.Select.Disable()
	this.Progress.SetValue(0)
	this.Reports = []LANReport{}
	this.ReportsList.Refresh()

	err := this.actualPinger()
	if err != nil {
		log.Println(err)
	}
	this.Select.Enable()
	this.IPEntry.Enable()
	this.PTNBtn.Enable()
	this.PTNBtn.Refresh()
	time.Sleep(time.Second * 2)
}

func (this *LANPinger) actualPinger() error {
	gen, err := ipnetgen.New(this.IPEntry.Text + this.Subnet)
	if err != nil {
		return err
	}

	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Println("Error on ListenPacket")
		return err
	}
	defer conn.Close()

	ctr := uint64(0)

	for ip2try := gen.Next(); ip2try != nil; ip2try = gen.Next() {
		this.Progress.SetValue(gen.GetProgress())
		this.Progress.Refresh()

		ip, err := net.ResolveIPAddr("ip4", ip2try.String())
		if err != nil {
			return err
		}

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: 1,
				Data: []byte(""),
			},
		}
		msg_bytes, err := msg.Marshal(nil)
		if err != nil {
			log.Println("Error on Marshal", msg_bytes)
			return err
		}

		// Write the message to the listening connection
		if _, err := conn.WriteTo(msg_bytes, &net.UDPAddr{IP: net.ParseIP(ip.String())}); err != nil {
			log.Println("Error on WriteTo %v", err)
			continue
		}

		ctr++
		if ctr%100 == 0 {
			this.readout(conn)
		}
	}

	/*
		pingok, err := SendPingRequest(ip2try.String())
				if err != nil {
					continue
				}

				if pingok {
					this.Reports = append(this.Reports, LANReport{
						IP: ip2try.String(),
					})
					this.ReportsList.Refresh()
				}
	*/
	this.readout(conn)

	return nil
}

func (this *LANPinger) readout(conn *icmp.PacketConn) error {
	err := conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		log.Println("Error on SetReadDeadline %v", err)
		return err
	}

	endtm := time.Now().Add(time.Second * 2)

	for time.Now().Before(endtm) {
		reply := make([]byte, 1500)
		n, from, err := conn.ReadFrom(reply)

		if err != nil {
			log.Println("Error on ReadFrom %v", err)
			continue
		}
		parsed_reply, err := icmp.ParseMessage(1, reply[:n])
		if parsed_reply.Code == 0 {
			//ok
			this.Reports = append(this.Reports, LANReport{
				IP: from.String(),
			})
			this.ReportsList.Refresh()
		}
	}
	return nil
}
