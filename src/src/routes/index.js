import Vue from 'vue';
import Router from 'vue-router';
import allRouter from './router';
Vue.use(Router);

const router = new Router({
  mode: 'history',
  routes: [
    ...allRouter
  ]
});
router.afterEach(route => {
});
export default router;
