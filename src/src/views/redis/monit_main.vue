<template>
  <div>
    <div style="width: 100%;float: left;margin-bottom: 20px;">
      <span class="each-chose" style="font-size: 20px;font-weight: border;">Refresh Every :</span>
      <a-button class="each-chose" :type="index[0]" @click="chose(0)">1 s</a-button>
      <a-button class="each-chose" :type="index[1]" @click="chose(1)">10 s</a-button>
      <a-button class="each-chose" :type="index[2]" @click="chose(2)">30 s</a-button>
      <a-button class="each-chose" :type="index[3]" @click="chose(3)">1 min</a-button>
      <a-button class="each-chose" :type="index[4]" @click="chose(4)">5 min</a-button>
      <a-button class="each-chose" :type="index[5]" @click="chose(5)">10 min</a-button>
    </div>
    <ve-line style="width: 50%;float: left;" :data="lineChartData" :settings="lineChartSettings"></ve-line>
    <ve-liquidfill
      style="width: 50%;float: left;"
      radius="50%;"
      :data="chartData"
      :settings="chartSettings"
    ></ve-liquidfill>
    <ve-line style="width: 50%;float: left;" :data="statsData" :settings="statsSettings"></ve-line>
  </div>
</template>
<script>
import hd from "../../lib/ws";
import { token } from "../../lib/token";
import WS from "../../lib/websocket";
let chartData = {};
let t = null;
let index = ["primary", "default", "default", "default", "default", "default"];
let timeData = [];
let statsDataT = [];
let interTime = 1000;
const PATH = "dev";
const ws = new WS(
  "ws://0.0.0.0:8081/v1/websocket?op=" +
    PATH +
    "&ot=" +
    token +
    "&ocid=nil"
);
export default {
  name: "setting_redis",
  data() {
    statsDataT = [];
    chartData = [];
    timeData = [];
    ws.Open();
    this.lineChartSettings = {
      area: true,
      scale: [true, true],
      yAxisName: ["M"],
      xAxisName: ["时间"]
    };
    this.statsSettings = {
      area: true,
      scale: [true, true],
      yAxisName: ["value"],
      xAxisName: ["时间"]
    };
    this.chartSettings = {
      seriesMap: {
        内存使用量: {
          radius: "40%",
          center: ["20%", "30%"],
          itemStyle: {
            opacity: 0.2
          },
          emphasis: {
            itemStyle: {
              opacity: 0.5
            }
          },
          backgroundStyle: {},
          label: {
            formatter(options) {
              const { seriesName, value } = options;
              return `${seriesName}\n${(value * 100).toFixed(2)}%`;
            },
            fontSize: 20
          }
        }
      }
    };
    const that = this;
    t = setInterval(() => {
      try {
        ws.SendObj({
          Func: "/redis/detail",
          Data: JSON.stringify({ id: that.$route.query.id })
        });
        ws.SendObj({
          Func: "/redis/stats",
          Data: JSON.stringify({ id: that.$route.query.id })
        });
      } catch (e) {
        console.error(e);
      }
    }, interTime);
    ws.Close();
    ws.OnData(
      hd(d => {
        if (d.Type === "/redis/detail") {
          let UsedMemoryTotal = 0;
          // let TotalSystemMemoryTotal = 0;
          let Maxmemory = 0;
          for (let i of d.Data) {
            UsedMemoryTotal += Number(i.UsedMemory);
            // TotalSystemMemoryTotal += Number(i.TotalSystemMemory);
            Maxmemory += Number(i.Maxmemory);
          }
          if (timeData.length >= 20) {
            timeData.shift();
          }
          timeData.push({
            t: that.$moment().format("hh:mm:ss"),
            memory_total: (UsedMemoryTotal / 1024 / 1024).toFixed(2)
          });
          chartData = {
            columns: ["key", "percent"],
            rows: [
              {
                key: "内存使用量",
                percent: (UsedMemoryTotal / Maxmemory).toFixed(4)
              }
            ]
          };
          that.chartData = chartData;
        }
        if (d.Type === "/redis/stats") {
          let InstantaneousInputKbps = 0;
          let InstantaneousOutputKbps = 0;
          let InstantaneousOpsPerSec = 0;
          for (let i of d.Data) {
            InstantaneousOutputKbps += Number(i.InstantaneousOutputKbps);
            InstantaneousInputKbps += Number(i.InstantaneousInputKbps);
            InstantaneousOpsPerSec += Number(i.InstantaneousOpsPerSec);
          }
          if (statsDataT.length >= 20) {
            statsDataT.shift();
          }
          statsDataT.push({
            t: that.$moment().format("hh:mm:ss"),
            output_Kbps: InstantaneousOutputKbps,
            input_Kbps: InstantaneousInputKbps,
            Ops: InstantaneousOpsPerSec
          });
        }
      })
    );

    return {
      chartData,
      interTime,
      lineChartData: {
        columns: ["t", "memory_total"],
        rows: timeData
      },
      index,
      statsData: {
        columns: ["t", "output_Kbps", "input_Kbps", "Ops"],
        rows: statsDataT
      }
    };
  },
  beforeDestroy() {
    ws.Close();
    if (t !== null) {
      window.clearInterval(t);
    }
  },
  methods: {
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
            this.$socket.sendObj({
              Func: "/redis/stats",
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
            this.$socket.sendObj({
              Func: "/redis/stats",
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
            this.$socket.sendObj({
              Func: "/redis/stats",
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
            this.$socket.sendObj({
              Func: "/redis/stats",
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
            this.$socket.sendObj({
              Func: "/redis/stats",
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
            this.$socket.sendObj({
              Func: "/redis/stats",
              Data: JSON.stringify({ id: that.$route.query.id })
            });
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
  }
};
</script>
<style lang="stylus" scoped>
.ant-table td {
  white-space: nowrap;
}

.each-chose {
  float: left;
  margin-right: 10px;
}
</style>
