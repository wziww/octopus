const RedisMonit = {
  path: '/setting/redis_monit',
  name: 'setting_redis_monit',
  component: () => import('@v/setting/redis/monit.vue'),
  meta: {
    title: '数据源',
    Index: '3'
  }
};
const RedisMonitMain = {
  path: '/setting/redis_monit_main',
  name: 'setting_redis_monit_main',
  component: () => import('@v/setting/redis/monit_main.vue'),
  meta: {
    title: '数据源',
    Index: '3'
  }
};
const RedisPage = {
  path: '/setting/redis',
  name: 'setting_redis',
  component: () => import('@v/setting/redis/redis.vue'),
  meta: {
    title: '数据源',
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
  RedisMonitMain
};
