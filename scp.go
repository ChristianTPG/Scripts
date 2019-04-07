package main

import (
	"golang.org/x/crypto/ssh"
	"net"
	"time"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type SSH struct {
	Ip string
	User string
	Cert string
	Port int
	session *ssh.Session
	client *ssh.Client
}

func (sshClient *SSH) readPublicKeyFile() ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(sshClient.Cert)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

func (sshClient *SSH) Connect() {
	auth := []ssh.AuthMethod{sshClient.readPublicKeyFile()}
	config := &ssh.ClientConfig{
		User: sshClient.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3, // Three seconds
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshClient.Ip, sshClient.Port), config)
	if err != nil {
		fmt.Println(err)
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		client.Close()
		return
	}

	sshClient.session = session
	sshClient.client = client
}

func (sshClient *SSH) RunCmd(cmd string) []byte {
	out, err := sshClient.session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(err)
	}
	return out
}

func (sshClient *SSH) Close() {
	sshClient.session.Close()
	sshClient.client.Close()
}

func main() {
	client := &SSH{
		Ip: "192.168.0.5",
		User: "osmc",
		Port: 22,
		Cert: "/Users/christian/.ssh/id_rsa",
	}

	client.Connect()
	filePath := strings.TrimSuffix(string(client.RunCmd("ls -1td \"$PWD/\"Downloads/*.png | head -1")), "\n")
	client.Close()

	addr := fmt.Sprintf("%s@%s", client.User, client.Ip)
	command := fmt.Sprintf("/usr/bin/scp -i %s %s:%s .", client.Cert, addr, filePath)

	if _, err := exec.Command("sh", "-c", command).CombinedOutput(); err != nil {
		log.Fatal(err)
	}
}