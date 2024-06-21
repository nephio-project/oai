package controller

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

type O1TelnetClient struct {
	url  string
	port string
	conn *net.TCPConn
}

/*
Opens a TCP connection to 'url:port' {Timeout of 5 sec}
*/
func (o *O1TelnetClient) openConnection() error {
	svc := o.url + ":" + o.port
	tcpServer, err := net.ResolveTCPAddr("tcp", svc)
	if err != nil {
		return err
	}
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", tcpServer.String())
	if err != nil {
		return err
	}

	o.conn = conn.(*net.TCPConn)
	return nil
}

func (o *O1TelnetClient) closeConnection() error {
	if o.conn == nil {
		// fmt.Println("Connection was already closed or never opened")
		return nil
	}
	o.conn.Close()
	return nil
}

/*
Writes to the tcp channel
*/
func (o *O1TelnetClient) writeToConnection(cmd string) error {
	_, err := o.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		return err
	}
	return nil
}

/*
Reads From the tcp channel
A Brief About the Logic:
It is observed that the telnet doesn't send EOF at the end of message due to following readers fails:
1. net.conn.Read
2. telent.conn.Read
3. ioutil.ReadAll
4. textproto.NewReader
Since, they keep reading until they don't find EOF (which they will never find :( )
As a work-around: we follow 2 cases:
1. It is seen that every telent message ends with 'softmodem_gnb', we will treat that as EOF
2. Set a ReadDeadline (generally 5 sec) to prevent infinite-loop
*/
func (o *O1TelnetClient) readFromConnection() (string, error) {
	var out, curLine string
	tmp := make([]byte, 1) // Read the charactor
	for {
		o.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err := o.conn.Read(tmp)
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "i/o timeout") {
				return out, err
			}
			break
		}
		curChar := string(tmp)
		if curChar == "\n" { // End of Line
			out += (curLine + "\n")
			curLine = ""
		} else {
			curLine += curChar
		}

		if curLine == "softmodem_gnb" { // It is treated as EOF
			break
		}
	}
	return out, nil
}

func (o *O1TelnetClient) RunCommand(cmd string) (string, error) {
	if err := o.openConnection(); err != nil {
		return "", errors.New("Unable to open-connection :: " + err.Error())
	}
	defer o.closeConnection()

	if err := o.writeToConnection(cmd); err != nil {
		return "", errors.New("Unable to write-Command :: " + err.Error())
	}

	out, err := o.readFromConnection()
	if err != nil {
		return "", errors.New("Unable to read-Command-Output :: " + err.Error())
	}
	return out, nil
}
