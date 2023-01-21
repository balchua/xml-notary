package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

type MemoryX509KeyStore struct {
	privateKey *rsa.PrivateKey
	cert       []byte
}

func (ks *MemoryX509KeyStore) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return ks.privateKey, ks.cert, nil
}

// func loadKeyPair(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey) {
func loadKeyPair(certFile, keyFile string) ([]byte, *rsa.PrivateKey) {
	cf, e := ioutil.ReadFile(certFile)
	if e != nil {
		fmt.Println("cfload:", e.Error())
		os.Exit(1)
	}

	kf, e := ioutil.ReadFile(keyFile)
	if e != nil {
		fmt.Println("kfload:", e.Error())
		os.Exit(1)
	}
	cpb, cr := pem.Decode(cf)

	fmt.Println(string(cr))
	kpb, kr := pem.Decode(kf)
	fmt.Println(string(kr))
	crt, e := x509.ParseCertificate(cpb.Bytes)
	if crt != nil {
		fmt.Println("certificate is good")
	}

	if e != nil {
		fmt.Println("parsex509:", e.Error())
		os.Exit(1)
	}
	key, e := x509.ParsePKCS1PrivateKey(kpb.Bytes)
	if e != nil {
		fmt.Println("parsekey:", e.Error())
		os.Exit(1)
	}
	return cpb.Bytes, key
}

func main() {

	certBytes, key := loadKeyPair("generated-certs/ca.pem", "generated-certs/ca-key.pem")

	ks := &MemoryX509KeyStore{
		privateKey: key,
		cert:       certBytes,
	}

	ctx := dsig.NewDefaultSigningContext(ks)
	elementToSign := &etree.Element{
		Tag: "ExampleElement",
	}
	elementToSign.CreateAttr("ID", "id1234")

	// Sign the element
	signedElement, err := ctx.SignEnveloped(elementToSign)
	if err != nil {
		panic(err)
	}

	// Serialize the signed element. It is important not to modify the element
	// after it has been signed - even pretty-printing the XML will invalidate
	// the signature.
	doc := etree.NewDocument()
	doc.SetRoot(signedElement)
	err = doc.WriteToFile("sample-signed.xml")
	if err != nil {
		panic(err)
	}

}
