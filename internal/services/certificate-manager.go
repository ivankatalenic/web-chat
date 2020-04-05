package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/ivankatalenic/web-chat/internal/config"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// SelfSignedCertFileName name of the file holding a self-signed certificate
const SelfSignedCertFileName = "self-signed-cert.pem"
// SelfSignedKeyFileName name of the file holding a private key for the self-signed certificate
const SelfSignedKeyFileName = "self-signed-key.pem"

// CertificateManager holds the file names for a certificate, and a private key
type CertificateManager struct {
	CertFilePath string
	KeyFilePath  string
}

// NewCertificateManager creates a new initialized certificate manager
func NewCertificateManager(config config.TLSConfig) *CertificateManager {
	if len(config.CertFilePath) == 0 || len(config.KeyFilePath) == 0 {
		_, certErr := os.Stat(SelfSignedCertFileName)
		_, keyErr := os.Stat(SelfSignedKeyFileName)
		if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
			createNewSelfSignedCertificate()
		}

		return &CertificateManager{
			CertFilePath: SelfSignedCertFileName,
			KeyFilePath:  SelfSignedKeyFileName,
		}
	}
	return &CertificateManager{
		CertFilePath: config.CertFilePath,
		KeyFilePath:  config.KeyFilePath,
	}
}

func createNewSelfSignedCertificate() {
	var err error

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	pub := priv.Public()
	if err != nil {
		log.Fatalf("Failed to generate a private key: %v", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(2 * 30 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate a serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Localhost Organization"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split("localhost,127.0.0.1", ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, priv)
	if err != nil {
		log.Fatalf("Failed to create a certificate: %v", err)
	}

	certOut, err := os.Create(SelfSignedCertFileName)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %v", SelfSignedCertFileName, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to %s: %v", SelfSignedCertFileName, err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing %s: %v", SelfSignedCertFileName, err)
	}
	log.Print("Wrote the certificate")

	keyOut, err := os.OpenFile(SelfSignedKeyFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %v", SelfSignedKeyFileName, err)
		return
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal the private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to %s: %v", SelfSignedKeyFileName, err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing %s: %v", SelfSignedKeyFileName, err)
	}
	log.Print("Wrote the key file\n")
}
