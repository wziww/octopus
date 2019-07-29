<template>
  <div>
    <a-list :grid="{ column: 4 }" :dataSource="Data" style="float: left;width: 100%;">
      <a-list-item slot="renderItem" slot-scope="item">
        <a-card :title="'Pod:' +item.name">
          <div class="each">
            <span style="float: left;">名称:{{item.data.Name}}</span>
          </div>
          <div class="each">
            <span style="float: left;">模式:{{item.data.Type}}</span>
          </div>
          <div class="each">
            <span style="float: left;">Host:{{item.data.Addrs}}</span>
          </div>
          <div class="each" style="border-bottom: 1px solid light; padding-bottom: 10px;">
            <span style="float: left;">Status:</span>
            <a-tag
              :color="item.data.Status==='ok'?'green':'red'"
              style="float: left;"
            >{{item.data.Status}}</a-tag>
          </div>
          <router-link :to="'/clusterSlots?id='+item.name" class="hd">
            <a-icon type="table" />
            <span>节点</span>
          </router-link>
          <router-link :to="'/redis_monit_main?id='+item.name" class="hd">
            <a-icon type="dashboard" />
            <span>监控</span>
          </router-link>
          <router-link :to="'/redis_dev?id='+item.name" class="hd">
            <a-icon type="code" />
            <span>dev</span>
          </router-link>
          <a-icon type="delete" class="hd" @click="del(item.name)" />
        </a-card>
      </a-list-item>
    </a-list>
  </div>
</template>
<script>
import Vue from "vue";
import hd from "../../lib/ws";
const vm = new Vue();
let Data = [];
let t = null;
const PATH = "monit";
export default {
  name: "setting_redis",
  data() {
    Data = [];
    vm.$connect(
      "ws://0.0.0.0:8081/v1/websocket?octopusPath=" +
        PATH +
        "&octopusToken=462426262a462a4a297c726f6f74"+
        "&octopusClusterID=nil",
      { format: "json" }
    );
    try {
      this.$socket.sendObj({
        Func: "/redis"
      });
    } catch (error) {
      console.error(error);
    }
    t = setInterval(() => {
      try {
        this.$socket.sendObj({
          Func: "/redis"
        });
      } catch (error) {
        console.error(error);
      }
    }, 1000 * 3);
    const that = this;
    this.$socket.onmessage = hd(d => {
      Data = [];
      switch (d.Type) {
        case "/redis": // 配置列表
          if (Object.keys(d.Data)) {
            for (let i of Object.keys(d.Data)) {
              Data.push({
                name: i,
                data: d.Data[i]
              });
            }
            that.Data = Data;
          }
          break;
      }
    });
    return {
      form: this.$form.createForm(this),
      visible: false,
      Data
    };
  },
  beforeDestroy() {
    this.Data = [];
    if (t !== null) {
      window.clearInterval(t);
    }
    vm.$disconnect();
  },
  methods: {
    del(id) {
      this.$socket.sendObj({
        Func: PATH + "/del",
        Data: JSON.stringify({ id })
      });
      this.$socket.sendObj({
        Func: PATH
      });
    }
  }
};
</script>
<style lang="stylus" scoped>
.hd {
  float: left;
  margin-right: 10px;
  cursor: pointer;
}

.each {
  width: 100%;
  margin-bottom: 10px;
  float: left;
}
</style>
