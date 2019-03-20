package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

// openssl x509 -in ca.crt -text -noout
func main() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	serial, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(159), nil))
	now := time.Now()

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: "Example Co., Ltd.",
			Country:    []string{"TH"},
			Province:   []string{"Bangkok"},
		},
		NotBefore:          now.UTC(),
		NotAfter:           now.AddDate(1, 0, 0).UTC(),
		IsCA:               true,
		KeyUsage:           x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		SignatureAlgorithm: x509.SHA512WithRSA,
	}
	crt, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		log.Fatal(err)
	}

	keyFp, err := os.Create("ca.key")
	if err != nil {
		log.Fatal(err)
	}
	defer keyFp.Close()
	pem.Encode(keyFp, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	crtFp, err := os.Create("ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	defer crtFp.Close()
	pem.Encode(crtFp, &pem.Block{Type: "CERTIFICATE", Bytes: crt})
}
