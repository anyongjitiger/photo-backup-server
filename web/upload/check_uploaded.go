package upload

import (
	"encoding/json"
	"net/http"

	"github.com/anyongjitiger/photo-backup-server/db/model"
	"github.com/anyongjitiger/photo-backup-server/log"
	"github.com/anyongjitiger/photo-backup-server/utils"
	"github.com/anyongjitiger/photo-backup-server/web/core/kit"
	"github.com/anyongjitiger/photo-backup-server/web/core/render"
)

func CheckUploaded(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	files := r.PostFormValue("files")
	var rtFiles []FileInfo;
	var _files []FileInfo
	json.Unmarshal([]byte(files), &_files)
	for _, r := range _files {
		log.Info(r.FileName)
		log.Info(r.FileSize)
		sha256Value := utils.GetTxtSha256(r.FileName + r.FileSize)
		temp := model.Resource{}
		temp.NameSha256 = sha256Value
		temp.Get()
		log.Info(temp.FileName)
		if temp.FileName != "" {
			continue
		}else{
			rtFiles = append(rtFiles, r)
		}
	}
	obj, _ := json.Marshal(rtFiles)
	log.Info(string(obj))
	ret := kit.GetCommonRet()
	if rtFiles != nil {
		ret.Data = string(obj)
		ret.Msg = "success"
		ret.State =  kit.RetStateOk
	}else {
		ret.Msg = "no data"
		ret.State =  kit.RetStateOk
	}
	render.RenderJson(w, ret)
}