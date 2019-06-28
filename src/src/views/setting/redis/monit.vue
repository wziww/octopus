<template>
  <div>
    <a-table :scroll="{ x: 1000 ,y:800}" :dataSource="data">
      <a-table-column title="address" data-index="address" key="address"/>
      <a-table-column title="id" data-index="id" key="id">
        <template slot-scope="id">
          <div v-for="each in split(id)" style="width: 100%;" :key="each">{{each}}</div>
        </template>
      </a-table-column>
      <a-table-column title="follow" data-index="follow" key="follow">
        <template slot-scope="follow">
          <div v-for="each in split(follow)" style="width: 100%;" :key="each">{{each}}</div>
        </template>
      </a-table-column>
      <a-table-column title="角色" data-index="role" key="role">
        <template slot-scope="role">
          <span>
            <a-tag v-for="each in role" :color="each.COLOR" :key="each.ROLE">{{each.ROLE}}</a-tag>
          </span>
        </template>
      </a-table-column>
      <a-table-column title="epoth 值" data-index="epoth" key="epoth"/>
      <a-table-column title="拥有 slot（槽点）" data-index="slot" key="slot">
        <template slot-scope="slot">
          <a-tag v-for="each in slot" :key="each" color="#042b36">{{each}}</a-tag>
        </template>
      </a-table-column>
      <a-table-column title="slot 拥有比例" data-index="slotPercent" key="slotPercent">
        <template slot-scope="slotPercent">
          <a-progress type="circle" :percent="parseInt(slotPercent * 100)" :width="80"/>
        </template>
      </a-table-column>
      <a-table-column title="占用内存" data-index="UsedMemory" key="UsedMemory">
        <template slot-scope="UsedMemory">{{UsedMemory}}M</template>
      </a-table-column>
      <a-table-column title="可用内存" data-index="Maxmemory" key="Maxmemory">
        <template slot-scope="Maxmemory">{{Maxmemory}}M</template>
      </a-table-column>
      <a-table-column title="系统总内存" data-index="TotalSystemMemory" key="TotalSystemMemory">
        <template slot-scope="TotalSystemMemory">{{TotalSystemMemory}}M</template>
      </a-table-column>
      <a-table-column title="内存占用比例" data-index="memoryPercent" key="memoryPercent">
        <template slot-scope="memoryPercent">
          <a-progress type="circle" :percent="parseInt(memoryPercent * 100)" :width="80"/>
        </template>
      </a-table-column>
      <a-table-column title="状态" data-index="state" key="state">
        <template slot-scope="state">
          <span>
            <a-tag v-for="each in state" :color="each.COLOR" :key="each.STATE">{{each.STATE}}</a-tag>
          </span>
        </template>
      </a-table-column>
      <a-table-column title="操作" data-index="operation" key="operation">
        <template slot-scope="operation">
          <a-tag color="#0ea7fb">{{"节点监控"||operation}}</a-tag>
        </template>
      </a-table-column>
    </a-table>
  </div>
</template>
<script>
let data = [];
export default {
  name: "setting_redis",
  data() {
    const that = this;
    this.$socket.sendObj({
      Func: "/config/redis/detail",
      Data: JSON.stringify({ id: that.$route.query.id })
    });
    this.$socket.onmessage = da => {
      const d = JSON.parse(da.data);
      if (d.Type === "/config/redis/detail") {
        data = [];
        for (let i of d.Data) {
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
              if (e.indexOf("-") !== -1) return e;
            }),
            slotPercent: (() => {
              let has = 0;
              for (let z of i.SLOT.split(" ")) {
                if (z.indexOf("-") !== -1) {
                  has += Number(z.split("-")[1]) - Number(z.split("-")[0]);
                }
              }
              return has / 16383;
            })(),
            memoryPercent: (() => {
              return Number(i.UsedMemory) / Number(i.Maxmemory);
            })(),
            state: [
              {
                STATE: i.STATE,
                COLOR: i.STATE === "connected" ? "#00c94d" : "RED"
              }
            ],
            Maxmemory: (i.Maxmemory / 1024 / 1024).toFixed(2),
            UsedMemory: (i.UsedMemory / 1024 / 1024).toFixed(2),
            TotalSystemMemory: (i.TotalSystemMemory / 1024 / 1024).toFixed(2),
            operation: []
          });
        }
        this.data = data;
      }
    };
    return {
      data
      // chartData
    };
  },
  beforeDestroy() {
    // console.log(1);
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
