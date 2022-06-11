package runner

import (
	"AScan/ai"
	"AScan/common"
	"AScan/common/utils"
	"strconv"
	"time"
)

func RunEnumeration(options common.Options) {
	var info common.Web
	var sn common.Sun
	info = new(common.WebInfo)
	fn := common.GetConfigDir()+"\\result\\website_"+strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
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