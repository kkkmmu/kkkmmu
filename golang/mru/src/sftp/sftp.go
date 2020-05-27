package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type Client struct {
	*sftp.Client
}

func NewClient(ip, port, user, pass string) (*Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(pass))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         1000 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr = fmt.Sprintf("%s:%s", ip, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return &Client{Client: sftpClient}, nil
}

func (cli *Client) Get(remote, local string) error {
	rFile, err := cli.Open(remote)
	if err != nil {
		return fmt.Errorf("Cannot download file %s with %s", remote, err)
	}
	defer rFile.Close()

	lFile, err := os.Create(local)
	if err != nil {
		return fmt.Errorf("Cannot download file %s with %s", remote, err)
	}
	defer lFile.Close()

	if _, err = rFile.WriteTo(lFile); err != nil {
		return fmt.Errorf("Cannot download file %s with %s", remote, err)
	}

	return nil
}

func (cli *Client) Put(local, remote string) error {
	rFile, err := os.Open(local)
	if err != nil {
		return fmt.Errorf("Cannot upload file %s with %s", local, err)
	}
	defer rFile.Close()

	lFile, err := cli.Create(remote)
	if err != nil {
		return fmt.Errorf("Cannot upload file %s with %s", local, err)
	}
	defer lFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := rFile.Read(buf)
		if n == 0 {
			break
		}
		lFile.Write(buf)
	}

	return nil
}
