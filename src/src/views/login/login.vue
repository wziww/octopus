<template>
  <div class="loginWrap fillcontain" :style="{ backgroundImage: `url(${bgImage})`}">
    <vue-canvas-nest class="nest" :config="{color:'64,169,255', count: 99}"></vue-canvas-nest>
    <transition name="form-fade" mode="in-out">
      <section class="form_contianer" v-show="showLogin">
        <div class="manage_tip">
          <p>octopus v0.0.1</p>
        </div>
        <a-form :model="loginForm" :rules="rules" ref="loginForm">
          <a-form-item prop="username">
            <a-input v-model="loginForm.username" placeholder="用户名">
              <span>dsfsf</span>
            </a-input>
          </a-form-item>
          <a-form-item prop="password">
            <a-input type="password" placeholder="密码" v-model="loginForm.password"></a-input>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="submitForm('loginForm')" class="submit_btn">登录</a-button>
          </a-form-item>
        </a-form>
      </section>
    </transition>
    <p class="coda">octopus v0.0.1</p>
  </div>
</template>

<script>
import { mapActions } from "vuex";
import vueCanvasNest from "vue-canvas-nest";
import moment from "moment";
import { message } from "ant-design-vue";
import bg1 from "../../../dist/img/login_bg_1.jpg";
import bg2 from "../../../dist/img/login_bg_2.jpg";
import WS from "../../lib/websocket";
import { TokenSet, PermissionSet } from "../../lib/token";
const PATH = "login";
let ws = new WS(
  "ws://0.0.0.0:8081/v1/websocket?op=" + PATH + "&ot=octopus"
);
export default {
  data() {
    ws.Open();
    return {
      loginForm: {
        username: "",
        password: ""
      },
      rules: {
        username: [
          { required: true, message: "请输入用户名", trigger: "blur" }
        ],
        password: [{ required: true, message: "请输入密码", trigger: "blur" }]
      },
      showLogin: false,
      bgImage: moment().hours() % 2 ? bg2 : bg1
    };
  },
  components: {
    vueCanvasNest
  },
  mounted() {
    this.showLogin = true;
  },
  computed: {},
  methods: {
    ...mapActions(["getAdminData"]),
    async submitForm(formName) {
      const that = this;
      ws.SendObj({
        Func: "/login",
        Data: JSON.stringify({
          username: this.loginForm.username,
          password: this.loginForm.password
        })
      });
      ws.OnData(function(d) {
        const data = JSON.parse(d.data);
        if (data.Type === "/login") {
          const t = JSON.parse(data.Data);
          if (t.token) {
            TokenSet(t.token);
            PermissionSet(t.permission);
            that.$router.push("/");
          } else {
            message.error("用户名或密码错误");
          }
        }
      });
    }
  },
  watch: {
    adminInfo: function(newValue) {}
  },
  beforeDestroy() {
    ws.Close();
  }
};
</script>

<style type="text/less" lang="less" scoped>
@import "../style/mixin";
.nest {
  z-index: 2 !important;
}
.loginWrap {
  background-color: #324057;
  height: 100vh;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  /* .bis('../assets/img/login_bg_1.jpg') */
}
.manage_tip {
  position: absolute;
  width: 100%;
  top: -100px;
  left: 0;
  p {
    font-size: 23px;
    color: #303133;
  }
}
.form_contianer {
  padding: 25px;
  border-radius: 5px;
  text-align: center;
  background-color: rgba(255, 255, 255, 0.5);
  .wh(320px, 210px);
  .ctp(320px, 210px);
}
.submit_btn {
  width: 100%;
  font-size: 16px;
}
.tip {
  font-size: 12px;
  color: red;
}
.form-fade-enter-active,
.form-fade-leave-active {
  transition: all 1s;
}
.form-fade-enter,
.form-fade-leave-active {
  transform: translate3d(0, -50px, 0);
  opacity: 0;
}
.coda {
  position: absolute;
  bottom: 5px;
  left: 50%;
  transform: translateX(-50%);
  color: darkgrey;
}
</style>
