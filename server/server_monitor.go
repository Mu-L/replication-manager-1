//go:build !clients
// +build !clients

// replication-manager - Replication Manager Monitoring and CLI for MariaDB and MySQL
// Copyright 2017-2021 SIGNAL18 CLOUD SAS
// Authors: Guillaume Lefranc <guillaume@signal18.io>
//          Stephane Varoqui  <svaroqui@gmail.com>
// This source code is licensed under the GNU General Public License, version 3.
// Redistribution/Reuse of this code is permitted under the GNU v3 license, as
// an additional term, ALL code must carry the original Author(s) credit in comment form.
// See LICENSE in this directory for the integral text.

package server

import (
	"fmt"
	"io/ioutil"

	mysqllog "log"

	"github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultFlagMap map[string]interface{}

func init() {

	conf.ProvOrchestrator = "local"
	var errLog = mysql.Logger(mysqllog.New(ioutil.Discard, "", 0))
	mysql.SetLogger(errLog)

	//monitorCmd.AddCommand(rootCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.AddCommand(configMergeCmd)
	RepMan.AddFlags(monitorCmd.Flags(), &conf)

	//configMergeCmd.PersistentFlags().StringVar(&cfgGroup, "cluster", "", "Cluster name (default is none)")
	//configMergeCmd.PersistentFlags().StringVar(&conf.ConfigFile, "config", "", "Configuration file (default is config.toml)")

	//cobra.OnInitialize()
	initLogFlags(monitorCmd)
	initLogFlags(configMergeCmd)

	//conf des defaults flag sans les paramètres en ligne de commande
	v := viper.GetViper()
	v.SetConfigType("toml")

	defaultFlagMap = make(map[string]interface{})

	for _, f := range v.AllKeys() {
		defaultFlagMap[f] = v.Get(f)
	}

	viper.BindPFlags(monitorCmd.Flags())
	viper.BindPFlags(configMergeCmd.Flags())

}

func initDeprecated() {
	//not needed use Alias in server.go

	monitorCmd.Flags().BoolVar(&conf.ProxysqlCopyGrants, "proxysql-copy-grants", true, "Deprecate copy grants from master")
	monitorCmd.Flags().MarkDeprecated("proxysql-copy-grants", "Deprecated for proxysql-bootstrap-users")
	monitorCmd.Flags().StringVar(&conf.BackupMyDumperPath, "mydumper-path", "/usr/bin/mydumper", "Deprecate Path to mydumper binary")
	monitorCmd.Flags().MarkDeprecated("mydumper-path", "Deprecated for backup-mydumper-path")
	monitorCmd.Flags().StringVar(&conf.BackupMyLoaderPath, "myloader-path", "/usr/bin/myloader", "Deprecate Path to myloader binary")
	monitorCmd.Flags().MarkDeprecated("myloader-path", "Deprecated for backup-myloader-path")
	monitorCmd.Flags().StringVar(&conf.BackupMysqldumpPath, "mysqldump-path", "", "Deprecate Path to mysqldump binary")

	monitorCmd.Flags().MarkDeprecated("mysqldump-path", "Deprecated for backup-mysqldump-path")
	monitorCmd.Flags().StringVar(&conf.BackupMysqlbinlogPath, "mysqlbinlog-path", "", "Deprecate Path to mysqlbinlog binary")
	monitorCmd.Flags().MarkDeprecated("mysqlbinlog-path", "Deprecated for backup-mysqlbinlog-path")
	monitorCmd.Flags().StringVar(&conf.BackupMysqlclientPath, "mysqlclient-path", "", "Deprecate Path to mysql client binary")
	monitorCmd.Flags().MarkDeprecated("mysqlclient-path", "Deprecated for backup-mysqlclient-path")
	monitorCmd.Flags().StringVar(&conf.HaproxyBinaryPath, "haproxy-binary-path", "/usr/sbin/haproxy", "HaProxy binary location")
	monitorCmd.Flags().StringVar(&conf.MasterConn, "replication-master-connection", "", "Connection name to use for multisource replication")
	monitorCmd.Flags().MarkDeprecated("replication-master-connection", "Depecrate for replication-source-name")
	monitorCmd.Flags().StringVar(&conf.LogFile, "logfile", "", "Write output messages to log file")
	monitorCmd.Flags().MarkDeprecated("logfile", "Deprecate for log-file")
	monitorCmd.Flags().Int64Var(&conf.SwitchWaitKill, "wait-kill", 5000, "Deprecate for switchover-wait-kill Wait this many milliseconds before killing threads on demoted master")
	monitorCmd.Flags().MarkDeprecated("wait-kill", "Deprecate for switchover-wait-kill Wait this many milliseconds before killing threads on demoted master")
	monitorCmd.Flags().StringVar(&conf.User, "user", "", "User for database login, specified in the [user]:[password] format")
	monitorCmd.Flags().MarkDeprecated("user", "Deprecate for db-servers-credential")

	monitorCmd.Flags().StringVar(&conf.ProvDBBinaryBasedir, "db-servers-binary-path", "/usr/local/mysql/bin", "Deprecate Path to mysqld binary for testing")
	monitorCmd.Flags().MarkDeprecated("db-servers-binary-path", "Deprecate for prov-db-binary-basedir")

	monitorCmd.Flags().StringVar(&conf.Hosts, "hosts", "", "List of database hosts IP and port (optional), specified in the host:[port] format and separated by commas")
	monitorCmd.Flags().MarkDeprecated("hosts", "Deprecate for db-servers-hosts")
	monitorCmd.Flags().StringVar(&conf.HostsTLSCA, "hosts-tls-ca-cert", "", "TLS authority certificate")
	monitorCmd.Flags().MarkDeprecated("hosts-tls-ca-cert", "Deprecate for db-servers-tls-ca-cert")
	monitorCmd.Flags().StringVar(&conf.HostsTlsCliKey, "hosts-tls-client-key", "", "TLS client key")
	monitorCmd.Flags().MarkDeprecated("hosts-tls-client-key", "Deprecate for db-servers-tls-client-key")
	monitorCmd.Flags().StringVar(&conf.HostsTlsCliCert, "hosts-tls-client-cert", "", "TLS client certificate")
	monitorCmd.Flags().MarkDeprecated("hosts-tls-client-cert", "Deprecate for db-servers-tls-client-cert")
	monitorCmd.Flags().IntVar(&conf.Timeout, "connect-timeout", 5, "Database connection timeout in seconds")
	monitorCmd.Flags().MarkDeprecated("connect-timeout", "Deprecate for db-servers-connect-timeout")
	monitorCmd.Flags().StringVar(&conf.RplUser, "rpluser", "", "Replication user in the [user]:[password] format")
	monitorCmd.Flags().MarkDeprecated("rpluser", "Deprecate for replication-credential")
	monitorCmd.Flags().StringVar(&conf.PrefMaster, "prefmaster", "", "Preferred candidate server for master failover, in host:[port] format")
	monitorCmd.Flags().MarkDeprecated("prefmaster", "Deprecate for db-servers-prefered-master")
	monitorCmd.Flags().StringVar(&conf.IgnoreSrv, "ignore-servers", "", "List of servers to ignore in slave promotion operations")
	monitorCmd.Flags().MarkDeprecated("ignore-servers", "Deprecate for db-servers-ignored-hosts")
	monitorCmd.Flags().StringVar(&conf.MasterConn, "master-connection", "", "Connection name to use for multisource replication")
	monitorCmd.Flags().MarkDeprecated("master-connection", "Deprecate for replication-master-connection")
	monitorCmd.Flags().IntVar(&conf.MasterConnectRetry, "master-connect-retry", 10, "Specifies how many seconds to wait between slave connect retries to master")
	monitorCmd.Flags().MarkDeprecated("master-connect-retry", "Deprecate for replication-master-connection-retry")
	monitorCmd.Flags().MarkDeprecated("api-user", "Deprecate for 	api-credential")
	monitorCmd.Flags().BoolVar(&conf.ReadOnly, "readonly", true, "Set slaves as read-only after switchover failover")
	monitorCmd.Flags().MarkDeprecated("readonly", "Deprecate for failover-readonly-state")
	monitorCmd.Flags().StringVar(&conf.MxsHost, "maxscale-host", "", "MaxScale host IP")
	monitorCmd.Flags().MarkDeprecated("maxscale-host", "Deprecate for maxscale-servers")
	monitorCmd.Flags().StringVar(&conf.MdbsProxyHosts, "mdbshardproxy-hosts", "127.0.0.1:3307", "MariaDB spider proxy hosts IP:Port,IP:Port")
	monitorCmd.Flags().MarkDeprecated("mdbshardproxy-hosts", "Deprecate for mdbshardproxy-servers")
	monitorCmd.Flags().BoolVar(&conf.MultiMaster, "multimaster", false, "Turn on multi-master detection")
	monitorCmd.Flags().MarkDeprecated("multimaster", "Deprecate for replication-multi-master")
	monitorCmd.Flags().BoolVar(&conf.MultiTierSlave, "multi-tier-slave", false, "Turn on to enable relay slaves in the topology")
	monitorCmd.Flags().MarkDeprecated("multi-tier-slaver", "Deprecate for replication-multi-tier-slave")
	monitorCmd.Flags().StringVar(&conf.PreScript, "pre-failover-script", "", "Path of pre-failover script")
	monitorCmd.Flags().MarkDeprecated("pre-failover-script", "Deprecate for failover-pre-script")
	monitorCmd.Flags().StringVar(&conf.PostScript, "post-failover-script", "", "Path of post-failover script")
	monitorCmd.Flags().MarkDeprecated("post-failover-script", "Deprecate for failover-post-script")
	monitorCmd.Flags().StringVar(&conf.RejoinScript, "rejoin-script", "", "Path of old master rejoin script")
	monitorCmd.Flags().MarkDeprecated("rejoin-script", "Deprecate for autorejoin-script")
	monitorCmd.Flags().StringVar(&conf.ShareDir, "share-directory", "/usr/share/replication-manager", "Path to HTTP monitor share files")
	monitorCmd.Flags().MarkDeprecated("share-directory", "Deprecate for monitoring-sharedir")
	monitorCmd.Flags().StringVar(&conf.WorkingDir, "working-directory", "/var/lib/replication-manager", "Path to HTTP monitor working directory")
	monitorCmd.Flags().MarkDeprecated("working-directory", "Deprecate for monitoring-datadir")
	monitorCmd.Flags().BoolVar(&conf.Interactive, "interactive", true, "Ask for user interaction when failures are detected")
	monitorCmd.Flags().MarkDeprecated("interactive", "Deprecate for failover-mode")
	monitorCmd.Flags().IntVar(&conf.MaxFail, "failcount", 5, "Trigger failover after N failures (interval 1s)")
	monitorCmd.Flags().MarkDeprecated("failcount", "Deprecate for failover-falsepositive-ping-counter")
	monitorCmd.Flags().IntVar(&conf.SwitchWaitWrite, "wait-write-query", 10, "Deprecate  Wait this many seconds before write query end to cancel switchover")
	monitorCmd.Flags().MarkDeprecated("wait-write-query", "Deprecate for switchover-wait-write-query")
	monitorCmd.Flags().Int64Var(&conf.SwitchWaitTrx, "wait-trx", 10, "Depecrate for switchover-wait-trx Wait this many seconds before transactions end to cancel switchover")
	monitorCmd.Flags().MarkDeprecated("wait-trx", "Deprecate for switchover-wait-trx")
	monitorCmd.Flags().BoolVar(&conf.SwitchGtidCheck, "gtidcheck", false, "Depecrate for failover-at-equal-gtid do not initiate failover unless slaves are fully in sync")
	monitorCmd.Flags().MarkDeprecated("gtidcheck", "Deprecate for switchover-at-equal-gtid")
	monitorCmd.Flags().Int64Var(&conf.FailMaxDelay, "maxdelay", 0, "Deprecate Maximum replication delay before initiating failover")
	monitorCmd.Flags().MarkDeprecated("maxdelay", "Deprecate for failover-max-slave-delay")
	monitorCmd.Flags().StringVar(&conf.APIUsers, "api-credential", "admin:repman", "Deprecate for api-credentials")
	monitorCmd.Flags().MarkDeprecated("api-credential", "Deprecate for failover-max-slave-delay")

}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Starts monitoring server",
	Long: `Starts replication-manager server in stateful monitor daemon mode.

For interacting with this daemon use,
- Interactive console client: "replication-manager-cli ".
- Command line clients: "replication-manager-cli switchover|failover|topology|test".
- HTTP dashboards on port 10001
- HTTPS dashboards on port 10005

`,
	Run: func(cmd *cobra.Command, args []string) {

		RepMan = new(ReplicationManager)
		RepMan.CommandLineFlag = GetCommandLineFlag(cmd)
		RepMan.DefaultFlagMap = defaultFlagMap
		RepMan.InitConfig(conf)
		RepMan.Run()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		// Close connections on exit.
		RepMan.Stop()
	},
}

var configMergeCmd = &cobra.Command{
	Use:   "config-merge",
	Short: "Merges the initial configuration file with the dynamic one",
	Long:  `Merges all parameters modified in dynamic mode with the original parameters (including immutable parameters) by merging the config files generated by the dynamic mode. Be careful, this command overwrites the original config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Start config-merge command !\n")
		fmt.Printf("Cluster: %s\n", cfgGroup)
		fmt.Printf("Config : %s\n", conf.ConfigFile)
		RepMan = new(ReplicationManager)
		RepMan.DefaultFlagMap = defaultFlagMap
		RepMan.InitConfig(conf)
		ImmFlagMap := RepMan.ImmuableFlagMaps[cfgGroup]
		err := conf.MergeConfig(conf.WorkingDir, cfgGroup, ImmFlagMap, RepMan.DefaultFlagMap, conf.ConfigFile)
		if err != nil {
			fmt.Printf("Config merge command fail: %s", err)
			return
		}
		fmt.Printf("Success of the configuration merge command!")
		//RepMan.DynamicFlagMaps
	},
}

func GetCommandLineFlag(cmd *cobra.Command) []string {
	var cmd_flag []string
	flag := viper.AllKeys()
	for _, f := range flag {
		if cmd.Flags().Changed(f) {
			cmd_flag = append(cmd_flag, f)
		}
	}
	return cmd_flag
}
