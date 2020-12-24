package taoweb

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/anyongjitiger/photo-backup-server/config"
	"github.com/anyongjitiger/photo-backup-server/db"
	"github.com/anyongjitiger/photo-backup-server/log"
)

var flags struct {
	addr, albumPath string
	dbAddr          string
	logto           string
	loglevel        string
}

func init() {
	fmt.Println(runtime.GOOS)
	flag.StringVar(&flags.addr, "addr", ":8000", "The TCP address to bind to")
	flag.StringVar(&flags.dbAddr, "dbAddr", "127.0.0.1:7398", "The TCP address to connect to taodb")
	flag.StringVar(&flags.logto, "log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	flag.StringVar(&flags.loglevel, "log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
	if runtime.GOOS == "windows" {
		flag.StringVar(&flags.albumPath, "albumPath", "D:/Album", "album save path")
	}else{
		flag.StringVar(&flags.albumPath, "albumPath", "/data/album", "album save path")
	}
}

func Main() {
	flag.Parse()
	proFile := config.Profile{}
	proFile.AlbumPath = flags.albumPath
	config.PFile = proFile
	log.LogTo(flags.logto, flags.loglevel)

	router := NewRouter()

	log.Info("start taodb ...")
	_, err := db.New(flags.dbAddr)

	if err != nil {
		log.Error("start taodb fail:", err)
		panic(err)
	}
	log.Info("start taodb success")

	if err := http.ListenAndServe(flags.addr, router); err != nil {
		log.Error("start fail:", err)
		panic(err)
	}
}
