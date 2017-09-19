package publisher

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/EGaaS/go-egaas-mvp/packages/crypto"
	"github.com/centrifugal/gocent"
)

var (
	clientsChannels = map[int64]string{}
	//admin centrifugo key
	centrifugoSecret  = "ALKSJDOI:Q@U(!"
	centrifugoURL     = "http://some_url/api"
	centrifugoTimeout = time.Second * 5
	publisher         *gocent.Client
)

func init() {
	publisher = gocent.NewClient(centrifugoURL, centrifugoSecret, centrifugoTimeout)
}

func GetHMACSign(userID int64) (string, error) {
	secret, err := crypto.HMAC(centrifugoSecret, strconv.FormatInt(userID, 64))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	result := hex.EncodeToString(secret)
	clientsChannels[userID] = result
	return result, nil
}

func Write(userID int64, data string) (bool, error) {
	ok, err := publisher.Publish("client&"+clientsChannels[userID], []byte(data))
	if err != nil {
		fmt.Println(err)
	}
	return ok, err
}
