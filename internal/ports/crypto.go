package ports

import "io"

// CryptoService defines the interface for encryption/decryption operations
type CryptoService interface {
	// DecryptPKCS7 decrypts PKCS7 encrypted data using a private key
	DecryptPKCS7(encryptedData []byte, privateKeyPath string) ([]byte, error)

	// DecryptPKCS7Reader decrypts from a reader and writes to a writer
	DecryptPKCS7Reader(encrypted io.Reader, privateKeyPath string, output io.Writer) error

	// ValidatePrivateKey checks if a private key file is valid
	ValidatePrivateKey(privateKeyPath string) error
}
