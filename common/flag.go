package common

import (
	"flag"
	"fmt"
)

const banner = `对ENScanPublic.exe进行二开,能对开业的分支公司、投资占比50%以上的公司进行递归查询!并对Domain和Title进行聚合整理!`
const author = `i11us0ry`
const Version = `0.0.1`

func Banner() {
	fmt.Println(fmt.Sprintf("\n%s\nauthor	:%s\nversion	:%s\n", banner, author, Version))
}

func Flag(Info *Options) {
	Banner()
	flag.StringVar(&Info.KeyWord, "n", "", "公司名称,最好是爱企查上公司对应的完整名称")
	flag.StringVar(&Info.InputFile, "f", "", "包含公司名称的文件，公司名按行存储")
	flag.BoolVar(&Info.IsGetBranch, "b", false, "是否递归查询开业状态的分支公司")
	flag.BoolVar(&Info.Invest, "s", false, "是否递归查询对外投资50%以上的开业公司")
	flag.Parse()
}
