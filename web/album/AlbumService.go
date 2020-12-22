package album

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/anyongjitiger/photo-backup-server/db/model"
)

func mapsToList(maps map[string]string) (model.ResourceSort, error){

	if maps!=nil && len(maps)>0 {
		list := make(model.ResourceSort,0)
		for _,value := range maps {
			if value !="" {
				res :=model.Resource{}
				json.Unmarshal([]byte(value),&res)
				list = append(list,res)
			}
		}
		sort.Sort(list)
		return list,nil
	}

	return nil, errors.New(" No data")
}