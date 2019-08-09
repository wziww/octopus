import { Error404 } from './base';
import common from '../common/common';
import { SettingPage, RedisPage, RedisMonit, RedisMonitMain, RedisDev, RedisOpcap } from './modules/setting';
export default [
  {
    path: '/login',
    name: 'login',
    component: () => import('@v/login/login.vue'),
    meta: {
      title: 'octopus'
    }
  },
  {
    path: '/',
    component: common,
    children: [
      SettingPage,
      RedisPage,
      RedisMonit,
      RedisMonitMain,
      RedisDev,
      RedisOpcap,
      Error404
    ]
  }
];
