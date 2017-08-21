// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package daylight

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/EGaaS/go-egaas-mvp/packages/api"
	"github.com/EGaaS/go-egaas-mvp/packages/config"
	"github.com/EGaaS/go-egaas-mvp/packages/config/syspar"
	"github.com/EGaaS/go-egaas-mvp/packages/consts"
	"github.com/EGaaS/go-egaas-mvp/packages/controllers"
	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	"github.com/EGaaS/go-egaas-mvp/packages/daemons"
	"github.com/EGaaS/go-egaas-mvp/packages/exchangeapi"
	"github.com/EGaaS/go-egaas-mvp/packages/language"
	"github.com/EGaaS/go-egaas-mvp/packages/model"
	"github.com/EGaaS/go-egaas-mvp/packages/parser"
	"github.com/EGaaS/go-egaas-mvp/packages/schema"
	"github.com/EGaaS/go-egaas-mvp/packages/static"
	"github.com/EGaaS/go-egaas-mvp/packages/stopdaemons"
	"github.com/EGaaS/go-egaas-mvp/packages/system"
	"github.com/EGaaS/go-egaas-mvp/packages/template"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
	"github.com/go-bindata-assetfs"
	"github.com/go-thrust/lib/bindings/window"
	"github.com/go-thrust/lib/commands"
	"github.com/go-thrust/thrust"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// FileAsset returns the body of the file
func FileAsset(name string) ([]byte, error) {
	if name := strings.Replace(name, "\\", "/", -1); name == `static/img/logo.`+utils.LogoExt {
		logofile := *utils.Dir + `/logo.` + utils.LogoExt
		if fi, err := os.Stat(logofile); err == nil && fi.Size() > 0 {
			return ioutil.ReadFile(logofile)
		}
	}
	return static.Asset(name)
}

func readConfig() {
	// read the config.ini
	if err := config.Read(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Info("Error reading config")
	}
	if *utils.TCPHost == "" {
		*utils.TCPHost = config.ConfigIni["tcp_host"]
	}
	if *utils.FirstBlockDir == "" {
		*utils.FirstBlockDir = config.ConfigIni["first_block_dir"]
	}
	if *utils.ListenHTTPPort == "" {
		*utils.ListenHTTPPort = config.ConfigIni["http_port"]
	}
	if *utils.Dir == "" {
		*utils.Dir = config.ConfigIni["dir"]
	}
	utils.OneCountry = converter.StrToInt64(config.ConfigIni["one_country"])
	utils.PrivCountry = config.ConfigIni["priv_country"] == `1` || config.ConfigIni["priv_country"] == `true`
	if len(config.ConfigIni["lang"]) > 0 {
		language.LangList = strings.Split(config.ConfigIni["lang"], `,`)
	}
}

func killOld() {
	path := *utils.Dir + "/daylight.pid"
	if _, err := os.Stat(path); err == nil {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			log.WithFields(logrus.Fields{"path": path, "error": err}).Warn("reading pid file")
		}
		var pidMap map[string]string
		err = json.Unmarshal(dat, &pidMap)
		if err != nil {
			log.WithFields(logrus.Fields{"path": path, "error": err}).Warn("Error unmarshalling pid file to json")
		}
		log.WithFields(logrus.Fields{"pid": pidMap["pid"]}).Info("Old process pid to kill")
		err = KillPid(pidMap["pid"])
		if nil != err {
			log.WithFields(logrus.Fields{"pid": pidMap["pid"], "error": err}).Warn("killing process with pid")
		}
		if fmt.Sprintf("%s", err) != "null" {
			// give 15 sec to end the previous process
			log.Info("Waiting for previous process to die")
			for i := 0; i < 15; i++ {
				if _, err := os.Stat(*utils.Dir + "/daylight.pid"); err == nil {
					time.Sleep(time.Second)
				} else { // if there is no daylight.pid, so it is finished
					break
				}
			}
		}
	}
}

func initLogs() error {
	if config.ConfigIni["log_output"] == "console" {
		f, err := os.OpenFile(*utils.Dir+"/dclog.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			log.WithFields(logrus.Fields{"path": *utils.Dir + "/dclog.txt", "error": err}).Warn("opening file for logging output")
			return err
		}
		logrus.SetOutput(f)
	}

	if *utils.LogLevel == "" {
		if level, ok := config.ConfigIni["log_level"]; ok {
			*utils.LogLevel = level
		} else {
			*utils.LogLevel = "INFO"
		}
	}

	switch *utils.LogLevel {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	}
	return nil
}

func savePid() error {
	pid := os.Getpid()
	PidAndVer, err := json.Marshal(map[string]string{"pid": converter.IntToStr(pid), "version": consts.VERSION})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Warn("Error marshalling process info to json")
		return err
	}
	if err := ioutil.WriteFile(*utils.Dir+"/daylight.pid", PidAndVer, 0644); err != nil {
		log.WithFields(logrus.Fields{"path": *utils.Dir + "/daylight.pid", "error": err}).Warn("Error saving pidfile")
		return err
	}
	return nil
}

func delPidFile() {
	path := filepath.Join(*utils.Dir, "daylight.pid")
	if err := os.Remove(path); err != nil {
		log.WithFields(logrus.Fields{"path": path, "error": err}).Warn("Can't remove pidfile")
	}
}

func rollbackToBlock(blockID int64) error {
	if err := template.LoadContracts(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error(`Load Contracts, while rollbacking to block`)
		return err
	}
	parser := new(parser.Parser)
	err := parser.RollbackToBlockID(*utils.RollbackToBlockID)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Rollback to block ID failed")
		return err
	}

	// we recieve the statistics of all tables
	allTable, err := model.GetAllTables()
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Get all tables error")
		return err

	}

	startData := map[string]int64{"install": 1, "config": 1, "queue_tx": 99999, "log_transactions": 1, "transactions_status": 99999, "block_chain": 1, "info_block": 1, "dlt_wallets": 1, "confirmations": 9999999, "full_nodes": 1, "system_parameters": 4, "my_node_keys": 99999, "transactions": 999999}
	for _, table := range allTable {
		query := `SELECT COUNT(*) FROM ` + converter.EscapeName(table)
		count, err := model.Single(`SELECT count(*) FROM ` + converter.EscapeName(table)).Int64()
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, query: query}).Error("Query failed")
			return err
		}
		if count > 0 && count > startData[table] {
			log.WithFields(logrus.Fields{"count": count, "normal_count": startData[table], "table": table}).Warn("count exceed normal count at start")
		}
	}
	return nil
}

func setRoute(route *httprouter.Router, path string, handle func(http.ResponseWriter, *http.Request), methods ...string) {
	for _, method := range methods {
		route.HandlerFunc(method, path, handle)
	}
}

func initRoutes(listenHost, browserHost string) string {
	route := httprouter.New()
	setRoute(route, `/`, controllers.Index, `GET`)
	setRoute(route, `/content`, controllers.Content, `GET`, `POST`)
	setRoute(route, `/template`, controllers.Template, `GET`, `POST`)
	setRoute(route, `/app`, controllers.App, `GET`, `POST`)
	setRoute(route, `/ajax`, controllers.Ajax, `GET`, `POST`)
	setRoute(route, `/wschain`, controllers.WsBlockchain, `GET`)
	setRoute(route, `/exchangeapi/:name`, exchangeapi.API, `GET`, `POST`)
	api.Route(route)
	route.Handler(`GET`, `/static/*filepath`, http.FileServer(&assetfs.AssetFS{Asset: FileAsset, AssetDir: static.AssetDir, Prefix: ""}))
	route.Handler(`GET`, `/.well-known/*filepath`, http.FileServer(http.Dir(*utils.TLS)))
	if len(*utils.TLS) > 0 {
		go http.ListenAndServeTLS(":443", *utils.TLS+`/fullchain.pem`, *utils.TLS+`/privkey.pem`, route)
	}

	httpListener(listenHost, browserHost, route)
	// for ipv6 server
	httpListenerV6(route)
	return browserHost
}

// Start starts the main code of the program
func Start(dir string, thrustWindowLoder *window.Window) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(logrus.Fields{"recover": r}).Panic("Panic with")
		}
	}()

	Exit := func(code int) {
		if thrustWindowLoder != nil {
			thrustWindowLoder.Close()
		}
		os.Exit(code)
	}

	if dir != "" {
		*utils.Dir = dir
	}

	readConfig()

	// create first block
	if *utils.GenerateFirstBlock == 1 {
		log.Info("Generating first block")
		utils.FirstBlock()
		os.Exit(0)

	}
	exchangeapi.InitAPI()
	log.Info("Exchange API Initialized")

	// kill previously run eGaaS
	if !utils.Mobile() {
		killOld()
	}

	controllers.SessInit()
	log.Info("Controllers Initialized")

	if fi, err := os.Stat(*utils.Dir + `/logo.png`); err == nil && fi.Size() > 0 {
		utils.LogoExt = `png`
	}

	err = initLogs()
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Logs init fail")
		Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// if there is OldFileName, so act on behalf dc.tmp and we have to restart on behalf the normal name
	if *utils.OldFileName != "" {
		if *utils.OldFileName != "" { //*utils.Dir+`/dc.tmp`
			err = utils.CopyFileContents(os.Args[0], *utils.OldFileName)
			if err != nil {
				log.WithFields(logrus.Fields{"from": *utils.OldFileName, "to": os.Args[0], "error": err}).Error("Error copying from old file no new")
			}
		}
		schema.Migration()

		if *utils.OldFileName != "" {
			err = model.DBConn.Close()
			if err != nil {
				log.WithFields(logrus.Fields{"error": err}).Error("error closing db")
			}
			path := filepath.Join(*utils.Dir, "daylight.pid")
			err = os.Remove(path)
			if err != nil {
				log.WithFields(logrus.Fields{"error": err, "path": path}).Error("error removing pidfile")
			}

			if thrustWindowLoder != nil {
				thrustWindowLoder.Close()
			}
			system.Finish()
			err = exec.Command(*utils.OldFileName, "-dir", *utils.Dir).Start()
			if err != nil {
				log.WithFields(logrus.Fields{"cmd": *utils.OldFileName + " -dir" + *utils.Dir, "error": err}).Error("error exec old file name")
			}
			os.Exit(1)
		}
	}

	// save the current pid and version
	if !utils.Mobile() {
		if err := savePid(); err != nil {
			log.WithFields(logrus.Fields{"error": err}).Fatal("can't save pid file")
		}
		defer delPidFile()
	}

	// database rollback to the specified block
	if *utils.RollbackToBlockID > 0 {
		rollbackToBlock(*utils.RollbackToBlockID)
		Exit(0)
	}

	dirpath := *utils.Dir + "/public"
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err = os.Mkdir(dirpath, 0755)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, "dir": dirpath}).Error("Error creating public dir")
			Exit(1)
		}
	}

	BrowserHTTPHost, _, ListenHTTPHost := GetHTTPHost()
	log.WithFields(logrus.Fields{"browser_host": BrowserHTTPHost, "listen_host": ListenHTTPHost}).Info("http hosts")
	if len(config.ConfigIni["db_type"]) > 0 {
		// The installation process is already finished (where user has specified DB and where wallet has been restarted)
		err = model.GormInit(config.ConfigIni["db_user"], config.ConfigIni["db_password"], config.ConfigIni["db_name"])
		if err != nil {
			log.WithFields(logrus.Fields{"db_user": config.ConfigIni["db_user"], "db_password": config.ConfigIni["db_password"],
				"db_name": config.ConfigIni["db_name"], "error": err}).Fatal("error initing gorm")
		}

		err = syspar.SysUpdate()
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Fatal("can't read system parameters")
		}

		log.WithFields(logrus.Fields{"workdir": *utils.Dir, "version": consts.VERSION}).Info("Start daemons")
		daemons.StartDaemons()
		daemonsTable := make(map[string]string)
		go func() {
			for {
				daemonNameAndTime := <-daemons.MonitorDaemonCh
				daemonsTable[daemonNameAndTime[0]] = daemonNameAndTime[1]
				if time.Now().Unix()%10 == 0 {
					log.Debug("daemonsTable: %v\n", daemonsTable)
				}
			}
		}()

		// signals for daemons to exit
		go stopdaemons.WaitStopTime()

		if err := template.LoadContracts(); err != nil {
			log.WithFields(logrus.Fields{"error": err}).Fatal("Error loading contracts, while starting project")
		}
		log.Info("Contracts loaded")

		tcpListener()
		log.Info("TCP Listener started")
		go controllers.GetChain()

	}

	stopdaemons.WaintForSignals()

	go func() {
		time.Sleep(time.Second)
		BrowserHTTPHost = initRoutes(ListenHTTPHost, BrowserHTTPHost)

		if *utils.Console == 0 && !utils.Mobile() {
			log.Info("Starting browser")
			time.Sleep(time.Second)
			if thrustWindowLoder != nil {
				thrustWindowLoder.Close()
				thrustWindow := thrust.NewWindow(thrust.WindowOptions{
					RootUrl: BrowserHTTPHost,
					Size:    commands.SizeHW{Width: 1024, Height: 700},
				})
				if *utils.DevTools != 0 {
					thrustWindow.OpenDevtools()
				}
				thrustWindow.HandleEvent("*", func(cr commands.EventResult) {
					log.WithFields(logrus.Fields{"event": cr}).Info("Handle event")
				})
				thrustWindow.HandleRemote(func(er commands.EventResult, this *window.Window) {
					if len(er.Message.Payload) > 7 && er.Message.Payload[:7] == `mailto:` && runtime.GOOS == `windows` {
						utils.ShellExecute(er.Message.Payload)
					} else if len(er.Message.Payload) > 7 && er.Message.Payload[:2] == `[{` {
						ioutil.WriteFile(filepath.Join(*utils.Dir, `accounts.txt`), []byte(er.Message.Payload), 0644)
					} else if er.Message.Payload == `ACCOUNTS` {
						accounts, _ := ioutil.ReadFile(filepath.Join(*utils.Dir, `accounts.txt`))
						this.SendRemoteMessage(string(accounts))
					} else {
						openBrowser(er.Message.Payload)
					}
					// Keep in mind once we have the message, lets say its json of some new type we made,
					// We can unmarshal it to that type.
					// Same goes for the other way around.
					//					this.SendRemoteMessage("boop")
				})
				thrustWindow.Show()
				thrustWindow.Focus()
			} else {
				openBrowser(BrowserHTTPHost)
			}
		}
	}()

	// waits for new records in chat, then waits for connect
	// (they are entered from the 'connections' daemon and from those who connected to the node by their own)
	// go utils.ChatOutput(utils.ChatNewTx)

	time.Sleep(time.Second * 3600 * 24 * 90)
}
