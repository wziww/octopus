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
      <a-tooltip placement="topLeft" title="id ip:端口[@总线端口] 角色 - " style="float: left;">
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
      <a-button class="eachHandle" type="primary">未分配 slots 计算</a-button>
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
    this.$socket.onmessage = da => {
      const d = JSON.parse(da.data);
      if (typeof d.Data !== "string") {
        return;
      }
      let z = [];
      z = d.Data.split("\n");
      if (d.Type === "/config/redis/clusterNodes") {
        data.push("// 节点信息");
      }
      if (d.Type === "/config/redis/clusterMeet") {
        data.push("// 节点添加");
      }
      if (d.Type === "/config/redis/clusterForget") {
        data.push("// 节点删除");
      }
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
      const container = this.$el.querySelector("#container");
      setTimeout(() => {
        container.scrollTop += 1000;
      }, 100);
    };
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
      repNodeID: ""
    };
  },
  methods: {
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
  width: 100%;
  height: 20vh;
  overflow: auto;
}

.main-box {
  box-sizing: border-box;
  width: 100%;
  height: 60vh;
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
  width: 100%;
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
