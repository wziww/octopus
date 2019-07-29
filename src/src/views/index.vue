<template>
  <router-link to="/redis">
    <a-card hoverable style="width: 300px">
      <img
        alt="redis"
        src="../../public/image/redis.png"
        slot="cover"
        style="width: 200px;margin-left: auto;margin-right: auto;"
      />
      <a-card-meta title="redis" description="redis 集群/单机 配置"></a-card-meta>
    </a-card>
  </router-link>
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
