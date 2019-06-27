import Vue from 'vue';
import Router from 'vue-router';
import allRouter from './router';
import { setTitle } from '@/tools/utils';
Vue.use(Router);

const router = new Router({
  mode: 'history',
  routes: [
    ...allRouter
  ]
});
router.afterEach(route => {
  setTitle(route.meta.title);
});
export default router;
