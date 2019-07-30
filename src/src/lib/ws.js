import { message } from 'ant-design-vue';

const hd = (fn) => {
  return (da) => {
    // 接受服务端数据
    try {
      const d = JSON.parse(da.data);
      d.Data = JSON.parse(d.Data);
      if (('' + d.Data.code).startsWith('403')) { // unauth
        return message.error(d.Data.message);
      }
      if (('' + d.Data.code).startsWith('404')) { // unauth
        return message.error(d.Data.message);
      }
      if (d.Data.code !== 200) {
        return message.error(d.Data.message);
      }
      d.Data = d.Data.message;
      fn(d);
    } catch (e) {
      console.error(e);
      message.error("sys error");
    }
  };
};
export default hd;
