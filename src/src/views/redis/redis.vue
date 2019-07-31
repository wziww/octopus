<template>
  <div>
    <a-list :grid="{ column: 4 }" :dataSource="Data" style="float: left;width: 100%;">
      <a-list-item slot="renderItem" slot-scope="item">
        <a-card :title="'' +item.data.Name">
          <!-- <div class="each">
            <span style="float: left;">名称:</span>
            <a-tag class="a-tag" color="blue">{{item.data.Name}}</a-tag>
          </div>-->
          <div class="each">
            <span style="float: left;">模式:</span>
            <a-tag class="a-tag" color="blue">{{item.data.Type}}</a-tag>
          </div>
          <div class="each">
            <span style="float: left;">地址:</span>
            <a-tag class="a-tag" color="blue">{{item.data.Addrs}}</a-tag>
          </div>
          <div class="each" style="border-bottom: 1px solid light; padding-bottom: 10px;">
            <span style="float: left;">状态:</span>
            <a-tag class="a-tag" :color="item.data.Status==='ok'?'green':'red'">{{item.data.Status}}</a-tag>
          </div>
          <router-link
            v-if="permission&permissionAll.PERMISSIONMONIT"
            :to="'/clusterSlots?id='+item.name"
            class="hd"
          >
            <a-icon type="table" />
            <span>节点</span>
          </router-link>
          <router-link
            v-if="permission&permissionAll.PERMISSIONMONIT"
            :to="'/redis_monit_main?id='+item.name"
            class="hd"
          >
            <a-icon type="dashboard" />
            <span>监控</span>
          </router-link>
          <router-link
            v-if="permission&permissionAll.PERMISSIONDEV"
            :to="'/redis_dev?id='+item.name"
            class="hd"
          >
            <a-icon type="code" />
            <span>dev</span>
          </router-link>
          <!-- <a-icon type="delete" class="hd" @click="del(item.name)" /> -->
        </a-card>
      </a-list-item>
    </a-list>
  </div>
</template>
<script>
import hd from "../../lib/ws";
import WS from "../../lib/websocket";
import { token, permission, permissionAll } from "../../lib/token";
import config from "../../config/index";
const PATH = "monit";
const ws = new WS(config.Host + "?op=" + PATH + "&ot=" + token + "&ocid=nil");
let Data = [];
let t = null;
export default {
  name: "setting_redis",
  data() {
    Data = [];
    ws.Open();
    const that = this;
    ws.OnData(
      hd(d => {
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
      })
    );
    ws.OnOpen(() => {
      try {
        ws.SendObj({
          Func: "/redis"
        });
      } catch (error) {
        console.error(error);
      }
      t = setInterval(() => {
        try {
          ws.SendObj({
            Func: "/redis"
          });
        } catch (error) {
          console.error(error);
        }
      }, 1000 * 3);
    });
    return {
      form: this.$form.createForm(this),
      visible: false,
      Data,
      permission,
      permissionAll
    };
  },
  beforeDestroy() {
    ws.Close();
    this.Data = [];
    if (t !== null) {
      window.clearInterval(t);
    }
  },
  methods: {}
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

.each > .a-tag {
  float: left;
  margin-left: 30px;
}
</style>
