// from https://github.com/dynport/gossh
package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type ClientConfig struct {
	User       string
	Host       string
	Port       int
	Password   string
	PrivateKey string
}

type Client struct {
	Conf  *ClientConfig
	Conn  *ssh.Client
	agent net.Conn
}

func NewClient(conf *ClientConfig) (*Client, error) {
	c := &Client{
		Conf: conf,
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
	if c.agent != nil {
		c.agent.Close()
	}
}

func (c *Client) Connect() (err error) {
	if c.Conf.Port == 0 {
		c.Conf.Port = 22
	}

	config := &ssh.ClientConfig{
		User:            c.Conf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	keys := []ssh.Signer{}
	if c.Conf.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(c.Conf.Password))
	}
	if c.agent, err = net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		signers, err := agent.NewClient(c.agent).Signers()
		if err == nil {
			keys = append(keys, signers...)
		}
	}

	if len(c.Conf.PrivateKey) != 0 {
		if pk, err := readPrivateKey(c.Conf.PrivateKey); err == nil {
			keys = append(keys, pk)
		}
	} else {
		if pk, err := readPrivateKey(os.ExpandEnv("$HOME/.ssh/id_rsa")); err == nil {
			keys = append(keys, pk)
		}
	}

	if len(keys) > 0 {
		config.Auth = append(config.Auth, ssh.PublicKeys(keys...))
	}

	c.Conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Conf.Host, c.Conf.Port), config)
	return err
}

func readPrivateKey(path string) (ssh.Signer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(b)
}

type Result struct {
	StdoutBuffer, StderrBuffer *bytes.Buffer
	Duration                   time.Duration
	Error                      error
	ExitStatus                 int
}

func (r *Result) Stdout() string {
	return r.StdoutBuffer.String()
}

func (r *Result) Stderr() string {
	return r.StderrBuffer.String()
}

func (self *Result) IsSuccess() bool {
	return self.ExitStatus == 0
}

func (r *Result) String() string {
	return fmt.Sprintf("stdout: %s\nstderr: %s\nduration: %f\nstatus: %d",
		r.StdoutBuffer.String(), r.StderrBuffer.String(), r.Duration.Seconds(), r.ExitStatus)
}

func (c *Client) Execute(s string) (r *Result, e error) {
	started := time.Now()
	ses, e := c.Conn.NewSession()
	if e != nil {
		return nil, e
	}
	defer ses.Close()

	tmodes := ssh.TerminalModes{
		53:  0,     // disable echoing
		128: 14400, // input speed = 14.4kbaud
		129: 14400, // output speed = 14.4kbaud
	}

	if e := ses.RequestPty("xterm", 80, 40, tmodes); e != nil {
		return nil, e
	}

	r = &Result{
		StdoutBuffer: bytes.NewBuffer(nil),
		StderrBuffer: bytes.NewBuffer(nil),
	}

	ses.Stdout = r.StdoutBuffer
	ses.Stderr = r.StderrBuffer

	r.Error = ses.Run(s)
	r.Duration = time.Now().Sub(started)

	if exitError, ok := r.Error.(*ssh.ExitError); ok {
		r.ExitStatus = exitError.ExitStatus()
	}

	if !r.IsSuccess() {
		r.Error = fmt.Errorf("process exited with %d", r.ExitStatus)
	}

	return r, r.Error
}
