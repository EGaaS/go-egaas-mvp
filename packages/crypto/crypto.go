package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	log "github.com/sirupsen/logrus"
)

type CryptoProvider int
type HashProvider int
type SignProvider int
type EllipticSize int
type ChecksumProvider int

const (
	MD5 HashProvider = iota
	SHA256
	DoubleSHA256
)

const (
	AESCFB CryptoProvider = iota
	AESCBC
)

const (
	ECDSA SignProvider = iota
)

const (
	Elliptic256 EllipticSize = iota
)

const (
	CRC64 ChecksumProvider = iota
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

var (
	table64 *crc64.Table
)

func init() {
	table64 = crc64.MakeTable(crc64.ECMA)
}

func HashString(msg string, prov HashProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn(HashingEmpty)
	}
	bytes := []byte(msg)
	return HashBytes(bytes, prov)
}

func HashBytes(msg []byte, prov HashProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn(HashingEmpty)
	}
	switch prov {
	case MD5:
		return hashMD5(msg), nil
	case SHA256:
		return hashSHA256(msg), nil
	case DoubleSHA256:
		return hashDoubleSHA256(msg), nil
	default:
		return nil, errors.New(UnknownProviderError)
	}
}

func EncryptString(msg string, key []byte, iv []byte, prov CryptoProvider) ([]byte, []byte, error) {
	if len(msg) == 0 {
		log.Warn(EncryptingEmpty)
	}
	bytes := []byte(msg)
	return EncryptBytes(bytes, key, iv, prov)
}

func EncryptBytes(msg []byte, key []byte, iv []byte, prov CryptoProvider) ([]byte, []byte, error) {
	if len(msg) == 0 {
		log.Warn(EncryptingEmpty)
	}
	switch prov {
	case AESCFB:
		return encryptCFB(msg, key, iv)
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

func DecryptString(msg string, key []byte, iv []byte, prov CryptoProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn(DecryptingEmpty)
	}
	bytes := []byte(msg)
	return DecryptBytes(bytes, key, iv, prov)
}

func DecryptBytes(msg []byte, key []byte, iv []byte, prov CryptoProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn(DecryptingEmpty)
	}
	switch prov {
	case AESCFB:
		return decryptCFB(msg, key, iv)
	case AESCBC:
		return decryptCBC(msg, key, iv)
	default:
		return nil, errors.New(UnknownProviderError)
	}
}

func Sign(privateKey string, data string, hashProv HashProvider, signProv SignProvider, size EllipticSize) ([]byte, error) {
	if len(data) == 0 {
		log.Warn(SigningEmpty)
	}
	switch signProv {
	case ECDSA:
		return signECDSA(privateKey, data, hashProv, size)
	default:
		return nil, errors.New(UnknownProviderError)
	}
}

func CheckSign(public []byte, data string, signature []byte, hashProv HashProvider, signProv SignProvider, size EllipticSize) (bool, error) {
	if len(public) == 0 {
		log.Warn(CheckingSignEmpty)
	}
	switch signProv {
	case ECDSA:
		return checkECDSA(public, data, signature, hashProv, size)
	default:
		return false, errors.New(UnknownProviderError)
	}
}

func signECDSA(privateKey string, data string, hashProv HashProvider, size EllipticSize) (ret []byte, err error) {
	var pubkeyCurve elliptic.Curve

	switch size {
	case Elliptic256:
		pubkeyCurve = elliptic.P256()
	default:
		log.Fatal(UnsupportedCurveSize)
	}

	b, err := hex.DecodeString(privateKey)
	if err != nil {
		return
	}
	bi := new(big.Int).SetBytes(b)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = pubkeyCurve
	priv.D = bi
	priv.PublicKey.X, priv.PublicKey.Y = pubkeyCurve.ScalarBaseMult(bi.Bytes())

	signhash, err := HashBytes([]byte(data), hashProv)
	if err != nil {
		log.Fatal(HashingError)
	}
	r, s, err := ecdsa.Sign(crand.Reader, priv, signhash)
	if err != nil {
		return
	}
	ret = append(converter.FillLeft(r.Bytes()), converter.FillLeft(s.Bytes())...)
	return
}

// TODO параметризировать, длина данных в зависимости от длины кривой
// CheckECDSA checks if forSign has been signed with corresponding to public the private key
func checkECDSA(public []byte, data string, signature []byte, hashProv HashProvider, size EllipticSize) (bool, error) {
	if len(data) == 0 || len(public) != 64 || len(signature) == 0 {
		return false, fmt.Errorf("invalid parameters")
	}
	var pubkeyCurve elliptic.Curve
	switch size {
	case Elliptic256:
		pubkeyCurve = elliptic.P256()
	default:
		log.Fatal(UnsupportedCurveSize)
	}

	hash, err := HashBytes([]byte(data), hashProv)
	if err != nil {
		log.Fatal(HashingError)
	}

	pubkey := new(ecdsa.PublicKey)
	pubkey.Curve = pubkeyCurve
	pubkey.X = new(big.Int).SetBytes(public[0:32])
	pubkey.Y = new(big.Int).SetBytes(public[32:])
	r, s, err := parseSign(hex.EncodeToString(signature))
	if err != nil {
		return false, err
	}
	verifystatus := ecdsa.Verify(pubkey, hash, r, s)
	if !verifystatus {
		return false, errors.New(IncorrectSign)
	}
	return true, nil
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

// TODO уточнить зачем два разных хеша
// Address gets int64 EGGAS address from the public key
func Address(pubKey []byte) int64 {
	h256 := sha256.Sum256(pubKey)
	h512 := sha512.Sum512(h256[:])
	crc := calcCRC64(h512[:])
	// replace the last digit by checksum
	num := strconv.FormatUint(crc, 10)
	val := []byte(strings.Repeat("0", 20-len(num)) + num)
	return int64(crc - (crc % 10) + uint64(checkSum(val[:len(val)-1])))
}

// PrivateToPublic returns the public key for the specified private key.
func PrivateToPublic(key []byte, size EllipticSize) ([]byte, error) {
	var pubkeyCurve elliptic.Curve
	switch size {
	case Elliptic256:
		pubkeyCurve = elliptic.P256()
	default:
		return nil, errors.New(UnsupportedCurveSize)
	}

	bi := new(big.Int).SetBytes(key)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = pubkeyCurve
	priv.D = bi
	priv.PublicKey.X, priv.PublicKey.Y = pubkeyCurve.ScalarBaseMult(bi.Bytes())
	return append(converter.FillLeft(priv.PublicKey.X.Bytes()), converter.FillLeft(priv.PublicKey.Y.Bytes())...), nil
}

// TODO убрать вместе с хексом
// PrivateToPublicHex returns the hex public key for the specified hex private key.
func PrivateToPublicHex(hexkey string) (string, error) {
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		return ``, errors.New("Decode hex error")
	}
	pubKey, err := PrivateToPublic(key, Elliptic256)
	if err != nil {
		return ``, err
	}
	return hex.EncodeToString(pubKey), nil
}

// TODO убрать отсюда
// KeyToAddress converts a public key to EGAAS address XXXX-...-XXXX.
func KeyToAddress(pubKey []byte) string {
	return converter.AddressToString(Address(pubKey))
}

func CalcChecksum(input []byte, checksumProv ChecksumProvider) (uint64, error) {
	switch checksumProv {
	case CRC64:
		return calcCRC64(input), nil
	default:
		return 0, errors.New(UnknownProviderError)
	}
}

// CRC64 returns crc64 sum
func calcCRC64(input []byte) uint64 {
	return crc64.Checksum(input, table64)
}

// CheckSum calculates the 0-9 check sum of []byte
func checkSum(val []byte) int {
	var one, two int
	for i, ch := range val {
		digit := int(ch - '0')
		if i&1 == 1 {
			one += digit
		} else {
			two += digit
		}
	}
	checksum := (two + 3*one) % 10
	if checksum > 0 {
		checksum = 10 - checksum
	}
	return checksum
}

// parseSign converts the hex signature to r and s big number
func parseSign(sign string) (*big.Int, *big.Int, error) {
	var (
		binSign []byte
		err     error
	)
	//	var off int
	parse := func(bsign []byte) []byte {
		blen := int(bsign[1])
		if blen > len(bsign)-2 {
			return nil
		}
		ret := bsign[2 : 2+blen]
		if len(ret) > 32 {
			ret = ret[len(ret)-32:]
		} else if len(ret) < 32 {
			ret = append(bytes.Repeat([]byte{0}, 32-len(ret)), ret...)
		}
		return ret
	}
	if len(sign) > 128 {
		binSign, err = hex.DecodeString(sign)
		if err != nil {
			return nil, nil, err
		}
		left := parse(binSign[2:])
		if left == nil || int(binSign[3])+6 > len(binSign) {
			return nil, nil, fmt.Errorf(`wrong left parsing`)
		}
		right := parse(binSign[4+binSign[3]:])
		if right == nil {
			return nil, nil, fmt.Errorf(`wrong right parsing`)
		}
		sign = hex.EncodeToString(append(left, right...))
	} else if len(sign) < 128 {
		return nil, nil, fmt.Errorf(`wrong len of signature %d`, len(sign))
	}
	all, err := hex.DecodeString(sign[:])
	if err != nil {
		return nil, nil, err
	}
	return new(big.Int).SetBytes(all[:32]), new(big.Int).SetBytes(all[len(all)-32:]), nil
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
		ciphertext := []byte(RandSeq(16))
		iv = ciphertext[:16]
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(plainText))
	encrypter.XORKeyStream(encrypted, plainText)

	return append(iv, encrypted...), iv, nil
}

// CBCEncrypt encrypts the text by using the key parameter. It uses CBC mode of AES.
func encryptCBC(text, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := PKCS7Padding(text, aes.BlockSize)
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
	if ret, err = PKCS7UnPadding(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// TODO в приватные
// PKCS7Padding realizes PKCS#7 encoding which is described in RFC 5652.
func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

//TODO в приватные
// PKCS7UnPadding realizes PKCS#7 decoding.
func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length < int(src[length-1]) {
		return nil, fmt.Errorf(`incorrect input of PKCS7UnPadding`)
	}
	return src[:length-int(src[length-1])], nil
}

// JSSignToBytes converts hex signature which has got from the browser to []byte
func JSSignToBytes(in string) ([]byte, error) {
	r, s, err := parseSign(in)
	if err != nil {
		return nil, err
	}
	return append(converter.FillLeft(r.Bytes()), converter.FillLeft(s.Bytes())...), nil
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
	val, _, err := EncryptBytes(shared, text, pub, cryptoProv)
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
	val, _, err := EncryptBytes(shared, ciphertext[64:], ciphertext[:aes.BlockSize], cryptoProv)
	return val, err
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
			key, err := HashBytes([]byte(hex.EncodeToString(x.Bytes())), hashProv)
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

// RandInt returns a random integer between min and max
func RandInt(min int, max int) int {
	if max-min <= 0 {
		return 1
	}
	return min + rand.Intn(max-min)
}
