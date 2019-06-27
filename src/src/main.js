import Vue from 'vue';
import App from './App.vue';
import router from './routes';
import store from './store';
import VCharts from 'v-charts';
// import './registerServiceWorker';
import '@a/css/common.styl';
import '@a/css/animate.styl';
// import plugins from './plugins';
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';
import VueNativeSock from 'vue-native-websocket';
import moment from 'vue-moment';
Vue.use(VCharts);
Vue.use(moment);
Vue.use(VueNativeSock, 'ws://0.0.0.0:8080/v1', {
  format: 'json', reconnection: true
});
Vue.config.productionTip = false;
Vue.use(Antd);
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');
