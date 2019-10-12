package gitsam

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/eyedeekay/eephttpd"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/gliderlabs/ssh"
	"github.com/sosedoff/gitkit"
)

//GitSAMTunnel is a structure which automatically configured the forwarding of
//a local service to i2p over the SAM API.
type GitSAMTunnel struct {
	*samforwarder.SAMForwarder
	*gitkit.SSH
	Conf       gitkit.Config
	OptPage    *eephttpd.EepHttpd
	PubKeyPath string
	SecurePath string
	up         bool
}

var err error

func (f *GitSAMTunnel) LookupKey() (*gitkit.PublicKey, error) {
	textkey, err := ioutil.ReadFile(f.PubKeyPath)
	if err != nil {
		return nil, err
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey(textkey)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(key.Marshal)
	print := base64.StdEncoding.EncodeToString(hash[:])
	return &gitkit.PublicKey{Fingerprint: print}, nil
}

func (f *GitSAMTunnel) GetType() string {
	return "gitsam"
}

func (f *GitSAMTunnel) ServeParent() {
	log.Println("Starting eepsite server", f.Base32())
	go f.OptPage.Serve()
	if err = f.SAMForwarder.Serve(); err != nil {
		f.Cleanup()
	}
}

//Serve starts the SAM connection and and forwards the local host:port to i2p
func (f *GitSAMTunnel) Serve() error {
	if f.Up() {
		go f.ServeParent()
		log.Println("Starting ssh server", f.Target())
		if err := f.SSH.ListenAndServe(f.Target()); err != nil {
			return err
		}
	}
	return nil
}

func (f *GitSAMTunnel) Up() bool {
	if f.SAMForwarder.Up() {
		return true
	}
	return false
}

//Close shuts the whole thing down.
func (f *GitSAMTunnel) Close() error {
	return f.SAMForwarder.Close()
}

func (s *GitSAMTunnel) Load() (samtunnel.SAMTunnel, error) {
	if !s.up {
		log.Println("Started putting tunnel up")
	}
	f, e := s.SAMForwarder.Load()
	if e != nil {
		return nil, e
	}
	if s.SecurePath == "" {
		s.SecurePath = filepath.Dir(s.Conf.Dir)
	}
	s.SAMForwarder = f.(*samforwarder.SAMForwarder)
	s.Conf.KeyDir = s.SecurePath
	s.SSH = gitkit.NewSSH(s.Conf)
	s.up = true
	log.Println("Finished putting tunnel up")
	return s, nil
}

//NewGitSAMTunnel makes a new SAM forwarder with default options, accepts host:port arguments
func NewGitSAMTunnel(host, port string) (*GitSAMTunnel, error) {
	return NewGitSAMTunnelFromOptions(SetHost(host), SetSSHPort(port))
}

//NewGitSAMTunnelFromOptions makes a new SAM forwarder with default options, accepts host:port arguments
func NewGitSAMTunnelFromOptions(opts ...func(*GitSAMTunnel) error) (*GitSAMTunnel, error) {
	var s GitSAMTunnel
	s.SAMForwarder = &samforwarder.SAMForwarder{}
	s.OptPage = &eephttpd.EepHttpd{}
	s.Conf = gitkit.Config{}
	s.SSH = &gitkit.SSH{}
	log.Println("Initializing gitsam")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	log.Println("Options loaded", s.Print())
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*GitSAMTunnel), nil
}
