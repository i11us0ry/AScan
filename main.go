package main

import (
	"AScan/common"
	"AScan/common/utils"
	"AScan/runner"
	"os"
)

func main() {
	var op common.Options
	conf := common.CheckConf()
	if !utils.FolderExists(conf.Output){
		os.Mkdir(conf.Output, 0777)
	}
	common.Flag(&op)
	common.Parse(&op)
	op.CookieInfo = conf.Cookies
	op.Output = conf.Output
	runner.RunEnumeration(op)
}