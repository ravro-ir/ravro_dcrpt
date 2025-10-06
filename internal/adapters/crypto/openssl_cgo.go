//go:build linux || darwin
// +build linux darwin

package crypto

/*
#cgo linux pkg-config: openssl
#cgo darwin pkg-config: openssl
#cgo darwin CFLAGS: -I/opt/homebrew/include -I/usr/local/include -DOPENSSL_API_COMPAT=0x30000000
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lssl -lcrypto

#include <openssl/bio.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <openssl/pkcs7.h>
#include <openssl/pem.h>
#include <stdlib.h>
#include <string.h>

// Helper function to decrypt PKCS7 data
static int decrypt_pkcs7(const unsigned char* encrypted_data, int encrypted_len,
                          const char* key_path,
                          unsigned char** decrypted_data, int* decrypted_len) {
    BIO *bio_in = NULL, *bio_key = NULL, *bio_out = NULL;
    EVP_PKEY *pkey = NULL;
    PKCS7 *p7 = NULL;
    int ret = 0;

    // Initialize OpenSSL
    #if OPENSSL_VERSION_NUMBER < 0x10100000L
        OpenSSL_add_all_algorithms();
        ERR_load_crypto_strings();
    #endif

    // Create BIO from encrypted data
    bio_in = BIO_new_mem_buf(encrypted_data, encrypted_len);
    if (!bio_in) {
        return -1;
    }

    // Read private key
    bio_key = BIO_new_file(key_path, "rb");
    if (!bio_key) {
        BIO_free(bio_in);
        return -2;
    }

    pkey = PEM_read_bio_PrivateKey(bio_key, NULL, NULL, NULL);
    if (!pkey) {
        BIO_free(bio_in);
        BIO_free(bio_key);
        return -3;
    }

    // Parse PKCS7 structure (DER format)
    p7 = d2i_PKCS7_bio(bio_in, NULL);
    if (!p7) {
        EVP_PKEY_free(pkey);
        BIO_free(bio_in);
        BIO_free(bio_key);
        return -4;
    }

    // Create output BIO
    bio_out = BIO_new(BIO_s_mem());
    if (!bio_out) {
        PKCS7_free(p7);
        EVP_PKEY_free(pkey);
        BIO_free(bio_in);
        BIO_free(bio_key);
        return -5;
    }

    // Decrypt
    if (!PKCS7_decrypt(p7, pkey, NULL, bio_out, 0)) {
        BIO_free(bio_out);
        PKCS7_free(p7);
        EVP_PKEY_free(pkey);
        BIO_free(bio_in);
        BIO_free(bio_key);
        return -6;
    }

    // Get decrypted data length
    *decrypted_len = BIO_pending(bio_out);

    // Allocate memory for decrypted data
    *decrypted_data = (unsigned char*)malloc(*decrypted_len);
    if (!*decrypted_data) {
        BIO_free(bio_out);
        PKCS7_free(p7);
        EVP_PKEY_free(pkey);
        BIO_free(bio_in);
        BIO_free(bio_key);
        return -7;
    }

    // Read decrypted data
    BIO_read(bio_out, *decrypted_data, *decrypted_len);

    // Cleanup
    BIO_free(bio_out);
    PKCS7_free(p7);
    EVP_PKEY_free(pkey);
    BIO_free(bio_in);
    BIO_free(bio_key);

    return 0;
}
*/
import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"ravro_dcrpt/internal/ports"
)

// OpenSSLCGOService implements CryptoService using OpenSSL via CGO
type OpenSSLCGOService struct{}

// NewOpenSSLCGOService creates a new OpenSSL CGO-based crypto service
func NewOpenSSLCGOService() ports.CryptoService {
	return &OpenSSLCGOService{}
}

// DecryptPKCS7 decrypts PKCS7 encrypted data using OpenSSL via CGO
func (s *OpenSSLCGOService) DecryptPKCS7(encryptedData []byte, privateKeyPath string) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("encrypted data is empty")
	}

	// Convert Go string to C string
	cKeyPath := C.CString(privateKeyPath)
	defer C.free(unsafe.Pointer(cKeyPath))

	// Prepare variables for C function
	var decryptedData *C.uchar
	var decryptedLen C.int

	// Call C function
	result := C.decrypt_pkcs7(
		(*C.uchar)(unsafe.Pointer(&encryptedData[0])),
		C.int(len(encryptedData)),
		cKeyPath,
		&decryptedData,
		&decryptedLen,
	)

	if result != 0 {
		return nil, fmt.Errorf("decryption failed with error code: %d", result)
	}

	// Convert C data to Go slice
	goData := C.GoBytes(unsafe.Pointer(decryptedData), decryptedLen)

	// Free C memory
	C.free(unsafe.Pointer(decryptedData))

	return goData, nil
}

// DecryptPKCS7Reader decrypts from a reader and writes to a writer
func (s *OpenSSLCGOService) DecryptPKCS7Reader(encrypted io.Reader, privateKeyPath string, output io.Writer) error {
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
func (s *OpenSSLCGOService) ValidatePrivateKey(privateKeyPath string) error {
	// Check if file exists
	if _, err := os.Stat(privateKeyPath); err != nil {
		return fmt.Errorf("key file not found: %w", err)
	}

	// Try to read the key
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	// Basic validation - check if it's a PEM file
	if len(keyData) < 20 {
		return fmt.Errorf("key file is too small")
	}

	keyStr := string(keyData)
	if !contains(keyStr, "BEGIN") || !contains(keyStr, "PRIVATE KEY") {
		return fmt.Errorf("invalid key format - not a PEM private key")
	}

	return nil
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != "" &&
		len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
