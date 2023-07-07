package scan

import (
	"bytes"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"golang.org/x/net/icmp"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	livewg sync.WaitGroup
)

func CheckLive(config cmd.Configs) []string {
	chanHosts := make(chan string, int(cmd.IPPoolsSize))
	go func() {
		for ip := range chanHosts {
			if _, ok := ExistHosts[ip]; !ok && cmd.InTarget(ip) {
				ExistHosts[ip] = true
				if config.Ping == false {
					fmt.Printf("(icmp) Target %-15s is alive\n", ip)
				} else {
					fmt.Printf("(ping) Target %-15s is alive\n", ip)
				}
				AliveHosts = append(AliveHosts, ip)
			}
			livewg.Done()
		}
	}()

	if config.Ping == true {
		//Use system ping command
		RunPing(chanHosts)
	} else {
		//First try icmp with listening
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err == nil {
			RunIcmp1(conn, chanHosts)
		} else {
			//Try icmp without listening
			fmt.Println("trying RunIcmp2")
			conn, err := net.DialTimeout("ip4:icmp", "127.0.0.1", 3*time.Second)
			defer func() {
				if conn != nil {
					conn.Close()
				}
			}()
			if err == nil {
				RunIcmp2(chanHosts)
			} else {
				//Use system ping command
				fmt.Println("The current user permissions unable to send icmp packets")
				fmt.Println("start ping")
				RunPing(chanHosts)
			}
		}
	}

	livewg.Wait()
	close(chanHosts)

	return AliveHosts
}

func RunIcmp1(conn *icmp.PacketConn, chanHosts chan string) {
	endflag := false
	go func() {
		for {
			if endflag == true {
				return
			}
			msg := make([]byte, 100)
			_, sourceIP, _ := conn.ReadFrom(msg)
			if sourceIP != nil {
				livewg.Add(1)
				chanHosts <- sourceIP.String()
			}
		}
	}()

	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		dst, _ := net.ResolveIPAddr("ip", host)
		IcmpByte := makemsg(host)
		conn.WriteTo(IcmpByte, dst)
	}
	//根据hosts数量修改icmp监听时间
	start := time.Now()
	for {
		if len(AliveHosts) == int(cmd.IPPoolsSize) {
			break
		}
		since := time.Now().Sub(start)
		var wait time.Duration
		switch {
		case int(cmd.IPPoolsSize) <= 256:
			wait = time.Second * 3
		default:
			wait = time.Second * 6
		}
		if since > wait {
			break
		}
	}
	endflag = true
	conn.Close()
}

func RunIcmp2(chanHosts chan string) {
	var num int64 = 1000
	if cmd.IPPoolsSize < num {
		num = cmd.IPPoolsSize
	}
	var wg sync.WaitGroup
	limiter := make(chan struct{}, num)
	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if icmpalive(host) {
				livewg.Add(1)
				chanHosts <- host
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
	close(limiter)
}

func icmpalive(host string) bool {
	startTime := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, 6*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return false
	}
	if err := conn.SetDeadline(startTime.Add(6 * time.Second)); err != nil {
		return false
	}
	msg := makemsg(host)
	if _, err := conn.Write(msg); err != nil {
		return false
	}

	receive := make([]byte, 60)
	if _, err := conn.Read(receive); err != nil {
		return false
	}

	return true
}

func RunPing(chanHosts chan string) {
	var bsenv = ""
	if OS != "windows" {
		bsenv = "/bin/bash"
	}
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 50)
	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if ExecCommandPing(host, bsenv) {
				livewg.Add(1)
				chanHosts <- host
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
}

func ExecCommandPing(ip string, bsenv string) bool {
	var command *exec.Cmd
	if OS == "windows" {
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	} else if OS == "linux" {
		command = exec.Command(bsenv, "-c", "ping -c 1 -w 1 "+ip+" >/dev/null && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	} else if OS == "darwin" {
		command = exec.Command(bsenv, "-c", "ping -c 1 -W 1 "+ip+" >/dev/null && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	}
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		return false
	}
	if err = command.Wait(); err != nil {
		return false
	} else {
		if strings.Contains(outinfo.String(), "true") {
			return true
		} else {
			return false
		}
	}
}

func makemsg(host string) []byte {
	msg := make([]byte, 40)
	id0, id1 := genIdentifier(host)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4], msg[5] = id0, id1
	msg[6], msg[7] = genSequence(1)
	check := checkSum(msg[0:40])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	return msg
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	answer := uint16(^sum)
	return answer
}

func genSequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}
