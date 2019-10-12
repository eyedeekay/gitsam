package gitsam

import (
	"log"

	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
    "github.com/eyedeekay/eephttpd"
)

//GitSAMTunnel is a structure which automatically configured the forwarding of
//a local service to i2p over the SAM API.
type GitSAMTunnel struct {
	*samforwarder.SAMForwarder
    OptPage     *eephttpd.EepHttpd
	//ServeDir string
	up       bool
}

var err error

func (f *GitSAMTunnel) GetType() string {
	return "gitsam"
}

func (f *GitSAMTunnel) ServeParent() {
	log.Println("Starting eepsite server", f.Base32())
	if err = f.SAMForwarder.Serve(); err != nil {
		f.Cleanup()
	}
}

//Serve starts the SAM connection and and forwards the local host:port to i2p
func (f *GitSAMTunnel) Serve() error {
	go f.ServeParent()
	if f.Up() {
		log.Println("Starting web server", f.Target())
		/*if err := http.ListenAndServe(f.Target(), f); err != nil {
			return err
		}*/
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
	s.SAMForwarder = f.(*samforwarder.SAMForwarder)
	s.up = true
	log.Println("Finished putting tunnel up")
	return s, nil
}

//NewGitSAMTunnel makes a new SAM forwarder with default options, accepts host:port arguments
func NewGitSAMTunnel(host, port string) (*GitSAMTunnel, error) {
	return NewGitSAMTunnelFromOptions(SetHost(host), SetPort(port))
}

//NewGitSAMTunnelFromOptions makes a new SAM forwarder with default options, accepts host:port arguments
func NewGitSAMTunnelFromOptions(opts ...func(*GitSAMTunnel) error) (*GitSAMTunnel, error) {
	var s GitSAMTunnel
	s.SAMForwarder = &samforwarder.SAMForwarder{}
    s.OptPage = &eephttpd.EepHttpd{}
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
