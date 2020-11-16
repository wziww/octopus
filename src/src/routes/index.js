import Vue from 'vue';
import Router from 'vue-router';
import allRouter from './router';
import { token, } from '../lib/token';
Vue.use(Router);

const router = new Router({
  mode: 'history',
  routes: [
    ...allRouter,
  ],
});
router.afterEach(route => {
  if (!token) {
    router.push({
      path: "/login",
    });
  }
});
export default router;
