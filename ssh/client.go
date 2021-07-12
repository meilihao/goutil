// from https://github.com/dynport/gossh
package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var (
	ErrNoAuth = errors.New("no ssh auth found")
)

type ClientConfig struct {
	User         string
	Host         string
	Port         int
	Password     string
	PrivateKey   string
	Passphrase   string // for decrypt PrivateKey
	UsePty       bool   // if true, will request a pty from the remote end
	DisableAgent bool
	Timeout      time.Duration
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
	if c.Conf.Timeout == 0 {
		c.Conf.Timeout = 5 * time.Second
	}

	config := &ssh.ClientConfig{
		Timeout:         c.Conf.Timeout,
		User:            c.Conf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: hostKeyCallBackFunc(c.Conf.Host),
	}

	keys := []ssh.Signer{}
	if !c.Conf.DisableAgent && os.Getenv("SSH_AUTH_SOCK") != "" {
		if c.agent, err = net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
			signers, err := agent.NewClient(c.agent).Signers()
			if err == nil {
				keys = append(keys, signers...)
			}
		}
	}
	if len(c.Conf.PrivateKey) != 0 {
		if pk, err := ReadPrivateKey(c.Conf.PrivateKey, c.Conf.Passphrase); err == nil {
			keys = append(keys, pk)
		}
	}
	if len(keys) > 0 {
		config.Auth = append(config.Auth, ssh.PublicKeys(keys...))
	}

	if c.Conf.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(c.Conf.Password))
	}

	if len(config.Auth) == 0 {
		return ErrNoAuth
	}

	c.Conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Conf.Host, c.Conf.Port), config)
	return err
}

func ReadPrivateKey(path, passphrase string) (ssh.Signer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if passphrase != "" {
		return ssh.ParsePrivateKeyWithPassphrase(b, []byte(passphrase))
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

func (r *Result) IsSuccess() bool {
	return r.ExitStatus == 0
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

	if c.Conf.UsePty {
		tmodes := ssh.TerminalModes{
			53:  0,     // disable echoing
			128: 14400, // input speed = 14.4kbaud
			129: 14400, // output speed = 14.4kbaud
		}

		if e := ses.RequestPty("xterm", 80, 40, tmodes); e != nil {
			return nil, e
		}
	}

	r = &Result{
		StdoutBuffer: bytes.NewBuffer(nil),
		StderrBuffer: bytes.NewBuffer(nil),
	}

	ses.Stdout = r.StdoutBuffer
	ses.Stderr = r.StderrBuffer

	r.Error = ses.Run(s)
	r.Duration = time.Since(started)

	if exitError, ok := r.Error.(*ssh.ExitError); ok {
		r.ExitStatus = exitError.ExitStatus()
	}

	if !r.IsSuccess() {
		if r.StderrBuffer.Len() > 0 {
			r.Error = errors.New(r.StdoutBuffer.String())
		}
	}

	return r, r.Error
}
