//go:build darwin
// +build darwin

package utils

// #cgo pkg-config: openssl
// #cgo CFLAGS: -DOPENSSL_API_COMPAT=0x30000000
// #include <openssl/bio.h>
// #include <openssl/err.h>
// #include <openssl/evp.h>
// #include <openssl/pkcs7.h>
// #include <openssl/ssl.h>
// #include <stdlib.h>
//
// char* smime_decrypt(const char* input_file, const char* key_file, const char* output_file) {
//     OpenSSL_add_all_algorithms();
//     ERR_load_crypto_strings();
//
//     BIO *in = NULL, *out = NULL, *key = NULL;
//     EVP_PKEY *pkey = NULL;
//     PKCS7 *p7 = NULL;
//     char* error_msg = NULL;
//
//     // Open input file
//     in = BIO_new_file(input_file, "rb");
//     if (!in) {
//         error_msg = strdup("Could not open input file");
//         goto err;
//     }
//
//     // Read private key
//     key = BIO_new_file(key_file, "rb");
//     if (!key) {
//         error_msg = strdup("Could not open key file");
//         goto err;
//     }
//
//     // Read private key
//     pkey = PEM_read_bio_PrivateKey(key, NULL, 0, NULL);
//     if (!pkey) {
//         error_msg = strdup("Could not read private key");
//         goto err;
//     }
//
//     // Open output file
//     out = BIO_new_file(output_file, "wb");
//     if (!out) {
//         error_msg = strdup("Could not open output file");
//         goto err;
//     }
//
//     // Read PKCS7 structure (DER format)
//     p7 = d2i_PKCS7_bio(in, NULL);
//     if (!p7) {
//         error_msg = strdup("Could not read PKCS7 structure");
//         goto err;
//     }
//
//     // Decrypt
//     if (!PKCS7_decrypt(p7, pkey, NULL, out, PKCS7_BINARY)) {
//         error_msg = strdup("Decryption failed");
//         goto err;
//     }
//
// err:
//     // Cleanup
//     if (in) BIO_free(in);
//     if (out) BIO_free(out);
//     if (key) BIO_free(key);
//     if (p7) PKCS7_free(p7);
//     if (pkey) EVP_PKEY_free(pkey);
//
//     ERR_free_strings();
//     EVP_cleanup();
//
//     return error_msg;
// }
import "C"
import (
	"fmt"
	"io/ioutil"
	"unsafe"
)

// SslDecrypt decrypts a file using OpenSSL via CGo
func SslDecrypt(name, filename, keyFixPath string) (out string, errOut error) {
	// Convert Go strings to C strings
	cName := C.CString(name)
	cKeyFile := C.CString(keyFixPath)
	cOutputFile := C.CString(filename)

	// Defer freeing C strings
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cKeyFile))
	defer C.free(unsafe.Pointer(cOutputFile))

	// Call CGo decryption function with BINARY flag
	errorMsg := C.smime_decrypt(cName, cKeyFile, cOutputFile)

	// Check for decryption errors
	if errorMsg != nil {
		defer C.free(unsafe.Pointer(errorMsg))
		return "", fmt.Errorf(C.GoString(errorMsg))
	}

	// Read the decrypted file
	outputBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read decrypted file: %v", err)
	}

	// Read the output (matching original behavior)
	out = string(outputBytes)

	return out, nil
}
