package main

import (
	"flag"

	"github.com/spf13/viper"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	"github.com/tidusant/c3m-common/mycrypto"
	rpch "github.com/tidusant/chadmin-repo/cuahang"
	rpsex "github.com/tidusant/chadmin-repo/session"
	//"io"
	"net"
	"net/http"
	//	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var mytoken string

func init() {

}

func main() {
	var port int
	var debug bool

	//fmt.Println(mycrypto.Encode("abc,efc", 5))
	flag.IntVar(&port, "port", 8084, "help message for flagname")
	flag.BoolVar(&debug, "debug", false, "Indicates if debug messages should be printed in log files")
	flag.StringVar(&mytoken, "token", ".fa1Xldsbe@", "Indicates if debug messages should be printed in log files")
	flag.Parse()

	//logLevel := log.DebugLevel
	if !debug {
		//logLevel = log.InfoLevel
		gin.SetMode(gin.ReleaseMode)
	}
	a := mycrypto.Encode4("abcdefgh")

	log.Printf("test decode 4:%s - %s", a, mycrypto.Decode4(a))
	// log.SetOutputFile(fmt.Sprintf("portal-"+strconv.Itoa(port)), logLevel)
	// defer log.CloseOutputFile()
	// log.RedirectStdOut()

	log.Infof("running with port:" + strconv.Itoa(port))

	//init config

	router := gin.Default()

	router.POST("/:param", func(c *gin.Context) {
		log.Debugf("header:%v", c.Request.Header)
		log.Debugf("Request:%v", c.Request)
		requestDomain := c.Request.Header.Get("Host")
		//allowDomain := c3mcommon.CheckDomain(requestDomain)
		param := c.Param("param")
		param = mycrypto.Decode3(param)
		strrt := ""
		c.Header("Access-Control-Allow-Origin", "*")
		if param != "" {
			//c.Header("Access-Control-Allow-Origin", allowDomain)
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,access-control-allow-credentials")
			c.Header("Access-Control-Allow-Credentials", "true")
			log.Debugf("check request:%s", c.Request.URL.Path)
			if rpsex.CheckRequest(c.Request.URL.Path, c.Request.UserAgent(), c.Request.Referer(), c.Request.RemoteAddr, "POST") {
				strrt = myRoute(c, param)

				strrt = mycrypto.Encode4(strrt)
				//log.Debugf("customer %s", strrt)
			} else {
				log.Debugf("check request error")
			}
		} else {
			log.Debugf("Not allow " + requestDomain)
		}
		if strrt == "" {
			strrt = c3mcommon.Fake64()
		}
		c.String(http.StatusOK, strrt)

	})

	router.Run(":" + strconv.Itoa(port))

}

func myRoute(c *gin.Context, param string) string {
	strrt := ""
	userIP, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	args := strings.Split(param, "|")
	if len(args) < 2 {
		return ""
	}
	RPCname := args[0]
	session := args[1]
	//function for not auth:
	if RPCname == "updateorderstatus" {
		if viper.GetString("config.whookmantoken") == session && len(args) >= 6 {
			partnercode := args[2]
			shopid := args[3]
			statusid := args[4]
			LabelID := args[5]
			status := rpch.GetStatusByPartnerStatus(partnercode, shopid, statusid)
			if status.ID.Hex() != "" {
				order := rpch.GetOrderByShipmentCode(LabelID, shopid)
				if order.ID.Hex() != "" {

					rpch.UpdateOrderStatusByShipmentCode(LabelID, status.ID.Hex(), shopid)
				}
			}
		}

	} else {
		//check login
		sessioninfo := rpch.GetLogin(session, userIP)
		if sessioninfo == "" {
			return ""
		}
		tmps := strings.Split(sessioninfo, "[+]")
		userid := tmps[0]
		shopid := tmps[1]
		log.Debugf("RPCname: %s", RPCname)
		if RPCname == "aut" {
			//check login
			log.Debugf("customer %s", sessioninfo)
			strrt = sessioninfo
		} else if RPCname == "cusexport" {
			cuss := rpch.GetAllCustomers(shopid)
			strphone := ""
			for _, v := range cuss {
				strphone += v.Phone + ","

			}
			return strphone

		} else if RPCname == "loadshopalbum" {
			shop := rpch.GetShopById(userid, shopid)
			strrt := "{\"\":\"\""

			for _, album := range shop.Albums {
				strrt += ",\"" + album.Slug + "\":\"" + album.Name + "\""
			}
			strrt += "}"
			return strrt
		}
	}
	return strrt
}
