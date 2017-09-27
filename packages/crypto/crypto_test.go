package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"
)

var hexCBCKey = "48F87EF4E749C70395AC5C3DC75BD628"
var hexCBCIV = "DA39A3EE5E6B4B0D3255BFEF95601890"
var hexMessage = "74657374206D65737361676500000000"

var hexCBCEncryptedWithPadding = "443e37da6b10527357f42bb3a4498af7f9f7c7618d2860e8a1586a358b0148df"
var ECDSASignedMessage = "5a926cacca5286c28663b4d1b22fa50a6f217881f7cdd59a997fe22c7bb9f7eb4cf7988d9cf7768dcfba75dce31e236a3a1421b8eb0d7d7ea7e89ef452ba3624"

type CBCcreds struct {
	Key     []byte
	IV      []byte
	Message []byte
}

func NewCBCcreds(key string, iv string, message string) (*CBCcreds, error) {
	result := &CBCcreds{}
	var err error
	result.Key, err = hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid key")
	}
	result.IV, err = hex.DecodeString(iv)
	if err != nil {
		return nil, fmt.Errorf("invalid iv")
	}

	result.Message, err = hex.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("invalid message")
	}
	return result, nil
}

func TestPaddingEqual(t *testing.T) {
	message, _ := hex.DecodeString(hexMessage)
	padded := _PKCS7Padding(message, consts.BlockSize)
	shouldBe, _ := hex.DecodeString("74657374206d6573736167650000000010101010101010101010101010101010")

	if bytes.Compare(padded, shouldBe) != 0 {
		t.Error("invalid padding want:\n", shouldBe, " have:\n", padded)
	}
}

func TestPaddingShort(t *testing.T) {
	message, _ := hex.DecodeString(hexMessage[4:])
	padded := _PKCS7Padding(message, consts.BlockSize)
	shouldBe, _ := hex.DecodeString("7374206d657373616765000000000202")

	if bytes.Compare(padded, shouldBe) != 0 {
		t.Error("invalid padding want:\n", shouldBe, " have:\n", padded)
	}
}

func TestPaddingGreater(t *testing.T) {
	message, _ := hex.DecodeString(hexMessage + hexMessage[:4])
	padded := _PKCS7Padding(message, consts.BlockSize)
	shouldBe, _ := hex.DecodeString("74657374206D6573736167650000000074650e0e0e0e0e0e0e0e0e0e0e0e0e0e")

	if bytes.Compare(padded, shouldBe) != 0 {
		t.Error("invalid padding want:\n", shouldBe, " have:\n", padded)
	}
}

func TestFakeUnpadding(t *testing.T) {
	message, _ := hex.DecodeString("74657374206d6573736167650000000011101010101010101010101010101010")
	_, err := _PKCS7UnPadding(message)
	if err == nil {
		t.Error("shouldn't be unpadded ", err)
	}
}

func TestFakeUnpadding1(t *testing.T) {
	message, _ := hex.DecodeString("74657374206d6573736167650000000010101010101010101010101010101011")
	_, err := _PKCS7UnPadding(message)
	if err == nil {
		t.Error("shouldn't be unpadded ", err)
	}
}

func TestUnpaddingEqual(t *testing.T) {
	message, _ := hex.DecodeString("74657374206d6573736167650000000010101010101010101010101010101010")
	unpadded, err := _PKCS7UnPadding(message)
	if err != nil {
		t.Error("unpadding error ", err)
	}
	shouldBe, _ := hex.DecodeString(hexMessage)

	if bytes.Compare(unpadded, shouldBe) != 0 {
		t.Error("invalid padding want:\n", shouldBe, " have:\n", unpadded)
	}
}

func TestUnpaddingShort(t *testing.T) {
	message, _ := hex.DecodeString("7374206d657373616765000000000202")
	unpadded, err := _PKCS7UnPadding(message)
	if err != nil {
		t.Error("unpadding error ", err)
	}
	shouldBe, _ := hex.DecodeString(hexMessage[4:])

	if bytes.Compare(unpadded, shouldBe) != 0 {
		t.Error("invalid padding want:\n", shouldBe, " have:\n", unpadded)
	}
}

func TestUnpaddingGreater(t *testing.T) {
	message, _ := hex.DecodeString("74657374206D6573736167650000000074650e0e0e0e0e0e0e0e0e0e0e0e0e0e")
	unpadded, err := _PKCS7UnPadding(message)
	if err != nil {
		t.Error("unpadding error ", err)
	}
	shouldBe, _ := hex.DecodeString(hexMessage + hexMessage[:4])

	if bytes.Compare(unpadded, shouldBe) != 0 {
		t.Error("invalid padding want:\n", shouldBe, " have:\n", unpadded)
	}
}

func TestCBCEnc(t *testing.T) {
	creds, err := NewCBCcreds(hexCBCKey, hexCBCIV, hexMessage)
	if err != nil {
		t.Error("creds error ", err)
	}

	shouldBe, _ := hex.DecodeString(hexCBCEncryptedWithPadding)
	shouldBe = append(creds.IV, shouldBe...)

	result, err := Encrypt(creds.Message, creds.Key, creds.IV)

	if len(result) != len(creds.IV)+len(creds.Message)*2 {
		t.Error("invalid encrypted len")
	}

	if bytes.Compare(result, shouldBe) != 0 {
		t.Error("invalid encrypting. want: \n", shouldBe, " have: \n", result)
	}
}

func TestCBCEncShortIV(t *testing.T) {
	creds, err := NewCBCcreds(hexCBCKey, hexCBCIV[2:], hexMessage)
	if err != nil {
		t.Error("creds error ", err)
	}

	_, err = Encrypt(creds.Message, creds.Key, creds.IV)

	if err == nil {
		t.Error("shoul be err")
	}
}

func TestCBCEncBigIV(t *testing.T) {
	creds, err := NewCBCcreds(hexCBCKey, hexCBCIV, hexMessage)
	if err != nil {
		t.Error("creds error ", err)
	}

	shouldBe, _ := hex.DecodeString(hexCBCEncryptedWithPadding)
	shouldBe = append(creds.IV, shouldBe...)

	result, err := Encrypt(creds.Message, creds.Key, append(creds.IV, creds.IV...))

	if len(result) != len(creds.IV)+len(creds.Message)*2 {
		t.Error("invalid encrypted len")
	}

	if bytes.Compare(result, shouldBe) != 0 {
		t.Error("invalid encrypting. want: \n", shouldBe, " have: \n", result)
	}
}

func TestCBCDec(t *testing.T) {
	creds, err := NewCBCcreds(hexCBCKey, hexCBCIV, hexCBCEncryptedWithPadding)
	if err != nil {
		t.Error("creds error ", err)
	}

	shouldBe, _ := hex.DecodeString(hexMessage)

	result, err := Decrypt(creds.Message, creds.Key, creds.IV)

	if len(result) != len(shouldBe) {
		t.Error("invalid decrypted len")
	}

	if bytes.Compare(result, shouldBe) != 0 {
		t.Error("invalid decrypting. want: \n", shouldBe, " have: \n", result)
	}
}

func TestCBCEncDec(t *testing.T) {
	creds, err := NewCBCcreds(hexCBCKey, hexCBCIV, hexMessage)
	if err != nil {
		t.Error("creds error ", err)
	}

	result, err := Encrypt(creds.Message, creds.Key, creds.IV)
	if err != nil {
		t.Error("encrypting error ", err)
	}

	result, err = Decrypt(result[consts.BlockSize:], creds.Key, creds.IV)
	if err != nil {
		t.Error("encrypting error ", err)
	}

	if bytes.Compare(result, creds.Message) != 0 {
		t.Error("encoding error: want: \n", result, " have: \n", creds.Message)
	}
}

func TestECDSASign(t *testing.T) {
	res, err := Sign(hexCBCKey, hexMessage)
	if err != nil {
		t.Error("sign error")
	}
	byteKey, _ := hex.DecodeString(hexCBCKey)
	pub, _ := PrivateToPublic(byteKey)
	verified, err := checkECDSA(pub, hexMessage, res)
	if err != nil {
		t.Error("sign check err ", err)
	}
	if !verified {
		t.Error("can't verify sign")
	}
}
