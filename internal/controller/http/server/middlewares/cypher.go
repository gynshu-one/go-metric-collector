package middlewares

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/rs/zerolog/log"
)

var privateKey *rsa.PrivateKey
var once = sync.Once{}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("could not decode PEM data")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PrivateKey), nil
}

func decryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
	once.Do(func() {
		var err error
		// Read the private key file
		privateKey, err = loadPrivateKey(config.GetConfig().CryptoKey)
		if err != nil {
			log.Error().Err(err).Msg("Could not load private key")
		}
	})

	parts := strings.Split(string(ciphertext), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid input string")
	}

	encryptedAESKey, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	ciphertext, err = base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// Decrypt the AES key with the RSA private key
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedAESKey, nil)
	if err != nil {
		return nil, err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Check the length of the ciphertext
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Split the nonce and the actual ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt the data using AES-GCM
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DecryptMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.GetConfig().CryptoKey != "" {
			// Read the encrypted data from the request
			encryptedData, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Error().Err(err).Msg("Could not read request body")
				return
			}

			// Decrypt the data
			decryptedData, err := decryptWithPrivateKey(encryptedData)
			if err != nil {
				log.Error().Err(err).Msg("Could not decrypt data")
				return
			}

			// Replace the request body with the decrypted data
			c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedData))

			c.Next()
		}
		c.Next()
	}
}
