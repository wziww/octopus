<template>
  <div id="main-box" class="main-box">
    <pre id="container" v-highlightjs="data">
    <code class="string">
    </code>
    </pre>
    <div class="handleButton">
      <a-divider orientation="left">通用</a-divider>
      <a-button class="eachHandle" @click="clearOutPut" type="primary">清屏</a-button>
      <router-link class="eachHandle" :to="'/setting/clusterSlots?id='+$route.query.id">
        <a-button type="primary">集群列表</a-button>
      </router-link>
      <a-divider orientation="left">节点操作</a-divider>
      <a-tooltip
        placement="topLeft"
        title="<id> <ip:port> <flags> <master> <ping-sent> <pong-recv> <config-epoch> <link-state> <slot> <slot> ... <slot>"
        style="float: left;"
      >
        <a-button class="eachHandle" @click="reloadClusterNodes" type="primary">刷新节点信息</a-button>
      </a-tooltip>
      <a-button class="eachHandle" @click="clusterMeet" type="primary">添加节点</a-button>
      <a-drawer
        title="添加节点"
        placement="top"
        @close="clusterMeetClose"
        :closable="true"
        :visible="clusterMeetShow"
      >
        <div class="each-input">
          <a-input placeholder="host" @change="inputHost" />
        </div>
        <div class="each-input">
          <a-input placeholder="port" @change="inputPort" />
        </div>
        <a-popconfirm @confirm="confirmAddNode" :title="'确认添加节点['+newhost+':'+newport+']进入集群么？'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button style="width: 30%;float: left;" class="submit" type="primary">提交</a-button>
        </a-popconfirm>
      </a-drawer>
      <a-button class="eachHandle" @click="clusterForget" type="primary">删除节点</a-button>
      <a-drawer
        title="删除节点"
        placement="top"
        @close="clusterForgetClose"
        :closable="true"
        :visible="clusterForgetShow"
      >
        <div class="each-input">
          <a-input placeholder="nodeid" @change="inputNodeID" />
        </div>
        <a-popconfirm @confirm="confirmClusterForget" :title="'确认删除节点['+nodeid+']么？'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button style="width: 30%;float: left;" class="submit" type="primary">提交</a-button>
        </a-popconfirm>
      </a-drawer>
      <a-button class="eachHandle" @click="clusterReplicate" type="primary">从节点分配</a-button>
      <a-drawer
        title="从节点分配"
        placement="top"
        @close="clusterReplicateClose"
        :closable="true"
        :visible="clusterReplicateShow"
        height="300"
      >
        <div class="each-input">
          <a-input placeholder="需要设置为从节点的 host" @change="inputRepHost" />
        </div>
        <div class="each-input">
          <a-input placeholder="需要设置为从节点的 port" @change="inputRepPort" />
        </div>
        <div class="each-input">
          <a-input placeholder="主节点 nodeid" @change="inputRepNodeID" />
        </div>
        <a-popconfirm @confirm="confirmClusterReplicate" :title="'确认设置从节点么？'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button style="width: 30%;float: left;" class="submit" type="primary">提交</a-button>
        </a-popconfirm>
      </a-drawer>
      <a-divider orientation="left">slots 操作</a-divider>
      <a-button class="eachHandle" type="primary" @click="slotsStats">未分配 slots 计算</a-button>
      <a-button class="eachHandle" @click="slotsSet" type="primary">slots 分配</a-button>
      <a-drawer
        title="slots 分配"
        placement="top"
        @close="slotsSetClose"
        :closable="true"
        :visible="slotsSetShow"
        height="400"
      >
        <div class="each-input">
          <a-input placeholder="需要设置为从节点的 host" @change="inputSlotsHost" />
        </div>
        <div class="each-input">
          <a-input placeholder="需要设置为从节点的 port" @change="inputSlotsPort" />
        </div>
        <div class="each-input">
          <a-input type="number" placeholder="slots 起（0-16383）" @change="inputSlotsStart" />
        </div>
        <div class="each-input">
          <a-input type="number" placeholder="slots 止（0-16383）" @change="inputSlotsEnd" />
        </div>
        <a-popconfirm @confirm="confirmSlotsSet" :title="'确认进行 slots 分配？'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button style="width: 30%;float: left;" class="submit" type="primary">提交</a-button>
        </a-popconfirm>
      </a-drawer>
      <!-- <a-button class="eachHandle" @click="slotsDel" type="primary">slots 删除</a-button>
      <a-drawer
        title="slots 删除"
        placement="top"
        @close="slotsDelClose"
        :closable="true"
        :visible="slotsDelShow"
        height="400"
      >
        <div class="each-input">
          <a-input placeholder="需要设置为从节点的 host" @change="inputSlotsDelHost" />
        </div>
        <div class="each-input">
          <a-input placeholder="需要设置为从节点的 port" @change="inputSlotsDelPort" />
        </div>
        <div class="each-input">
          <a-input type="number" placeholder="slots 起（0-16383）" @change="inputSlotsDelStart" />
        </div>
        <div class="each-input">
          <a-input type="number" placeholder="slots 止（0-16383）" @change="inputSlotsDelEnd" />
        </div>
        <a-popconfirm @confirm="confirmSlotsDel" :title="'确认进行 slots 删除'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button style="width: 30%;float: left;" class="submit" type="primary">提交</a-button>
        </a-popconfirm>
      </a-drawer>-->
      <a-button class="eachHandle" @click="slotsMig" type="primary">slots 迁移</a-button>
      <a-drawer
        title="slots 迁移"
        placement="top"
        @close="slotsMigClose"
        :closable="true"
        :visible="slotsMigShow"
        height="400"
      >
        <div class="each-input">
          <a-input placeholder="来源节点" @change="inputSlotsMigSource" />
        </div>
        <div class="each-input">
          <a-input placeholder="目标节点" @change="inputSlotsMigTarget" />
        </div>
        <div class="each-input">
          <a-input type="number" placeholder="slots 起（0-16383）" @change="inputSlotsMigStart" />
        </div>
        <div class="each-input">
          <a-input type="number" placeholder="slots 止（0-16383）" @change="inputSlotsMigEnd" />
        </div>
        <a-popconfirm @confirm="confirmSlotsMig" :title="'确认进行 slots 迁移'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button style="width: 30%;float: left;" class="submit" type="primary">提交</a-button>
        </a-popconfirm>
      </a-drawer>
    </div>
  </div>
</template>
<script>
let data = [];
let t = null;
export default {
  name: "setting_redis",
  data() {
    const that = this;
    this.$socket.sendObj({
      Func: "/config/redis/clusterNodes",
      Data: JSON.stringify({ id: that.$route.query.id })
    });
    const handMessage = function(da) {
      const d = JSON.parse(da.data);
      let z = [];
      if (d.Type === "/config/redis/clusterNodes") {
        data.push(
          "// 节点信息 <id> <ip:port> <role> <follow-node-id> <ping-sent> <pong-recv> <config-epoch> <link-state> <slots> ..."
        );
      }
      if (d.Type === "/config/redis/clusterMeet") {
        data.push("// 节点添加");
      }
      if (d.Type === "/config/redis/clusterForget") {
        data.push("// 节点删除");
      }
      if (d.Type === "/config/redis/setSlots") {
        data.push("// slots 添加");
      }
      if (d.Type === "/config/redis/delSlots") {
        data.push("// slots 删除");
      }
      if (d.Type === "/config/redis/slots/migrating") {
        data.push("// slots 迁移");
      }
      if (d.Type === "/config/redis/clusterSlots") {
        const max = 16384;
        let availableSlots = [];
        let usedSlots = [];
        data.push("// slots 统计");
        let used = 0;
        const tmpArray = [];
        for (let i = 0; i < d.Data.length; i++) {
          const t = d.Data[i];
          usedSlots.push(t.Start, t.End);
          used += t.End - t.Start + 1;
          tmpArray.push(
            `节点：${t.Nodes[0].Id} (${t.Nodes[0].Addr})  拥有 slots: ${t.Start} - ${t.End}`
          );
        }
        usedSlots.sort((a, b) => {
          return a - b;
        });
        if (usedSlots[0] > 0) {
          availableSlots = [0, usedSlots[0] - 1];
        }
        for (let i = 1; i < usedSlots.length + 1; i += 2) {
          availableSlots.push(usedSlots[i] + 1);
          if (usedSlots[i + 1]) {
            availableSlots.push(usedSlots[i + 1] - 1);
          }
        }
        if (availableSlots[availableSlots.length - 1] < 16383) {
          availableSlots.push(16383);
        }
        for (let i = 0; i < availableSlots.length - 1; i += 2) {
          if (availableSlots[i] > availableSlots[i + 1]) continue;
          tmpArray.push(
            `slots: ${availableSlots[i]} - ${availableSlots[i + 1]} 待分配`
          );
        }
        tmpArray.push(
          "共有：" +
            used +
            " 个 slots 被占用，剩余需分配 slots 总数：" +
            (max - used)
        );
        d.Data = tmpArray.join("\n");
      }
      if (typeof d.Data !== "string") {
        return;
      }
      z = d.Data.split("\n");
      data.push(
        ...z
          .filter(e => {
            if (e.replace(/\r|\n/g, "") !== "") {
              return e;
            }
          })
          .map(e => {
            return e;
          })
      );
      while (data.length > 200) {
        data.shift();
      }
      that.data = data.join("\n");
      const container = that.$el.querySelector("#container");
      setTimeout(() => {
        container.scrollTop += 1000;
      }, 100);
    };
    this.$socket.onmessage = handMessage;
    t = setInterval(() => {
      if (that.$socket.onmessage !== handMessage) {
        that.$socket.onmessage = handMessage;
      }
    }, 1000);
    return {
      data,
      clusterMeetShow: false,
      newhost: "",
      newport: "",
      clusterForgetShow: false,
      nodeid: "",
      clusterReplicateShow: false,
      repHost: "",
      repPort: "",
      repNodeID: "",
      slotsSetShow: false,
      slotsHost: "",
      slotsPort: "",
      slotsStart: "",
      slotsEnd: "",
      slotsMigShow: false,
      slotsMigSource: "",
      slotsMigTarget: "",
      slotsMigStart: "",
      slotsMigEnd: ""
    };
  },
  methods: {
    // slots add
    slotsSet() {
      this.slotsSetShow = true;
    },
    slotsSetClose() {
      this.slotsSetShow = false;
    },
    inputSlotsHost(e) {
      this.slotsHost = e.target.value;
    },
    inputSlotsPort(e) {
      this.slotsPort = e.target.value;
    },
    inputSlotsStart(e) {
      this.slotsStart = e.target.value;
    },
    inputSlotsEnd(e) {
      this.slotsEnd = e.target.value;
    },
    confirmSlotsSet() {
      this.$socket.sendObj({
        Func: "/config/redis/setSlots",
        Data: JSON.stringify({
          id: this.$route.query.id,
          host: this.slotsHost,
          port: this.slotsPort,
          start: Number(this.slotsStart),
          end: Number(this.slotsEnd)
        })
      });
      this.slotsSetShow = false;
    },
    // slots Mig
    slotsMig() {
      this.slotsMigShow = true;
    },
    slotsMigClose() {
      this.slotsMigShow = false;
    },
    inputSlotsMigSource(e) {
      this.slotsMigSource = e.target.value;
    },
    inputSlotsMigTarget(e) {
      this.slotsMigTarget = e.target.value;
    },
    inputSlotsMigStart(e) {
      this.slotsMigStart = e.target.value;
    },
    inputSlotsMigEnd(e) {
      this.slotsMigEnd = e.target.value;
    },
    confirmSlotsMig() {
      this.$socket.sendObj({
        Func: "/config/redis/slots/migrating",
        Data: JSON.stringify({
          id: this.$route.query.id,
          sourceId: this.slotsMigSource,
          TargetID: this.slotsMigTarget,
          slotsStart: Number(this.slotsMigStart),
          slotsEnd: Number(this.slotsMigEnd)
        })
      });
      this.slotsMigShow = false;
    },
    // slots stats
    slotsStats() {
      this.$socket.sendObj({
        Func: "/config/redis/clusterSlots",
        Data: JSON.stringify({
          id: this.$route.query.id
        })
      });
    },
    // cluster
    confirmClusterReplicate() {
      this.$socket.sendObj({
        Func: "/config/redis/clusterReplicate",
        Data: JSON.stringify({
          id: this.$route.query.id,
          host: this.repHost,
          port: this.repPort,
          nodeid: this.repNodeID
        })
      });
      this.clusterReplicateShow = false;
    },
    confirmClusterForget() {
      this.$socket.sendObj({
        Func: "/config/redis/clusterForget",
        Data: JSON.stringify({
          id: this.$route.query.id,
          nodeid: this.nodeid
        })
      });
      this.clusterForgetShow = false;
    },
    confirmAddNode() {
      this.$socket.sendObj({
        Func: "/config/redis/clusterMeet",
        Data: JSON.stringify({
          id: this.$route.query.id,
          host: this.newhost,
          port: this.newport
        })
      });
      this.clusterMeetShow = false;
    },
    inputRepNodeID(e) {
      this.repNodeID = e.target.value;
    },
    inputRepHost(e) {
      this.repHost = e.target.value;
    },
    inputRepPort(e) {
      this.repPort = e.target.value;
    },
    inputNodeID(e) {
      this.nodeid = e.target.value;
    },
    inputHost(e) {
      this.newhost = e.target.value;
    },
    inputPort(e) {
      this.newport = e.target.value;
    },
    clusterForgetClose() {
      this.clusterForgetShow = false;
    },
    clusterMeetClose() {
      this.clusterMeetShow = false;
    },
    clusterReplicateClose() {
      this.clusterReplicateShow = false;
    },
    clusterReplicate() {
      this.clusterReplicateShow = true;
    },
    clusterForget() {
      this.clusterForgetShow = true;
    },
    clusterMeet() {
      this.clusterMeetShow = true;
    },
    clearOutPut() {
      data = [];
      this.data = "";
    },
    reloadClusterNodes() {
      this.$socket.sendObj({
        Func: "/config/redis/clusterNodes",
        Data: JSON.stringify({ id: this.$route.query.id })
      });
    }
  },
  beforeDestroy() {
    data = [];
    this.data = "";
    if (t !== null) {
      window.clearInterval(t);
    }
  }
};
</script>
<style lang="stylus" scoped>
.handleButton {
  width: 24%;
  height: 78vh;
  overflow: auto;
  float: right;
}

.main-box {
  box-sizing: border-box;
  width: 100%;
  height: 78vh;
}

.each-input {
  width: 100%;
  margin-bottom: 20px;
  float: left;
}

.eachHandle {
  float: left;
  margin-right: 10px;
  margin-bottom: 10px;
}

#container {
  box-sizing: border-box;
  float: left;
  width: 74%;
  z-index: 1;
  height: 100%;
  background-color: #2b2b2b;
  box-sizing: border-box;
  margin: 0 0;
}

code {
  wdith: 100%;
  text-align: left;
  box-sizing: border-box;
  margin: 0 0;
  float: left;
  padding: 0 0;
}
</style>
