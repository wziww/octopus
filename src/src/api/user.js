import HttpReq from '@/lib/https';
// 登录
function loginApi(data) {
  return HttpReq.formPost({
    url: 'login',
    data,
    notLogin: true
  });
}
export {
  loginApi // 登录
};
