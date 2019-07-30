const CLOSED = 0;
const CONNECTED = 1;
const ERROR = -1;
class WS {
  constructor(url, obj = {
    reconnect: true
  }) {
    this.$socket = null;
    this.$url = url;
    this.$socketStatus = CLOSED;
    this.$reconnect = obj.reconnect;
    this._onclose = () => {
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
    this.$socket = new WebSocket(this.$url);
    this.$socketStatus = CONNECTED;
    this.$socket.onclose = this._initOnClose();
    this.$socket.onerror = this._initOnError();
    this.$socket.onmessage = this._initOnMessage();
  }
  Close(fn) {
    if (!this.$socket) return;
    this.$socket.close();
  }
  OnClose(fn) {
    if (!this.$socket) return;
    this.$onclose.push(fn);
  }
  OnOpen(fn) {
    if (!this.$socket) return;
    this.$onopen.push(fn);
  }
  OnData(fn) {
    if (!this.$socket) return;
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
    this.$socket.send(JSON.stringify(d));
  }
  _initOnClose() {
    this.$socketStatus = CLOSED;
    const that = this;
    return () => {
      that._onclose();
      for (let i = 0; i < that.$onclose.length; i++) {
        if (typeof that.$onclose[i] === 'function') {
          that.$onclose[i]();
        }
      }
    };
  }
  _initOnError(e) {
    this.$socketStatus = ERROR;
    const that = this;
    return () => {
      for (let i = 0; i < that.$onerror.length; i++) {
        if (typeof that.$onerror[i] === 'function') {
          that.$onerror[i](e);
        }
      }
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
