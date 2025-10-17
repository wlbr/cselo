package elo

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-gcfg/gcfg"

	"github.com/wlbr/commons"
	"github.com/wlbr/commons/log"
)

type Config struct {
	commons.CommonConfig
	ConfigFileName string
	Elo            struct {
		ImportFileName   string
		RecorderFileName string
		OutputFileName   string
		OutputFile       *os.File
		Port             string
		ForceOverwrite   bool
		ExportDatafiles  bool
		DiscordWebhook   string
		BufferSize       int
	}
	PostgreSQL struct {
		Host     string
		Port     string
		Database string
		User     string
		Password string
	}
	InfluxDB struct {
		Host   string
		Port   string
		Token  string
		Bucket string
		Org    string
	}
}

func (c *Config) String() string {
	pw := "\"\""
	if c.PostgreSQL.Password != "" {
		pw = "***"
	}
	tok := "\"\""
	if c.InfluxDB.Token != "" {
		tok = "***"
	}
	return fmt.Sprintf("%s\tPort: %s\n"+
		"\tForceoverWrite: %t\n"+
		"\tExportDatafiles: %t\n"+
		"\tOutputFileName: %s\n"+
		"\tLogFileName: %s\n"+
		"\tRecorderFileName: %s\n"+
		"\tDiscordWebhook: %s\n"+
		"\tBufferSize:%d\n"+
		"\tConfigFileName: %s\n\n"+
		"\tPostgreSQL:\n"+
		"\t\tHost: %s\n"+
		"\t\tPort: %s\n"+
		"\t\tDatabase: %s\n"+
		"\t\tUser: %s\n"+
		"\t\tPassword: %s\n"+
		"\tInfluxDB: \n"+
		"\t\tHost: %s\n"+
		"\t\tPort: %s\n"+
		"\t\tToken: %s\n"+
		"\t\tBucket: %s\n"+
		"\t\tOrg: %s\n",

		c.CommonConfig, c.Elo.Port, c.Elo.ForceOverwrite, c.Elo.ExportDatafiles, c.Elo.OutputFileName, c.Elo.ImportFileName, c.Elo.RecorderFileName, c.Elo.DiscordWebhook, c.Elo.BufferSize,
		c.ConfigFileName, c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.Database, c.PostgreSQL.User, pw, c.InfluxDB.Host, c.InfluxDB.Port, tok, c.InfluxDB.Bucket, c.InfluxDB.Org)
}

func (cfg Config) CheckForceOverwrite() {
	log.Info("Checking forceoverwrite parameter = %t", cfg.Elo.ForceOverwrite)
	if cfg.Elo.OutputFileName != "" && cfg.Elo.OutputFileName != "<STDOUT>" {
		info, err := os.Stat(cfg.Elo.OutputFileName)
		if err == nil {
			if info.IsDir() {
				log.Fatal("Outputfile is a directory, not overwriting.")
				cfg.FatalExit()
			}
			if !cfg.Elo.ForceOverwrite {
				//log.Fatal("Outputfile exists, not overwriting. Use -f to force overwrite.")
				//cfg.FatalExit()
				log.Warn("Outputfile exists, appending")
			}
		}
	}
}

func (cfg *Config) FlagDefinition() {
	cfg.CommonConfig.FlagDefinition()
	flag.StringVar(&cfg.Elo.Port, "port", "42820", "The port to listen on.")
	flag.StringVar(&cfg.Elo.ImportFileName, "import", "", "Use a captured logfile instead of listening to the net.")
	flag.StringVar(&cfg.Elo.RecorderFileName, "rec", "", "Save captured data to file.")
	flag.BoolVar(&cfg.Elo.ForceOverwrite, "f", false, "Overwrite all output files.")
	flag.StringVar(&cfg.ConfigFileName, "cfg", "", "File to read (partial) config from.")
	flag.IntVar(&cfg.Elo.BufferSize, "size", 5000, "Buffersize.")
}

func (cfg *Config) Initialize(version string, buildtimestamp string) *Config {
	if !flag.Parsed() {
		cfg.FlagDefinition()
		flag.Parse()
	}
	if cfg.ConfigFileName != "" {
		if err := gcfg.ReadFileInto(cfg, cfg.ConfigFileName); err != nil {
			log.Error("Error reading ini-file: %s", err)
		}
	}
	cfg.CommonConfig.Initialize(version, buildtimestamp)

	// Outputfile
	cfg.Elo.OutputFileName = flag.Arg(0)
	if cfg.Elo.OutputFileName == "" || strings.ToUpper(cfg.Elo.OutputFileName) == "STDOUT" {
		log.Info("No output file file given as argument. Using <STDOUT>.")
		cfg.Elo.OutputFileName = "<STDOUT>"
		cfg.Elo.OutputFile = os.Stdout
	} else {
		cfg.CheckForceOverwrite()
		var err error
		//cfg.Elo.OutputFile, err = os.Create(cfg.Elo.OutputFileName)
		cfg.Elo.OutputFile, err = os.OpenFile(cfg.Elo.OutputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatal("Cannot create output file %s. Error: %s", cfg.Elo.OutputFileName, err)
		}
		cfg.AddCleanUpFn(cfg.Elo.OutputFile.Close)
	}

	if cfg.ActiveLogLevel >= log.INFO {
		fmt.Println(cfg.GetInspectData())
		fmt.Println(cfg.String())
	}
	return cfg
}
