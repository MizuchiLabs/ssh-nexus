package data

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/MizuchiLabs/ssh-nexus/tools/util"
)

// GenerateServerCA generates a CA certificate and private key pair.
func GenerateServerCA() error {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// Prepare CA certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "CA ssh-nexus"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Create CA certificate
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		publicKey,
		privateKey,
	)
	if err != nil {
		return err
	}

	// Write CA certificate to file
	certFile, err := os.Create(ServerCaCert)
	if err != nil {
		return err
	}
	defer certFile.Close()
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return err
	}

	// Write CA private key to file
	keyFile, err := os.Create(ServerCaKey)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	return pem.Encode(keyFile, &pem.Block{Type: "ED25519 PRIVATE KEY", Bytes: privateKeyBytes})
}

// GenerateKeyPair generates a certificate and private key pair signed by the given CA.
func GenerateKeyPair(domain string) error {
	caCertPEM, err := os.ReadFile(ServerCaCert)
	if err != nil {
		return err
	}
	caKeyPEM, err := os.ReadFile(ServerCaKey)
	if err != nil {
		return err
	}
	caCert, _ := pem.Decode(caCertPEM)
	caKey, _ := pem.Decode(caKeyPEM)
	caCertParsed, err := x509.ParseCertificate(caCert.Bytes)
	if err != nil {
		return err
	}
	caKeyParsed, err := x509.ParsePKCS8PrivateKey(caKey.Bytes)
	if err != nil {
		return err
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// Prepare certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "ssh-nexus"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
			util.GetOutboundIP(),
		},
		DNSNames:              []string{domain},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// Create certificate signed by CA
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		caCertParsed,
		publicKey,
		caKeyParsed,
	)
	if err != nil {
		return err
	}

	// Write certificate to file
	certFile, err := os.Create(ServerCert)
	if err != nil {
		return err
	}
	defer certFile.Close()
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return err
	}

	// Write private key to file
	keyFile, err := os.Create(ServerKey)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	return pem.Encode(keyFile, &pem.Block{Type: "ED25519 PRIVATE KEY", Bytes: privateKeyBytes})
}

func GetPublicServerCA() ([]byte, error) {
	caCertPEM, err := os.ReadFile(ServerCaCert)
	if err != nil {
		return nil, err
	}
	return caCertPEM, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GenerateTLS(host string) error {
	// check if domain is valid
	if util.IsValidDomain(host) {
		return fmt.Errorf("invalid domain: %s", host)
	}

	if !fileExists(ServerCaCert) || !fileExists(ServerCaKey) {
		return GenerateServerCA()
	}

	if !fileExists(ServerCert) || !fileExists(ServerKey) {
		return GenerateKeyPair(host)
	}

	// Check CA cert expiration
	caCert, err := tls.LoadX509KeyPair(ServerCaCert, ServerCaKey)
	if err != nil {
		return err
	}
	x509CaCert, err := x509.ParseCertificate(caCert.Certificate[0])
	if err != nil {
		return err
	}
	if time.Now().After(x509CaCert.NotAfter) {
		return GenerateServerCA()
	}

	// Check server cert expiration
	serverCert, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	if err != nil {
		return err
	}
	x509Cert, err := x509.ParseCertificate(serverCert.Certificate[0])
	if err != nil {
		return err
	}
	if time.Now().After(x509Cert.NotAfter) {
		return GenerateKeyPair(host)
	}

	slog.Info("Certificates and keys generated successfully")
	return nil
}

func RegenerateKeys(host string) error {
	if err := os.Remove(ServerCert); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(ServerKey); err != nil && !os.IsNotExist(err) {
		return err
	}

	return GenerateKeyPair(host)
}
