package main

//"io"
import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/tidusant/c3m-common/mycrypto"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/chadmin-repo/models"

	"github.com/spf13/viper"
	rpch "github.com/tidusant/chadmin-repo/cuahang"
	//	"os"
)

func posSync(userid, shopid, posname string) string {
	//sync product

	PData := struct {
		Employees []models.Employee
		Products  []models.Product
		Cats      []models.ProdCat
		PosClient models.PosClient
	}{}
	//get posclient
	pos := rpch.GetPosClient(userid, shopid, posname)
	if pos.UserId == userid && pos.ShopId == shopid {
		//get all not sync product
		prods := rpch.GetSyncProds(userid, shopid)
		PData.Products = prods
		//get all cat
		cats := rpch.GetSyncCats(userid, shopid)
		PData.Cats = cats
		//get all cat
		empls := rpch.GetSyncEmployees(userid, shopid)
		PData.Employees = empls
		if pos.IsSync {
			PData.PosClient = pos
		}
		info, _ := json.Marshal(PData)
		return string(info)
	}
	return ""
}

func posSyncUpdate(userid, shopid, params string) string {
	//sync product
	data := mycrypto.Base64Decode(params)
	PData := struct {
		Employees []string
		Products  []string
		Cats      []string
		PosClient models.PosClient
	}{}
	err := json.Unmarshal([]byte(data), &PData)
	if c3mcommon.CheckError("Fail to parse syncupdate", err) {
		//get posclient
		pos := rpch.GetPosClient(userid, shopid, PData.PosClient.Name)
		if pos.UserId == userid && pos.ShopId == shopid {
			if len(PData.Cats) > 0 {
				rpch.UpdateSyncCats(userid, shopid, PData.Cats)
			}
			if len(PData.Products) > 0 {
				rpch.UpdateSyncProds(userid, shopid, PData.Products)
			}
			if len(PData.Employees) > 0 {
				rpch.UpdateSyncEmployees(userid, shopid, PData.Employees)
			}
			if pos.IsSync == false {
				rpch.UpdateSyncPosClient(userid, shopid, pos.Name)
			}
		}
	}
	return ""
}
func posGetImageThumb(userid, shopid, imageid string) string {
	//sync product
	//get all not sync product
	imagepath := viper.GetString("config.imagefolder") + "/" + shopid + "/thumb_"
	if _, err := os.Stat(imagepath + imageid); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			return "image not found!"
		} else {
			// other error
			return "unknow error!"
		}
	}
	b, _ := ioutil.ReadFile(imagepath + imageid)
	imgBase64Str := base64.StdEncoding.EncodeToString(b)
	return string(imgBase64Str)
}
func posGetShopInfo(userid, shopid string) string {
	//sync product
	//get all not sync product
	shop := rpch.GetShopById(userid, shopid)
	info, _ := json.Marshal(shop)
	return string(info)
}
