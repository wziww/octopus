package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"octopus/config"
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
	var c config.RedisDetail
	for _, v := range config.C.Redis {
		if v.Name == names[0] {
			c = v
			goto next
		}
	}
	w.Write([]byte("404"))
	return
next:
	id := fmt.Sprintf("%x", md5.Sum([]byte(c.Name)))
	all := myredis.GetStatsObj(id)
	key := "cache_hit_rate_" + c.Name
	var exportData bytes.Buffer
	exportData.WriteString(fmt.Sprintf("%s%s%s\n", "# HELP ", key, " The cache hit rate of the entire of cluster."))
	exportData.WriteString(fmt.Sprintf("%s%s%s\n", "# TYPE ", key, " gauge"))
	var hit int
	var total int
	for _, v := range all {
		h, _ := strconv.Atoi(myredis.Trim(v.KeyspaceHits))
		m, _ := strconv.Atoi(myredis.Trim(v.KeyspaceMisses))
		hit += h
		total += h + m
		if h+m == 0 {
			exportData.WriteString(fmt.Sprintf("%s%s%s%s\n", "{type=\"each\",host=\"", v.ADDR, "\"} ", "100"))
		} else {
			exportData.WriteString(fmt.Sprintf("%s%s%s%d\n", "{type=\"each\",host=\"", v.ADDR, "\"} ", h*100/(h+m)))
		}
	}
	if total == 0 {
		exportData.WriteString("{type=\"total\",host=\"*\"} 100\n")
	} else {
		exportData.WriteString(fmt.Sprintf("%s%d\n", "{type=\"total\",host=\"*\"} ", hit*100/total))
	}
	w.Write(exportData.Bytes())
}

// 慢日志统计项
func slowlog(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	if len(names) < 1 {
		w.Write([]byte("404"))
		return
	}
	var c config.RedisDetail
	for _, v := range config.C.Redis {
		if v.Name == names[0] {
			c = v
			goto next
		}
	}
	w.Write([]byte("404"))
	return
next:
	id := fmt.Sprintf("%x", md5.Sum([]byte(c.Name)))
	all := myredis.GetSlowLogObj(id)
	key := "slowlog_" + c.Name
	var exportData bytes.Buffer
	exportData.WriteString(fmt.Sprintf("%s%s%s", "# HELP ", key, " The slowlog count of the entire of cluster.\n"))
	exportData.WriteString(fmt.Sprintf("%s%s%s", "# TYPE ", key, " gauge\n"))
	var memoryT int
	for _, v := range all {
		memoryT += v.Count
		exportData.WriteString(fmt.Sprintf("%s%s%s%s%d%s", key, "{type=\"each\",host=\"", v.Addr, "\"} ", v.Count, " \n"))
	}
	exportData.WriteString(fmt.Sprintf("%s%s%d%s", key, "{type=\"total\",host=\"*\"} ", memoryT, " \n"))
	w.Write(exportData.Bytes())
}

// 内存统计项
func memoryTotal(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	w.Write(_memoryTotal(names))
}
func _memoryTotal(names []string) []byte {
	if len(names) < 1 {
		return []byte("404")
	}
	var c config.RedisDetail
	for _, v := range config.C.Redis {
		if v.Name == names[0] {
			c = v
			goto next
		}
	}
	return []byte("404")
next:
	id := fmt.Sprintf("%x", md5.Sum([]byte(c.Name)))
	all := myredis.GetDetailObj(id)
	key := "memory_" + c.Name
	var exportData bytes.Buffer
	exportData.WriteString(fmt.Sprintf("%s%s%s", "# HELP ", key, " The memory usage situation of the entire of cluster.\n"))
	exportData.WriteString(fmt.Sprintf("%s%s%s", "# TYPE ", key, " gauge\n"))
	var memoryT int64
	for _, v := range all {
		t, _ := strconv.ParseInt(myredis.Trim(v.Memory.UsedMemory), 10, 0)
		memoryT += t
		exportData.WriteString(fmt.Sprintf("%s%s%s%s%s%s", key, "{type=\"each\",host=\"", v.ADDR, "\"} ", myredis.Trim(v.Memory.UsedMemory), " \n"))
	}
	exportData.WriteString(fmt.Sprintf("%s%s%d%s", key, "{type=\"total\",host=\"*\"} ", memoryT, " \n"))
	return exportData.Bytes()
}

// 流量 / ops 统计项
func opsTotal(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	w.Write(_opsTotal(names))
}
func _opsTotal(names []string) []byte {
	if len(names) < 1 {
		return []byte("404")
	}
	var c config.RedisDetail
	for _, v := range config.C.Redis {
		if v.Name == names[0] {
			c = v
			goto next
		}
	}
	return []byte("404")
next:
	id := fmt.Sprintf("%x", md5.Sum([]byte(c.Name)))
	all := myredis.GetStatsObj(id)
	key := "ops_" + c.Name
	var exportData bytes.Buffer
	exportData.WriteString(fmt.Sprintf("%s%s%s", "# HELP ", key, " The ops & outputKbps & inputKbps situation  of the entire of cluster.\n"))
	exportData.WriteString(fmt.Sprintf("%s%s%s", "# TYPE ", key, " gauge\n"))
	var InstantaneousOutputKbps, InstantaneousInputKbps, InstantaneousOpsPerSec float64
	for _, v := range all {
		// InstantaneousOutputKbps
		t, _ := strconv.ParseFloat(myredis.Trim(v.InstantaneousOutputKbps), 64)
		InstantaneousOutputKbps += t
		exportData.WriteString(fmt.Sprintf("%s%s%s%s%s%s", key, "{type=\"each-okbps\",host=\"", v.ADDR, "\"} ", myredis.Trim(v.InstantaneousOutputKbps), " \n"))
		// InstantaneousInputKbps
		t2, _ := strconv.ParseFloat(myredis.Trim(v.InstantaneousInputKbps), 64)
		InstantaneousInputKbps += t2
		exportData.WriteString(fmt.Sprintf("%s%s%s%s%s%s", key, "{type=\"each-ikbps\",host=\"", v.ADDR, "\"} ", myredis.Trim(v.InstantaneousInputKbps), " \n"))
		// InstantaneousOpsPerSec
		t3, _ := strconv.ParseFloat(myredis.Trim(v.InstantaneousOpsPerSec), 64)
		InstantaneousOpsPerSec += t3
		exportData.WriteString(fmt.Sprintf("%s%s%s%s%s%s", key, "{type=\"each-ikbps\",host=\"", v.ADDR, "\"} ", myredis.Trim(v.InstantaneousOpsPerSec), " \n"))
	}
	exportData.WriteString(fmt.Sprintf("%s%s%.2f%s", key, "{type=\"total-okbps\",host=\"*\"} ", InstantaneousOutputKbps, " \n"))
	exportData.WriteString(fmt.Sprintf("%s%s%.2f%s", key, "{type=\"total-ikbps\",host=\"*\"} ", InstantaneousInputKbps, " \n"))
	exportData.WriteString(fmt.Sprintf("%s%s%.2f%s", key, "{type=\"total-ops\",host=\"*\"} ", InstantaneousOpsPerSec, " \n"))
	return exportData.Bytes()
}
