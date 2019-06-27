<template>
  <a-card hoverable style="width: 300px">
    <img
      alt="redis"
      src="../../../public/image/redis.png"
      slot="cover"
      style="width: 200px;margin-left: auto;margin-right: auto;"
    >
    <template class="ant-card-actions" slot="actions">
      <router-link to="/setting/redis">
        <a-icon type="setting"/>
      </router-link>
      <a-icon type="edit"/>
      <a-icon type="ellipsis"/>
    </template>
    <a-card-meta title="redis" description="redis 集群/单机 配置"></a-card-meta>
  </a-card>
</template>

<script>
export default {
  name: "setting",

  data() {
    return {};
  },

  components: {},

  computed: {},

  methods: {},

  mounted() {}
};
</script>
<style lang="stylus" scoped></style>
