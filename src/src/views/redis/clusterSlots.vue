<template>
  <div>
    <a-tooltip
      placement="topLeft"
      title="此处仅展示与 slot 有关的节点信息"
      :getPopupContainer="getPopupContainer"
      style="float: left;"
    >
      <a-button type="primary">说明</a-button>
    </a-tooltip>
    <router-link
      style="float: left;margin-left: 10px;"
      :to="'/redis_dev?id='+$route.query.id"
      class="hd"
    >
      <a-button v-if="permission&permissionAll.PERMISSIONDEV" type="danger">dev</a-button>
    </router-link>
    <div style="width: 100%;float: left;margin-bottom: 20px;">
      <a-alert
        message="存在迁移中的 slots"
        type="warning"
        showIcon
        :style="{display: reShardingSlots,maxWidth:'400px',float: 'left'}"
      />
      <a-alert
        message="存在节点未配置「最大可用内存」 : maxmemory"
        type="warning"
        showIcon
        :style="{display: maxmemoryWarning.indexOf('show')===-1?'none':'',width:'400px',float: 'left'}"
      />
      <a-alert
        message="存在节点内存使用达 90% +"
        type="warning"
        showIcon
        :style="{display: usedmemoryWarning,width:'400px',float: 'left'}"
      />
    </div>
    <div style="width: 100%;float: left;margin-bottom: 20px;">
      <span class="each-chose" style="font-size: 20px;font-weight: border;">Refresh Every :</span>
      <a-button class="each-chose" :type="index[0]" @click="chose(0)">1 s</a-button>
      <a-button class="each-chose" :type="index[1]" @click="chose(1)">10 s</a-button>
      <a-button class="each-chose" :type="index[2]" @click="chose(2)">30 s</a-button>
      <a-button class="each-chose" :type="index[3]" @click="chose(3)">1 min</a-button>
      <a-button class="each-chose" :type="index[4]" @click="chose(4)">5 min</a-button>
      <a-button class="each-chose" :type="index[5]" @click="chose(5)">10 min</a-button>
    </div>
    <a-table :scroll="{ x: 1000 }" :dataSource="data" style="float: left;width: 100%;">
      <a-table-column title="address" data-index="address" key="address">
        <template slot-scope="address">
          <div style="width: 100px;word-wrap:break-word;" :key="address">{{address}}</div>
        </template>
      </a-table-column>
      <a-table-column title="version" data-index="version" key="version">
        <template slot-scope="version">
          <a-tag
            :color="version.value===''?'red':version.color"
          >{{version.value===""?"未知":version.value}}</a-tag>
        </template>
      </a-table-column>
      <a-table-column title="id" data-index="id" key="id" v-if="type==='cluster'">
        <template slot-scope="id">
          <div
            v-for="each in split(id)"
            style="width: 100px;word-wrap:break-word;"
            :key="each"
          >{{each}}</div>
        </template>
      </a-table-column>
      <a-table-column title="follow" data-index="follow" key="follow" v-if="type==='cluster'">
        <template slot-scope="follow">
          <div
            v-for="each in split(follow)"
            style="width: 100px;word-wrap:break-word;"
            :key="each"
          >{{each}}</div>
        </template>
      </a-table-column>
      <a-table-column title="角色" data-index="role" key="role" v-if="type==='cluster'">
        <template slot-scope="role">
          <span>
            <a-tag v-for="each in role" :color="each.COLOR" :key="each.ROLE">{{each.ROLE}}</a-tag>
          </span>
        </template>
      </a-table-column>
      <a-table-column title="epoth 值" data-index="epoth" key="epoth" v-if="type==='cluster'" />
      <a-table-column title="拥有 slot（槽点）" data-index="slot" key="slot" v-if="type==='cluster'">
        <template slot-scope="slot">
          <a-tag
            @click="slotsClick(each)"
            v-for="each in slot"
            :key="each"
            :color="each.indexOf('->-')!==-1||each.indexOf('-<-')!==-1?'#ff001d':'#042b36'"
          >{{each}}</a-tag>
        </template>
      </a-table-column>
      <a-table-column
        title="slot 拥有比例"
        data-index="slotPercent"
        key="slotPercent"
        v-if="type==='cluster'"
      >
        <template slot-scope="slotPercent">
          <a-progress
            type="circle"
            :percent="parseInt(slotPercent * 100)"
            :width="80"
            :status="slotPercent===1?'exception':'success'"
            :format="(e) => {
              return e===100?'100%':e + '%';
            }"
          />
        </template>
      </a-table-column>
      <a-table-column title="占用内存" data-index="UsedMemory" key="UsedMemory">
        <template slot-scope="UsedMemory">{{UsedMemory}}M</template>
      </a-table-column>
      <a-table-column title="可用内存" data-index="Maxmemory" key="Maxmemory">
        <template slot-scope="Maxmemory">
          <span :style="{color: Maxmemory.color}">{{Maxmemory.value}}M</span>
        </template>
      </a-table-column>
      <a-table-column title="系统总内存" data-index="TotalSystemMemory" key="TotalSystemMemory">
        <template slot-scope="TotalSystemMemory">{{TotalSystemMemory}}M</template>
      </a-table-column>
      <a-table-column title="内存占用比例" data-index="memoryPercent" key="memoryPercent">
        <template slot-scope="memoryPercent">
          <a-progress type="circle" :percent="parseInt(memoryPercent * 100)" :width="80" />
        </template>
      </a-table-column>
      <a-table-column title="状态" data-index="state" key="state">
        <template slot-scope="state">
          <span>
            <a-tag
              v-for="each in state"
              :color="each.COLOR"
              :key="each.STATE"
            >{{each.STATE?each.STATE:"disconnected"}}</a-tag>
          </span>
        </template>
      </a-table-column>
      <a-table-column title="opcap" data-index="operation" key="operation">
        <template slot-scope="operation">
          <a-tag :color="operation=='节点监控'?'#0ea7fb':'#2c2c2c'">{{operation}}</a-tag>
        </template>
      </a-table-column>
    </a-table>
  </div>
</template>
<script>
import hd from "../../lib/ws";
import WS from "../../lib/websocket";
import { token, permission, permissionAll } from "../../lib/token";
import config from "../../config/index";
const PATH = "monit";
const ws = new WS(config.Host + "?op=" + PATH + "&ot=" + token + "&ocid=nil");
let data = [];
let type = "cluster";
let interTime = 1000;
let t = null;
let t2 = null;
let index = ["primary", "default", "default", "default", "default", "default"];
export default {
  name: "setting_redis",
  data() {
    data = [];
    let maxmemoryWarning = ["none"];
    let usedmemoryWarning = "none";
    let reShardingSlots = "none";
    const that = this;
    ws.Open();
    ws.OnOpen(() => {
      t = setInterval(() => {
        try {
          ws.SendObj({
            Func: "/redis/detail",
            Data: JSON.stringify({ id: that.$route.query.id })
          });
        } catch (e) {
          console.error(e);
        }
      }, interTime);
    });
    const handMessage = hd(d => {
      if (d.Type === "/redis/detail") {
        that.maxmemoryWarning = [];
        data = [];
        for (let i of d.Data) {
          that.type = i.Type;
          data.push({
            key: i.ID,
            id: i.ID,
            address: i.ADDR,
            follow: i.FOLLOW,
            role: [
              {
                ROLE: i.ROLE.indexOf("master") !== -1 ? "master" : "slave",
                COLOR: i.ROLE.indexOf("master") !== -1 ? "blue" : "green"
              }
            ],
            epoth: i.EPOTH,
            slot: i.SLOT.split(" ").filter(e => {
              if (e.indexOf("-") !== -1 && Number.isInteger(Number(e[0]))) {
                return e;
              } else if ("" + e === "" + Number(e)) {
                return e;
              } else if (e.indexOf("->-") !== -1 || e.indexOf("-<-") !== -1) {
                // Migrating - Importing
                that.reShardingSlots = "";
                return e;
              }
            }),
            slotPercent: (() => {
              let has = 0;
              for (let z of i.SLOT.split(" ")) {
                if (z.indexOf("-") !== -1 && Number.isInteger(Number(z[0]))) {
                  has += Number(z.split("-")[1]) - Number(z.split("-")[0]);
                } else if ("" + z === "" + Number(z)) {
                  has++;
                }
              }
              return (has + 1) / 16384;
            })(),
            memoryPercent: (() => {
              const percent = Number(i.UsedMemory) / Number(i.Maxmemory);
              if (Infinity !== percent && percent >= 0.9) {
                if (parseInt(Math.random() * 100) > 80) {
                  that.usedmemoryWarning = "";
                }
              }
              return percent;
            })(),
            state: [
              {
                STATE: i.STATE,
                COLOR: i.STATE === "connected" ? "#00c94d" : "RED"
              }
            ],
            version: (() => {
              return {
                color: i.VERSION.startsWith("4.") ? "#00c94d" : "#2593fc",
                value: i.VERSION
              };
            })(),
            Maxmemory: (() => {
              const tmp = i.Maxmemory / 1024 / 1024;
              if (tmp === 0) {
                that.maxmemoryWarning.push("show");
              } else {
                that.maxmemoryWarning.push("none");
              }
              return {
                value: tmp.toFixed(2),
                color: tmp < 1 ? "red" : ""
              };
            })(),
            UsedMemory: (i.UsedMemory / 1024 / 1024).toFixed(2),
            TotalSystemMemory: (i.TotalSystemMemory / 1024 / 1024).toFixed(2),
            operation: i.OpcapOnline ? "节点监控" : "不可用"
          });
        }
        data = data
          .sort((a, b) => {
            return Number(a.epoth) - Number(b.epoth);
          })
          .sort((a, b) => {
            if (a.role[0].ROLE > b.role[0].ROLE) return 1;
            return -1;
          });
        that.data = data;
      }
    });
    ws.OnData(handMessage);
    return {
      data,
      index,
      type,
      maxmemoryWarning,
      usedmemoryWarning,
      reShardingSlots,
      permission,
      permissionAll
    };
  },
  methods: {
    getPopupContainer(trigger) {
      return trigger.parentElement;
    },
    slotsClick(e) {
      console.log(e);
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
            this.$socket.sendObj({
              Func: "/redis/detail",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
          }, interTime);
          break;
        case 1:
          interTime = 1000 * 10;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            this.$socket.sendObj({
              Func: "/redis/detail",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
          }, interTime);
          break;
        case 2:
          interTime = 1000 * 30;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            this.$socket.sendObj({
              Func: "/redis/detail",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
          }, interTime);
          break;
        case 3:
          interTime = 1000 * 60;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            this.$socket.sendObj({
              Func: "/redis/detail",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
          }, interTime);
          break;
        case 4:
          interTime = 1000 * 60 * 5;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            this.$socket.sendObj({
              Func: "/redis/detail",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
          }, interTime);
          break;
        case 5:
          interTime = 1000 * 60 * 10;
          this.interTime = interTime;
          window.clearInterval(t);
          t = setInterval(() => {
            this.$socket.sendObj({
              Func: "/redis/detail",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
          }, interTime);
          break;
      }
    },
    split: str => {
      if (typeof str !== "string") return [];
      const arr = [str];
      return arr;
    }
  },
  beforeDestroy() {
    ws.Close();
    if (t !== null) {
      window.clearInterval(t);
    }
    if (t2 !== null) {
      window.clearInterval(t2);
    }
  }
};
</script>
<style lang="stylus" scoped>
.each-chose {
  float: left;
  margin-right: 10px;
}

.ant-table td {
  white-space: nowrap;
}
</style>
