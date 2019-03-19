package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

func loadPem(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	block, _ := pem.Decode(b)
	return block.Bytes
}

func main() {
	// load ca
	caKey, _ := x509.ParsePKCS1PrivateKey(loadPem("ca.key"))
	caCrt, _ := x509.ParseCertificate(loadPem("ca.crt"))

	// load csr
	csr, _ := x509.ParseCertificateRequest(loadPem("server.csr"))

	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	now := time.Now()
	crt, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		Subject:      csr.Subject,
		Issuer:       caCrt.Subject,
		SerialNumber: serial,
		NotBefore:    now.UTC(),
		NotAfter:     now.AddDate(1, 0, 0).UTC(),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     csr.DNSNames,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}, caCrt, csr.PublicKey, caKey)
	if err != nil {
		log.Fatal(err)
	}

	crtFp, err := os.Create("server.crt")
	if err != nil {
		log.Fatal(err)
	}
	defer crtFp.Close()
	pem.Encode(crtFp, &pem.Block{Type: "CERTIFICATE", Bytes: crt})
}
