package scan

import (
	"bytes"
	"github.com/MayMistery/noscan/cmd"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func CheckLive(host string) bool {

	if cmd.Config.Ping == true {
		//Use system ping command
		return ExecCommandPing(host)
	} else {
		//Try icmp without listening
		conn, err := net.DialTimeout("ip4:icmp", "127.0.0.1", 3*time.Second)
		defer func() {
			if conn != nil {
				err := conn.Close()
				if err != nil {
					cmd.ErrLog("connect close error %v", err)
					return
				}
			}
		}()
		if err == nil {
			return IcmpAlive(host)
		} else {
			//Use system ping command
			//TODO add it
			//fmt.Println("The current user permissions unable to send icmp packets")
			//fmt.Println("start ping")
			return ExecCommandPing(host)
		}
	}
}

func IcmpAlive(host string) bool {
	startTime := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, 6*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	//TODO log timeout or not
	if err != nil {
		return false
	}
	if err := conn.SetDeadline(startTime.Add(6 * time.Second)); err != nil {
		return false
	}
	msg := makeMsg(host)
	if _, err := conn.Write(msg); err != nil {
		return false
	}

	receive := make([]byte, 60)
	if _, err := conn.Read(receive); err != nil {
		return false
	}

	return true
}

func ExecCommandPing(host string) bool {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+host+" && echo true || echo false")
	case "linux":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+host+" >/dev/null && echo true || echo false")
	case "darwin":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+host+" >/dev/null && echo true || echo false")
	case "freebsd":
		command = exec.Command("ping", "-c", "1", "-W", "200", host)
	case "openbsd":
		command = exec.Command("ping", "-c", "1", "-w", "200", host)
	case "netbsd":
		command = exec.Command("ping", "-c", "1", "-w", "2", host)
	default:
		command = exec.Command("ping", "-c", "1", host)
	}
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		cmd.ErrLog("Fail to ping by bash %v", err)
		return false
	}
	if err = command.Wait(); err != nil {
		cmd.ErrLog("Fail to ping by bash (wait) %v", err)
		return false
	} else {
		if strings.Contains(outinfo.String(), "true") {
			return true
		} else {
			return false
		}
	}
}

func makeMsg(host string) []byte {
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
