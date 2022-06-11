# 参数说明
对ENScanPublic.exe进行二开(只保留爱企查接口),能对开业的分支公司、投资占比50%以上的公司进行递归查询!并对Domain和Title进行聚合整理!

***author  :i11us0ry***
***version :0.0.1***

Usage of AScan.exe:
  -b    bool    是否递归查询开业状态的分支公司
  -f     string  包含公司名称的文件，公司名按行存储
  -n    string  公司名称,最好是爱企查上公司对应的`完整名称`
  -s    bool     是否递归查询对外投资50%以上的开业公司
  
# 使用说明
  首次使用先运行AScan.exe，会自动生成配置文件和保存目录。生成完配置文件后手动添加cookie，程序提示cookie过期或提示需要更新cookie时也需要更新cookie
  
  ![image](https://github.com/i11us0ry/AScan/blob/main/img/Pasted%20image%2020220401194624.png)
  
  运行时根据需求输入参数即可
  
  ![[Pasted image 20220401195416.png]]

# 保存结果
AScan会将查询结果分为两部分。
第一部分以查询公司、分支公司及对外投资公司为单位的xlsx表格。表格内容包括基本信息、网站备案、分支机构、微信公众号、对外投资、软件著作权信息(APP)，方便红队和攻防中做备用计划
  
![[Pasted image 20220401195020.png]]

第二部分为查询公司、分支公司、备案信息及对外投资50%以上公司的域名信息整合，方便快速扫描或做下一步信息收集

![[Pasted image 20220401195229.png]]
