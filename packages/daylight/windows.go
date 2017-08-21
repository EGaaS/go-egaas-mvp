// +build windows

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
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/EGaaS/go-egaas-mvp/packages/model"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

// KillPid kills the process with the specified pid
func KillPid(pid string) error {
	if model.DBConn != nil {
		sd := &model.StopDaemon{StopTime: time.Now().Unix()}
		err := sd.Create()
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, "pid": pid}).Error("Error inserting into stop_daemons, when killing process")
			return err
		}
	}
	cmd := "tasklist /fi PID eq" + pid
	res, err := exec.Command("tasklist", "/fi", "PID eq "+pid).Output()
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "pid": pid, "command": cmd}).Error("error executing command, when killing process")
		return err
	}
	if string(res) == "" {
		log.WithFields(logrus.Fields{"pid": pid, "command": cmd}).Error("command returned no result, when killing process")
		return fmt.Errorf("null")
	}
	log.WithFields(logrus.Fields{"pid": pid, "command": cmd}).Info("command returned result")
	if ok, _ := regexp.MatchString(`(?i)PID`, string(res)); !ok {
		log.WithFields(logrus.Fields{"pid": pid, "command": cmd, "output": res}).Error("command returned incorrect result, when killing process")
		return fmt.Errorf("null")
	}
	return nil
}

func tray() {

}
