package main

import (
	"crypto/md5"
	"fmt"
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
	} else if easyRouter(r.RequestURI, "/prometheus/ioo") { // inputkbps outputkbps ops
		opsTotal(w, r)
		return
	}
}

// 内存统计项
func memoryTotal(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	if len(names) < 1 {
		w.Write([]byte("404"))
		return
	}
	id := fmt.Sprintf("%x", md5.Sum([]byte(names[0])))
	fmt.Println(id)
	all := myredis.GetDetailObj(id)
	fmt.Println(all)
	key := "memory_" + names[0]
	exportData := "# HELP " + key + " The memory usage situation of the entire of cluster.\n"
	exportData += "# TYPE " + key + " gauge\n"
	var memoryT int64
	for _, v := range all {
		t, _ := strconv.ParseInt(myredis.Trim(v.Memory.UsedMemory), 10, 0)
		memoryT += t
		exportData += key + "{type=\"each\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.Memory.UsedMemory) + " \n"
	}
	exportData += key + "{type=\"total\",host=\"*\"} " + strconv.FormatInt(memoryT, 10) + " \n"
	fmt.Println(exportData)
	w.Write([]byte(exportData))
}

// 流量 / ops 统计项
func opsTotal(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	if len(names) < 1 {
		w.Write([]byte("404"))
		return
	}
	id := fmt.Sprintf("%x", md5.Sum([]byte(names[0])))
	all := myredis.GetStatsObj(id)
	key := "ops_" + names[0]
	exportData := "# HELP " + key + " The ops & outputKbps & inputKbps situation  of the entire of cluster.\n"
	exportData += "# TYPE " + key + " gauge\n"
	var InstantaneousOutputKbps, InstantaneousInputKbps, InstantaneousOpsPerSec int64
	for _, v := range all {
		// InstantaneousOutputKbps
		t, _ := strconv.ParseInt(myredis.Trim(v.InstantaneousOutputKbps), 10, 0)
		InstantaneousOutputKbps += t
		exportData += key + "{type=\"each-okbps\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.InstantaneousOutputKbps) + " \n"
		// InstantaneousInputKbps
		t2, _ := strconv.ParseInt(myredis.Trim(v.InstantaneousInputKbps), 10, 0)
		InstantaneousInputKbps += t2
		exportData += key + "{type=\"each-ikbps\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.InstantaneousInputKbps) + " \n"
		// InstantaneousOpsPerSec
		t3, _ := strconv.ParseInt(myredis.Trim(v.InstantaneousOpsPerSec), 10, 0)
		InstantaneousOpsPerSec += t3
		exportData += key + "{type=\"each-ikbps\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.InstantaneousOpsPerSec) + " \n"
	}
	exportData += key + "{type=\"total-okbps\",host=\"*\"} " + strconv.FormatInt(InstantaneousOutputKbps, 10) + " \n"
	exportData += key + "{type=\"total-ikbps\",host=\"*\"} " + strconv.FormatInt(InstantaneousInputKbps, 10) + " \n"
	exportData += key + "{type=\"total-ops\",host=\"*\"} " + strconv.FormatInt(InstantaneousInputKbps, 10) + " \n"
	w.Write([]byte(exportData))
}
