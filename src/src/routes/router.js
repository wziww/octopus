import { Error404 } from './base';
import { CasePage } from './modules/case';
import common from '../common/common';
import { System } from './modules/system';
import { SettingPage, RedisPage, RedisMonit, RedisMonitMain, RedisSlowLog } from './modules/setting';
export default [
  {
    path: '/',
    alias: '',
    component: common,
    children: [
      CasePage,
      SettingPage,
      RedisPage,
      RedisMonit,
      RedisMonitMain,
      RedisSlowLog,
      System,
      Error404
    ]
  }
];
