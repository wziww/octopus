import Vue from 'vue';
import App from './App.vue';
import router from './routes';
import VCharts from 'v-charts';
// import './registerServiceWorker';
// import plugins from './plugins';
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';
import 'highlight.js/styles/a11y-dark.css';
import VueHighlightJS from 'vue-highlightjs';
import moment from 'vue-moment';
Vue.use(VCharts);
Vue.use(VueHighlightJS);
Vue.use(moment);
Vue.config.productionTip = false;
Vue.use(Antd);
new Vue({
  router,
  render: h => h(App)
}).$mount('#app');
