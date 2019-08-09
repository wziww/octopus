<template></template>
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
            Func: "/opcap",
            Data: JSON.stringify({ address: that.$route.query.address })
          });
        } catch (e) {
          console.error(e);
        }
      }, interTime);
    });
    const handMessage = hd(d => {
      console.log(d);
      if (d.Type === "/opcap") {
      }
    });
    ws.OnData(handMessage);
    return {
      data,
    };
  },
  methods: {
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
