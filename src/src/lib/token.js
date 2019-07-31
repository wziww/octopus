let token = localStorage.getItem("token");
let permission = localStorage.getItem("permission");
function TokenSet(t) {
  token = t;
  localStorage.setItem("token", t);
}
function PermissionSet(arr) {
  permission = arr;
  localStorage.setItem("permission", arr);
}
const permissionAll = {
  // PSRMISSIONMONIT 查看监控界面
  PSRMISSIONMONIT: 1 << 0,
  // PERMISSIONDEV dev 运维模式
  PERMISSIONDEV: 1 << 1,
  // PERMISSIONEXEC 在线操作模式
  PERMISSIONEXEC: 1 << 2
};
function clear() {
  localStorage.removeItem("token");
  localStorage.removeItem("permission");
}
export {
  token, TokenSet, permission, PermissionSet, permissionAll, clear
};
