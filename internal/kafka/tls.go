package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// NewTLSConfig creates a TLS config from a PEM encoded CA certificate string.
// Used for connecting to Aiven Kafka with SASL_SSL authentication.
func NewTLSConfig(caCert string) (*tls.Config, error) {
	// Try base64 decode first (for production env vars)
	decoded, err := base64.StdEncoding.DecodeString(caCert)
	if err != nil {
		// If base64 decode fails, assume it's raw PEM (local development)
		decoded = []byte(caCert)
	}

	caCertPool := x509.NewCertPool()
	block, _ := pem.Decode(decoded)
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
