<template>
  <!-- <router-link to="/redis"> -->
  <div>
    <a-card @click="jumpTo('redis')" hoverable class="each-card">
      <img
        alt="redis"
        src="../../public/image/redis.png"
        slot="cover"
        style="width: 200px; margin-left: auto; margin-right: auto"
      />
      <a-card-meta
        title="redis"
        description="redis 集群/单机 管理"
      ></a-card-meta>
    </a-card>
  </div>
  <!-- </router-link> -->
</template>

<script>
export default {
  name: "setting",
  data() {
    return {};
  },
  methods: {
    jumpTo(url) {
      this.$router.push(url);
    },
  },
};
</script>
<style lang="stylus" scoped>
.each-card {
  float: left;
  width: 300px;
  margin-right: 10px;
  height: 300px;
}
</style>
