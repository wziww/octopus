import { getToken, setToken, removeToken } from '@/tools/utils'
import { loginApi } from '@/api'
export default {
  state: {
    token: getToken(),
    userId: '',
    userName: '',
    userHeadImage: '',
    coinBalance: '0'
  },
  mutations: {
    // 设置用户头像
    USER_SETUSERHEADIMAGE_MUTATE (state, headImgPath) {
      state.userHeadImage = headImgPath
    },
    // 设置用户id
    USER_SETUSERID_MUTATE (state, userId) {
      state.userId = userId
    },
    // 设置用户名字
    USER_SETUSERNAME_MUTATE (state, userName) {
      state.userName = userName
    },
    // 设置剩余币数
    USER_SETCOINBALANCE_MUTATE (state, coin) {
      state.coinBalance = coin
    },
    // 设置token
    USER_SETTOKEN_MUTATE (state, token) {
      state.token = token
      if (token) {
        setToken(token)
      } else {
        removeToken()
      }
    }
  },
  actions: {
    // 获取用户信息
    // async USER_GETUSERINFO_ACTION ({ commit, state }) {
    //   let res = await getUserInfoApi()
    //   if (res && res.return_code === 0) {
    //     const { coin } = res.data
    //     commit('USER_SETCOINBALANCE_MUTATE', coin)
    //   }
    //   return res
    // },
    // 登录
    async USER_LOGIN_ACTION ({ commit, dispatch }, { account, password, pwdLen }) {
      const res = await loginApi({
        account,
        password,
        pwd_len: pwdLen
      })
      if (res.return_code === '0') {
        commit('USER_SETTOKEN_MUTATE', res.data.token)
      }
      // let res = await dispatch('USER_GETUSERINFO_ACTION')
      return res
    },
    // 登出去除相关信息
    async USER_LOGOUT_ACTION ({ state, commit }) {
      commit('CLEARSTATE_MUTATE')
      commit('USER_SETTOKEN_MUTATE')
      // location.href = `${location.href.split('#')[0]}#/login` // 登录页
      return {}
    }
  }
}
