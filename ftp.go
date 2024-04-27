package ftp

//https://github.com/secsy/goftp
import (
	"crypto/tls"
	"github.com/civet148/log"
	"github.com/secsy/goftp"
	"net/url"
	"os"
	"time"
)

type Option struct {
	// Maximum number of FTP connections to open per-host. Defaults to 5. Keep in
	// mind that FTP servers typically limit how many connections a single user
	// may have open at once, so you may need to lower this if you are doing
	// concurrent transfers.
	ConnectionsPerHost int

	// Timeout for opening connections, sending control commands, and each read/write
	// of data transfers. Defaults to 5 seconds.
	TimeoutSeconds int

	// TLS Config used for FTPS. If provided, it will be an error if the server
	// does not support TLS. Both the control and data connection will use TLS.
	TLSConfig *tls.Config

	// FTPS mode. TLSExplicit means connect non-TLS, then upgrade connection to
	// TLS via "AUTH TLS" command. TLSImplicit means open the connection using
	// TLS. Defaults to TLSExplicit.
	TLSMode goftp.TLSMode

	// This flag controls whether to use IPv6 addresses found when resolving
	// hostnames. Defaults to false to prevent failures when your computer can't
	// IPv6. If the hostname(s) only resolve to IPv6 addresses, Dial() will still
	// try to use them as a last ditch effort. You can still directly give an
	// IPv6 address to Dial() even with this flag off.
	IPv6Lookup bool

	// Time zone of the FTP server. Used when parsing mtime from "LIST" output if
	// server does not support "MLST"/"MLSD". Defaults to UTC.
	ServerLocation *time.Location

	// Enable "active" FTP data connections where the server connects to the client to
	// establish data connections (does not work if client is behind NAT). If TLSConfig
	// is specified, it will be used when listening for active connections.
	ActiveTransfers bool

	// Set the host:port to listen on for active data connections. If the host and/or
	// port is empty, the local address/port of the control connection will be used. A
	// port of 0 will listen on a random port. If not specified, the default behavior is
	// ":0", i.e. listen on the local control connection host and a random port.
	ActiveListenAddr string

	// Disables EPSV in favour of PASV. This is useful in cases where EPSV connections
	// neither complete nor downgrade to PASV successfully by themselves, resulting in
	// hung connections.
	DisableEPSV bool
}

type Client struct {
	strHost string
	fc      *goftp.Client
}

// NewFtpClient create an FTP client
// strUrl => ftp://user:password@127.0.0.1:21
func NewFtpClient(strUrl string, opts ...*Option) *Client {
	var opt Option
	ui, err := url.Parse(strUrl)
	if err != nil {
		log.Panic("invalid URL %s", strUrl)
		return nil
	}
	if len(opts) != 0 {
		opt = *opts[0]
	}
	var strUser, strPassword string
	strUser = ui.User.Username()
	strPassword, _ = ui.User.Password()
	strHost := ui.Host

	var fc *goftp.Client
	fc, err = goftp.DialConfig(goftp.Config{
		User:               strUser,
		Password:           strPassword,
		ConnectionsPerHost: opt.ConnectionsPerHost,
		Timeout:            time.Duration(opt.TimeoutSeconds) * time.Second,
		TLSConfig:          opt.TLSConfig,
		TLSMode:            opt.TLSMode,
		IPv6Lookup:         opt.IPv6Lookup,
		ServerLocation:     opt.ServerLocation,
		ActiveTransfers:    opt.ActiveTransfers,
		ActiveListenAddr:   opt.ActiveListenAddr,
		DisableEPSV:        opt.DisableEPSV,
	}, strHost)
	if err != nil {
		log.Panic("connect ftp [%s] error [%s]", strHost, err.Error())
		return nil
	}
	return &Client{
		fc:      fc,
		strHost: strHost,
	}
}

func (m *Client) Mkdir(dir string) (err error) {
	_, err = m.fc.Mkdir(dir)
	if err != nil {
		return log.Errorf("mkdir [%s] error [%s]", dir, err)
	}
	return nil
}

func (m *Client) Download(strFilePath, strStorePath string) error {
	writer, err := os.Open(strStorePath)
	if err != nil {
		return log.Errorf("open file [%s] error [%s]", strStorePath, err)
	}
	err = m.fc.Retrieve(strFilePath, writer)
	if err != nil {
		return log.Errorf("write file [%s] error [%s]", strStorePath, err)
	}
	return nil
}

func (m *Client) Upload(strFilePath, strStorePath string) error {
	reader, err := os.Open(strFilePath)
	if err != nil {
		return log.Errorf("read file [%s] error [%s]", strFilePath, err)
	}
	err = m.fc.Store(strStorePath, reader)
	if err != nil {
		return log.Errorf("upload file [%s] to [%s] error [%s]", strFilePath, strStorePath, err)
	}
	return nil
}

func (m *Client) ReadDir(dir string) (fis []os.FileInfo, err error) {
	fis, err = m.fc.ReadDir(dir)
	if err != nil {
		return nil, log.Errorf("list work dir [%s] files error [%s]", dir, err)
	}

	return fis, nil
}

func (m *Client) Delete(path string) error {
	fi, err := m.fc.Stat(path)
	if err != nil {
		return log.Errorf("stat path [%s] error [%s]", path, err)
	}
	if fi.IsDir() {
		return m.fc.Rmdir(path)
	}
	return m.fc.Delete(path)
}

func (m *Client) Stat(path string) (fi os.FileInfo, err error) {
	return m.fc.Stat(path)
}

func (m *Client) Rename(from, to string) (err error) {
	return m.fc.Rename(from, to)
}

func (m *Client) GetWorkDir() (string, error) {
	return m.fc.Getwd()
}

func (m *Client) Close() error {
	return m.fc.Close()
}
