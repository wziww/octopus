<template>
  <div>
    <a-button type="primary" @click="showDrawer" style="float: left;margin-bottom: 10px;">
      <a-icon type="plus"/>新增配置
    </a-button>
    <a-drawer
      title="添加监控"
      :width="720"
      @close="onClose"
      :visible="visible"
      :wrapStyle="{height: 'calc(100% - 108px)',overflow: 'auto',paddingBottom: '108px'}"
    >
      <a-form :form="form" layout="vertical" hideRequiredMark>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Name">
              <a-input
                v-decorator="['name', {
                  rules: [{ required: true, message: 'Please enter name' }]
                }]"
                placeholder="Please enter user name"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Url(多个 ; 隔开)">
              <a-input
                v-decorator="['url', {
                  rules: [{ required: true, message: 'please enter url' }]
                }]"
                style="width: 100%"
                placeholder="please enter url"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="模式属性">
              <a-select
                v-decorator="['type', {
                  rules: [{ required: true, message: 'Please choose the instance type' }]
                }]"
                placeholder="Please choose the instance type"
              >
                <a-select-option value="cluster">集群模式(cluster)</a-select-option>
                <!-- <a-select-option value="single">单机模式</a-select-option> -->
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <div
        :style="{
          position: 'absolute',
          left: 0,
          bottom: 0,
          width: '100%',
          borderTop: '1px solid #e9e9e9',
          padding: '10px 16px',
          background: '#fff',
          textAlign: 'right',
        }"
      >
        <a-button :style="{marginRight: '8px'}" @click="onClose">Cancel</a-button>
        <a-button @click="onSubmit" type="primary">Submit</a-button>
      </div>
    </a-drawer>
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
            <a-tag color="green" style="float: left;">{{item.data.Status}}</a-tag>
          </div>
          <router-link :to="'/setting/redis_monit?id='+item.name">
            <a-icon type="align-left" class="hd"/>
          </router-link>
          <router-link :to="'/setting/redis_monit_main?id='+item.name">
            <a-icon type="dashboard" class="hd"/>
          </router-link>
          <a-icon type="delete" class="hd" @click="del(item.name)"/>
        </a-card>
      </a-list-item>
    </a-list>
  </div>
</template>
<script>
let Data = [];
let t = null;
export default {
  name: "setting_redis",
  data() {
    try {
      this.$socket.sendObj({
        Func: "/config/redis"
      });
    } catch (error) {
      console.error(error);
    }
    t = setInterval(() => {
      try {
        this.$socket.sendObj({
          Func: "/config/redis"
        });
      } catch (error) {
        console.error(error);
      }
    }, 1000 * 3);
    const that = this;
    this.$socket.onmessage = data => {
      const d = JSON.parse(data.data);
      console.log(d);
      Data = [];
      switch (d.Type) {
        case "/config/redis": // 配置列表
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
    };
    return {
      form: this.$form.createForm(this),
      visible: false,
      Data
    };
  },
  beforeDestroy() {
    if (t !== null) {
      window.clearInterval(t);
    }
  },
  methods: {
    showDrawer() {
      this.visible = true;
    },
    onClose() {
      this.visible = false;
    },
    onSubmit() {
      let obj = this.form.getFieldsValue();
      this.$socket.sendObj({
        Func: "/config/redis/add",
        Data: JSON.stringify(obj)
      });
      this.visible = false;
      this.$socket.sendObj({
        Func: "/config/redis"
      });
    },
    del(id) {
      this.$socket.sendObj({
        Func: "/config/redis/del",
        Data: JSON.stringify({ id })
      });
      this.$socket.sendObj({
        Func: "/config/redis"
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
