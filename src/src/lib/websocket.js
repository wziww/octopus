import { token } from './token';
import { message, } from 'ant-design-vue';

const CLOSED = 0;
const CONNECTED = 1;
const ERROR = -1;
const CONNECTING = 2;
let handle;
class WS {
  constructor(url, obj = {
    reconnect: true,
  }) {
    if (url.startsWith("ws")) {
      console.log("yes");
    } else {
      const host = window.location.host;
      const protocol = window.location.protocol;
      switch (protocol.startsWith("https")) {
        case true:
          url = 'wss://' + host + url;
          break;
        default:
          url = 'ws://' + host + url;
          break;
      }
    }
    this.$socket = null;
    this.$url = url;
    this.$socketStatus = CLOSED;
    this.$reconnect = obj.reconnect;
    this._onclose = (e) => {
      if (this.$reconnect && this.$socketStatus !== CONNECTED) {
        setTimeout(() => {
          this.Open();
        }, 1000);
      }
    };
    this.$onclose = [];
    this.$onerror = [];
    this.$onopen = [];
    this.$onmessage = [];
  }
  Open() {
    if (this.$socketStatus === CONNECTED || this.$socketStatus === CONNECTING) return;
    handle = setTimeout(() => {
      if (this.$socketStatus === CONNECTING) {
        this.$socket.close();
        this.$socketStatus = ERROR;
        this.Open();
      }
    }, 1000*5);
    this.$socketStatus = CONNECTING
    this.$socket = new WebSocket(this.$url);
    this.$socket.onopen = this._initOnOpen();
    this.$socket.onclose = this._initOnClose();
    this.$socket.onerror = this._initOnError();
    this.$socket.onmessage = this._initOnMessage();
  }
  _clean() {
    this.$onclose = [];
    this.$onerror = [];
    this.$onopen = [];
    this.$onmessage = [];
  }
  Close(fn) {
    if (!this.$socket) return;
    this.$socket.close();
  }
  OnClose(fn) {
    this.$onclose.push(fn);
  }
  OnOpen(fn) {
    this.$onopen.push(fn);
  }
  OnData(fn) {
    this.$onmessage.push(fn);
  }
  OnError(fn) {
    this.$onerror.push(fn);
  }
  Send(d) {
    if (!this.$socket) return;
    this.$socket.send(d);
  }
  SendObj(d) {
    if (!this.$socket) return;
    try {
      this.$socket.send(JSON.stringify(d));
    } catch (e) {
      this.Open();
    }
  }
  _initOnClose(e) {
    const that = this;
    return (e) => {
      if (this.$socketStatus !== CONNECTED) { return; }
      message.error("ws 已断开!");
      this.$socketStatus = CLOSED;
      that._onclose(e);
      for (let i = 0; i < that.$onclose.length; i++) {
        if (typeof that.$onclose[i] === 'function') {
          that.$onclose[i](e);
        }
      }
      // that._clean();
    };
  }
  _initOnOpen() {
    const that = this;
    return () => {
      this.$socketStatus = CONNECTED;
      this.SendObj({
        Func: "token",
        Data: JSON.stringify({
          token,
        }),
      });
      message.success("ws 成功连接!");
      for (let i = 0; i < that.$onopen.length; i++) {
        if (typeof that.$onopen[i] === 'function') {
          that.$onopen[i]();
        }
      }
    };
  }
  _initOnError(e) {
    const that = this;
    return () => {
      if (this.$socketStatus == CONNECTING) return;
      this.$socketStatus = ERROR;
      for (let i = 0; i < that.$onerror.length; i++) {
        if (typeof that.$onerror[i] === 'function') {
          that.$onerror[i](e);
        }
      }
      this.Open();
    };
  }
  _initOnMessage(d) {
    const that = this;
    return (d) => {
      for (let i = 0; i < that.$onmessage.length; i++) {
        if (typeof that.$onmessage[i] === 'function') {
          that.$onmessage[i](d);
        }
      }
    };
  }
};
export default WS;
