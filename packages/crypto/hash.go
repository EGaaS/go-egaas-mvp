package crypto

import (
	"crypto/sha256"
	"errors"

	log "github.com/sirupsen/logrus"
)

type HashProvider int

const (
	SHA256 HashProvider = iota
	DoubleSHA256
)

func Hash(msg []byte, prov HashProvider) ([]byte, error) {
	if len(msg) == 0 {
		log.Warn(HashingEmpty)
	}
	switch prov {
	case SHA256:
		return hashSHA256(msg), nil
	case DoubleSHA256:
		return hashDoubleSHA256(msg), nil
	default:
		return nil, errors.New(UnknownProviderError)
	}
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
