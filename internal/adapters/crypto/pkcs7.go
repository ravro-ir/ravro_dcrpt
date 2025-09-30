package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"ravro_dcrpt/internal/ports"

	"go.mozilla.org/pkcs7"
)

// PKCS7Service implements the CryptoService interface using Pure Go
type PKCS7Service struct{}

// NewPKCS7Service creates a new PKCS7 service
func NewPKCS7Service() ports.CryptoService {
	return &PKCS7Service{}
}

// DecryptPKCS7 decrypts PKCS7 encrypted data using a private key
func (s *PKCS7Service) DecryptPKCS7(encryptedData []byte, privateKeyPath string) ([]byte, error) {
	// Load private key
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Parse PKCS7 structure
	p7, err := pkcs7.Parse(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS7 data: %w", err)
	}

	// Decrypt the data
	decrypted, err := p7.Decrypt(nil, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decrypted, nil
}

// DecryptPKCS7Reader decrypts from a reader and writes to a writer
func (s *PKCS7Service) DecryptPKCS7Reader(encrypted io.Reader, privateKeyPath string, output io.Writer) error {
	// Read all encrypted data
	encryptedData, err := io.ReadAll(encrypted)
	if err != nil {
		return fmt.Errorf("failed to read encrypted data: %w", err)
	}

	// Decrypt
	decrypted, err := s.DecryptPKCS7(encryptedData, privateKeyPath)
	if err != nil {
		return err
	}

	// Write to output
	_, err = output.Write(decrypted)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data: %w", err)
	}

	return nil
}

// ValidatePrivateKey checks if a private key file is valid
func (s *PKCS7Service) ValidatePrivateKey(privateKeyPath string) error {
	_, err := loadPrivateKey(privateKeyPath)
	return err
}

// loadPrivateKey loads an RSA private key from a PEM file
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Read the key file
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse private key
	var privateKey *rsa.PrivateKey

	// Try PKCS1 format first
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		privateKey = key
	} else {
		// Try PKCS8 format
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		var ok bool
		privateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an RSA private key")
		}
	}

	return privateKey, nil
}
