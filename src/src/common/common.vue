<template>
  <a-layout id="components-layout-demo-custom-trigger">
    <a-layout-sider
      :trigger="null"
      collapsible
      :style="{ height: '100vh', overflow: 'auto', position: 'fixed', left: 0 }"
    >
      <div class="logo">
        <img :style="{ width: 'auto', height: '100%' }" />
      </div>
      <a-menu
        theme="dark"
        mode="vertical"
        :defaultSelectedKeys="[$route.meta.Index]"
      >
        <a-menu-item key="1">
          <router-link to="/">
            <a-icon type="cloud-o" />
            <span>数据源</span>
          </router-link>
        </a-menu-item>
        <a-menu-item
          selectable="false"
          @click="
            clear();
            logout();
          "
        >
          <span>登出</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout
      style="box-sizing: border-box"
      :style="{ marginLeft: '200px', minHeight: '100vh' }"
    >
      <a-layout-content :style="{ margin: '24px 16px 0' }">
        <div :style="{ padding: '24px' }">
          <router-view></router-view>
        </div>
      </a-layout-content>
      <a-layout-footer style="textalign: center"
        >octopus v0.0.2</a-layout-footer
      >
    </a-layout>
  </a-layout>
</template>

<script>
import { clear } from "../lib/token";
export default {
  name: "common",
  data() {
    return {
      clear,
      // collapsed: false
    };
  },
  // collapsed() {
  //   this.collapsed = !this.collapsed;
  // },
  methods: {
    logout() {
      this.$router.push({
        path: "/login",
      });
    },
  },
};
</script>
