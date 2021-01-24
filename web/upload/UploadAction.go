package upload

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	fileName := r.PostFormValue("fileName")
	log.Info("FileName: %s", fileName)
	file, _, err := r.FormFile("uploadFile")
	if err == nil {
		_, err := os.Stat(albumPath + deviceName)
		if err != nil {
			log.Error("Read dir error:%v", err)
			err = os.MkdirAll(albumPath+deviceName, 0765)
			if err != nil {
				log.Error("Mkdir error:%v", err)
			} else {
				log.Info("Mkdir success:%s", albumPath+deviceName)
			}
		}
	}
	// 读文件错误
	if err != nil {
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
	defer file.Close()

	// 获取文件的扩展名
	extName := utils.GetFileExt(fileName)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("Read file error : %v", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	// 生成文件sha256码
	// sha256Value := utils.GetByteSha256(data)
	sha256Value := utils.GetTxtSha256(fileName + fileSize)
	// log.Info("sha256:%s", sha256Value)

	// 获取DB中是否已经保存该文件
	temp := model.Resource{}
	temp.NameSha256 = sha256Value
	temp.Get()
	if temp.FileName != "" {
		log.Error(" 文件已经存在，文件名=%s\n", temp.FileName)
		return
	}

	tempFile := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	// tempFileName := tempFile + "." + extName

	// previewFileName := tempFile + "_" + utils.PhotoPreviewSizeStr + "." + extName
	tempStoreFile := albumPath + deviceName + "/" + fileName
	if exist, _ := PathExists(tempStoreFile); exist {
		tempStoreFile = albumPath + deviceName + "/" + strings.TrimSuffix(fileName, extName) + tempFile + "."  + extName
	}
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
	res.FileName = fileName
	res.FileSize = fileInfo.Size()
	res.NameSha256 = sha256Value
	res.FileType = utils.GetFileType(extName)
	
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

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
	return false, err
}
