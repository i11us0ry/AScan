package common

import (
	"AScan/common/utils/gologger"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Options struct {
	KeyWord     string // Keyword of Search
	CompanyID   string // Company ID
	InputFile   string // Scan Input File
	Output      string
	CookieInfo  string
	ScanType    string
	IsGetBranch bool
	Invest      bool
	GetFlags    string
	Version     bool
}

type Sun struct {
	Name 		[]string
}

type Conf struct {
	Version 	string 			`yaml:"version"`
	Output 		string 			`yaml:"output"`
	Cookies 	string 			`yaml:"cookies"`
}

type WebInfo struct {
	FileName string
	Info  []Info
}

type Info struct{
	Domain	string
	Title 	string
}

type Web interface {
	Check(w, t string)
	SetFileName(fn string)
}

func CheckConf() *Conf{
	config := GetConfigDir() + "\\conf.yml"
	_, exist := os.Stat(config)
	// 文件不存在
	if os.IsNotExist(exist) {
		return writeConf(config)
		gologger.Printf("已自动生成配置文件 conf.yml")
	}
	return readConf()
}

func writeConf(fileName string) *Conf{
	conf := &Conf{}
	conf.Version = Version
	conf.Output = GetConfigDir() + "\\result"
	conf.Cookies = "cookie"
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer file.Close()
	enc := yaml.NewEncoder(file)
	err := enc.Encode(conf)
	if err != nil {
		gologger.Labelf("%v" ,err)
		os.Exit(1)
	}
	return conf
}

func readConf() *Conf{
	conf := &Conf{}
	yamlFile, err := ioutil.ReadFile(GetConfigDir() + "\\conf.yml")
	if err != nil {
		gologger.Labelf("conf.yml read err!\n")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		fmt.Println()
		gologger.Labelf("conf.yml read err!\n")
		os.Exit(1)
	}
	return conf
}

func GetConfigDir() string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		gologger.Labelf("%v\n",err)
		os.Exit(1)
	}
	//strings.Replace(dir, "\\", "/", -1)
	return fmt.Sprintf("%v",dir)
}

func (info *WebInfo)SetFileName (fn string){
	info.FileName = fn
}

func (info *WebInfo)Check (d, t string){
	flag := false
	for _, v := range(info.Info){
		if v.Domain == d {
			flag = true
			break
		}
	}
	if !flag {
		info.Info = append(info.Info,Info{d,t})
		writeFile(info.FileName,d,t)
	}
}

// 写入domain_title
func writeFile(fileName string, d, t string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0666) //打开文件
	defer f.Close()
	if err != nil {
		gologger.Labelf("%v 打开失败!",fileName)
		return
	}
	// 将文件写进去
	if _, err = io.WriteString(f, fmt.Sprintf("Domain:%-60v Title:%v\n",d,t)); err != nil {
		gologger.Labelf("%v 写入失败!",fileName)
		return
	}
}

func WriteSun(fileName string,sn []string){
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0666) //打开文件
	defer f.Close()
	if err != nil {
		gologger.Labelf("%v 打开失败! %v",fileName,err)
		return
	}
	// 将文件写进去
	if _, err = io.WriteString(f, fmt.Sprintf("\n%s\n","子公司：")); err != nil {
		gologger.Labelf("%v 写入失败! %v",fileName,err)
		return
	}
	for _,v := range(sn){
		if _, err = io.WriteString(f, fmt.Sprintf("%s\n",v)); err != nil {
			gologger.Labelf("%v 写入失败! %v",fileName,err)
			return
		}
	}
}