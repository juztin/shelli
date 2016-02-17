package shelli

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// ConfigFunc is just a function signature for retrieving ssh configuration functions.
type ConfigFunc func() (*ssh.ClientConfig, error)

func getKey(file string) (ssh.Signer, error) {
	var signer ssh.Signer
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return signer, err
	}
	signer, err = ssh.ParsePrivateKey(buf)
	return signer, err
}

// ConfigForCert returns a configuration func for the given username and certificate.
func ConfigForCert(user, cert string) ConfigFunc {
	return func() (*ssh.ClientConfig, error) {
		key, err := getKey(cert)
		if err != nil {
			return nil, err
		}
		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
		}
		return config, nil
	}
}

// ConfigForPassword returns a configuration func for the given username and password.
func ConfigForPassword(user, password string) ConfigFunc {
	return func() (*ssh.ClientConfig, error) {
		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
		}
		return config, nil
	}
}

// NewClient returns a new SSH client for the given host and configuration.
func NewClient(host string, config ConfigFunc) (*ssh.Client, error) {
	cfg, err := config()
	if err != nil {
		return nil, err
	}
	return ssh.Dial("tcp", host, cfg)
}
