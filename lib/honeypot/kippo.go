package honeypot

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"regexp"
	"time"
)

func isKippoHoneypot(ip string, port int) (*ssh.Client, bool) {
	// Change this to your own SSH client config
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "SSH-1.9-OpenSSH_5.9p1",
		Timeout:         3 * time.Second,
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	//log.Printf("%v", err)
	if err != nil {
		match, _ := regexp.MatchString(`bad version`, err.Error())
		//if err.Error() == "ssh: handshake failed: ssh: disconnect, reason 8: bad version 1.9" {
		if match {
			// If we get this specific error message, it's probably a kippo honeypot
			//log.Printf("kippo")
			return nil, true
		}
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, false
		}
	}

	return conn, false
}
