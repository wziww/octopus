import Vue from 'vue'
import Vuex from 'vuex'

import app from './modules/app'
import user from './modules/user'

Vue.use(Vuex)
// 数据初始化保存
const initData = {
  app: app.state,
  user: user.state
}
export default new Vuex.Store({
  state: {
  },
  mutations: {
    CLEARSTATE_MUTATE (state) {
      for (let key in state) {
        for (let name in state[key]) {
          state[key][name] = initData[key][name]
        }
      }
    }
  },
  actions: {
    //
  },
  modules: {
    app,
    user
  }
})
