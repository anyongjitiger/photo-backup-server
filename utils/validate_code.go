package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/anyongjitiger/photo-backup-server/db"
)

const PrefixPIN = "current pin"

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	str := sb.String()
	db := db.GetDb()
	db.Set(PrefixPIN, []byte(str))
	return str
}

func GetCurrentPIN() string {
	/* if pin,err := db.GetDb().Get(PrefixPIN); err != nil {
		log.Println(err)
		return string(pin)
	} */
	pin, _ := db.GetDb().Get(PrefixPIN)
	return string(pin)
}
