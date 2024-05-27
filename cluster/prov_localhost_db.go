// replication-manager - Replication Manager Monitoring and CLI for MariaDB and MySQL
// Copyright 2017-2021 SIGNAL18 CLOUD SAS
// Authors: Guillaume Lefranc <guillaume@signal18.io>
//          Stephane Varoqui  <svaroqui@gmail.com>
// This source code is licensed under the GNU General Public License, version 3.

package cluster

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/signal18/replication-manager/config"
)

func (cluster *Cluster) LocalhostUnprovisionDatabaseService(server *ServerMonitor) error {
	cluster.LocalhostStopDatabaseService(server)
	cmd := exec.Command("rm", "-rf", server.Datadir)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s", err)
		cluster.errorChan <- err
		return err
	}
	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Remove datadir done: %s", out.Bytes())
	cluster.master = nil
	cluster.errorChan <- nil
	return nil
}

func (cluster *Cluster) LocalhostProvisionGetVersionFromMysqld(server *ServerMonitor) string {
	out := &bytes.Buffer{}
	versionCmd := exec.Command(cluster.Conf.ProvDBBinaryBasedir+"/mysqld", "--version")
	versionCmd.Stdout = out

	err := versionCmd.Run()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "mysqld version err: %s", out.Bytes())
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s", err)
		return ""
	}
	return strings.ToLower(string(out.Bytes()))
}

func (cluster *Cluster) LocalhostProvisionDatabaseService(server *ServerMonitor) error {

	out := &bytes.Buffer{}
	path := server.Datadir + "/var"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		nofile, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		nofile.Close()
	}

	/*
		//os.RemoveAll(path)

		cmd := exec.Command("rm", "-rf", path)

		cmd.Stdout = out
		err := cmd.Run()
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,config.LvlErr, "%s", err)
			cluster.errorChan <- err
			return err
		}
	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,LvlInfo, "Remove datadir done: %s", out.Bytes())*/
	server.GetDatabaseConfig()
	///	os.Symlink(server.Datadir+"/init/data", path)

	/*cmd = exec.Command("cp", "-rp", cluster.Conf.ShareDir+"/tests/data"+cluster.Conf.ProvDatadirVersion, path)

	// Attach buffer to command
	cmd.Stdout = out
	err = cmd.Run()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,config.LvlErr, "%s", err)
		return err
	}
	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,LvlInfo, "Copy fresh datadir done: %s", out.Bytes())
	*/
	cmd := exec.Command("cp", "-rp", server.Datadir+"/init/data/.", path)
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "cp -rp %s %s failed %s ", server.Datadir+"/init/data/.system", path, err)
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "init fresh datadir err: %s", out.Bytes())
		return err
	}
	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "copy datadir done: %s", out.Bytes())

	var sysCmd *exec.Cmd
	err = errors.New("No database version found")
	version := cluster.LocalhostProvisionGetVersionFromMysqld(server)
	if version == "" {
		cluster.errorChan <- err
		return err
	}
	if strings.Contains(version, "mariadb") {
		sysCmd = exec.Command(cluster.Conf.ProvDBClientBasedir+"/mysql_install_db", "--defaults-file="+server.Datadir+"/init/etc/mysql/my.cnf", "--datadir="+server.Datadir+"/var", "--basedir="+cluster.Conf.ProvDBBinaryBasedir+"/../", "--force")
	} else {
		sysCmd = exec.Command(cluster.Conf.ProvDBBinaryBasedir+"/mysqld", "--defaults-file="+server.Datadir+"/init/etc/mysql/my.cnf", "--datadir="+server.Datadir+"/var", "--basedir="+cluster.Conf.ProvDBBinaryBasedir+"/../", "--initialize", "--initialize-insecure")
	}
	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "%s", sysCmd.String())
	sysCmd.Stdout = out
	err = sysCmd.Run()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "init fresh datadir err: %s", out.Bytes())
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s", err)
		cluster.errorChan <- err
		return err
	}

	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "init fresh datadir done: %s", out.Bytes())
	if server.Id == "" {
		_, err := os.Stat(server.Id)
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, "TEST", "Found no os process continue with start ")
		}

	}

	/*	err := os.RemoveAll(path + "/" + server.Id + ".pid")
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,config.LvlErr, "%s", err)
			return err
		}*/

	err = cluster.LocalhostStartDatabaseServiceFistTime(server)
	if err != nil {
		cluster.errorChan <- err
		return err

	}
	cluster.errorChan <- nil
	return nil
}

func (cluster *Cluster) LocalhostStopDatabaseService(server *ServerMonitor) error {
	server.StopSlave()
	return server.Shutdown()
}

func (cluster *Cluster) LocalhostStartDatabaseServiceFistTime(server *ServerMonitor) error {

	if server.Id == "" {
		_, err := os.Stat(server.Id)
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, "TEST", "Found no os process continue with start ")
		}

	}
	path := server.Datadir + "/var"
	/*	err := os.RemoveAll(path + "/" + server.Id + ".pid")
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,config.LvlErr, "%s", err)
			return err
		}*/
	usr, err := user.Current()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s", err)
		return err
	}
	user := usr.Username
	version := cluster.LocalhostProvisionGetVersionFromMysqld(server)
	if version == "" {
		return errors.New("mysqld --version not found ")
	}
	time.Sleep(time.Millisecond * 2000)
	if strings.Contains(version, "mariadb") {
		user = "root"
	}
	mariadbdCmd := exec.Command(cluster.Conf.ProvDBBinaryBasedir+"/mysqld", "--defaults-file="+server.Datadir+"/init/etc/mysql/my.cnf", "--port="+server.Port, "--server-id="+server.Port, "--datadir="+path, "--socket="+server.GetDatabaseSocket(), "--user="+user, "--bind-address=0.0.0.0", "--pid_file="+path+"/"+server.Id+".pid")

	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "%s %s", mariadbdCmd.Path, mariadbdCmd.Args)

	var out bytes.Buffer
	mariadbdCmd.Stdout = &out

	go func() {
		err = mariadbdCmd.Run()
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s ", err)
		}
		fmt.Printf("Command finished with error: %v", err)
	}()
	exitloop := 0
	time.Sleep(time.Millisecond * 4000)
	for exitloop < 30 {
		haveerror := false
		time.Sleep(time.Millisecond * 2000)
		//cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,LvlInfo, "Waiting database startup ")
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Waiting database first start   .. %s", out)

		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "Can't get replication-manager process user: %s", err)
		}
		dsn := user + ":@unix(" + server.GetDatabaseSocket() + ")/?timeout=15s"
		conn, err2 := sqlx.Open("mysql", dsn)
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlDbg, "DNS start prov localhost first time : %s\n", dsn)
		if err2 == nil {
			defer conn.Close()
			_, err := conn.Exec("set sql_log_bin=0")
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", "set sql_log_bin=0", err)
			}

			_, err = conn.Exec("delete from mysql.user where password='' and user!='mariadb.sys'")
			if err != nil {
				//	haveerror = true
				// don't trigger error for mysql 5.7 and mariadb 10.4 that does not have password column

				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlWarn, " %s %s ", "delete from mysql.user where password=''", err)
			}
			grants := "grant all on *.* to '" + server.User + "'@'localhost' identified by '" + server.Pass + "' WITH GRANT OPTION"
			_, err = conn.Exec(grants)
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", grants, err)
			}
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "%s", grants)
			grants = "grant all on *.* to '" + server.User + "'@'%' identified by '" + server.Pass + "'"
			_, err = conn.Exec(grants)
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", grants, err)
			}
			grants = "grant all on *.* to '" + server.User + "'@'127.0.0.1' identified by '" + server.Pass + "'"
			_, err = conn.Exec(grants)
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", grants, err)
			}
			grants = "grant all on *.* to '" + server.ClusterGroup.GetRplUser() + "'@'localhost' identified by '" + server.ClusterGroup.GetRplPass() + "'"
			_, err = conn.Exec(grants)
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", grants, err)
			}
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "%s", grants)
			grants = "grant all on *.* to '" + server.ClusterGroup.GetRplUser() + "'@'%' identified by '" + server.ClusterGroup.GetRplPass() + "'"
			_, err = conn.Exec(grants)
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", grants, err)
			}
			grants = "grant all on *.* to '" + server.ClusterGroup.GetRplUser() + "'@'127.0.0.1' identified by '" + server.ClusterGroup.GetRplPass() + "'"
			_, err = conn.Exec(grants)
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", grants, err)
			}
			_, err = conn.Exec("flush privileges")
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", "flush privileges", err)
			}

			_, err = conn.Exec("reset master")
			if err != nil {
				haveerror = true
				cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, " %s %s ", "reset master", err)
			}

			if !haveerror {
				exitloop = 100
			}

		} else {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "Database connection to init user  %s ", err2)
		}
		exitloop++

	}
	if exitloop == 101 {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Database started.")

	} else {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Database timeout.")
		return errors.New("Failed to start")
	}

	//	mariadbdCmd.Process.Release()

	return nil
}

func (cluster *Cluster) LocalhostStartDatabaseService(server *ServerMonitor) error {
	server.GetDatabaseConfig()
	if server.Id == "" {
		_, err := os.Stat(server.Id)
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, "TEST", "Found no os process continue with start ")
		}

	}
	path := server.Datadir + "/var"
	/*	err := os.RemoveAll(path + "/" + server.Id + ".pid")
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,config.LvlErr, "%s", err)
			return err
		}*/
	usr, err := user.Current()
	if err != nil {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s", err)
		return err
	}
	//	mariadbdCmd := exec.Command(cluster.Conf.ProvDBBinaryBasedir+"/mysqld", "--defaults-file="+server.Datadir+"/init/etc/mysql/my.cnf --port="+server.Port, "--server-id="+server.Port, "--datadir="+path, "--socket="+server.Datadir+"/"+server.Id+".sock", "--user="+usr.Username, "--bind-address=0.0.0.0", "--general_log=1", "--general_log_file="+path+"/"+server.Id+".log", "--pid_file="+path+"/"+server.Id+".pid", "--log-error="+path+"/"+server.Id+".err")
	time.Sleep(time.Millisecond * 2000)
	mariadbdCmd := exec.Command(cluster.Conf.ProvDBBinaryBasedir+"/mysqld", "--defaults-file="+server.Datadir+"/init/etc/mysql/my.cnf", "--port="+server.Port, "--server-id="+server.Port, "--datadir="+path, "--socket="+server.GetDatabaseSocket(), "--user="+usr.Username, "--bind-address=0.0.0.0", "--pid_file="+path+"/"+server.Id+".pid")
	cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "%s %s", mariadbdCmd.Path, mariadbdCmd.Args)

	var out bytes.Buffer
	mariadbdCmd.Stdout = &out

	go func() {
		err = mariadbdCmd.Run()
		if err != nil {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "%s ", err)
		}
		fmt.Printf("Command finished with error: %v", err)
	}()

	exitloop := 0
	time.Sleep(time.Millisecond * 4000)
	for exitloop < 30 {

		time.Sleep(time.Millisecond * 2000)
		//cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator,LvlInfo, "Waiting database startup ")
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Waiting database startup %d: %s", exitloop, out.String())
		conn, err2 := sqlx.Open("mysql", server.DSN)
		if err2 == nil {
			defer conn.Close()
			exitloop = 100

		} else {
			cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlErr, "Database connection to init user  %s ", err2)
		}
		exitloop++

	}
	if exitloop == 101 {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Database started.")

	} else {
		cluster.LogModulePrintf(cluster.Conf.Verbose, config.ConstLogModOrchestrator, config.LvlInfo, "Database timeout.")
		return errors.New("Failed to start")
	}
	server.Process = mariadbdCmd.Process
	//	mariadbdCmd.Process.Release()

	return nil
}
