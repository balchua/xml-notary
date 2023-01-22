package notary

import (
	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

type Notary struct {
	signingCtx *dsig.SigningContext
}

func New(ks dsig.X509KeyStore) (*Notary, error) {
	signCtx := dsig.NewDefaultSigningContext(ks)
	return &Notary{
		signingCtx: signCtx,
	}, nil
}

func (n *Notary) SignEnvelope(elementToSign *etree.Element) (*etree.Element, error) {
	// Sign the element
	signedElement, err := n.signingCtx.SignEnveloped(elementToSign)
	if err != nil {
		return nil, err
	}

	return signedElement, nil
}
