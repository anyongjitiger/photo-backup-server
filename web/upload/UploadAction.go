package upload

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anyongjitiger/photo-backup-server/config"
	"github.com/anyongjitiger/photo-backup-server/db/model"
	"github.com/anyongjitiger/photo-backup-server/log"
	"github.com/anyongjitiger/photo-backup-server/utils"
	"github.com/anyongjitiger/photo-backup-server/web/common"
	"github.com/anyongjitiger/photo-backup-server/web/core/kit"
	"github.com/anyongjitiger/photo-backup-server/web/core/render"
	"github.com/julienschmidt/httprouter"
)

const baseFormat = "2006-01-02 15:04:05"

var albumPath = ""

type Controller struct {
}

func (Controller) Upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	albumPath = config.PFile.AlbumPath + "/"
	deviceName := r.PostFormValue("device")
	fileSize := r.PostFormValue("fileSize")
	file, fHead, err := r.FormFile("uploadFile")
	// 读文件错误
	if err != nil {
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
	defer file.Close()

	// 获取文件的扩展名
	extName := utils.GetFileExt(fHead.Filename)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("Read file error : %v", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	// 生成文件sha256码
	// sha256Value := utils.GetByteSha256(data)
	sha256Value := utils.GetTxtSha256(fHead.Filename + fileSize)
	// log.Info("sha256:%s", sha256Value)

	// 获取DB中是否已经保存该文件
	temp := model.Resource{}
	temp.NameSha256 = sha256Value
	temp.Get()
	if temp.FileName != "" {
		log.Error(" 文件已经存在，文件名=%s\n", temp.FileName)
		return
	}

	fileName := fHead.Filename
	log.Info("originFileName: %s", fileName)

	tempFile := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	tempFileName := tempFile + "." + extName

	// previewFileName := tempFile + "_" + utils.PhotoPreviewSizeStr + "." + extName
	tempStoreFile := albumPath + tempFile + "." + extName
	err = ioutil.WriteFile(tempStoreFile, data, 0664)

	if err != nil {
		log.Error("Write file error: %v", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
	fileInfo, err := os.Stat(tempStoreFile)
	if err != nil {
		log.Error(" get fileInfo error: %v", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
	res := model.Resource{}
	// cTime, err := utils.Photo{}.GetDate(tempStoreFile)
	res.FileName = tempFileName
	res.FileSize = fileInfo.Size()
	res.NameSha256 = sha256Value
	res.FileType = utils.GetFileType(extName)
	if err == nil {
		// res.CTime = strconv.FormatInt(cTime.Unix(), 10)
		// res.FilePath = deviceName + "/" +utils.GetDateYYYYMM(cTime) 
		res.FilePath = deviceName
		netFileName := albumPath + res.FilePath + "/" + tempFileName
		_, err := os.Stat(albumPath + res.FilePath)
		if err != nil {
			log.Error("Read dir error:%v", err)
			err = os.MkdirAll(albumPath+res.FilePath, 0765)
			//err = os.MkdirAll(albumPath+res.FilePath, 0664)
			if err != nil {
				log.Error("Mkdir error:%v", err)
			} else {
				log.Info("Mkdir success:%s", albumPath+res.FilePath)
			}
		}
		err = os.Rename(tempStoreFile, netFileName)
		if err != nil {
			log.Error("rename error:%v", err)
		}
		// previewFileName
		/* if res.FileType == "video/mp4" || res.FileType == "video/mov" {
			res.Preview = ""
		}else{
			err = utils.Photo{}.CreatePreviewImg(netFileName, albumPath+res.FilePath+"/"+previewFileName)
			if err != nil {
				log.Error("rename error:%v", err)
			} else {
				res.Preview = previewFileName
			}
		} */
	}
	log.Info("save res.")
	// save to taodb
	res.Save()
	ret := kit.GetCommonRet()
	ret.State = kit.RetStateOk
	bean := Bean{}
	bean.FileName = fileName
	bean.State = 1
	ret.Data = bean
	render.RenderJson(w, ret)

}
