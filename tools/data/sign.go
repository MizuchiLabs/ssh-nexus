package data

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func SignHostCertificate(
	publicKey string,
	hostname string,
	expiration time.Duration,
) ([]byte, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	if expiration == 0 {
		return nil, fmt.Errorf("invalid expiration time")
	}

	cert := &ssh.Certificate{
		Key:         pub,
		Serial:      uint64(time.Now().Unix()),
		CertType:    ssh.HostCert,
		KeyId:       hostname + "@ssh-nexus",
		ValidAfter:  uint64(time.Now().Unix()),
		ValidBefore: uint64(time.Now().Add(expiration).Unix()),
		Permissions: ssh.Permissions{
			CriticalOptions: make(map[string]string),
			Extensions:      make(map[string]string),
		},
	}

	signer, err := GetHostSigner()
	if err != nil {
		return nil, err
	}

	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(cert), nil
}

func SignUserCertificate(
	publicKey string,
	principal string,
	expiration time.Duration,
) ([]byte, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	if expiration == 0 {
		return nil, fmt.Errorf("invalid expiration time")
	}

	cert := &ssh.Certificate{
		Key:             pub,
		Serial:          uint64(time.Now().Unix()),
		CertType:        ssh.UserCert,
		KeyId:           principal + "@ssh-nexus",
		ValidPrincipals: []string{principal},
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Add(expiration).Unix()),
		Permissions: ssh.Permissions{
			CriticalOptions: make(map[string]string),
			Extensions: map[string]string{
				"permit-X11-forwarding":   "",
				"permit-agent-forwarding": "",
				"permit-port-forwarding":  "",
				"permit-pty":              "",
				"permit-user-rc":          "",
			},
		},
	}

	signer, err := GetUserSigner()
	if err != nil {
		return nil, err
	}

	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		return nil, err
	}

	certBytes := ssh.MarshalAuthorizedKey(cert)
	return certBytes, nil
}
