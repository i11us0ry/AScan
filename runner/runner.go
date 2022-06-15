package runner

import (
	"AScan/ai"
	"AScan/common"
	"AScan/common/utils"
	"AScan/common/utils/gologger"
	"os"
	"runtime"
	"strconv"
	"time"
)

func RunEnumeration(options common.Options) {
	var info common.Web
	var sn common.Sun
	info = new(common.WebInfo)
	sysType := runtime.GOOS
	fn := ""
	if sysType == "windows"{
		fn = common.GetConfigDir()+"\\result\\website_"+strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	} else {
		fn = common.GetConfigDir()+"/result/website_"+strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
		f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777) //打开文件
		if err!=nil{
			gologger.Errorf("file create err!")
		}
		f.Close()
	}
	info.SetFileName(fn)
	ai.Init(info)
	if options.InputFile!=""{
		ns := utils.FileReadByline(options.InputFile)
		for _, n := range(ns){
			options.KeyWord = n
			ai.GetEnInfoByPid(options,&sn)
		}
	} else {
		ai.GetEnInfoByPid(options,&sn)
	}
	// 输出所有子公司名称
	if len(sn.Name)!=0{
		common.WriteSun(fn,sn.Name)
	}
}