<template>
  <!-- <div> -->
  <ve-pie style="width: 50%;float: left;" :data="data"></ve-pie>
  <!-- </div> -->
</template>
<script>
import hd from "../../lib/ws";
import WS from "../../lib/websocket";
import { token } from "../../lib/token";
import config from "../../config/index";
const PATH = "monit";
const ws = new WS(config.Host + "?op=" + PATH + "&ot=" + token + "&ocid=nil");
let data = {};
let interTime = 1000;
let t = null;
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
        that.data = {
          columns: ["命令", "次数"],
          rows: tmpD
        };
      }
    });
    ws.OnData(handMessage);
    return {
      data
    };
  },
  methods: {},
  beforeDestroy() {
    ws.Close();
    if (t !== null) {
      window.clearInterval(t);
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
