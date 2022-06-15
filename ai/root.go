package ai

import (
	"AScan/common"
	"AScan/common/utils"
	"AScan/common/utils/gologger"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xuri/excelize/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Info common.Web
)

func Init(info common.Web){
	Info = info
}

// pageParseJson 提取页面中的JSON字段
func pageParseJson(content string) gjson.Result {
	tag1 := "window.pageData ="
	tag2 := "window.isSpider ="
	//tag2 := "/* eslint-enable */</script><script data-app"
	idx1 := strings.Index(content, tag1)
	idx2 := strings.Index(content, tag2)
	if idx2 > idx1 {
		str := content[idx1+len(tag1) : idx2]
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, " ", "", -1)
		str = str[:len(str)-1]
		return gjson.Get(string(str), "result")
	}
	return gjson.Result{}
}

// 第一步 GetEnInfoByPid 根据PID获取公司信息
func GetEnInfoByPid(options common.Options,sn *common.Sun) {
	pid := ""
	enop := common.Options{}
	if options.CompanyID == "" {
		_, enop = SearchName(options)
	} else {
		enop = options
	}
	pid = enop.CompanyID
	if pid == "" {
		gologger.Errorf("没有获取到PID\n")
		return
	}
	gologger.Infof("查询PID %s\n", pid)

	//获取公司信息
	res := getCompanyInfoById(pid, true, enop,sn)
	outPutExcelByEnInfo(res, enop, sn)

}

func outPutExcelByEnInfo(enInfo EnInfo, options common.Options,sn *common.Sun) {
	f := excelize.NewFile()
	//Base info
	baseHeaders := []string{"信息", "值"}
	baseData := [][]interface{}{
		{"PID", enInfo.Pid},
		{"企业名称", enInfo.EntName},
		{"法人代表", enInfo.legalPerson},
		{"开业状态", enInfo.openStatus},
		{"电话", enInfo.telephone},
		{"邮箱", enInfo.email},
		{"网址",enInfo.webSite},
	}
	if enInfo.webSite!="-" && enInfo.webSite!=""{
		Info.Check(enInfo.webSite,enInfo.EntName)
	}
	f, _ = utils.ExportExcel("基本信息", baseHeaders, baseData, f)

	for k, s := range enInfo.ensMap {
		if s.total > 0 && s.api != "" {
			//gologger.Infof("正在导出%s\n", s.name)
			headers := s.keyWord
			var data [][]interface{}
			for _, y := range enInfo.infos[k] {
				results := gjson.GetMany(y.Raw, s.field...)
				var str []interface{}
				for _, s1 := range results {
					str = append(str, s1.String())
					// 新增 聚合备案和基本信息中的website
					if s.name == "网站备案"{
						if results[0].String()!="-" && results[0].String()!=""{
							Info.Check(results[0].String(),results[1].String())
						}
					}
				}
				data = append(data, str)
			}
			f, _ = utils.ExportExcel(s.name, headers, data, f)
		}
	}

	f.DeleteSheet("Sheet1")
	// Save spreadsheet by the given path.
	sysType := runtime.GOOS
	savaPath := ""
	if sysType == "windows"{
		savaPath = options.Output + "\\" +
			time.Now().Format("2006-01-02") +
			enInfo.EntName + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	} else {
		savaPath = options.Output + "/" +
			time.Now().Format("2006-01-02") +
			enInfo.EntName + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	}

	// 将所有子公司名字收集
	sn.Name = append(sn.Name, enInfo.EntName)
	if err := f.SaveAs(savaPath); err != nil {
		gologger.Fatalf("导出失败：%s", err)
	}
	//gologger.Infof("导出成功路径： %s\n", savaPath)

}

// getCompanyInfoById 获取公司基本信息
// pid 公司id
// isSearch 是否递归搜索信息【分支机构、对外投资信息】
// options options
func getCompanyInfoById(pid string, isSearch bool, options common.Options,sn *common.Sun) EnInfo {
	var enInfo EnInfo
	enInfo.infos = make(map[string][]gjson.Result)
	urls := "https://aiqicha.baidu.com/company_detail_" + pid
	content := common.GetReq(urls, options)
	res := pageParseJson(string(content))
	//获取企业基本信息情况
	enInfo.Pid = res.Get("pid").String()
	enInfo.EntName = res.Get("entName").String()
	enInfo.legalPerson = res.Get("legalPerson").String()
	enInfo.openStatus = res.Get("openStatus").String()
	enInfo.telephone = res.Get("telephone").String()
	enInfo.email = res.Get("email").String()
	enInfo.webSite = res.Get("website").String()
	//gologger.Infof("企业基本信息\n")
	data := [][]string{
		{"PID", enInfo.Pid},
		{"企业名称", enInfo.EntName},
		{"法人代表", enInfo.legalPerson},
		{"开业状态", enInfo.openStatus},
		{"电话", enInfo.telephone},
		{"邮箱", enInfo.email},
		{"网址",enInfo.webSite},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(data)
	table.Render()

	//判断企业状态，不然就可以跳过了
	if enInfo.openStatus == "注销" || enInfo.openStatus == "吊销" {
		return enInfo
	}

	//获取企业信息
	enInfoUrl := "https://aiqicha.baidu.com/compdata/navigationListAjax?pid=" + pid
	enInfoRes := common.GetReq(enInfoUrl, options)
	ensInfoMap := make(map[string]*EnsGo)
	if gjson.Get(string(enInfoRes), "status").String() == "0" {
		data := gjson.Get(string(enInfoRes), "data").Array()
		for _, s := range data {
			for _, t := range s.Get("children").Array() {
				ensInfoMap[t.Get("id").String()] = &EnsGo{
					t.Get("name").String(),
					t.Get("total").Int(),
					t.Get("avaliable").Int(),
					"",
					[]string{},
					[]string{},
				}
			}
		}
	}

	//赋值API数据
	ensInfoMap["webRecord"].api = "detail/icpinfoAjax"
	ensInfoMap["webRecord"].field = []string{"domain", "siteName", "homeSite", "icpNo"}
	ensInfoMap["webRecord"].keyWord = []string{"域名", "站点名称", "首页", "ICP备案号"}

	//ensInfoMap["appinfo"].api = "c/appinfoAjax"
	//ensInfoMap["appinfo"].field = []string{"name", "classify", "logoWord", "logoBrief", "entName"}
	//ensInfoMap["appinfo"].keyWord = []string{"APP名称", "分类", "LOGO文字", "描述", "所属公司"}

	//ensInfoMap["microblog"].api = "c/microblogAjax"
	//ensInfoMap["microblog"].field = []string{"nickname", "weiboLink", "logo"}
	//ensInfoMap["microblog"].keyWord = []string{"微博昵称", "链接", "LOGO"}

	ensInfoMap["wechatoa"].api = "c/wechatoaAjax"
	ensInfoMap["wechatoa"].field = []string{"wechatName", "wechatId", "wechatIntruduction", "wechatLogo", "qrcode", "entName"}
	ensInfoMap["wechatoa"].keyWord = []string{"名称", "ID", "描述", "LOGO", "二维码", "归属公司"}

	//ensInfoMap["enterprisejob"].api = "c/enterprisejobAjax"
	//ensInfoMap["enterprisejob"].field = []string{"jobTitle", "location", "education", "publishDate", "desc"}
	//ensInfoMap["enterprisejob"].keyWord = []string{"职位名称", "工作地点", "学历要求", "发布日期", "招聘描述"}

	ensInfoMap["copyright"].api = "detail/copyrightAjax"
	ensInfoMap["copyright"].field = []string{"softwareName", "shortName", "softwareType", "typeCode", "regDate"}
	ensInfoMap["copyright"].keyWord = []string{"软件名称", "软件简介", "分类", "行业", "日期"}

	//ensInfoMap["supplier"].api = "c/supplierAjax"
	//ensInfoMap["supplier"].field = []string{"supplier", "source", "principalNameClient", "cooperationDate"}
	//ensInfoMap["supplier"].keyWord = []string{"供应商名称", "来源", "所属公司", "日期"}

	ensInfoMap["invest"].api = "detail/investajax" //对外投资
	ensInfoMap["invest"].field = []string{"entName", "openStatus", "regRate", "data"}
	ensInfoMap["invest"].keyWord = []string{"公司名称", "状态", "投资比例", "数据信息"}

	ensInfoMap["branch"].api = "detail/branchajax" //分支机构
	ensInfoMap["branch"].field = []string{"entName", "openStatus", "data"}
	ensInfoMap["branch"].keyWord = []string{"公司名称", "状态", "数据信息"}

	enInfo.ensMap = ensInfoMap

	//获取数据
	for k, s := range ensInfoMap {
		if s.total > 0 && s.api != "" {
			gologger.Infof("正在查询 %s\n", s.name)
			t := getInfoList(res.Get("pid").String(), s.api, options)

			//判断下网站备案，然后提取出来，留个坑看看有没有更好的解决方案
			if k == "webRecord" {
				var tmp []gjson.Result
				for _, y := range t {
					for _, o := range y.Get("domain").Array() {
						value, _ := sjson.Set(y.Raw, "domain", o.String())
						tmp = append(tmp, gjson.Parse(value))
					}
				}
				t = tmp
			}
			//保存数据
			enInfo.infos[k] = t

			//命令输出展示
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(ensInfoMap[k].keyWord)
			for _, y := range t {
				results := gjson.GetMany(y.Raw, ensInfoMap[k].field...)
				var str []string
				for _, s := range results {
					str = append(str, s.String())
				}
				table.Append(str)
			}
			table.Render()

		}
	}

	// 查询对外投资详细信息
	// 对外投资>0 && 是否递归 && 参数投资信息大于0
	if ensInfoMap["invest"].total > 0 && isSearch && options.Invest {
		enInfo.investInfos = make(map[string]EnInfo)
		for _, t := range enInfo.infos["invest"] {
			openStatus := t.Get("openStatus").String()
			if openStatus == "注销" || openStatus == "吊销" {
				continue
			}
			investNum := 0.00
			if t.Get("regRate").String() == "-" {
				investNum = -1
			} else {
				str := strings.Replace(t.Get("regRate").String(), "%", "", -1)
				investNum, _ = strconv.ParseFloat(str, 2)
			}
			// 50%以上控股递归查询
			if investNum >= 50 {
				gologger.Infof("企业名称：%s 投资占比：%v\n", t.Get("entName"), investNum)
				options.CompanyID = t.Get("pid").String()
				options.KeyWord = ""
				GetEnInfoByPid(options,sn)
				//n := getCompanyInfoById(t.Get("pid").String(), true, options)
				//enInfo.investInfos[t.Get("pid").String()] = n
			}
		}
	}

	// 查询分支机构公司详细信息
	// 分支机构大于0 && 是否递归模式 && 参数是否开启查询
	if ensInfoMap["branch"].total > 0 && isSearch && options.IsGetBranch {
		enInfo.branchInfos = make(map[string]EnInfo)
		for _, t := range enInfo.infos["branch"] {
			if t.Get("openStatus").String() == "开业"{
				gologger.Infof("分支名称：%s 状态：%s\n", t.Get("entName"), t.Get("openStatus"))
				options.CompanyID = t.Get("pid").String()
				options.KeyWord = ""
				GetEnInfoByPid(options,sn)
				//n := getCompanyInfoById(t.Get("pid").String(), false, options)
				//enInfo.branchInfos[t.Get("pid").String()] = n
			}
		}
	}
	return enInfo
}

// getInfoList 获取信息列表
func getInfoList(pid string, types string, options common.Options) []gjson.Result {
	urls := "https://aiqicha.baidu.com/" + types + "?size=100&pid=" + pid
	content := common.GetReq(urls, options)
	var listData []gjson.Result
	if gjson.Get(string(content), "status").String() == "0" {
		data := gjson.Get(string(content), "data")
		//判断一个获取的特殊值
		if types == "relations/relationalMapAjax" {
			data = gjson.Get(string(content), "data.investRecordData")
		}
		//判断是否多页，遍历获取所有数据
		pageCount := data.Get("pageCount").Int()
		if pageCount > 1 {
			for i := 1; int(pageCount) >= i; i++ {
				gologger.Infof("当前：%s,%d\n", types, i)
				reqUrls := urls + "&p=" + strconv.Itoa(i)
				content = common.GetReq(reqUrls, options)
				listData = append(listData, gjson.Get(string(content), "data.list").Array()...)
			}
		} else {
			listData = data.Get("list").Array()
		}
	}
	return listData

}

// SearchName 根据企业名称搜索信息
func SearchName(options common.Options) ([]gjson.Result,common.Options) {
	//fmt.Println(options.KeyWord)
	name := options.KeyWord
	urls := "https://aiqicha.baidu.com/s?q=" + name + "&t=0"
	content := common.GetReq(urls, options)
	enList1 := pageParseJson(string(content))
	enList := enList1.Get("resultList").Array()

	if len(enList) == 0 {
		gologger.Errorf("没有查询到关键词 “%s”\n", name)
		return enList,options
	} else {
		gologger.Infof("关键词：“%s” 查询到 %d 个结果，默认选择第一个 \n", name, len(enList))
	}
	options.CompanyID = enList[0].Get("pid").String()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PID", "企业名称", "法人代表"})
	for _, v := range enList {
		table.Append([]string{v.Get("pid").String(), v.Get("titleName").String(), v.Get("titleLegal").String()})
	}
	table.Render()
	return enList,options
}
