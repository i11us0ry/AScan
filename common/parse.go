package common

import (
	"flag"
	"AScan/common/utils/gologger"
	"os"
)

func Parse(options *Options) {
	if options.Version {
		gologger.Infof("Current Version: %s\n", Version)
		os.Exit(0)
	}
	if options.KeyWord == "" && options.CompanyID == "" && options.InputFile == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

}
