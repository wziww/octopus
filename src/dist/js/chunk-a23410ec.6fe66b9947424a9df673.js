(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-a23410ec"],{1169:function(t,e,n){var r=n("2d95");t.exports=Array.isArray||function(t){return"Array"==r(t)}},"11e9":function(t,e,n){var r=n("52a7"),i=n("4630"),a=n("6821"),o=n("6a99"),s=n("69a8"),c=n("c69a"),u=Object.getOwnPropertyDescriptor;e.f=n("9e1e")?u:function(t,e){if(t=a(t),e=o(e,!0),c)try{return u(t,e)}catch(n){}if(s(t,e))return i(!r.f.call(t,e),t[e])}},"145e":function(t,e,n){"use strict";n("f559");var r=n("f64c"),i=function(t){return function(e){try{var n=JSON.parse(e.data);if(n.Data=JSON.parse(n.Data),(""+n.Data.code).startsWith("403"))return r["a"].warn(n.Data.message);if((""+n.Data.code).startsWith("404"))return r["a"].error(n.Data.message);if(200!==n.Data.code)return r["a"].error(n.Data.message);n.Data=n.Data.message,t(n)}catch(i){console.error(i),r["a"].error("sys error")}}};e["a"]=i},2771:function(t,e,n){"use strict";n.r(e);var r=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",[n("div",{staticStyle:{width:"100%",float:"left","margin-bottom":"20px"}},[n("span",{staticClass:"each-chose",staticStyle:{"font-size":"20px","font-weight":"border"}},[t._v("Refresh Every :")]),n("a-button",{staticClass:"each-chose",attrs:{type:t.index[0]},on:{click:function(e){return t.chose(0)}}},[t._v("1 s")]),n("a-button",{staticClass:"each-chose",attrs:{type:t.index[1]},on:{click:function(e){return t.chose(1)}}},[t._v("10 s")]),n("a-button",{staticClass:"each-chose",attrs:{type:t.index[2]},on:{click:function(e){return t.chose(2)}}},[t._v("30 s")]),n("a-button",{staticClass:"each-chose",attrs:{type:t.index[3]},on:{click:function(e){return t.chose(3)}}},[t._v("1 min")]),n("a-button",{staticClass:"each-chose",attrs:{type:t.index[4]},on:{click:function(e){return t.chose(4)}}},[t._v("5 min")]),n("a-button",{staticClass:"each-chose",attrs:{type:t.index[5]},on:{click:function(e){return t.chose(5)}}},[t._v("10 min")])],1),n("ve-line",{staticStyle:{width:"50%",float:"left"},attrs:{data:t.lineChartData,settings:t.lineChartSettings}}),n("ve-liquidfill",{staticStyle:{width:"50%",float:"left"},attrs:{radius:"50%;",data:t.chartData,settings:t.chartSettings}}),n("ve-line",{staticStyle:{width:"50%",float:"left"},attrs:{data:t.statsData,settings:t.statsSettings}})],1)},i=[],a=(n("c5f6"),n("ac4d"),n("8a81"),n("ac6a"),n("145e")),o=n("46b8"),s=n("7cc5"),c=n("f121"),u={},f=null,l=["primary","default","default","default","default","default"],h=[],d=[],y=1e3,p="dev",b=new s["a"](c["a"].Host+"?op="+p+"&ot="+o["f"]+"&ocid=nil"),v={name:"setting_redis",data:function(){d=[],u=[],h=[],b.Open(),this.lineChartSettings={area:!0,scale:[!0,!0],yAxisName:["M"],xAxisName:["时间"]},this.statsSettings={area:!0,scale:[!0,!0],yAxisName:["value"],xAxisName:["时间"]},this.chartSettings={seriesMap:{"内存使用量":{radius:"40%",center:["20%","30%"],itemStyle:{opacity:.2},emphasis:{itemStyle:{opacity:.5}},backgroundStyle:{},label:{formatter:function(t){var e=t.seriesName,n=t.value;return"".concat(e,"\n").concat((100*n).toFixed(2),"%")},fontSize:20}}}};var t=this;return f=setInterval(function(){try{b.SendObj({Func:"/redis/detail",Data:JSON.stringify({id:t.$route.query.id})}),b.SendObj({Func:"/redis/stats",Data:JSON.stringify({id:t.$route.query.id})})}catch(e){console.error(e)}},y),b.Close(),b.OnData(Object(a["a"])(function(e){if("/redis/detail"===e.Type){var n=0,r=0,i=!0,a=!1,o=void 0;try{for(var s,c=e.Data[Symbol.iterator]();!(i=(s=c.next()).done);i=!0){var f=s.value;n+=Number(f.UsedMemory),r+=Number(f.Maxmemory)}}catch(k){a=!0,o=k}finally{try{i||null==c.return||c.return()}finally{if(a)throw o}}h.length>=20&&h.shift(),h.push({t:t.$moment().format("hh:mm:ss"),memory_total:(n/1024/1024).toFixed(2)}),u={columns:["key","percent"],rows:[{key:"内存使用量",percent:(n/r).toFixed(4)}]},t.chartData=u}if("/redis/stats"===e.Type){var l=0,y=0,p=0,b=!0,v=!1,m=void 0;try{for(var g,S=e.Data[Symbol.iterator]();!(b=(g=S.next()).done);b=!0){var O=g.value;y+=Number(O.InstantaneousOutputKbps),l+=Number(O.InstantaneousInputKbps),p+=Number(O.InstantaneousOpsPerSec)}}catch(k){v=!0,m=k}finally{try{b||null==S.return||S.return()}finally{if(v)throw m}}d.length>=20&&d.shift(),d.push({t:t.$moment().format("hh:mm:ss"),output_Kbps:y,input_Kbps:l,Ops:p})}})),{chartData:u,interTime:y,lineChartData:{columns:["t","memory_total"],rows:h},index:l,statsData:{columns:["t","output_Kbps","input_Kbps","Ops"],rows:d}}},beforeDestroy:function(){b.Close(),null!==f&&window.clearInterval(f)},methods:{chose:function(t){var e=this,n=this;switch(l=["default","default","default","default","default","default"],l[t]="primary",this.index=l,t){case 0:y=1e3,this.interTime=y,window.clearInterval(f),f=setInterval(function(){e.$socket.sendObj({Func:"/redis/detail",Data:JSON.stringify({id:n.$route.query.id})}),e.$socket.sendObj({Func:"/redis/stats",Data:JSON.stringify({id:n.$route.query.id})})},y);break;case 1:y=1e4,this.interTime=y,window.clearInterval(f),f=setInterval(function(){e.$socket.sendObj({Func:"/redis/detail",Data:JSON.stringify({id:n.$route.query.id})}),e.$socket.sendObj({Func:"/redis/stats",Data:JSON.stringify({id:n.$route.query.id})})},y);break;case 2:y=3e4,this.interTime=y,window.clearInterval(f),f=setInterval(function(){e.$socket.sendObj({Func:"/redis/detail",Data:JSON.stringify({id:n.$route.query.id})}),e.$socket.sendObj({Func:"/redis/stats",Data:JSON.stringify({id:n.$route.query.id})})},y);break;case 3:y=6e4,this.interTime=y,window.clearInterval(f),f=setInterval(function(){e.$socket.sendObj({Func:"/redis/detail",Data:JSON.stringify({id:n.$route.query.id})}),e.$socket.sendObj({Func:"/redis/stats",Data:JSON.stringify({id:n.$route.query.id})})},y);break;case 4:y=3e5,this.interTime=y,window.clearInterval(f),f=setInterval(function(){e.$socket.sendObj({Func:"/redis/detail",Data:JSON.stringify({id:n.$route.query.id})}),e.$socket.sendObj({Func:"/redis/stats",Data:JSON.stringify({id:n.$route.query.id})})},y);break;case 5:y=6e5,this.interTime=y,window.clearInterval(f),f=setInterval(function(){e.$socket.sendObj({Func:"/redis/detail",Data:JSON.stringify({id:n.$route.query.id})}),e.$socket.sendObj({Func:"/redis/stats",Data:JSON.stringify({id:n.$route.query.id})})},y);break}},split:function(t){if("string"!==typeof t)return[];for(var e=t.length,n=[],r=0;r<e;r+=10)n.push(t.substr(r,10));return n}}},m=v,g=(n("7c6b"),n("2877")),S=Object(g["a"])(m,r,i,!1,null,"4d4a4efb",null);e["default"]=S.exports},"37c8":function(t,e,n){e.f=n("2b4c")},"3a72":function(t,e,n){var r=n("7726"),i=n("8378"),a=n("2d00"),o=n("37c8"),s=n("86cc").f;t.exports=function(t){var e=i.Symbol||(i.Symbol=a?{}:r.Symbol||{});"_"==t.charAt(0)||t in e||s(e,t,{value:o.f(t)})}},5147:function(t,e,n){var r=n("2b4c")("match");t.exports=function(t){var e=/./;try{"/./"[t](e)}catch(n){try{return e[r]=!1,!"/./"[t](e)}catch(i){}}return!0}},"5dbc":function(t,e,n){var r=n("d3f4"),i=n("8b97").set;t.exports=function(t,e,n){var a,o=e.constructor;return o!==n&&"function"==typeof o&&(a=o.prototype)!==n.prototype&&r(a)&&i&&i(t,a),t}},"67ab":function(t,e,n){var r=n("ca5a")("meta"),i=n("d3f4"),a=n("69a8"),o=n("86cc").f,s=0,c=Object.isExtensible||function(){return!0},u=!n("79e5")(function(){return c(Object.preventExtensions({}))}),f=function(t){o(t,r,{value:{i:"O"+ ++s,w:{}}})},l=function(t,e){if(!i(t))return"symbol"==typeof t?t:("string"==typeof t?"S":"P")+t;if(!a(t,r)){if(!c(t))return"F";if(!e)return"E";f(t)}return t[r].i},h=function(t,e){if(!a(t,r)){if(!c(t))return!0;if(!e)return!1;f(t)}return t[r].w},d=function(t){return u&&y.NEED&&c(t)&&!a(t,r)&&f(t),t},y=t.exports={KEY:r,NEED:!1,fastKey:l,getWeak:h,onFreeze:d}},"7bbc":function(t,e,n){var r=n("6821"),i=n("9093").f,a={}.toString,o="object"==typeof window&&window&&Object.getOwnPropertyNames?Object.getOwnPropertyNames(window):[],s=function(t){try{return i(t)}catch(e){return o.slice()}};t.exports.f=function(t){return o&&"[object Window]"==a.call(t)?s(t):i(r(t))}},"7c6b":function(t,e,n){"use strict";var r=n("b0b4"),i=n.n(r);i.a},"7cc5":function(t,e,n){"use strict";n("f559");function r(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}var i=n("85f2"),a=n.n(i);function o(t,e){for(var n=0;n<e.length;n++){var r=e[n];r.enumerable=r.enumerable||!1,r.configurable=!0,"value"in r&&(r.writable=!0),a()(t,r.key,r)}}function s(t,e,n){return e&&o(t.prototype,e),n&&o(t,n),t}var c=n("f64c"),u=n("46b8"),f=0,l=1,h=-1,d=function(){function t(e){var n=this,i=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{reconnect:!0};if(r(this,t),e.startsWith("ws"))console.log("yes");else{var a=window.location.host,o=window.location.protocol;switch(o.startsWith("https")){case!0:e="wss://"+a+e;break;default:e="ws://"+a+e;break}}this.$socket=null,this.$url=e,this.$socketStatus=f,this.$reconnect=i.reconnect,this._onclose=function(t){n.$reconnect&&n.$socketStatus!==l&&setTimeout(function(){n.Open()},1e3)},this.$onclose=[],this.$onerror=[],this.$onopen=[],this.$onmessage=[]}return s(t,[{key:"Open",value:function(){this.$socket=new WebSocket(this.$url),this.$socket.onclose=this._initOnClose(),this.$socket.onopen=this._initOnOpen(),this.$socket.onerror=this._initOnError(),this.$socket.onmessage=this._initOnMessage()}},{key:"_clean",value:function(){this.$onclose=[],this.$onerror=[],this.$onopen=[],this.$onmessage=[]}},{key:"Close",value:function(t){this.$socket&&this.$socket.close()}},{key:"OnClose",value:function(t){this.$onclose.push(t)}},{key:"OnOpen",value:function(t){this.$onopen.push(t)}},{key:"OnData",value:function(t){this.$onmessage.push(t)}},{key:"OnError",value:function(t){this.$onerror.push(t)}},{key:"Send",value:function(t){this.$socket&&this.$socket.send(t)}},{key:"SendObj",value:function(t){this.$socket&&this.$socket.send(JSON.stringify(t))}},{key:"_initOnClose",value:function(t){this.$socketStatus=f;var e=this;return function(t){e._onclose(t);for(var n=0;n<e.$onclose.length;n++)"function"===typeof e.$onclose[n]&&e.$onclose[n](t);e._clean()}}},{key:"_initOnOpen",value:function(){var t=this,e=this;return function(){t.$socketStatus=l,t.SendObj({Func:"token",Data:JSON.stringify({token:u["f"]})}),c["a"].success("ws 成功连接!");for(var n=0;n<e.$onopen.length;n++)"function"===typeof e.$onopen[n]&&e.$onopen[n]()}}},{key:"_initOnError",value:function(t){var e=this,n=this;return function(){e.$socketStatus=h;for(var r=0;r<n.$onerror.length;r++)"function"===typeof n.$onerror[r]&&n.$onerror[r](t)}}},{key:"_initOnMessage",value:function(t){var e=this;return function(t){for(var n=0;n<e.$onmessage.length;n++)"function"===typeof e.$onmessage[n]&&e.$onmessage[n](t)}}}]),t}();e["a"]=d},"85f2":function(t,e,n){t.exports=n("454f")},"8a81":function(t,e,n){"use strict";var r=n("7726"),i=n("69a8"),a=n("9e1e"),o=n("5ca1"),s=n("2aba"),c=n("67ab").KEY,u=n("79e5"),f=n("5537"),l=n("7f20"),h=n("ca5a"),d=n("2b4c"),y=n("37c8"),p=n("3a72"),b=n("d4c0"),v=n("1169"),m=n("cb7c"),g=n("d3f4"),S=n("4bf8"),O=n("6821"),k=n("6a99"),w=n("4630"),$=n("2aeb"),N=n("7bbc"),_=n("11e9"),D=n("2621"),x=n("86cc"),I=n("0d58"),E=_.f,F=x.f,T=N.f,j=r.Symbol,C=r.JSON,L=C&&C.stringify,A="prototype",J=d("_hidden"),P=d("toPrimitive"),M={}.propertyIsEnumerable,q=f("symbol-registry"),V=f("symbols"),G=f("op-symbols"),R=Object[A],W="function"==typeof j&&!!D.f,K=r.QObject,H=!K||!K[A]||!K[A].findChild,Y=a&&u(function(){return 7!=$(F({},"a",{get:function(){return F(this,"a",{value:7}).a}})).a})?function(t,e,n){var r=E(R,e);r&&delete R[e],F(t,e,n),r&&t!==R&&F(R,e,r)}:F,z=function(t){var e=V[t]=$(j[A]);return e._k=t,e},U=W&&"symbol"==typeof j.iterator?function(t){return"symbol"==typeof t}:function(t){return t instanceof j},X=function(t,e,n){return t===R&&X(G,e,n),m(t),e=k(e,!0),m(n),i(V,e)?(n.enumerable?(i(t,J)&&t[J][e]&&(t[J][e]=!1),n=$(n,{enumerable:w(0,!1)})):(i(t,J)||F(t,J,w(1,{})),t[J][e]=!0),Y(t,e,n)):F(t,e,n)},B=function(t,e){m(t);var n,r=b(e=O(e)),i=0,a=r.length;while(a>i)X(t,n=r[i++],e[n]);return t},Q=function(t,e){return void 0===e?$(t):B($(t),e)},Z=function(t){var e=M.call(this,t=k(t,!0));return!(this===R&&i(V,t)&&!i(G,t))&&(!(e||!i(this,t)||!i(V,t)||i(this,J)&&this[J][t])||e)},tt=function(t,e){if(t=O(t),e=k(e,!0),t!==R||!i(V,e)||i(G,e)){var n=E(t,e);return!n||!i(V,e)||i(t,J)&&t[J][e]||(n.enumerable=!0),n}},et=function(t){var e,n=T(O(t)),r=[],a=0;while(n.length>a)i(V,e=n[a++])||e==J||e==c||r.push(e);return r},nt=function(t){var e,n=t===R,r=T(n?G:O(t)),a=[],o=0;while(r.length>o)!i(V,e=r[o++])||n&&!i(R,e)||a.push(V[e]);return a};W||(j=function(){if(this instanceof j)throw TypeError("Symbol is not a constructor!");var t=h(arguments.length>0?arguments[0]:void 0),e=function(n){this===R&&e.call(G,n),i(this,J)&&i(this[J],t)&&(this[J][t]=!1),Y(this,t,w(1,n))};return a&&H&&Y(R,t,{configurable:!0,set:e}),z(t)},s(j[A],"toString",function(){return this._k}),_.f=tt,x.f=X,n("9093").f=N.f=et,n("52a7").f=Z,D.f=nt,a&&!n("2d00")&&s(R,"propertyIsEnumerable",Z,!0),y.f=function(t){return z(d(t))}),o(o.G+o.W+o.F*!W,{Symbol:j});for(var rt="hasInstance,isConcatSpreadable,iterator,match,replace,search,species,split,toPrimitive,toStringTag,unscopables".split(","),it=0;rt.length>it;)d(rt[it++]);for(var at=I(d.store),ot=0;at.length>ot;)p(at[ot++]);o(o.S+o.F*!W,"Symbol",{for:function(t){return i(q,t+="")?q[t]:q[t]=j(t)},keyFor:function(t){if(!U(t))throw TypeError(t+" is not a symbol!");for(var e in q)if(q[e]===t)return e},useSetter:function(){H=!0},useSimple:function(){H=!1}}),o(o.S+o.F*!W,"Object",{create:Q,defineProperty:X,defineProperties:B,getOwnPropertyDescriptor:tt,getOwnPropertyNames:et,getOwnPropertySymbols:nt});var st=u(function(){D.f(1)});o(o.S+o.F*st,"Object",{getOwnPropertySymbols:function(t){return D.f(S(t))}}),C&&o(o.S+o.F*(!W||u(function(){var t=j();return"[null]"!=L([t])||"{}"!=L({a:t})||"{}"!=L(Object(t))})),"JSON",{stringify:function(t){var e,n,r=[t],i=1;while(arguments.length>i)r.push(arguments[i++]);if(n=e=r[1],(g(e)||void 0!==t)&&!U(t))return v(e)||(e=function(t,e){if("function"==typeof n&&(e=n.call(this,t,e)),!U(e))return e}),r[1]=e,L.apply(C,r)}}),j[A][P]||n("32e9")(j[A],P,j[A].valueOf),l(j,"Symbol"),l(Math,"Math",!0),l(r.JSON,"JSON",!0)},"8b97":function(t,e,n){var r=n("d3f4"),i=n("cb7c"),a=function(t,e){if(i(t),!r(e)&&null!==e)throw TypeError(e+": can't set as prototype!")};t.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(t,e,r){try{r=n("9b43")(Function.call,n("11e9").f(Object.prototype,"__proto__").set,2),r(t,[]),e=!(t instanceof Array)}catch(i){e=!0}return function(t,n){return a(t,n),e?t.__proto__=n:r(t,n),t}}({},!1):void 0),check:a}},9093:function(t,e,n){var r=n("ce10"),i=n("e11e").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return r(t,i)}},aa77:function(t,e,n){var r=n("5ca1"),i=n("be13"),a=n("79e5"),o=n("fdef"),s="["+o+"]",c="​",u=RegExp("^"+s+s+"*"),f=RegExp(s+s+"*$"),l=function(t,e,n){var i={},s=a(function(){return!!o[t]()||c[t]()!=c}),u=i[t]=s?e(h):o[t];n&&(i[n]=u),r(r.P+r.F*s,"String",i)},h=l.trim=function(t,e){return t=String(i(t)),1&e&&(t=t.replace(u,"")),2&e&&(t=t.replace(f,"")),t};t.exports=l},aae3:function(t,e,n){var r=n("d3f4"),i=n("2d95"),a=n("2b4c")("match");t.exports=function(t){var e;return r(t)&&(void 0!==(e=t[a])?!!e:"RegExp"==i(t))}},ac4d:function(t,e,n){n("3a72")("asyncIterator")},ac6a:function(t,e,n){for(var r=n("cadf"),i=n("0d58"),a=n("2aba"),o=n("7726"),s=n("32e9"),c=n("84f2"),u=n("2b4c"),f=u("iterator"),l=u("toStringTag"),h=c.Array,d={CSSRuleList:!0,CSSStyleDeclaration:!1,CSSValueList:!1,ClientRectList:!1,DOMRectList:!1,DOMStringList:!1,DOMTokenList:!0,DataTransferItemList:!1,FileList:!1,HTMLAllCollection:!1,HTMLCollection:!1,HTMLFormElement:!1,HTMLSelectElement:!1,MediaList:!0,MimeTypeArray:!1,NamedNodeMap:!1,NodeList:!0,PaintRequestList:!1,Plugin:!1,PluginArray:!1,SVGLengthList:!1,SVGNumberList:!1,SVGPathSegList:!1,SVGPointList:!1,SVGStringList:!1,SVGTransformList:!1,SourceBufferList:!1,StyleSheetList:!0,TextTrackCueList:!1,TextTrackList:!1,TouchList:!1},y=i(d),p=0;p<y.length;p++){var b,v=y[p],m=d[v],g=o[v],S=g&&g.prototype;if(S&&(S[f]||s(S,f,h),S[l]||s(S,l,v),c[v]=h,m))for(b in r)S[b]||a(S,b,r[b],!0)}},b0b4:function(t,e,n){},c5f6:function(t,e,n){"use strict";var r=n("7726"),i=n("69a8"),a=n("2d95"),o=n("5dbc"),s=n("6a99"),c=n("79e5"),u=n("9093").f,f=n("11e9").f,l=n("86cc").f,h=n("aa77").trim,d="Number",y=r[d],p=y,b=y.prototype,v=a(n("2aeb")(b))==d,m="trim"in String.prototype,g=function(t){var e=s(t,!1);if("string"==typeof e&&e.length>2){e=m?e.trim():h(e,3);var n,r,i,a=e.charCodeAt(0);if(43===a||45===a){if(n=e.charCodeAt(2),88===n||120===n)return NaN}else if(48===a){switch(e.charCodeAt(1)){case 66:case 98:r=2,i=49;break;case 79:case 111:r=8,i=55;break;default:return+e}for(var o,c=e.slice(2),u=0,f=c.length;u<f;u++)if(o=c.charCodeAt(u),o<48||o>i)return NaN;return parseInt(c,r)}}return+e};if(!y(" 0o1")||!y("0b1")||y("+0x1")){y=function(t){var e=arguments.length<1?0:t,n=this;return n instanceof y&&(v?c(function(){b.valueOf.call(n)}):a(n)!=d)?o(new p(g(e)),n,y):g(e)};for(var S,O=n("9e1e")?u(p):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),k=0;O.length>k;k++)i(p,S=O[k])&&!i(y,S)&&l(y,S,f(p,S));y.prototype=b,b.constructor=y,n("2aba")(r,d,y)}},d2c8:function(t,e,n){var r=n("aae3"),i=n("be13");t.exports=function(t,e,n){if(r(e))throw TypeError("String#"+n+" doesn't accept regex!");return String(i(t))}},d4c0:function(t,e,n){var r=n("0d58"),i=n("2621"),a=n("52a7");t.exports=function(t){var e=r(t),n=i.f;if(n){var o,s=n(t),c=a.f,u=0;while(s.length>u)c.call(t,o=s[u++])&&e.push(o)}return e}},f121:function(t,e,n){"use strict";e["a"]={Host:"/v1/websocket"}},f559:function(t,e,n){"use strict";var r=n("5ca1"),i=n("9def"),a=n("d2c8"),o="startsWith",s=""[o];r(r.P+r.F*n("5147")(o),"String",{startsWith:function(t){var e=a(this,t,o),n=i(Math.min(arguments.length>1?arguments[1]:void 0,e.length)),r=String(t);return s?s.call(e,r,n):e.slice(n,n+r.length)===r}})},fdef:function(t,e){t.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);