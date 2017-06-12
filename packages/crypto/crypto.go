package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

type CryptoProvider int
type HashProvider int

const (
	MD5 HashProvider = iota
	SHA256
	DoubleSHA256
)

const (
	UnknownProvider = "Unknown provider"
)

const (
	AESCFB CryptoProvider = iota
)

func HashString(msg string, prov HashProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn("Hashing of empty string")
	}
	bytes := []byte(msg)
	return HashBytes(bytes, prov)
}

func HashBytes(msg []byte, prov HashProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn("Hashing of empty array")
	}
	switch prov {
	case MD5:
		return hashMD5(msg), nil
	case SHA256:
		return hashSHA256(msg), nil
	case DoubleSHA256:
		return hashDoubleSHA256(msg), nil
	default:
		return nil, errors.New(UnknownProvider)
	}
}

func EncryptString(msg string, key []byte, iv []byte, prov CryptoProvider) ([]byte, []byte, error) {
	if len(msg) == 0 {
		log.Warn("Encrypting of empty string")
	}
	bytes := []byte(msg)
	return EncryptBytes(bytes, key, iv, prov)
}

func EncryptBytes(msg []byte, key []byte, iv []byte, prov CryptoProvider) ([]byte, []byte, error) {
	if len(msg) == 0 {
		log.Warn("Encrypting of empty array")
	}
	switch prov {
	case AESCFB:
		return encryptCFB(msg, key, iv)
	default:
		return nil, nil, errors.New(UnknownProvider)
	}
}

func DecryptString(msg string, key []byte, iv []byte, prov CryptoProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn("Decrypting of empty string")
	}
	bytes := []byte(msg)
	return DecryptBytes(bytes, key, iv, prov)
}

func DecryptBytes(msg []byte, key []byte, iv []byte, prov CryptoProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn("Decrypting of empty array")
	}
	switch prov {
	case AESCFB:
		return decryptCFB(msg, key, iv)
	default:
		return nil, errors.New(UnknownProvider)
	}
}

func hashMD5(msg []byte) []byte {
	hash := md5.Sum(msg)
	return hash[:]
}

func hashSHA256(msg []byte) []byte {
	hash := sha256.Sum256(msg)
	return hash[:]
}

//In the previous version of this function (api v 1.0) this func worked in another way.
//First, hash has been calculated from input data
//Second, obtained hash has been converted to hex
//Third, hex value has been hashed once more time
//In this variant second step is omited.
func hashDoubleSHA256(msg []byte) []byte {
	firstHash := sha256.Sum256(msg)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:]
}

func encryptCFB(plainText, key, iv []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	if len(iv) == 0 {
		ciphertext := []byte(randSeq(16))
		iv = ciphertext[:16]
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(plainText))
	encrypter.XORKeyStream(encrypted, plainText)

	return append(iv, encrypted...), iv, nil
}

//In previous version of this func the parameters order was another (iv, encrypted, key)
func decryptCFB(cipherText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(cipherText))
	decrypter.XORKeyStream(decrypted, cipherText)
	return decrypted, nil
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
