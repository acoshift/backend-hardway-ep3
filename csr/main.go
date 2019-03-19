package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"os"
)

// openssl req -new -newkey rsa:2048 -nodes -keyout server.key -out server.csr
// openssl req -in server.csr -noout -text
func main() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		DNSNames: []string{
			"localhost",
			"www.localhost",
		},
		SignatureAlgorithm: x509.SHA512WithRSA,
	}, key)
	if err != nil {
		log.Fatal(err)
	}

	keyFp, err := os.Create("server.key")
	if err != nil {
		log.Fatal(err)
	}
	defer keyFp.Close()
	pem.Encode(keyFp, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	csrFp, err := os.Create("server.csr")
	if err != nil {
		log.Fatal(err)
	}
	defer csrFp.Close()
	pem.Encode(csrFp, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})
}
