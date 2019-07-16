package cluster

import "strconv"

// AddSlotsLua ...
func AddSlotsLua(start, end int64) string {
	return "for i=" + strconv.FormatInt(start, 10) + "," + strconv.FormatInt(end, 10) + ", 1 do redis.call(\"CLUSTER\",\"ADDSLOTS\",i)  end return 'success'"
}

// // DelSlotsLua ...
// func DelSlotsLua(start, end int64) string {
// 	return "for i=" + strconv.FormatInt(start, 10) + "," + strconv.FormatInt(end, 10) + ", 1 do redis.call(\"CLUSTER\",\"DELSLOTS\",i)  end return 'success'"
// }
