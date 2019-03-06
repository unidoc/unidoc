package sighandler

import (
	"bytes"
	"crypto/x509"

	"github.com/gunnsth/pkcs7"
	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/model"
)

type emptyPKCS7Detached struct {
	sigLen      int
	certificate *x509.Certificate
}

// NewEmptyPKCS7Detached creates a new Adobe.PPKMS/Adobe.PPKLite adbe.pkcs7.detached
// signature handler. The generated signature is empty and of size sigLen.
// The certificate parameter may be nil for the signature validation.
func NewEmptyPKCS7Detached(sigLen int, certificate *x509.Certificate) model.SignatureHandler {
	return &emptyPKCS7Detached{
		sigLen:      sigLen,
		certificate: certificate,
	}
}

// InitSignature initialises the PdfSignature.
func (e *emptyPKCS7Detached) InitSignature(sig *model.PdfSignature) error {
	handler := *e
	sig.Handler = &handler
	sig.Filter = core.MakeName("Adobe.PPKLite")
	sig.SubFilter = core.MakeName("adbe.pkcs7.detached")
	sig.Reference = nil
	if sig.Cert != nil {
		sig.Cert = core.MakeString(string(handler.certificate.Raw))
	}

	digest, err := handler.NewDigest(sig)
	if err != nil {
		return err
	}

	return handler.Sign(sig, digest)
}

func (e *emptyPKCS7Detached) getCertificate(sig *model.PdfSignature) (*x509.Certificate, error) {
	certificate := e.certificate
	if certificate == nil {
		certData := sig.Cert.(*core.PdfObjectString).Bytes()
		certs, err := x509.ParseCertificates(certData)
		if err != nil {
			return nil, err
		}
		certificate = certs[0]
	}
	return certificate, nil
}

// NewDigest creates a new digest.
func (e *emptyPKCS7Detached) NewDigest(sig *model.PdfSignature) (model.Hasher, error) {
	return bytes.NewBuffer(nil), nil
}

// Validate validates PdfSignature.
func (e *emptyPKCS7Detached) Validate(sig *model.PdfSignature, digest model.Hasher) (model.SignatureValidationResult, error) {
	p7, err := pkcs7.Parse(sig.Contents.Bytes())
	if err != nil {
		return model.SignatureValidationResult{}, err
	}

	buffer := digest.(*bytes.Buffer)
	p7.Content = buffer.Bytes()
	if err = p7.Verify(); err != nil {
		return model.SignatureValidationResult{}, err
	}

	return model.SignatureValidationResult{
		IsSigned:   true,
		IsVerified: true,
	}, nil
}

// Sign sets the Contents fields for the PdfSignature.
func (e *emptyPKCS7Detached) Sign(sig *model.PdfSignature, digest model.Hasher) error {
	sigLen := e.sigLen
	if e.sigLen <= 0 {
		sigLen = 8092
	}

	sig.Contents = core.MakeHexString(string(make([]byte, sigLen)))
	return nil
}

// IsApplicable returns true if the signature handler is applicable for the PdfSignature.
func (e *emptyPKCS7Detached) IsApplicable(sig *model.PdfSignature) bool {
	if sig == nil || sig.Filter == nil || sig.SubFilter == nil {
		return false
	}
	return (*sig.Filter == "Adobe.PPKMS" || *sig.Filter == "Adobe.PPKLite") && *sig.SubFilter == "adbe.pkcs7.detached"
}
