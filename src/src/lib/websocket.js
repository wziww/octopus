class WS {
  constructor(url, obj = {
    reconnect: true
  }) {
    this.$socket = null;
    this.$url = url;
    this.$reconnect = obj.reconnect;
  }
  Open() {
    this.$socket = new WebSocket(this.$url);
  }
  Close(fn) {
    this.$socket.onclose = fn;
    if (this.$reconnect) {
      this.$socket.Open();
    }
  }
  OnOpen(fn) {
    this.$socket.onopen = fn;
  }
  OnData(fn) {
    this.$socket.onmessage = fn;
  }
  OnError(fn) {
    this.$socket.onclose = fn;
  }
  Send(d) {
    this.$socket.send(d);
  }
  SendObj(d) {
    this.$socket.send(JSON.stringify(d));
  }
};
export default WS;
