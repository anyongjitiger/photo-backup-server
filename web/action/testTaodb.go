package action

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anyongjitiger/photo-backup-server/db"
	"github.com/anyongjitiger/photo-backup-server/log"
	mux "github.com/julienschmidt/httprouter"
)

func TestTaodb(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	log.Info("TestTaodb.")
	db := db.GetDb()
	for i := 0; i < 100000; i++ {
		db.Set(fmt.Sprintf("testtaodb%d", i), []byte(fmt.Sprintf("测试taodbvalue%d", i)))
		fmt.Println("set key=" + fmt.Sprintf("testtaodb%d", i))
	}

	for i := 0; i < 100000; i++ {
		value, _ := db.Get(fmt.Sprintf("testtaodb%d", i))
		fmt.Println("value:", string(value[:]))
	}
	fmt.Fprintf(w, "<h1>Hello, TestTaodb</h1>")
}



func TestAlbumList(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	log.Info("TestTaodb.")
	db := db.GetDb()

	list, _ := db.Iterator("album")
	fmt.Println(list)
	jsonData ,_:=json.Marshal(list)
	fmt.Fprintf(w, "%v",string(jsonData))
}

