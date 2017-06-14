package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	log "github.com/sirupsen/logrus"
)

type CryptoProvider int
type EllipticSize int

const (
	AESCBC CryptoProvider = iota
)

const (
	Elliptic256 EllipticSize = iota
)

const (
	HashingError         = "Hashing error"
	EncryptingError      = "Encoding error"
	DecryptingError      = "Decrypting error"
	UnknownProviderError = "Unknown provider"
	HashingEmpty         = "Hashing empty value"
	EncryptingEmpty      = "Encrypting empty value"
	DecryptingEmpty      = "Decrypting empty value"
	SigningEmpty         = "Signing empty value"
	CheckingSignEmpty    = "Cheking sign of empty"
	IncorrectSign        = "Incorrect sign"
	UnsupportedCurveSize = "Unsupported curve size"
)

func Encrypt(msg []byte, key []byte, iv []byte, prov CryptoProvider) ([]byte, []byte, error) {
	if len(msg) == 0 {
		log.Warn(EncryptingEmpty)
	}
	switch prov {
	case AESCBC:
		res, err := encryptCBC(msg, key, iv)
		if err != nil {
			return nil, nil, err
		}
		return res, nil, nil
	default:
		return nil, nil, errors.New(UnknownProviderError)
	}
}

func Decrypt(msg []byte, key []byte, iv []byte, prov CryptoProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn(DecryptingEmpty)
	}
	switch prov {
	case AESCBC:
		return decryptCBC(msg, key, iv)
	default:
		return nil, errors.New(UnknownProviderError)
	}
}

// SharedEncrypt creates a shared key and encrypts text. The first 32 characters are the created public key.
// The cipher text can be only decrypted with the original private key.
func SharedEncrypt(public, text []byte, hashProv HashProvider, signProv SignProvider, cryptoProv CryptoProvider, ellipticSize EllipticSize) ([]byte, error) {
	priv, pub, err := GenBytesKeys(ellipticSize)
	if err != nil {
		return nil, err
	}
	shared, err := getSharedKey(priv, public, hashProv, signProv, ellipticSize)
	if err != nil {
		return nil, err
	}
	val, _, err := Encrypt(shared, text, pub, cryptoProv)
	return val, err
}

// SharedDecrypt decrypts the ciphertext by using private key.
func SharedDecrypt(private, ciphertext []byte, hashProv HashProvider, signProv SignProvider, cryptoProv CryptoProvider, ellipticSize EllipticSize) ([]byte, error) {
	if len(ciphertext) <= 64 {
		return nil, fmt.Errorf(`too short cipher %d`, len(ciphertext))
	}
	shared, err := getSharedKey(private, ciphertext[:64], hashProv, signProv, ellipticSize)
	if err != nil {
		return nil, err
	}
	val, _, err := Encrypt(shared, ciphertext[64:], ciphertext[:aes.BlockSize], cryptoProv)
	return val, err
}

// GenBytesKeys generates a random pair of ECDSA private and public binary keys.
// TODO параметризировать fillLeft
func GenBytesKeys(size EllipticSize) ([]byte, []byte, error) {
	var curve elliptic.Curve
	switch size {
	case Elliptic256:
		curve = elliptic.P256()
	default:
		return nil, nil, errors.New(UnsupportedCurveSize)
	}
	private, err := ecdsa.GenerateKey(curve, crand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return private.D.Bytes(), append(converter.FillLeft(private.PublicKey.X.Bytes()), converter.FillLeft(private.PublicKey.Y.Bytes())...), nil
}

// GenHexKeys generates a random pair of ECDSA private and public hex keys.
func GenHexKeys(size EllipticSize) (string, string, error) {
	priv, pub, err := GenBytesKeys(size)
	if err != nil {
		return ``, ``, err
	}
	return hex.EncodeToString(priv), hex.EncodeToString(pub), nil
}

// CBCEncrypt encrypts the text by using the key parameter. It uses CBC mode of AES.
func encryptCBC(text, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := _PKCS7Padding(text, aes.BlockSize)
	if iv == nil {
		iv = make([]byte, aes.BlockSize, aes.BlockSize+len(plaintext))
		if _, err := io.ReadFull(crand.Reader, iv); err != nil {
			return nil, err
		}
	}
	if len(iv) < aes.BlockSize {
		return nil, fmt.Errorf(`wrong size of iv %d`, len(iv))
	}
	mode := cipher.NewCBCEncrypter(block, iv[:aes.BlockSize])
	encrypted := make([]byte, len(plaintext))
	mode.CryptBlocks(encrypted, plaintext)
	return append(iv, encrypted...), nil
}

// CBCDecrypt decrypts the text by using key. It uses CBC mode of AES.
func decryptCBC(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize || len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf(`Wrong size of cipher %d`, len(ciphertext))
	}
	if iv == nil {
		iv = ciphertext[:aes.BlockSize]
		ciphertext = ciphertext[aes.BlockSize:]
	}
	ret := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv[:aes.BlockSize]).CryptBlocks(ret, ciphertext)
	if ret, err = _PKCS7UnPadding(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// TODO в приватные
// PKCS7Padding realizes PKCS#7 encoding which is described in RFC 5652.
func _PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

//TODO в приватные
// PKCS7UnPadding realizes PKCS#7 decoding.
func _PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length < int(src[length-1]) {
		return nil, fmt.Errorf(`incorrect input of PKCS7UnPadding`)
	}
	return src[:length-int(src[length-1])], nil
}

// GetSharedKey creates and returns the shared key = private * public.
// public must be the public key from the different private key.
func getSharedKey(private, public []byte, hashProv HashProvider, signProv SignProvider, ellipticSize EllipticSize) (shared []byte, err error) {
	var pubkeyCurve elliptic.Curve
	switch ellipticSize {
	case Elliptic256:
		pubkeyCurve = elliptic.P256()
	default:
		return nil, errors.New(UnknownProviderError)
	}

	switch signProv {
	case ECDSA:
		private = converter.FillLeft(private)
		public = converter.FillLeft(public)
		pub := new(ecdsa.PublicKey)
		pub.Curve = pubkeyCurve
		pub.X = new(big.Int).SetBytes(public[0:32])
		pub.Y = new(big.Int).SetBytes(public[32:])

		bi := new(big.Int).SetBytes(private)
		priv := new(ecdsa.PrivateKey)
		priv.PublicKey.Curve = pubkeyCurve
		priv.D = bi
		priv.PublicKey.X, priv.PublicKey.Y = pubkeyCurve.ScalarBaseMult(bi.Bytes())

		if priv.Curve.IsOnCurve(pub.X, pub.Y) {
			x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, priv.D.Bytes())
			key, err := Hash([]byte(hex.EncodeToString(x.Bytes())), hashProv)
			if err != nil {
				return nil, errors.New(UnknownProviderError)
			}
			shared = key
		} else {
			err = fmt.Errorf("Not IsOnCurve")
		}
	default:
		return nil, errors.New(UnknownProviderError)
	}

	return
}
