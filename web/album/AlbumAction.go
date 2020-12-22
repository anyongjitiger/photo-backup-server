package album

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anyongjitiger/photo-backup-server/config"
	"github.com/anyongjitiger/photo-backup-server/db"
	"github.com/anyongjitiger/photo-backup-server/web/common"
	"github.com/anyongjitiger/photo-backup-server/web/core/kit"
	"github.com/anyongjitiger/photo-backup-server/web/core/render"
	mux "github.com/julienschmidt/httprouter"
)

/* func List(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("/albums")
	prePath :=ps.ByName("prePath")
	if "all" == prePath {
		prePath = ""
	}
	db := db.GetDb()
	dataMap, _ := db.Iterator("album:"+prePath)
	ret := kit.GetCommonRet()
	if dataMap != nil{
		list ,_:= mapsToList(dataMap)
		fmt.Println(list)
		ret.Data = list
		ret.State =  kit.RetStateOk
	}
	render.RenderJson(w, ret)
} */

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("/albums")
	if ps, ok := common.ParamsFromContext(r.Context()); ok {
		prePath := ps.ByName("prePath")
		if "all" == prePath {
			prePath = ""
		}
		db := db.GetDb()
		dataMap, _ := db.Iterator("album:"+prePath)
		ret := kit.GetCommonRet()
		if dataMap != nil{
			list ,_:= mapsToList(dataMap)
			fmt.Println(list)
			ret.Data = list
			ret.State =  kit.RetStateOk
		}
		render.RenderJson(w, ret)
	}
}

func Show(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	fileName := ps.ByName("fileName")
	filePath := ps.ByName("filePath")
	fmt.Println("fileName:"+fileName)
	if fileName !=""{
		albumPath := config.PFile.AlbumPath+"/"+filePath+"/"
		data, err := os.Open(albumPath+fileName)
		if err !=nil{
			log.Printf("Read file error : %v", err)
			common.SendErrorResponse(w, http.StatusNotFound, "Not found file.")
			return
		}
		if data!=nil{
			w.Header().Set("Content-Type", "image/jpeg")
			fmt.Println("fileName:"+fileName)
			http.ServeContent(w, r, "", time.Now().UTC(), data)
		}
	}

}
