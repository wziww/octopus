const RedisDev = {
  path: '/setting/redis_dev',
  name: 'setting_redis_dev',
  component: () => import('@v/setting/redis/dev.vue'),
  meta: {
    title: '数据源-运维模式',
    Index: '3'
  }
};
const RedisMonit = {
  path: '/setting/clusterSlots',
  name: 'setting_redis_monit',
  component: () => import('@v/setting/redis/clusterSlots.vue'),
  meta: {
    title: '数据源-节点列表',
    Index: '3'
  }
};
const RedisMonitMain = {
  path: '/setting/redis_monit_main',
  name: 'setting_redis_monit_main',
  component: () => import('@v/setting/redis/monit_main.vue'),
  meta: {
    title: '数据源-实时监控',
    Index: '3'
  }
};
const RedisPage = {
  path: '/setting/redis',
  name: 'setting_redis',
  component: () => import('@v/setting/redis/redis.vue'),
  meta: {
    title: '数据源-redis-列表',
    Index: '3'
  }
  // children: [RedisList]
};
const SettingPage = {
  path: '/setting',
  name: 'setting',
  component: () => import('@v/setting/index.vue'),
  meta: {
    title: '数据源',
    Index: '3'
  }
};
export {
  SettingPage,
  RedisPage,
  RedisMonit,
  RedisMonitMain,
  RedisDev
};
