package daylight

import "flag"
import "github.com/EGaaS/go-egaas-mvp/packages/controllers"

func install() error {
	dir := *flag.String("dir", "./", "installation directory")
	genFirstBlock := *flag.Int64("gen_first_block", 1, "should first block be generated")
	tcpHost := *flag.String("host", "7079", "installation host")
	httpPort := *flag.String("port", "localhost", "installation port")
	logLevel := *flag.String("log_level", "ERROR", "Log level. Should be DEBUG or ERRROR")
	if logLevel != "DEBUG" && logLevel != "ERROR" {
		logLevel = "ERROR"
	}

	firstLoadBlockchainURL := *flag.String("first_load_url", "", "Blockchain load url")
	firstLoad := *flag.String("first_load", "", "")
	dbType := "postgres"
	dbHost := *flag.String("db_host", "5432", "postgres host")
	dbPort := *flag.String("db_port", "5432", "postgres port")
	dbName := *flag.String("db_name", "", "database name")
	dbUsername := *flag.String("username", "", "database username")
	dbPassword := *flag.String("password", "", "database password")

	c := &controllers.Controller{}
	err := c.LocalInstall(dir, genFirstBlock, tcpHost, httpPort, logLevel, firstLoadBlockchainURL, firstLoad, dbType, dbHost,
		dbPort, dbName, dbUsername, dbPassword)
	return err
}
