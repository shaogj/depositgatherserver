package main

import (
	"2019NNZXProj10/depositgatherserver/service"
	//"encoding/base64"

	//"backend/support/libraries/loggers"
	"flag"
	"fmt"
	"net/http"
	//"time"

	//"2019NNZXProj10/depositgatherserver"
	"os"

	//sgj 1217add

	"2019NNZXProj10/depositgatherserver/config"

	"github.com/gorilla/mux"
	mylog "github.com/mkideal/log"

	//"2019NNZXProj10/depositgatherserver/worker"
	//"2019NNZXProj10/depositgatherserver/cryptoutil"

	//0615doing

	//"io"
	"path/filepath"
	//"strings"
)

var flConfig string

func init() {
	//flag.StringVar(&flConfig, "c", "./config.conf", "config filepath")
	//	flag.StringVar(&flConfig, "cmy", "/Users/gejians/go/src/2019NNZXProj10/depositgatherserver/config.conf", "config filepath")
	flag.StringVar(&flConfig, "cmy", "config.conf", "config filepath")
}

//sgj 1105 add for settlecenter:

///113testing:
//Encrypt

func main() {
	err := config.NewConfigTools(flConfig)
	if nil != err {
		mylog.Error("Can't load config error=%v", err)
		os.Exit(0)
	}
	cfgtools := &config.HConf

	err = config.InitConfigInfo()
	if nil != err {
		mylog.Error("from config.json,get json conf err!")
		os.Exit(0)
	}
	gbConf := &config.GbConf
	mylog.Info("--sgj==>get conf info is %v", cfgtools.CurDSCConf)

	err = service.InitMysqlDB(gbConf.MySqlCfg)
	if nil != err {
		mylog.Error("cur InitMysqlDB() to conn err!,err is :%v", err)
		os.Exit(0)

	}

	//1103add
	service.GSettleAccessKey = gbConf.SettleAccessKey.AccessPrivKey

	//一个默认的合作商Key
	getKeystr := []byte("1234567812345678")
	/**/
	//2)

	//12.17doing,,for encry to db:
	//curPrivkey :="L5DoNfVEtdwEuPkmTYQT11p7dLsmnsnMpKLsD4mZbn3ozEizdv37"
	curPrivkey := "KznoaLNGcSzJUWkBk7FLXFRbNBqkL21SVGn6CMZrEUE2qJRX3SFf"
	//sgj20200612doingTmp
	curPrivkey = "b463b84f5af1243661267441f0e8daa6a57dfde5b9059ad14e8ddc2a6617c4e0"

	//sgj 0615doing for WGCaddr:
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	//file, err := os.Open(exPath+"/"+filename)
	fmt.Println("exPath is :%s, getKeystr is :%s,curPrivkey is:%s", exPath, getKeystr, curPrivkey)

	router := mux.NewRouter().StrictSlash(true)

	//sgj 1017 adding:	签名服务器: 成生新的地址(批量)
	router.HandleFunc("/remote/getnewaddress", service.RemoteSignCreateAddress)

	//sgj 1114add for guiji:

	//sgj 1217cadd for guiji:
	//router.HandleFunc("/remote/monitorwalletwdc", service.RemoteMonitorWalletAddress)
	//sgj 1218doing

	strHost := fmt.Sprintf(":%d", gbConf.WebPort)
	mylog.Info("strHost is :%s", strHost)

	err = http.ListenAndServe(strHost, router)
	if nil != err {
		mylog.Error("%+v", err)
		os.Exit(0)
	}
	//select {}

}
