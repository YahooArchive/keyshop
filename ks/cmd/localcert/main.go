// localcert is a small utility that generates a CA, signs a certificate
// with it, and then throws away the key.

// NOTE: Original version derived from Go distribution.

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
)

var (
	curve = elliptic.P256()
)

const (
	expireDays = 30
)

func serialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	return serialNumber
}

func writePem(fn string, t string, derBytes []byte) {
	certOut, err := os.Create(fn)
	if err != nil {
		log.Fatalf("failed to open %s for writing: %s", fn, err)
	}
	pem.Encode(certOut, &pem.Block{Type: t, Bytes: derBytes})
	certOut.Close()
	log.Infof("written %s\n", fn)
}

func generate(c *cli.Context) {

}

func main() {
	app := cli.NewApp()
	app.Name = "localcert"
	app.Usage = "Generate a temporary self-signed cert."

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "validfor",
			Value: 30,
		},
		&cli.StringFlag{
			Name:  "ou",
			Value: "You're clearly not reading the source",
		},
		&cli.StringFlag{
			Name:  "organization",
			Value: "if your certificate looks this silly.",
		},
		&cli.IntFlag{
			Name:  "certv",
			Value: 1 << 24,
		},
		&cli.StringFlag{
			Name:  "prefix",
			Value: "",
		},
	}

	app.Action = func(c *cli.Context) {
		domains := c.Args()
		if len(domains) == 0 {
			return
		}
		ou := strings.Split(c.String("ou"), ",")
		organization := strings.Split(c.String("organization"), ",")
		version := c.Int("certv")
		expireDays := c.Int("validfor")
		if expireDays == 0 {
			expireDays = 30
		}
		prefix := c.String("prefix")
		if prefix == "" {
			prefix = domains[0] + "."
		}

		nameTemplate := &pkix.Name{
			Country:            []string{"US"},
			Organization:       organization,
			OrganizationalUnit: ou,
			SerialNumber:       strconv.Itoa(version),
		}

		notBefore := time.Now().UTC()
		notAfter := notBefore.Add(time.Duration(expireDays*24) * time.Hour)
		fmt.Printf("notBefore: %v, notAfter: %v", notBefore, notAfter)

		caName := *nameTemplate
		caName.CommonName = fmt.Sprintf("TEST-CA: %s (%d)", domains[0], version)

		caTemplate := &x509.Certificate{
			// Root certificate.
			Issuer:  caName,
			Subject: caName,

			// Basic constraints.
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
			IsCA:                  true,
			MaxPathLen:            0,
			MaxPathLenZero:        true,
			BasicConstraintsValid: true,

			// Extended key usage.
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageOCSPSigning},

			// DNS name constraint; alas, not recognized by OS X.
			PermittedDNSDomains: domains,
			//PermittedDNSDomainsCritical: true,

			// Common fields:
			NotBefore:    notBefore,
			NotAfter:     notAfter,
			SerialNumber: serialNumber(),
			Version:      version,
		}

		serverName := *nameTemplate
		serverName.CommonName = fmt.Sprintf("%s", domains[0])

		serverTemplate := &x509.Certificate{
			KeyUsage: x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			// Server certificate.
			Issuer:  caName,
			Subject: serverName,

			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth, // needed for WebSockets
			},
			// SNI
			DNSNames: domains,

			// Common fields:
			NotBefore:    notBefore,
			NotAfter:     notAfter,
			SerialNumber: serialNumber(),
			Version:      version,
		}

		// Generate the two private keys.
		caPriv, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			log.Fatalf("error generating private key: %s", err)
		}
		serverPriv, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			log.Fatalf("error generating private key: %s", err)
		}

		caDer, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPriv.PublicKey, caPriv)
		if err != nil {
			log.Fatalf("%s", err)
		}
		writePem(prefix+"ca.pem", "CERTIFICATE", caDer)

		ca, err := x509.ParseCertificate(caDer)
		if err != nil {
			log.Fatalf("%s", err)
		}

		serverDer, err := x509.CreateCertificate(rand.Reader, serverTemplate, ca, &serverPriv.PublicKey, caPriv)
		if err != nil {
			log.Fatalf("%s", err)
		}
		writePem(prefix+"server.pem", "CERTIFICATE", serverDer)

		b, err := x509.MarshalECPrivateKey(serverPriv)
		if err != nil {
			log.Fatalf("%s", err)
		}
		writePem(prefix+"privatekey.pem", "EC PRIVATE KEY", b)
		return
	}
	app.Run(os.Args)
}
