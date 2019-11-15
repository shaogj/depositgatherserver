package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mkideal/log"
	"io/ioutil"

	//"strconv"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/go-xorm/xorm"
	"github.com/robfig/config"
)
var   (
	HConf ConfigTools

)
var GbConf ConfigInfomation = ConfigInfomation{}
const (
	CoinBitcoin string ="BTC"
	CoinEthereum string ="ETH"

	//sgj 1015 add
	CoinDASH string="DASH"
	CoinWDC string ="WDC"
)
type DSCConf struct{
	//RPC for DSC:
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string
	RPCTestNet		int

}

//10.30,目前预留
type WDCNodeConf struct{
	//RPC for WDC:
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string

}

type  MongoDatabaseInfo struct {
	Db string   	`josn:"dbase"`  //数据库
	Coll string   	`josn:"coll"` //数据集
	User string		`josn:"user"` //用户名
	Pass  string	`josn:"pass"` //密码
}
type MySqlConfig struct  {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Dbname string `json:"dbname"`
}

type SettleApiReq struct {
	SettlApiQuery string `json:"settleapiquery"`
	SettlApiUpdate string `json:"settleapiupdate"`
}

type SettleAccessKey struct {
	AccessComePubKey string `json:"AccessPubKey"`
	AccessPrivKey string `json:"AccessPrivKey"`
}
//sgj 1019 add
type ConfigInfomation struct {
	MySqlCfg        MySqlConfig         `json:"MySqlConfig"`
	SettleApiReq    SettleApiReq         `json:"SettleApiReq"`
	SettleApiQuery string `json:"SettleApiQuery"`
	SettleApiUpdate string `json:"SettleApiUpdate"`
	//1112 add
	SettleApiDepositQuery string `json:"SettleApiDepositQuery"`
	//1019 add
	WebPort			int `json:"WebPort"`
	JavaSDKUrl    string         `json:"JavaSDKUrl"`
	WDCTransUrl    string         `json:"WDCTransUrl"`
	WDCNodeUrl    string         `json:"WDCNodeUrl"`

	WDCConf			WDCNodeConf		`json:"SettleApiReq"`
	SettleAccessKey    SettleAccessKey         `json:"SettleAccessKey"`
	//WDC提现的大账户地址，可选
	WDCTransferOutAddress    string         `json:"WDCTransferOutAddress"`
	//WDC归集的大账户地址，可选
	WDCGatterToAddress    string         `json:"WDCGatterToAddress"`
	//WDC归集的获取配置接口，可选
	WDCGatterConfigUrl    string         `json:"WDCGatterConfigUrl"`

}

func InitWithProviders(providers, dir string) error {
	return log.Init(providers, log.M{
		"rootdir":     dir,
		"suffix":      ".txt",
		"date_format": "%04d-%02d-%02d",
	})
}

func defaultInt(ptr *int, dft int) {
	if *ptr == 0 {
		*ptr = dft
	}
}

func defaultString(ptr *string, dft string) {
	if len(*ptr) == 0 {
		*ptr = dft
	}
}
type ConfigTools struct {
	//sjg 0924 add for BCH RPC
	CurDSCConf DSCConf

	// log
	LogProviders string
	LogLevel     string
	Logpath      string
	MgoAddrData		MongoDatabaseInfo   /*用户addr数据*/
	
}

func NewConfigTools(configpath string) error {
	//ct := ConfigTools{}
	HConf= ConfigTools{}
	ct:=&HConf
	fmt.Printf("dsfd---tsint")
	// 初始配置获取，获取失败则直接抛异常
	cfg, err := config.ReadDefault(configpath)
	if err != nil {
		panic(fmt.Sprintf("config.ReadDefault error:%v", err))
	}
	sectionName := "log"
	ct.LogProviders, _ = cfg.String(sectionName, "log.providers")
	defaultString(&ct.LogProviders, "multifile/console")
	ct.Logpath, _ = cfg.String(sectionName, "log.path")
	ct.LogLevel, _ = cfg.String(sectionName, "log.level")
	if err := InitWithProviders(ct.LogProviders, ct.Logpath); err != nil {
		panic("init log error: " + err.Error())
	}
	log.Info("log level: %v", log.SetLevelFromString(ct.LogLevel))
	sectionName = "DSCRPC"
	// 当前BTC's RPC服务端口，job供监听是否正常运行也有个对外web
	ct.CurDSCConf.RPCPort, _ = cfg.Int(sectionName, "port")
	defaultInt(&ct.CurDSCConf.RPCPort, 9332)
	//add 
	ct.CurDSCConf.RPCHostPort, _ = cfg.String(sectionName, "rpchostport")
	defaultString(&ct.CurDSCConf.RPCHostPort, "127.0.0.1:9332")
	log.Info("----------------DSC-PRC RPCHostPort is:---------%d------",ct.CurDSCConf.RPCHostPort)
	//
	ct.CurDSCConf.RPCUser, _ = cfg.String(sectionName, "rpcuser")
	defaultString(&ct.CurDSCConf.RPCUser, "shaogj")
	log.Info("---------------DSC--PRC RPCUser is:---------%s------",ct.CurDSCConf.RPCUser)
	
	ct.CurDSCConf.RPCPassWord, _ = cfg.String(sectionName, "rpcpassword")
	defaultString(&ct.CurDSCConf.RPCPassWord, "123456")

	//sgj 0522 add,是否testnet开关
	ct.CurDSCConf.RPCTestNet, _ = cfg.Int(sectionName, "testnet")
	defaultInt(&ct.CurDSCConf.RPCTestNet, 0)
		
	return nil
}

//sgj 1019 add
func InitConfigInfo() error {
	//*good conf:
	//log.SetFlags(log.Lshortfile | log.Ltime)
	var strConf string
	flag.StringVar(&strConf, "conf", "config.json", "config <file>")
	flag.Parse()
	byData, err := ioutil.ReadFile(strConf)
	if nil != err {
		log.Error("Read config file :::%v", err)
		return err
	}
	err = json.Unmarshal(byData, &GbConf)
	if nil != err {
		log.Error("Unmarshal config file :::%v", err)
		return err
	}
	log.Info("ConfigInfo:::%+v", GbConf)
	return nil
}
