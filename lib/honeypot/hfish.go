package honeypot

import (
	"golang.org/x/crypto/ssh"
	"log"
	"regexp"
)

//func isKippoHoneypot(ip string, port int) bool {
// // Change this to your own SSH client config
// //config := &ssh.ClientConfig{
// // User:            "username",
// // Auth:            []ssh.AuthMethod{ssh.Password("password")},
// // HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// // ClientVersion:   "SSH-1.9-OpenSSH_5.9p1",
// // Timeout:         3 * time.Second,
// //}
//
// config := &ssh.ClientConfig{
// User:            "root",
// Auth:            []ssh.AuthMethod{ssh.Password("")}, // No password
// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// ClientVersion:   "SSH-1.9-OpenSSH_5.9p1",
// Timeout:         1 * time.Second,
// }
//
// conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
//
// //_, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
// //log.Printf("%v", err)
// if err != nil {
// match, _ := regexp.MatchString(`bad version`, err.Error())
// //if err.Error() == "ssh: handshake failed: ssh: disconnect, reason 8: bad version 1.9" {
// if match {
// // If we get this specific error message, it's probably a kippo honeypot
// log.Printf("kippo")
// return true
// }
// if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
// return false
// }
// }
// session, err := conn.NewSession()
// if session == nil || err != nil {
// return false
// }
// defer func(session *ssh.Session) {
// err := session.Close()
// if err != nil {
// }
// }(session)
//
// out, err := session.CombinedOutput("\r\n")
// if err != nil {
// return false
// }
//
// match, _ := regexp.MatchString(`test`, string(out))
// if match {
// log.Printf("hfish")
// return true
// }
//
// return false
//}

func isHfishHoneypot(conn *ssh.Client) bool {
	session, err := conn.NewSession()
	if err != nil {
		return false
	}
	defer func() {
		err := session.Close()
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}()

	out, err := session.CombinedOutput("\r\n")
	if err != nil {
		return false
	}

	match, err := regexp.MatchString(`test`, string(out))
	if err != nil {
		log.Printf("Error matching string: %v", err)
		return false
	}
	if match {
		log.Printf("hfish")
		return true
	}

	return false
}
