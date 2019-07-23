package main

import (
	"net/http"
	"octopus/myredis"
	"strconv"
)

func easyRouter(str1, str2 string) bool {
	if len(str1) >= len(str2) && str1[:len(str2)] == str2 {
		return true
	}
	return false
}
func httprouter(w http.ResponseWriter, r *http.Request) {
	if easyRouter(r.RequestURI, "/prometheus/memory") {
		memoryTotal(w, r)
		return
	}
}
func memoryTotal(w http.ResponseWriter, r *http.Request) {
	// timeNow := time.Now().UnixNano() / 1e9
	params := r.URL.Query()
	id := params["id"]
	if len(id) < 1 {
		w.Write([]byte("404"))
		return
	}
	all := myredis.GetDetail(id[0])
	key := "memory_" + id[0]
	exportData := "# HELP " + key + " The memory usage situation  of the entire of cluster.\n"
	exportData += "# TYPE " + key + " gauge\n"
	var memoryT int64
	for _, v := range all {
		t, _ := strconv.ParseInt(myredis.Trim(v.Memory.UsedMemory), 10, 0)
		memoryT += t
		exportData += key + "{type=\"each\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.Memory.UsedMemory) + " \n"
	}
	exportData += key + "{type=\"total\",host=\"*\"} " + strconv.FormatInt(memoryT, 10) + " \n"
	w.Write([]byte(exportData))
}
