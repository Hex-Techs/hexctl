package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// default remote info
const (
	DefaultPort    = 22
	DefaultUser    = "root"
	DefaultKeyFile = "~/.ssh/id_rsa"
)

// NewRemoteOption return a new RemoteOption
func NewRemoteOption(host, port, user, password, key, netdevice string) (*RemoteOption, error) {
	auth := make([]ssh.AuthMethod, 0)
	if key != "" {
		pemBytes, err := os.ReadFile(key)
		if err != nil {
			return nil, err
		}
		var signer ssh.Signer
		signer, err = ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	} else {
		auth = append(auth, ssh.Password(password))
	}

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 5 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), clientConfig)
	if err != nil {
		return nil, err
	}
	remote := &RemoteOption{
		Host:      host,
		Port:      port,
		User:      user,
		KeyFile:   &key,
		Password:  &password,
		NetDevice: netdevice,
	}
	if remote.Session, err = client.NewSession(); err != nil {
		return nil, err
	}
	if remote.SftpClient, err = sftp.NewClient(client); err != nil {
		return nil, err
	}
	return remote, nil
}

type RemoteOption struct {
	Host       string
	Port       string
	User       string
	KeyFile    *string
	Password   *string
	WorkDir    string
	Command    *Command
	File       *TransferFile
	Session    *ssh.Session
	SftpClient *sftp.Client
	NetDevice  string
}

type TransferFile struct {
	SrcFile string
	DstFile string
}

type Command struct {
	Cmd   string
	Shell string
}

func (r *RemoteOption) RunCommand() (string, error) {
	b, err := r.Session.CombinedOutput(r.Command.Cmd)
	r.Session.Close()
	return string(b), err
}

func (r *RemoteOption) CopyFileToRmote() error {
	srcFile, err := os.Open(r.File.SrcFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := r.SftpClient.Create(r.File.DstFile)
	if err != nil {
		return err
	}
	defer destFile.Close()
	buf, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}
	_, err = destFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *RemoteOption) CopyFileFromRemote() error {
	srcFile, err := r.SftpClient.Open(r.File.SrcFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(r.File.DstFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	return nil
}

func (r *RemoteOption) Close() {
	r.Session.Close()
	r.SftpClient.Close()
}
