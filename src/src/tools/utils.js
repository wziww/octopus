import Cookies from 'js-cookie';
import config from '@/config';
const { cookieTokenName, cookieExpires, title } = config;
/**
 * 设置token
 * @param {string} token 登录信息的token
 */
function setToken(token) {
  Cookies.set(cookieTokenName, token, { expires: cookieExpires || 1 });
}
/**
 * 获取token
 * @returns {string|boolean} 返回token或false
 */
function getToken() {
  const token = Cookies.get(cookieTokenName);
  if (token) return token;
  return false;
}
/**
 * 移除token
 */
function removeToken() {
  Cookies.remove(cookieTokenName);
}
/**
 * 动态设置页面title
 * @param {string|void} pageTitle 页面title
 */
function setTitle(pageTitle) {
  document.title = pageTitle ? `${title} - ${pageTitle}` : title;
}
/**
 * 判断是否有属性children
 * @param {JSON} item  对象
 * @return {boolean}
 */
// function hasChildren (item) {
//   return item.children && item.children.length > 0
// }
export {
  setToken, // 设置token
  getToken, // 获取token
  removeToken, // 移除token
  setTitle // 设置title
};
