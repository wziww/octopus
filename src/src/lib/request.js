import axios from 'axios'
import qs from 'qs'
import store from '@/store'
// import { APP_API_NAME } from '@/api'
// import { checkObjectTool } from './tools'
// import { getToken, setToken } from './utils'
// 需求： 1、请求插件 2、接口出错收集
function addErrorLog (info) {
  const { statusText, status, request: { responseURL } } = info
  let data = {
    type: 'ajax',
    code: status,
    mes: statusText,
    url: responseURL
  }
  // 判断是否是发送报错日志的接口 防止进入死循环
  // if (!responseURL.includes(APP_API_NAME.SAVE_ERROR_LOG))
  store.dispatch('APP_ADDERRORLOG_ACTION', data)
}
// 请求时处理
axios.interceptors.request.use(request => {
  let format = request.format
  delete request.format
  // loading
  switch (format) {
    case 1:
      // Indicator.open('加载中')
      return request
    case 2:
      // Indicator.open('正在提交')
      return request
    default:
      return request
  }
}, error => {
  return Promise.reject(error)
})
// 返回结束处理

axios.interceptors.response.use(response => {
  // 返回值处理 如果已经登出那么
  const status = response.status
  switch (status) {
    case 200:
      let { data } = response
      const code = data.return_code
      // 未登录
      if (code === '-200') {
        store.dispatch('USER_LOGOUT_ACTION')
      }
      return data
    default:
      // doSomething 收集信息
      return {}
  }
}, error => {
  // 网络异常 收集信息
  if (error && error.response) {
    addErrorLog(error.response)
  } else {
    // 网络出问题了
  }
  return Promise.reject(error)
  // Toast('网络异常')
  // return Promise.resolve(error.response)
})
class HttpRequest {
  constructor (baseUrl) {
    this.baseURL = baseUrl
    this.queue = {}
  }
  // 根据有无传接口地址
  getOptionChange (options) {
    if (!options.baseURL) {
      // 相对路径
      if (options.url.indexOf('/') !== 0) {
        options.baseURL = this.baseURL
      }
    }
    return options
  }
  // 判断是否需要登录等处理 notLogin true 为不需要登录
  askLogin (notLogin) {
    // 有token
    if (!notLogin && store.state.user.token) return true
    return notLogin
  }
  // 过滤整合axios配置
  getParamConfig (options) {
    options = this.getOptionChange(options)
    let method = options.method
    let baseURL = options.baseURL
    let url = options.url || ''
    let headers = options.headers || {}
    let head = options.head || {}
    let data = options.data || {}
    let format = options.format
    headers = { ...head, ...headers }
    if (!options.notLogin) {
      headers.token = store.state.user.token
      // data.token = store.state.user.token
    }
    let config = {
      method,
      baseURL,
      url,
      headers,
      format
    }
    if (options.timeout) {
      config.timeout = options.timeout
    }
    if (method === 'post') {
      data = headers['Content-Type'] === 'application/x-www-form-urlencoded' ? qs.stringify(data) : data
      config.data = data
    } else if (method === 'get') {
      config.params = data
    }
    return config
  }
  // 添加内部配置
  getInsideConfig () {
    const config = {
      withCredentials: false, // 允许携带cookie
      timeout: 10000
    }
    return config
  }
  async request (options) {
    // 默认登录（判断是否有token）有token返回true，不登录返回true
    if (!this.askLogin(options.notLogin)) {
    // 需要登录（没有登录信息）
      await store.dispatch('USER_LOGOUT_ACTION')
      return false
    }
    const instance = axios
    let config = Object.assign(this.getInsideConfig(), this.getParamConfig(options))
    let res = instance(config)
    return res
  }
  async formPost (options) {
    let config = {
      method: 'post',
      head: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Accept-Language': 'charset=utf-8'
      },
      ...options
      // url: options.url,
      // data: options.data ? qs.stringify(options.data) : null,
      // headers: options
    }
    return this.request(config)
  }
  async jsonPost (options) {
    let config = {
      method: 'post',
      head: {
        'Content-Type': 'application/json',
        'Accept-Language': 'charset=utf-8'
      },
      ...options
      // url: options.url,
      // data: options.data,
      // headers: options
    }
    return this.request(config)
  }
  async urlGet (options) {
    let config = {
      method: 'get',
      head: {
        'Accept-Language': 'charset=utf-8'
      },
      ...options
      // url: options.url,
      // data: options.data ? qs.stringify(options.data) : null,
      // headers: options
    }
    return this.request(config)
  }
  async uploadPost (options) {
    let config = {
      method: 'post',
      head: {
        'Content-Type': 'multipart/form-data',
        'Accept-Language': 'charset=utf-8'
      },
      ...options
      // url: options.url,
      // data: options.data ? qs.stringify(options.data) : null,
      // headers: options
    }
    return this.request(config)
  }
}
export default HttpRequest
