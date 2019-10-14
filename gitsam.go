package gitsam

import (
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
    "os/exec"
	"path/filepath"
	"strings"

	"github.com/eyedeekay/eephttpd"
	"github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/config/helpers"
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
	GitConf    gitkit.Config
	OptPage    *eephttpd.EepHttpd
	PubKeyPath string
	SecurePath string
	PagePort   string
	page       bool
	up         bool
	prex       bool
}

var err error

func (s *GitSAMTunnel) PRBytes() []byte {
	r := "#!/bin/sh"
	r += "GIT_WORK_TREE=" + s.GitConf.Dir + " git checkout -f"
	return []byte(r)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info != nil {
		return !info.IsDir()
	}
	return false
}

func (s *GitSAMTunnel) ListDirs() ([]os.FileInfo, error){
    unclassedfiles, err := ioutil.ReadDir(s.GitConf.Dir)
	if err != nil {
		return nil, err
	}
    var directories []os.FileInfo
    for _, unclassedfile := range unclassedfiles {
        if unclassedfile.IsDir() {
            log.Println("FoundDir", unclassedfile.Name())
            directories = append(directories, unclassedfile)
        }
    }
	gitdir, err := os.Stat(s.GitConf.Dir)
	if err != nil {
		return nil, err
	}
    if gitdir.IsDir() {
        directories = append(directories, gitdir)
    }
    return directories, nil
}

func (s *GitSAMTunnel) AssurePostRecieve() error {
	if !s.page {
		return nil
	}
    files, err := s.ListDirs()
    if err != nil {
        return err
    }
	for _, dir := range files {
		info, err := os.Stat(filepath.Join(s.GitConf.Dir, dir.Name(), ".git"))
		if err == nil {
			if info.IsDir() {
                log.Println(filepath.Join(s.GitConf.Dir, dir.Name(), ".git"), "is a directory")
                dirpath := filepath.Join(s.GitConf.Dir, dir.Name())
                if strings.HasSuffix(s.GitConf.Dir, dir.Name()){
                    dirpath = s.GitConf.Dir
                }

                cmd := exec.Command("git", "update-server-info")
                cmd.Dir = dirpath
                if err := cmd.Run(); err != nil {
                    return err
                }
                log.Println("Updated git server info for static copy", cmd.Dir)
				if err := os.MkdirAll(filepath.Join(s.GitConf.Dir, "/hooks"), 0755); err != nil {
					return err
				} else {
					if !FileExists(filepath.Join(s.GitConf.Dir, dir.Name(), "hooks", "post-recieve")) {
						if err := ioutil.WriteFile(filepath.Join(s.GitConf.Dir, "hooks", "post-recieve"), s.PRBytes(), 0755); err != nil {
							s.prex = true
							return err
						}
					}
				}
			}else{
                log.Println(filepath.Join(s.GitConf.Dir, dir.Name(), ".git"), "is not a directory")
            }
		}else if os.IsNotExist(err) {
            log.Println(filepath.Join(s.GitConf.Dir, dir.Name(), ".git"), "does not exist")
        }else{
            return err
        }
	}
	return nil
}

func (s *GitSAMTunnel) DeletePostRecieve() error {
	if s.prex {
		files, err := ioutil.ReadDir(s.GitConf.Dir)
		if err != nil {
			return err
		}
		gitdir, err := os.Stat(s.GitConf.Dir)
		if err != nil {
			return err
		}
		files = append(files, gitdir)
		for _, dir := range files {
			if err := os.Remove(filepath.Join(s.GitConf.Dir, dir.Name(), "hooks", "post-recieve")); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *GitSAMTunnel) AssureGitIgnore() error {
	PubKeyPath, err := filepath.Rel(s.GitConf.Dir, s.PubKeyPath)
	if err != nil {
		return err
	}
	if filepath.Dir(s.PubKeyPath) == s.GitConf.Dir {
		if FileExists(filepath.Join(s.GitConf.Dir, "/.gitignore")) {
			if bytes, err := ioutil.ReadFile(filepath.Join(s.GitConf.Dir, "/.gitignore")); err == nil {
				if !strings.Contains(string(bytes), s.PubKeyPath) {
					f, err := os.OpenFile(filepath.Join(s.GitConf.Dir, "/.gitignore"), os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					defer f.Close()
					if _, err = f.WriteString(PubKeyPath); err != nil {
						return err
					}
				}
			} else {
				return err
			}
		} else {
			if err := ioutil.WriteFile(filepath.Join(s.GitConf.Dir, "/.gitignore"), []byte(PubKeyPath), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *GitSAMTunnel) LookupKey(content string) (*gitkit.PublicKey, error) {
	textkey, err := ioutil.ReadFile(f.PubKeyPath)
	if err != nil {
		return nil, err
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey(textkey)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(key.Marshal())
	print := base64.StdEncoding.EncodeToString(hash[:])
	return &gitkit.PublicKey{Fingerprint: print}, nil
}

func (f *GitSAMTunnel) GetType() string {
	return "gitsam"
}

func (f *GitSAMTunnel) ServeParent() {
	log.Println("Starting eepssh server", f.Base32())
	if err = f.SAMForwarder.Serve(); err != nil {
		log.Println(err)
		f.Cleanup()
	}
}

func (f *GitSAMTunnel) ServePage() {
	if err = f.OptPage.Serve(); err != nil {
		log.Println(err)
		f.Cleanup()
	}
}

//Serve starts the SAM connection and and forwards the local host:port to i2p
func (f *GitSAMTunnel) Serve() error {
	if f.Up() {
		go f.ServePage()
		go f.ServeParent()
		log.Println("Starting ssh server", f.Target())
		if err := f.SSH.ListenAndServe(f.Target()); err != nil {
			return err
		}
	}
	return nil
}

func (f *GitSAMTunnel) Up() bool {
	return f.up
}

//Close shuts the whole thing down.
func (f *GitSAMTunnel) Close() error {
	if err := f.DeletePostRecieve(); err != nil {
		return err
	}
	return f.SAMForwarder.Close()
}

func (s *GitSAMTunnel) Load() (samtunnel.SAMTunnel, error) {
	if !s.Up() {
		log.Println("Started putting tunnel up", s.Conf.SaveDirectory)
	}
	f, e := s.SAMForwarder.Load()
	if e != nil {
		return nil, e
	}
	s.Conf.ServeDirectory = s.GitConf.Dir
	s.SAMForwarder = f.(*samforwarder.SAMForwarder)
	s.GitConf.KeyDir = s.SecurePath
	s.SSH = gitkit.NewSSH(s.GitConf)
	s.SSH.PublicKeyLookupFunc = s.LookupKey
	if err := s.AssureGitIgnore(); err != nil {
		return nil, err
	}
	if err := s.AssurePostRecieve(); err != nil {
		return nil, err
	}
	log.Println("Finished putting tunnel up")
	s.up = true
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
	s.Conf = &i2ptunconf.Conf{}
	s.OptPage = &eephttpd.EepHttpd{}
	//s.OptPage.SAMForwarder = &samforwarder.SAMForwarder{}
	s.GitConf = gitkit.Config{}
	s.SSH = &gitkit.SSH{}
	s.page = true
	s.up = false
	log.Println("Initializing gitsam")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	var err error
	if s.SecurePath == s.GitConf.Dir {
		s.SecurePath = filepath.Dir(s.GitConf.Dir) + "/.gitsam_secure"
		s.SAMForwarder.Config().SaveDirectory = s.SecurePath
		s.SAMForwarder.Config().FilePath = s.SecurePath
	}
	conf := *s.Conf
	conf.CloseIdleTime = 6000000
	conf.TargetPort = s.PagePort
	conf.TunName = s.ID() + "-eephttpd"
	conf.SaveDirectory = s.Conf.SaveDirectory
	log.Println("Setting up secure path", s.Conf.SaveDirectory, conf.SaveDirectory)
	if s.OptPage, err = i2ptunhelper.NewEepHttpdFromConf(&conf); err != nil {
		return nil, err
	}
	l, e := s.Load()
	//log.Println("Options loaded", s.Print())
	if e != nil {
		return nil, e
	}
	return l.(*GitSAMTunnel), nil
}
