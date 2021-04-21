package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var(
	targets[] string
	filename string
	outputfile string
	help bool
)

const (
	URL   = "https://api.ip138.com/ip/"
	TOKEN = "cf8efae9550c25033e085bb89b0bac69"
)

//----------------------------------
// iP地址调用示例代码
//----------------------------------

// xml struct
type xmlinfo struct {
	Ret  string          `xml:"ret"`
	Ip   string          `xml:"ip"`
	Data locationxmlInfo `xml:"data"`
}

type locationxmlInfo struct {
	Country string `xml:"country"`
	Region  string `xml:"region"`
	City    string `xml:"city"`
	Isp     string `xml:"isp"`
	Zip     string `xml:"zip"`
	Zone    string `xml:zone`
}

//json struct
type jsoninfo struct {
	Ret  string    `json:"ret"`
	Ip   string    `json:"ip"`
	Data [6]string `json:"data"`
}

func init(){
	flag.BoolVar(&help,"h, --help",false,"help, 帮助命令")
	flag.StringVar(&filename,"i","","要读取的文件")
	flag.StringVar(&outputfile,"o","","要输出的文件")
	flag.Usage = usage
	flag.Parse()

}

func usage(){
	flag.PrintDefaults()
}


func main() {
	//ipLocation("8.8.8.8","json")
	var datas []jsoninfo
	if filename != "" && outputfile != ""{
		targetfile, err := os.OpenFile(filename,os.O_RDONLY,1)
		outfile, err := os.OpenFile(outputfile,os.O_WRONLY,2)
		defer targetfile.Close()
		if err != nil{
			log.Println(err)
		}

		tmps := bufio.NewScanner(targetfile)
		for tmps.Scan(){
			tmp := tmps.Text()
			if strings.Contains(tmp,"https://") || strings.Contains(tmp,"http://"){
				tmp2 := strings.Split(tmp,"//")
				tmpip := tmp2[1]
				if strings.Contains(tmpip,":") {
					tmpip2 := strings.Split(tmpip,":")[0]
					targets = append(targets, tmpip2)
				}else {
					targets = append(targets, tmpip)
				}
			}else {
				tmp2 := strings.Split(tmp,":")
				tmpip := tmp2[0]
				targets = append(targets, tmpip)
			}
		}
		for _,ip := range targets{

			data := ipLocation(ip,"jsonp")
			datas = append(datas, data)
		}
		for _,v := range datas{
			outfile.WriteString(v.Ip + " " +v.Ret)
			for _,v1 := range v.Data{
				outfile.WriteString(" "+v1)
			}
			outfile.WriteString("\n")
		}
	}else {
		flag.Usage()
	}

}

func ipLocation(ip string,dataType string) jsoninfo{
	var info jsoninfo
	queryUrl := fmt.Sprintf("%s?ip=%s&datatype=%s",URL,ip,dataType)
	client := &http.Client{}
	reqest, err := http.NewRequest("GET",queryUrl,nil)

	if err != nil {
		fmt.Println("Fatal error ",err.Error())
	}

	reqest.Header.Add("token",TOKEN)
	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Fatal error ",err.Error())
	}
	if response.StatusCode == 200 {
		bodyByte, _ := ioutil.ReadAll(response.Body)

		if dataType == "jsonp" {

			json.Unmarshal(bodyByte,&info)
			fmt.Println(info.Ip)
			return info
		} else if dataType == "xml" {
			var info xmlinfo
			xml.Unmarshal(bodyByte,&info)
			fmt.Println(info.Ip)
		}
	}

	return info
}