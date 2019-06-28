<template>
  <div>
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
let chartData = {};
let t = null;
let timeData = [];
let statsDataT = [];
export default {
  name: "setting_redis",
  data() {
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
        集群内存占用: {
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
              return `${seriesName}\n${value * 100}%`;
            },
            fontSize: 20
          }
        }
      }
    };
    const that = this;
    this.$socket.sendObj({
      Func: "/config/redis/detail",
      Data: JSON.stringify({ id: that.$route.query.id })
    });
    this.$socket.sendObj({
      Func: "/redis/stats",
      Data: JSON.stringify({ id: that.$route.query.id })
    });
    t = setInterval(() => {
      this.$socket.sendObj({
        Func: "/config/redis/detail",
        Data: JSON.stringify({ id: that.$route.query.id })
      });
      this.$socket.sendObj({
        Func: "/redis/stats",
        Data: JSON.stringify({ id: that.$route.query.id })
      });
    }, 1000);
    this.$socket.onmessage = da => {
      const d = JSON.parse(da.data);
      if (d.Type === "/config/redis/detail") {
        let UsedMemoryTotal = 0;
        let TotalSystemMemoryTotal = 0;
        let Maxmemory = 0;
        for (let i of d.Data) {
          UsedMemoryTotal += Number(i.UsedMemory);
          TotalSystemMemoryTotal += Number(i.TotalSystemMemory);
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
              key: "集群内存占用",
              percent: (UsedMemoryTotal / Maxmemory).toFixed(4)
            }
          ]
        };
        that.chartData = chartData;
      }
      if (d.Type === "/redis/stats") {
        console.log(d);
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
    };
    return {
      chartData,
      lineChartData: {
        columns: ["t", "memory_total"],
        rows: timeData
      },
      statsData: {
        columns: ["t", "output_Kbps", "input_Kbps", "Ops"],
        rows: statsDataT
      }
    };
  },
  beforeDestroy() {
    if (t !== null) {
      window.clearInterval(t);
    }
  },
  methods: {
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
}</style>
