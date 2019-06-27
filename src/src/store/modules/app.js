export default {
  state: {
    errorList: [] // 日错误志列表
  },
  mutations: {
    // 错误日志列表添加
    APP_ADDERRORLOG_MUTATE (state, error) {
      state.errorList.unshift(error)
    }
  },
  actions: {
    // 添加错误日志
    APP_ADDERRORLOG_ACTION ({ rootState, commit }, info) {
      const { user: { userId, userName, userHeadImage } } = rootState
      let data = {
        ...info,
        userId,
        time: new Date().getTime(),
        userName,
        userHeadImage
      }
      commit('APP_ADDERRORLOG_MUTATE', data)
    }
  }
}
