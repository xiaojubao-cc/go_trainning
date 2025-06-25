package main

import (
	"fmt"
	"regexp"
	"strings"
)

/*基本数据类型赋值：int、float、bool、string、数组、结构体属于值复制*/
func main() {
	///*本质上是一个字节数组*/
	//var str string = "hello world"
	//for _, value := range str {
	//	fmt.Printf("value=%c\n", value)
	//}
	//str1 := "石头人"
	//for _, value := range str1 {
	//	fmt.Printf("value=%c\n", value)
	//}
	//for i := 0; i < len([]rune(str1)); i++ {
	//	fmt.Printf("value=%c\n", []rune(str1)[i])
	//}
	///*字符串的拼接使用*/
	//builder := strings.Builder{}
	//builder.WriteString("hello")
	//builder.WriteString("world")
	//fmt.Println(builder.String())
	logLine := `2025/06/19 14:53:00.844781 [C7E90D00E34E6E6C] --- DEBUG[redis.EXEC:Retn=0000,Desc=OK,REQ=string=expireat FQ000301202506191453008390982201 1750316000; ANSWER:table={"Result":1,"RowNum":1,"Desc":"OK","Retn":"0000"}; EXEC: 0.0ms][adapterMid.lua:1167]`
	// 提取 redis.EXEC 段落 + 字段
	re := regexp.MustCompile(`\[redis\.EXEC:Retn=(\d+),Desc=([^,]+),REQ=string=([^;]+);\s*ANSWER:table=({[^}]+});\s*EXEC:\s*(\d+\.?\d*ms)\]`)
	match := re.FindStringSubmatch(logLine)
	if len(match) > 0 {
		split := strings.Split(strings.ReplaceAll(strings.ReplaceAll(match[0], "[", ""), "]", ""), ";")
		if len(split) == 3 {
			i := strings.Split(split[0], ",")
			re := regexp.MustCompile(`Retn=(\d+)`)
			match := re.FindString(i[0])
			fmt.Println(match)
		}
	} else {
		fmt.Println("未找到匹配项")
	}

	logLine = `2025/06/19 14:53:00.846988 [C7E90D00E34E6E6C] --- DEBUG[_DATASVR.SqlRun:Retn=0000,Desc=2:insert 1 rows(3840818249).31,37,99,save asyn OK.,SQL=insert into mcn_mid_userData301(unumber1,midnumber,unumber2,cityid,areacode,requestid,datastate,exptime,record,international,midgroupid,bindtime,bindid,gnflag,userData,audio,audioanswer,xappbizid,lastgnumber,authorisedstatus) values ('057156068911','18461966141','15069148890',936,'936','0a4ff6bf4cda11f0897e7e10b3528a1c','0','20250619145320','1','1','FQ0003','20250619145300','01202506191453008390982201','00','','||','|||','','','1'); EXEC: 2ms][adapterMid.lua:1208]`
	// 提取 redis.EXEC 段落 + 字段
	re = regexp.MustCompile(`\[_DATASVR\.(\w+):Retn=(\d+),Desc=([^,]+?\.\d+,\d+,\d+,[^,]+?),SQL=([^;]+); EXEC: (\d+\.?\d*ms)\]`)
	match = re.FindStringSubmatch(logLine)
	if len(match) > 0 {
		split := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(match[0], "[", ""), "]", ""), "EXEC:", ""), ";")
		for i := range split {
			fmt.Println(strings.TrimSpace(split[i]))
		}
	} else {
		fmt.Println("未找到匹配项")
	}
}
