<template>
  <div id="main-box" class="main-box">
    <pre id="container" v-highlightjs="data">
    <code class="string">
    </code>
    </pre>
    <div class="handleButton">
      <a-divider orientation="left">通用</a-divider>
      <a-button class="eachHandle" @click="clearOutPut" type="primary"
        >清屏</a-button
      >
      <router-link
        class="eachHandle"
        :to="'/clusterSlots?id=' + $route.query.id"
      >
        <a-button type="primary">集群列表</a-button>
      </router-link>
      <a-divider orientation="left">节点操作</a-divider>
      <a-tooltip
        placement="topLeft"
        title="<id> <ip:port> <flags> <master> <ping-sent> <pong-recv> <config-epoch> <link-state> <slot> <slot> ... <slot>"
        style="float: left"
      >
        <a-button class="eachHandle" @click="reloadClusterNodes" type="primary"
          >刷新节点信息</a-button
        >
      </a-tooltip>
      <a-button class="eachHandle" @click="clusterMeet" type="primary"
        >添加节点</a-button
      >
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
        <a-popconfirm
          @confirm="confirmAddNode"
          :title="'确认添加节点[' + newhost + ':' + newport + ']进入集群么？'"
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
      <a-button class="eachHandle" @click="clusterForget" type="primary"
        >删除节点</a-button
      >
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
        <a-popconfirm
          @confirm="confirmClusterForget"
          :title="'确认删除节点[' + nodeid + ']么？'"
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
      <a-button class="eachHandle" @click="clusterReplicate" type="primary"
        >从节点分配</a-button
      >
      <a-drawer
        title="从节点分配"
        placement="top"
        @close="clusterReplicateClose"
        :closable="true"
        :visible="clusterReplicateShow"
        height="300"
      >
        <div class="each-input">
          <a-input placeholder="需要分配的节点 host" @change="inputRepHost" />
        </div>
        <div class="each-input">
          <a-input placeholder="需要分配的节点 port" @change="inputRepPort" />
        </div>
        <div class="each-input">
          <a-input placeholder="主节点 nodeid" @change="inputRepNodeID" />
        </div>
        <a-popconfirm
          @confirm="confirmClusterReplicate"
          :title="'确认设置从节点么？'"
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
      <a-button class="eachHandle" @click="safeDel" type="primary"
        >KEY 安全删除操作</a-button
      >
      <a-drawer
        title="KEY 安全删除操作"
        placement="top"
        @close="safeDelClose"
        :closable="true"
        :visible="safeDelShow"
      >
        <div class="each-input">
          <a-input placeholder="address" @change="inputSafeDelKeyAddress" />
        </div>
        <div class="each-input">
          <a-input placeholder="key" @change="inputSafeDelKey" />
        </div>
        <a-popconfirm
          @confirm="confirmSafeDel"
          :title="
            '确认进行[' +
            safeDelAddress +
            '节点的 key:' +
            safeDelKey +
            ']删除操作？'
          "
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>

      <a-button class="eachHandle" @click="nodeDebug" type="primary"
        >节点内存分析</a-button
      >
      <a-drawer
        title="节点内存分析"
        placement="top"
        @close="nodeDebugClose"
        :closable="true"
        :visible="nodeDebugShow"
      >
        <div class="each-input">
          <a-input placeholder="host" @change="inputDebugHost" />
        </div>
        <div class="each-input">
          <a-input placeholder="port" @change="inputDebugPort" />
        </div>
        <a-popconfirm
          @confirm="confirmDebugNode"
          :title="'确认查看节点[' + debugHost + ':' + debugPort + ']内存信息？'"
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
      <a-divider orientation="left">rdb 操作</a-divider>
      <a-button class="eachHandle" @click="rdbLS" type="primary"
        >rdb ls</a-button
      >
      <a-button class="eachHandle" @click="rdbOpen" type="primary"
        >rdb 文件分析</a-button
      >
      <a-drawer
        title="rdb 文件分析"
        placement="top"
        @close="rdbClose"
        :closable="true"
        :visible="rdbShow"
        height="350"
      >
        <div class="each-input">
          <a-input placeholder="文件名" @change="inputRdbFileName" />
        </div>
        <div class="each-input">
          <a-input placeholder="统计数量" @change="inputCount" />
        </div>
        <div class="each-input">
          <a-input placeholder="大小偏移量限制" @change="inputOffset" />
        </div>
        <div class="each-input">
          <a-input placeholder="子级数量限制" @change="inputChildSize" />
        </div>
        <a-popconfirm
          @confirm="confirmRdb"
          :title="'确认对文件[' + rdbFileName + ']进行分析？'"
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
      <a-upload
        class="eachHandle"
        name="file"
        :multiple="true"
        :action="uploadAddress"
        :headers="headers"
        @change="handleChange"
      >
        <a-button> <a-icon type="upload" />点击上传文件</a-button>
      </a-upload>

      <a-divider orientation="left">slots 操作</a-divider>

      <a-button class="eachHandle" type="primary" @click="slotsStats"
        >未分配 slots 计算</a-button
      >
      <a-button class="eachHandle" @click="slotsSet" type="primary"
        >slots 分配</a-button
      >
      <a-drawer
        title="slots 分配"
        placement="top"
        @close="slotsSetClose"
        :closable="true"
        :visible="slotsSetShow"
        height="400"
      >
        <div class="each-input">
          <a-input
            placeholder="需要设置为从节点的 host"
            @change="inputSlotsHost"
          />
        </div>
        <div class="each-input">
          <a-input
            placeholder="需要设置为从节点的 port"
            @change="inputSlotsPort"
          />
        </div>
        <div class="each-input">
          <a-input
            type="number"
            placeholder="slots 起（0-16383）"
            @change="inputSlotsStart"
          />
        </div>
        <div class="each-input">
          <a-input
            type="number"
            placeholder="slots 止（0-16383）"
            @change="inputSlotsEnd"
          />
        </div>
        <a-popconfirm
          @confirm="confirmSlotsSet"
          :title="'确认进行 slots 分配？'"
        >
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
      <a-button class="eachHandle" @click="slotsMig" type="primary"
        >slots 迁移</a-button
      >
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
          <a-input
            type="number"
            placeholder="slots 起（0-16383）"
            @change="inputSlotsMigStart"
          />
        </div>
        <div class="each-input">
          <a-input
            type="number"
            placeholder="slots 止（0-16383）"
            @change="inputSlotsMigEnd"
          />
        </div>
        <a-popconfirm @confirm="confirmSlotsMig" :title="'确认进行 slots 迁移'">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-button
            style="width: 30%; float: left"
            class="submit"
            type="primary"
            >提交</a-button
          >
        </a-popconfirm>
      </a-drawer>
    </div>
  </div>
</template>
<script>
import hd from "../../lib/ws";
import { token } from "../../lib/token";
import WS from "../../lib/websocket";
import config from "../../config/index";
const PATH = "dev";
let ws = null;
let data = [];
let t = null;
let lastTime = new Date().getTime();
export default {
  name: "setting_redis",
  data() {
    ws = new WS(
      config.Host +
        "?op=" +
        PATH +
        "&ot=" +
        token +
        "&ocid=" +
        this.$route.query.id
    );
    const that = this;
    ws.OnOpen(() => {
      ws.SendObj({
        Func: "namespace",
        Data: JSON.stringify({
          namespace: `dev-${that.$route.query.id}`,
        }),
      });
    });
    const handMessage = hd((d) => {
      // 接受服务端数据
      let z = [];
      if (["token", "namespace"].indexOf(d.Type) !== -1) {
        return;
      }
      if (d.Type === "/redis/clusterNodes") {
        data.push(
          "// 节点信息 <id> <ip:port> <role> <follow-node-id> <ping-sent> <pong-recv> <config-epoch> <link-state> <slots> ..."
        );
      }
      if (d.Type === "/redis/clusterMeet") {
        data.push("// 节点添加");
      }
      if (d.Type === "/redis/clusterForget") {
        data.push("// 节点删除");
      }
      if (d.Type === "/redis/setSlots") {
        data.push("// slots 添加");
      }
      if (d.Type === "/redis/delSlots") {
        data.push("// slots 删除");
      }
      if (d.Type === "/redis/slots/migrating") {
        data.push("// slots 迁移");
      }
      if (d.Type === "/redis/rdb/analyze/0") {
        const percent = Number(d.Data);
        let str = "rdb 分析中:\n[";
        const starsNum = parseInt((percent * 50) / 100);
        for (let i = 0; i < starsNum; i++) {
          str += "*";
        }
        for (let i = 0; i < 50 - starsNum; i++) {
          str += "_";
        }
        str += "] " + percent + " %";
        data = [str];
        str += "\n";
        d.Data = str;
        if (lastTime + 100 < new Date().getTime() || status[1] === status[0]) {
          lastTime = new Date().getTime();
          that.data = d.Data;
        }
        return;
      }
      if (d.Type === "/redis/rdb/analyze" && typeof d.Data !== "string") {
        let data = [];
        data.push("rdb 分析结果:\n");
        data.push("总键值对:" + d.Data.TotalNums + "\n");
        data.push("过期键总数:" + d.Data.Expires + "\n");
        data.push("已过期键总数:" + d.Data.AlreadyExpired + "\n");
        data.push("缓存脚本数:" + d.Data.LuaNums + "\n");
        data.push("大键 top " + d.Data.Count + " :(存储空间)\n");
        for (let v of d.Data.OffSetLog) {
          data.push(`key: ${v.Key}\nval: ${v.Val}\n`);
        }
        data.push("大键 top " + d.Data.Count + " :(子级数量)\n");
        for (let v of d.Data.ChildLog) {
          data.push(`key: ${v.Key}\nval: ${v.Val}\n`);
        }
        console.log(d.Data);
        that.data = data.join("\n");
        return;
      }
      if (d.Type === "/redis/slots/migrating/0") {
        const status = d.Data.split(" ");
        status[0]++;
        status[1]++;
        let str = "slots 迁移中:\n";
        str +=
          "自 " +
          status[2] +
          " 迁移 " +
          status[4] +
          " - " +
          status[5] +
          " slots 去 " +
          status[3] +
          "\n[";
        for (let i = 0; i < (status[1] / status[0]) * 50; i++) {
          str += "*";
        }
        for (let i = 0; i < 50 - (status[1] / status[0]) * 50; i++) {
          str += "_";
        }
        str += "] " + ((status[1] / status[0]) * 100).toFixed(2) + " %";
        data = [str];
        str += "\n";
        d.Data = str;
        if (lastTime + 100 < new Date().getTime() || status[1] === status[0]) {
          lastTime = new Date().getTime();
          that.data = d.Data;
        }
        return;
      }
      if (d.Type === "/redis/clusterSlots") {
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
      if (d.Type === "/redis/debug/htstats") {
        data.push("// debug htstats 0 参数说明:");
        data.push("// [Dictionary HT] : 数据字典哈希表");
        data.push("// [Expires HT]    : 过期 KV 字典哈希表");
        data.push("// 每张哈希表至多附带两个子表，编号为 0，1");
      }
      if (d.Type === "/redis/rdb/ls") {
        let data = ["rdb 文件列表:"];
        if (d.Data) {
          for (let v of d.Data) {
            data.push(v);
          }
        }
        that.data = data.join("\n");
        return;
      }
      if (d.Type === "/redis/safe/del") {
        data.push("// key 删除操作");
      }
      if (typeof d.Data !== "string") {
        return;
      }

      z = d.Data.split("\n");
      data.push(
        ...z
          .filter((e) => {
            if (e.replace(/\r|\n/g, "") !== "") {
              return e;
            }
          })
          .map((e) => {
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
    });
    ws.Open();
    ws.OnData(handMessage);
    return {
      uploadAddress: config.Upload,
      // rdb
      rdbFileName: "",
      rdbCount: "",
      rdbOffset: "",
      childSize: "",
      // rdb
      data,
      clusterMeetShow: false,
      nodeDebugShow: false,
      safeDelKey: "",
      safeDelAddress: "",
      safeDelShow: false,
      rdbShow: false,
      debugHost: "",
      debugPort: "",
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
      slotsMigEnd: "",
    };
  },
  methods: {
    handleChange(info) {
      if (info.file.status !== "uploading") {
        console.log(info.file, info.fileList);
      }
      if (info.file.status === "done") {
        this.$message.success(`${info.file.name} file uploaded successfully`);
      } else if (info.file.status === "error") {
        this.$message.error(`${info.file.name} file upload failed.`);
      }
    },
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
      ws.SendObj({
        Func: "/redis/setSlots",
        Data: JSON.stringify({
          id: this.$route.query.id,
          host: this.slotsHost,
          port: this.slotsPort,
          start: Number(this.slotsStart),
          end: Number(this.slotsEnd),
        }),
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
      ws.SendObj({
        Func: "/redis/slots/migrating",
        Data: JSON.stringify({
          id: this.$route.query.id,
          sourceId: this.slotsMigSource,
          TargetID: this.slotsMigTarget,
          slotsStart: Number(this.slotsMigStart),
          slotsEnd: Number(this.slotsMigEnd),
        }),
      });
      this.slotsMigShow = false;
    },
    rdbLS() {
      ws.SendObj({
        Func: "/redis/rdb/ls",
      });
    },
    // slots stats
    slotsStats() {
      ws.SendObj({
        Func: "/redis/clusterSlots",
        Data: JSON.stringify({
          id: this.$route.query.id,
        }),
      });
    },
    confirmRdb() {
      ws.SendObj({
        Func: "/redis/rdb/analyze",
        Data: JSON.stringify({
          filename: this.rdbFileName,
          count: Number(this.rdbCount),
          offsetSize: Number(this.rdbOffset),
          childSize: Number(this.childSize),
        }),
      });
      this.rdbShow = false;
    },
    // debug htstats 0
    confirmSafeDel() {
      ws.SendObj({
        Func: "/redis/safe/del",
        Data: JSON.stringify({
          address: this.safeDelAddress,
          key: this.safeDelKey,
        }),
      });
      this.safeDelShow = false;
    },
    // debug htstats 0
    confirmDebugNode() {
      ws.SendObj({
        Func: "/redis/debug/htstats",
        Data: JSON.stringify({
          host: this.debugHost,
          port: this.debugPort,
        }),
      });
      this.nodeDebugShow = false;
    },

    // cluster
    confirmClusterReplicate() {
      ws.SendObj({
        Func: "/redis/clusterReplicate",
        Data: JSON.stringify({
          id: this.$route.query.id,
          host: this.repHost,
          port: this.repPort,
          nodeid: this.repNodeID,
        }),
      });
      this.clusterReplicateShow = false;
    },
    confirmClusterForget() {
      ws.SendObj({
        Func: "/redis/clusterForget",
        Data: JSON.stringify({
          id: this.$route.query.id,
          nodeid: this.nodeid,
        }),
      });
      this.clusterForgetShow = false;
    },
    confirmAddNode() {
      ws.SendObj({
        Func: "/redis/clusterMeet",
        Data: JSON.stringify({
          id: this.$route.query.id,
          host: this.newhost,
          port: this.newport,
        }),
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
    inputSafeDelKeyAddress(e) {
      this.safeDelAddress = e.target.value;
    },
    inputSafeDelKey(e) {
      this.safeDelKey = e.target.value;
    },
    // rdb
    inputRdbFileName(e) {
      this.rdbFileName = e.target.value;
    },
    inputCount(e) {
      this.rdbCount = e.target.value;
    },
    inputOffset(e) {
      this.rdbOffset = e.target.value;
    },
    inputChildSize(e) {
      this.childSize = e.target.value;
    },
    // rdb
    inputDebugHost(e) {
      this.debugHost = e.target.value;
    },
    inputDebugPort(e) {
      this.debugPort = e.target.value;
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
    safeDelClose() {
      this.safeDelShow = false;
    },
    rdbOpen() {
      this.rdbShow = true;
    },
    rdbClose() {
      this.rdbShow = false;
    },
    nodeDebugClose() {
      this.nodeDebugShow = false;
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
    safeDel() {
      this.safeDelShow = true;
    },
    nodeDebug() {
      this.nodeDebugShow = true;
    },
    clusterMeet() {
      this.clusterMeetShow = true;
    },
    clearOutPut() {
      data = [];
      this.data = "";
    },
    reloadClusterNodes() {
      ws.SendObj({
        Func: "/redis/clusterNodes",
        Data: JSON.stringify({ id: this.$route.query.id }),
      });
    },
  },
  beforeDestroy() {
    data = [];
    this.data = "";
    if (t !== null) {
      window.clearInterval(t);
    }
    ws.Close();
  },
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
