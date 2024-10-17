package data

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"

	"golang.org/x/crypto/ssh"
)

// NewSigner generates a new ed25519 keypair
func NewSigner(path, comment string) error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	privateKey, err := ssh.MarshalPrivateKey(priv, comment)
	if err != nil {
		return err
	}
	privatePem := pem.EncodeToMemory(privateKey)

	err = os.WriteFile(path, privatePem, 0600)
	if err != nil {
		return err
	}

	publicKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+".pub", ssh.MarshalAuthorizedKey(publicKey), 0600)
	if err != nil {
		return err
	}

	return nil
}

func GetUserSigner() (ssh.Signer, error) {
	privateKey, err := os.ReadFile(UserKey)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func GetPublicUserKey() ([]byte, error) {
	signer, err := GetUserSigner()
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(signer.PublicKey()), nil
}

func GetHostSigner() (ssh.Signer, error) {
	privateKey, err := os.ReadFile(HostCAKey)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func GetPublicHostKey() ([]byte, error) {
	signer, err := GetHostSigner()
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(signer.PublicKey()), nil
}

func GenerateSSHKeys(rotate bool) error {
	if rotate {
		if err := os.Remove(HostCAKey); err != nil {
			return err
		}
		if err := os.Remove(UserKey); err != nil {
			return err
		}
	}

	if _, err := os.Stat(HostCAKey); os.IsNotExist(err) {
		if err = NewSigner(HostCAKey, "ca@ssh-nexus"); err != nil {
			return err
		}
	}

	if _, err := os.Stat(UserKey); os.IsNotExist(err) {
		if err = NewSigner(UserKey, "user@ssh-nexus"); err != nil {
			return err
		}
	}
	return nil
}
