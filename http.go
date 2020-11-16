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
	} else if easyRouter(r.RequestURI, "/prometheus/slowlog") {
		slowlog(w, r)
		return
	} else if easyRouter(r.RequestURI, "/prometheus/cacheRate") {
		cacheRate(w, r)
		return
	}
}

// 慢日志统计项
func cacheRate(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	if len(names) < 1 {
		w.Write([]byte("404"))
		return
	}
	id := fmt.Sprintf("%x", md5.Sum([]byte(names[0])))
	all := myredis.GetStatsObj(id)
	key := "cache_hit_rate_" + names[0]
	exportData := "# HELP " + key + " The cache hit rate of the entire of cluster.\n"
	exportData += "# TYPE " + key + " gauge\n"
	var hit int
	var total int
	for _, v := range all {
		h, _ := strconv.Atoi(myredis.Trim(v.KeyspaceHits))
		m, _ := strconv.Atoi(myredis.Trim(v.KeyspaceMisses))
		hit += h
		total += h + m
		if h+m == 0 {
			exportData += key + "{type=\"each\",host=\"" + v.ADDR + "\"} " + "100" + " \n"
		} else {
			exportData += key + "{type=\"each\",host=\"" + v.ADDR + "\"} " + strconv.Itoa(h*100/(h+m)) + " \n"
		}
	}
	if total == 0 {
		exportData += key + "{type=\"total\",host=\"*\"} " + "100" + " \n"
	} else {
		exportData += key + "{type=\"total\",host=\"*\"} " + strconv.Itoa(hit*100/total) + " \n"
	}
	w.Write([]byte(exportData))
}

// 慢日志统计项
func slowlog(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	if len(names) < 1 {
		w.Write([]byte("404"))
		return
	}
	id := fmt.Sprintf("%x", md5.Sum([]byte(names[0])))
	all := myredis.GetSlowLogObj(id)
	key := "slowlog_" + names[0]
	exportData := "# HELP " + key + " The slowlog count of the entire of cluster.\n"
	exportData += "# TYPE " + key + " gauge\n"
	var memoryT int
	for _, v := range all {
		memoryT += v.Count
		exportData += key + "{type=\"each\",host=\"" + v.Addr + "\"} " + strconv.Itoa(v.Count) + " \n"
	}
	exportData += key + "{type=\"total\",host=\"*\"} " + strconv.Itoa(memoryT) + " \n"
	w.Write([]byte(exportData))
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
	all := myredis.GetDetailObj(id)
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
	var InstantaneousOutputKbps, InstantaneousInputKbps, InstantaneousOpsPerSec float64
	for _, v := range all {
		// InstantaneousOutputKbps
		t, _ := strconv.ParseFloat(myredis.Trim(v.InstantaneousOutputKbps), 64)
		InstantaneousOutputKbps += t
		exportData += key + "{type=\"each-okbps\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.InstantaneousOutputKbps) + " \n"
		// InstantaneousInputKbps
		t2, _ := strconv.ParseFloat(myredis.Trim(v.InstantaneousInputKbps), 64)
		InstantaneousInputKbps += t2
		exportData += key + "{type=\"each-ikbps\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.InstantaneousInputKbps) + " \n"
		// InstantaneousOpsPerSec
		t3, _ := strconv.ParseFloat(myredis.Trim(v.InstantaneousOpsPerSec), 64)
		InstantaneousOpsPerSec += t3
		exportData += key + "{type=\"each-ikbps\",host=\"" + v.ADDR + "\"} " + myredis.Trim(v.InstantaneousOpsPerSec) + " \n"
	}
	exportData += key + "{type=\"total-okbps\",host=\"*\"} " + fmt.Sprintf("%.2f", InstantaneousOutputKbps) + " \n"
	exportData += key + "{type=\"total-ikbps\",host=\"*\"} " + fmt.Sprintf("%.2f", InstantaneousInputKbps) + " \n"
	exportData += key + "{type=\"total-ops\",host=\"*\"} " + fmt.Sprintf("%.2f", InstantaneousOpsPerSec) + " \n"
	w.Write([]byte(exportData))
}
