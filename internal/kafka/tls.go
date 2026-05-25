package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// NewTLSConfig creates a TLS config from a PEM encoded CA certificate string.
// Used for connecting to Aiven Kafka with SASL_SSL authentication.
func NewTLSConfig(caCert string) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	block, _ := pem.Decode([]byte(caCert))
	if block == nil {
		return nil, fmt.Errorf("failed to decode CA certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}
	caCertPool.AddCert(cert)
	return &tls.Config{RootCAs: caCertPool}, nil
}
