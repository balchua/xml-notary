package certmgr

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type FileBasedX509KeyStore struct {
	privateKey *rsa.PrivateKey
	cert       []byte
	certfile   string
	keyfile    string
}

func loadKeyPair(certFile, keyFile string) ([]byte, *rsa.PrivateKey, error) {
	cf, e := ioutil.ReadFile(certFile)
	if e != nil {
		return nil, nil, fmt.Errorf("unable to load certificate: %v", e)
	}

	kf, e := ioutil.ReadFile(keyFile)
	if e != nil {
		return nil, nil, fmt.Errorf("unable to load private key: %v", e)
	}
	cpb, _ := pem.Decode(cf)

	kpb, _ := pem.Decode(kf)

	key, e := x509.ParsePKCS1PrivateKey(kpb.Bytes)
	if e != nil {
		return nil, nil, fmt.Errorf("unable to parse private key: %v", e)
	}
	return cpb.Bytes, key, nil
}

func New(certfile string, keyfile string) (*FileBasedX509KeyStore, error) {
	certBytes, privateKey, err := loadKeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	return &FileBasedX509KeyStore{
		certfile:   certfile,
		keyfile:    keyfile,
		privateKey: privateKey,
		cert:       certBytes,
	}, nil
}

func (ks *FileBasedX509KeyStore) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return ks.privateKey, ks.cert, nil
}
