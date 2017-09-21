package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln(err)
	}

	bin := x509.MarshalPKCS1PrivateKey(privateKey)

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bin,
	}
	pem.Encode(os.Stdout, block)

	bin, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bin,
	}
	pem.Encode(os.Stdout, block)
}
