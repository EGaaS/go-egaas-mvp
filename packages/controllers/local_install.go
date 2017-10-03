package controllers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"
	"github.com/EGaaS/go-egaas-mvp/packages/lib"
	"github.com/EGaaS/go-egaas-mvp/packages/static"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
	"github.com/astaxie/beego/config"
)

func (c *Controller) LocalInstall(directory string,
	genFirstBlock int64,
	tcpHost string,
	httpPort string,
	ll string,
	firstLoadURL string,
	firstLoad string,
	dbType string,
	dbHost string,
	dbPort string,
	dbName string,
	dbUsername string,
	dbPassword string) error {
	if _, err := os.Stat(*utils.Dir + "/config.ini"); err == nil {
		return errors.New("config.ini exists")
	}

	dir := directory
	*utils.Dir = dir
	*utils.GenerateFirstBlock = genFirstBlock

	installType := "standart"
	*utils.TcpHost = tcpHost
	*utils.ListenHttpPort = httpPort
	logLevel := ll

	firstLoadBlockchainURL := firstLoadURL

	if len(firstLoadBlockchainURL) == 0 {
		firstLoadBlockchainURL = consts.BLOCKCHAIN_URL
	}

	if _, err := os.Stat(*utils.Dir + "/config.ini"); os.IsNotExist(err) {
		ioutil.WriteFile(*utils.Dir+"/config.ini", []byte(``), 0644)
	}
	confIni, err := config.NewConfig("ini", *utils.Dir+"/config.ini")
	confIni.Set("log_level", logLevel)
	confIni.Set("install_type", installType)
	confIni.Set("dir", *utils.Dir)
	confIni.Set("tcp_host", *utils.TcpHost)
	confIni.Set("http_port", *utils.ListenHttpPort)
	confIni.Set("first_block_dir", *utils.FirstBlockDir)
	confIni.Set("db_type", dbType)
	confIni.Set("db_user", dbUsername)
	confIni.Set("db_host", dbHost)
	confIni.Set("db_port", dbPort)
	confIni.Set("db_password", dbPassword)
	confIni.Set("db_name", dbName)

	err = confIni.SaveConfigFile(*utils.Dir + "/config.ini")
	if err != nil {
		dropConfig()
		return err
	}

	configIni, err = confIni.GetSection("default")
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}
	utils.DB, err = utils.NewDbConnect(configIni)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}
	c.DCDB = utils.DB
	if c.DCDB.DB == nil {
		err = fmt.Errorf("utils.DB == nil")
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	err = c.DCDB.ExecSql(`DROP SCHEMA public CASCADE`)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	err = c.DCDB.ExecSql(`CREATE SCHEMA public`)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	schema, err := static.Asset("static/schema.sql")
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	err = c.DCDB.ExecSql(string(schema))
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	err = c.DCDB.ExecSql("INSERT INTO config (first_load_blockchain, first_load_blockchain_url, auto_reload) VALUES (?, ?, ?)", firstLoad, firstLoadBlockchainURL, 259200)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	err = c.DCDB.ExecSql(`INSERT INTO install (progress) VALUES ('complete')`)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	log.Debug("GenerateFirstBlock", *utils.GenerateFirstBlock)

	if _, err := os.Stat(*utils.FirstBlockDir + "/1block"); len(*utils.FirstBlockDir) > 0 && os.IsNotExist(err) {

		fmt.Println("not exists " + *utils.FirstBlockDir + "/1block")

		// If there is no key, this is the first run and the need to create them in the working directory.
		if _, err := os.Stat(*utils.Dir + "/PrivateKey"); os.IsNotExist(err) {
			fmt.Println("not exists " + *utils.Dir + "/PrivateKey")
			if len(*utils.FirstBlockPublicKey) == 0 {
				priv, pub := lib.GenKeys()
				fmt.Println("WriteFile " + *utils.Dir + "/PrivateKey")
				err := ioutil.WriteFile(*utils.Dir+"/PrivateKey", []byte(priv), 0644)
				if err != nil {
					log.Error("%v", utils.ErrInfo(err))
				}
				*utils.FirstBlockPublicKey = pub
			}
		}

		if _, err := os.Stat(*utils.Dir + "/NodePrivateKey"); os.IsNotExist(err) {
			fmt.Println("not exists " + *utils.Dir + "/NodePrivateKey")
			if len(*utils.FirstBlockNodePublicKey) == 0 {
				priv, pub := lib.GenKeys()
				fmt.Println("WriteFile " + *utils.Dir + "/NodePrivateKey")
				err := ioutil.WriteFile(*utils.Dir+"/NodePrivateKey", []byte(priv), 0644)
				if err != nil {
					log.Error("%v", utils.ErrInfo(err))
				}
				*utils.FirstBlockNodePublicKey = pub
			}
		}

		*utils.GenerateFirstBlock = 1
		utils.FirstBlock(false)
	}
	log.Debug("1block")

	NodePrivateKey, _ := ioutil.ReadFile(*utils.Dir + "/NodePrivateKey")
	NodePrivateKeyStr := strings.TrimSpace(string(NodePrivateKey))
	npubkey := lib.PrivateToPublicHex(NodePrivateKeyStr)
	err = c.DCDB.ExecSql(`INSERT INTO my_node_keys (private_key, public_key, block_id) VALUES (?, [hex], ?)`, NodePrivateKeyStr, npubkey, 1)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}

	if *utils.DltWalletId == 0 {
		PrivateKey, _ := ioutil.ReadFile(*utils.Dir + "/PrivateKey")
		PrivateHex, _ := hex.DecodeString(string(PrivateKey))
		PublicKeyBytes2 := lib.PrivateToPublic(PrivateHex)
		log.Debug("dlt_wallet_id %d", int64(lib.Address(PublicKeyBytes2)))
		*utils.DltWalletId = int64(lib.Address(PublicKeyBytes2))
	}

	err = c.DCDB.ExecSql(`UPDATE config SET dlt_wallet_id = ?`, *utils.DltWalletId)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		dropConfig()
		return err
	}
	return nil
}
