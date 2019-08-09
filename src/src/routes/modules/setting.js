const RedisOpcap = {
  path: '/opcap',
  name: 'setting_node_opcap',
  component: () => import('@v/redis/opcap.vue'),
  meta: {
    title: 'redis-opcap',
    Index: '1'
  }
};
const RedisDev = {
  path: '/redis_dev',
  name: 'setting_redis_dev',
  component: () => import('@v/redis/dev.vue'),
  meta: {
    title: '数据源-运维模式',
    Index: '1'
  }
};
const RedisMonit = {
  path: '/clusterSlots',
  name: 'setting_redis_monit',
  component: () => import('@v/redis/clusterSlots.vue'),
  meta: {
    title: '数据源-节点列表',
    Index: '1'
  }
};
const RedisMonitMain = {
  path: '/redis_monit_main',
  name: 'setting_redis_monit_main',
  component: () => import('@v/redis/monit_main.vue'),
  meta: {
    title: '数据源-实时监控',
    Index: '1'
  }
};
const RedisPage = {
  path: '/redis',
  name: 'setting_redis',
  component: () => import('@v/redis/redis.vue'),
  meta: {
    title: '数据源-redis-列表',
    Index: '1'
  }
};
const SettingPage = {
  path: '',
  name: 'setting',
  component: () => import('@v/index.vue'),
  meta: {
    title: '数据源',
    Index: '1'
  }
};
export {
  SettingPage,
  RedisPage,
  RedisMonit,
  RedisMonitMain,
  RedisDev,
  RedisOpcap
};
