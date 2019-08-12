<template>
  <div>
    <div style="width: 80%;float: left;margin-bottom: 20px;">
      <span class="each-chose" style="font-size: 20px;font-weight: border;">Refresh Every :</span>
      <a-button class="each-chose" :type="index[0]" @click="chose(0)">1 s</a-button>
      <a-button class="each-chose" :type="index[1]" @click="chose(1)">10 s</a-button>
      <a-button class="each-chose" :type="index[2]" @click="chose(2)">30 s</a-button>
      <a-button class="each-chose" :type="index[3]" @click="chose(3)">1 min</a-button>
      <a-button class="each-chose" :type="index[4]" @click="chose(4)">5 min</a-button>
      <a-button class="each-chose" :type="index[5]" @click="chose(5)">10 min</a-button>
      <a-button type="primary" @click="showDrawer">设置</a-button>
    </div>
    <div class="each-chart">
      <ve-line :data="lineChartData" :settings="lineChartSettings"></ve-line>
      <a-tag color="#f50">命令趋势</a-tag>
    </div>
    <div class="each-chart">
      <ve-pie :data="data"></ve-pie>
      <a-tag color="#f50">5 分钟内命令统计</a-tag>
    </div>
    <a-drawer
      width="500"
      title="监控命令设置"
      placement="right"
      :closable="false"
      @close="onClose"
      :visible="visible"
    >
      <a-checkbox-group :options="plainOptions" v-model="checkedList" @change="onChange" />
    </a-drawer>
  </div>
</template>
<script>
import hd from "../../lib/ws";
import WS from "../../lib/websocket";
import { token } from "../../lib/token";
import config from "../../config/index";
const PATH = "monit";
let timeData = [];
const ws = new WS(config.Host + "?op=" + PATH + "&ot=" + token + "&ocid=nil");
let data = {};
let interTime = 1000;
let index = ["primary", "default", "default", "default", "default", "default"];
let t = null;
let checkedList = JSON.parse(localStorage.getItem("checkedList")) || [
  "get",
  "set"
];
export default {
  name: "setting_redis",
  data() {
    const that = this;
    ws.Open();
    ws.OnOpen(() => {
      t = setInterval(() => {
        try {
          ws.SendObj({
            Func: "/opcap",
            Data: JSON.stringify({ address: that.$route.query.address })
          });
        } catch (e) {
          console.error(e);
        }
      }, interTime);
    });
    const handMessage = hd(d => {
      if (d.Type === "/opcap") {
        let tmpD = [];
        if (d.Data && typeof d.Data === "string") {
          d.Data = d.Data.split("_");
        }
        for (let i = 0; i < d.Data.length; i += 2) {
          tmpD.push({
            命令: d.Data[i],
            次数: d.Data[i + 1]
          });
        }
        if (timeData.length >= 20) {
          timeData.shift();
        }
        const columns = ["t"];
        const t = that.$moment().format("hh:mm:ss");
        const obj = {
          t
        };
        let _tmpD = [];
        for (let i = 0; i < tmpD.length; i++) {
          if (checkedList.indexOf(tmpD[i].命令) !== -1) {
            _tmpD.push(tmpD[i]);
            obj[tmpD[i].命令] = tmpD[i].次数;
            columns.push(tmpD[i].命令);
          }
        }
        timeData.push(obj);
        that.lineChartData = {
          columns,
          rows: timeData
        };
        that.data = {
          columns: ["命令", "次数"],
          rows: _tmpD
        };
      }
    });
    ws.OnData(handMessage);
    this.lineChartSettings = {
      area: false,
      scale: [true, true],
      yAxisName: ["value"],
      xAxisName: ["时间"]
    };
    return {
      lineChartData: {},
      data,
      checkedList,
      plainOptions: [
        "set",
        "setnx",
        "setex",
        "psetex",
        "get",
        "getset",
        "strlen",
        "append",
        "setrange",
        "getrange",
        "incr",
        "incrby",
        "incrbyfloat",
        "decr",
        "decrby",
        "mset",
        "msetnx",
        "mget",
        // hash table
        "hset",
        "hsetnx",
        "hget",
        "hexists",
        "hdel",
        "hlen",
        "hstrlen",
        "hincrby",
        "hincrbyfloat",
        "hmset",
        "hmget",
        "hkeys",
        "hvals",
        "hgetall",
        "hscan",
        // list
        "lpush",
        "lpushx",
        "rpush",
        "rpushx",
        "lpop",
        "rpop",
        "rpoplpush",
        "lrem",
        "llen",
        "lindex",
        "linsert",
        "lset",
        "lrange",
        "ltrim",
        "blpop",
        "brpop",
        "brpoplpush",
        // set
        "sadd",
        "sismember",
        "spop",
        "srandmember",
        "srem",
        "smove",
        "scard",
        "sasmembersdd",
        "sscan",
        "sinter",
        "sinterstore",
        "sunion",
        "sunionstore",
        "sdiff",
        "sdiffstore",
        // zip list
        "zadd",
        "zscore",
        "zincrby",
        "zcard",
        "zcount",
        "zrange",
        "zrevrange",
        "zrangebyscore",
        "zrevrangebyscore",
        "zrank",
        "zrevrank",
        "zrem",
        "zremrangebyrank",
        "zremrangebyscore",
        "zrangebylex",
        "zlexcount",
        "zremrangebylex",
        "zscan",
        "zunionstore",
        "zinterstore",
        // bit map
        "setbit",
        "getbit",
        "bitcount",
        "bitpos",
        "bitop",
        "bitfield",
        // db
        "exists",
        "type",
        "rename",
        "renamenx",
        "move",
        "del",
        "randomkey",
        "dbsize",
        "keys",
        "scan",
        // ttl
        "expire",
        "expireat",
        "ttl"
      ],
      visible: false,
      index
    };
  },
  methods: {
    onChange(checkedValues) {
      checkedList = checkedValues;
      this.checkedList = checkedList;
      localStorage.setItem("checkedList", JSON.stringify(checkedList));
    },
    showDrawer() {
      this.visible = true;
    },
    onClose() {
      this.visible = false;
    },
    chose(x) {
      const that = this;
      index = [
        "default",
        "default",
        "default",
        "default",
        "default",
        "default"
      ];
      index[x] = "primary";
      this.index = index;
      switch (x) {
        case 0:
          interTime = 1000;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            try {
              ws.SendObj({
                Func: "/opcap",
                Data: JSON.stringify({ address: that.$route.query.address })
              });
            } catch (e) {
              console.error(e);
            }
          }, interTime);
          break;
        case 1:
          interTime = 1000 * 10;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            try {
              ws.SendObj({
                Func: "/opcap",
                Data: JSON.stringify({ address: that.$route.query.address })
              });
            } catch (e) {
              console.error(e);
            }
          }, interTime);
          break;
        case 2:
          interTime = 1000 * 30;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            try {
              ws.SendObj({
                Func: "/opcap",
                Data: JSON.stringify({ address: that.$route.query.address })
              });
            } catch (e) {
              console.error(e);
            }
          }, interTime);
          break;
        case 3:
          interTime = 1000 * 60;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            try {
              ws.SendObj({
                Func: "/opcap",
                Data: JSON.stringify({ address: that.$route.query.address })
              });
            } catch (e) {
              console.error(e);
            }
          }, interTime);
          break;
        case 4:
          interTime = 1000 * 60 * 5;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            try {
              ws.SendObj({
                Func: "/opcap",
                Data: JSON.stringify({ address: that.$route.query.address })
              });
            } catch (e) {
              console.error(e);
            }
          }, interTime);
          break;
        case 5:
          interTime = 1000 * 60 * 10;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            try {
              ws.SendObj({
                Func: "/opcap",
                Data: JSON.stringify({ address: that.$route.query.address })
              });
            } catch (e) {
              console.error(e);
            }
          }, interTime);
          break;
      }
    },
    split: str => {
      if (typeof str !== "string") return [];
      const len = str.length;
      const arr = [];
      for (let i = 0; i < len; i += 10) {
        arr.push(str.substr(i, 10));
      }
      return arr;
    }
  },
  beforeDestroy() {
    ws.Close();
    if (t !== null) {
      window.clearInterval(t);
    }
  }
};
</script>
<style lang="stylus" scoped>
.each-chart {
  width: 50%;
  float: left;
  min-height: 440px;
}

.each-chose {
  float: left;
  margin-right: 10px;
}

.ant-table td {
  white-space: nowrap;
}
</style>
