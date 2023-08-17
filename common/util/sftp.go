package util

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
)

// SftpConnect is a controller that contains all the sftp server address information for sftp connector.
type SftpConnect struct {
	Host       string       //ip
	Port       int          //port
	Username   string       //user name
	Password   string       //password
	sshClient  *ssh.Client  //ssh client
	sftpClient *sftp.Client //sftp client
	File       *sftp.File
}

// SftpConfig for simplify the input of sftp connection parameters
type SftpConfig struct {
	Host     string //ip
	Port     int    //port
	Username string //user name
	Password string //password
}

// ConnectSftp connect sftp server using SftpConfig
func ConnectSftp(config *SftpConfig) (*SftpConnect, error) {
	connect := new(SftpConnect)
	host := config.Host
	port := config.Port
	userName := config.Username
	password := config.Password

	err := connect.CreateClient(host, port, userName, password)
	if err != nil {
		return nil, err
	}
	return connect, nil
}

// ConnectSftpByParam connect sftp server using the parameters
func ConnectSftpByParam(host string, port int, userName string, password string) (*SftpConnect, error) {
	connect := new(SftpConnect)
	err := connect.CreateClient(host, port, userName, password)
	if err != nil {
		return nil, errors.Errorf(constant.SystemInternalError, "Connect sftp %s:%d failed, err=%++v", host, port, err)
	}
	return connect, nil
}

// SftpHandler defines the interface of the sftp controller should have
type SftpHandler interface {
	CreateClient(host string, port int, username, password string) error
	CreateClientUseKey(host string, port int, username string, key []byte) error
	ClientClose() error
	Upload(srcPath, dstPath string) error
	UploadFromReader(srcFile io.ReadCloser, dstPath string) error
	Download(srcPath, dstPath string) error
	MoveTo(srcPath, dstPath string) error
	GetSftpClient() *sftp.Client
	FileOpen(path string) (*sftp.File, error)
	FileClose() error
}

// NewSftpHandler implements all the SftpHandler interface
func NewSftpHandler() SftpHandler {
	return &SftpConnect{}
}

// CreateClient create sftp client
func (cliConf *SftpConnect) CreateClient(host string, port int, username, password string) error {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password

	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create ssh client failed, err=%++v", err)
	}
	cliConf.sshClient = sshClient

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create sftp client failed, err=%++v", err)
	}
	cliConf.sftpClient = sftpClient
	return nil
}

// CreateClientUseKey create sftp client
func (cliConf *SftpConnect) CreateClientUseKey(host string, port int, username string, key []byte) error {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Parse private Key failed. error: %v", err)
	}

	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create ssh client failed, err=%++v", err)
	}
	cliConf.sshClient = sshClient

	//At this time, sshClient is obtained, and sshClient is used to build sftpClient below
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create sftp client failed, err=%++v", err)
	}
	cliConf.sftpClient = sftpClient
	return nil
}

// ClientClose close sftp client
func (cliConf SftpConnect) ClientClose() error {
	if cliConf.sftpClient != nil {
		err := cliConf.sftpClient.Close()
		if err != nil {
			return errors.Errorf(constant.SystemInternalError, "Close sftp client failed, err=%++v", err)
		}
	}
	if cliConf.sshClient != nil {
		err := cliConf.sshClient.Close()
		if err != nil {
			return errors.Errorf(constant.SystemInternalError, "Close ssh client failed, err=%++v", err)
		}
	}
	return nil
}

// Upload upload file to sftp server
func (cliConf *SftpConnect) Upload(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Open local file %s failed, err=%++v", srcPath, err)
	}
	dstFile, err := cliConf.sftpClient.Create(dstPath)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create sftp file %s failed, err=%++v", dstPath, err)
	}
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			log.Errorsf("read data file %s failed, err=%++v", srcPath, err)
			if err != io.EOF {
				return err
			}
			if n != 0 {
				count, err := dstFile.Write(buf[:n])
				if err != nil {
					return errors.Errorf(constant.SystemInternalError, "Write data to sftp file %s failed, err=%++v", dstPath, err)
				}
				log.Debugsf("Write %d bytes to sftp file %s success", count, dstPath)
			}
			break
		}
		log.Debugsf("read %d bytes data from reader success", n)
		count, err := dstFile.Write(buf[:n])
		if err != nil {
			return errors.Errorf(constant.SystemInternalError, "Write data to sftp file %s failed, err=%++v", dstPath, err)
			//return err
		}
		log.Debugsf("Write %d bytes to sftp file %s success", count, dstPath)
	}
	log.Debugsf("Upload file %s to %s success", srcPath, dstPath)
	return nil
}

// UploadFromReader upload file from reader
func (cliConf *SftpConnect) UploadFromReader(srcFile io.ReadCloser, dstPath string) error {
	dstFile, err := cliConf.sftpClient.Create(dstPath)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create sftp file %s failed, err=%++v", dstPath, err)
	}
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			fmt.Errorf("read data failed, err=%++v", err)
			if err != io.EOF {
				return err
			}
			if n != 0 {
				count, err := dstFile.Write(buf[:n])
				if err != nil {
					return errors.Errorf(constant.SystemInternalError, "Write data to sftp file %s failed, err=%++v", dstPath, err)
					//return err
				}
				log.Debugsf("Write %d bytes to sftp file %s success", count, dstPath)
			}
			break
		}
		count, err := dstFile.Write(buf[:n])
		if err != nil {
			return errors.Errorf(constant.SystemInternalError, "Write data to sftp file %s failed, err=%++v", dstPath, err)
		}
		log.Debugsf("Write %d bytes to sftp file %s success", count, dstPath)
	}
	log.Debugsf("Upload file %s success", dstPath)
	return nil
}

//Download download file from sftp server
func (cliConf *SftpConnect) Download(srcPath, dstPath string) error {
	srcFile, err := cliConf.sftpClient.Open(srcPath)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Open sftp file %s failed, err=%++v", srcPath, err)
	}
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create local file %s failed, err=%++v", dstPath, err)
	}
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return errors.Errorf(constant.SystemInternalError, "Write data to local file %s failed, err=%++v", dstPath, err)
	}
	log.Debugsf("Download file %s to %s success", srcPath, dstPath)
	return nil
}

// MoveTo move the file to another path
func (cliConf SftpConnect) MoveTo(srcPath, dstPath string) error {
	srcFile, err := cliConf.sftpClient.OpenFile(srcPath, os.O_RDONLY)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Open file %s failed, err=%++v", srcPath, err)
	}
	dstFile, err := cliConf.sftpClient.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Create file %s failed, err=%++v", dstPath, err)
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Copy file %s to %s failed, err=%++v", srcPath, dstPath, err)
	}
	log.Debugsf("Copy file %s to %s success", srcPath, dstPath)
	for i := 0; i < 3; i++ {
		err := cliConf.sftpClient.Remove(srcPath)
		if err != nil {
			log.Errorsf("Remove file %s failed, err=%++v", srcPath, err)
			continue
		}
		break
	}
	log.Debugsf("Move file %s to %s success", srcPath, dstPath)
	return nil
}

// FileOpen open the file in sftp server
func (cliConf SftpConnect) FileOpen(path string) (*sftp.File, error) {

	file, err := cliConf.GetSftpClient().Open(path)
	if err != nil {
		return nil, err
	}

	cliConf.File = file

	return file, nil
}

// FileClose close the file
func (cliConf SftpConnect) FileClose() error {

	if cliConf.File == nil {
		return fmt.Errorf("SftpConnect(for File) hasn't been initialized, please use FileOpen to initialize the sftpFile")
	}
	return cliConf.File.Close()
}

// GetSftpClient get sftp client from connector
func (cliConf SftpConnect) GetSftpClient() *sftp.Client {
	return cliConf.sftpClient
}
