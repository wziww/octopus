(function(e){function t(t){for(var a,r,o=t[0],u=t[1],d=t[2],i=0,f=[];i<o.length;i++)r=o[i],s[r]&&f.push(s[r][0]),s[r]=0;for(a in u)Object.prototype.hasOwnProperty.call(u,a)&&(e[a]=u[a]);l&&l(t);while(f.length)f.shift()();return c.push.apply(c,d||[]),n()}function n(){for(var e,t=0;t<c.length;t++){for(var n=c[t],a=!0,r=1;r<n.length;r++){var o=n[r];0!==s[o]&&(a=!1)}a&&(c.splice(t--,1),e=u(u.s=n[0]))}return e}var a={},r={app:0},s={app:0},c=[];function o(e){return u.p+"js/"+({}[e]||e)+"."+{"chunk-05522dae":"2eb0da34fbff31daa7d6","chunk-12efb352":"8b2a56cefe47cdb79770","chunk-2d0ab4d9":"60068b72dde6a159eff2","chunk-3db9aeea":"f3f6c6c963e9db2a06f5","chunk-66716589":"4c74aa70b91d7ce274a4","chunk-7c55d574":"90992023bdf315217054","chunk-89528e28":"a6cd0b4fbe21ff196498","chunk-a23410ec":"6fe66b9947424a9df673"}[e]+".js"}function u(t){if(a[t])return a[t].exports;var n=a[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,u),n.l=!0,n.exports}u.e=function(e){var t=[],n={"chunk-05522dae":1,"chunk-12efb352":1,"chunk-3db9aeea":1,"chunk-7c55d574":1,"chunk-89528e28":1,"chunk-a23410ec":1};r[e]?t.push(r[e]):0!==r[e]&&n[e]&&t.push(r[e]=new Promise(function(t,n){for(var a="css/"+({}[e]||e)+"."+{"chunk-05522dae":"4102bab1","chunk-12efb352":"1bae2e7a","chunk-2d0ab4d9":"31d6cfe0","chunk-3db9aeea":"01db269c","chunk-66716589":"31d6cfe0","chunk-7c55d574":"a684a7a0","chunk-89528e28":"2c5e559f","chunk-a23410ec":"5f35a743"}[e]+".css",s=u.p+a,c=document.getElementsByTagName("link"),o=0;o<c.length;o++){var d=c[o],i=d.getAttribute("data-href")||d.getAttribute("href");if("stylesheet"===d.rel&&(i===a||i===s))return t()}var f=document.getElementsByTagName("style");for(o=0;o<f.length;o++){d=f[o],i=d.getAttribute("data-href");if(i===a||i===s)return t()}var l=document.createElement("link");l.rel="stylesheet",l.type="text/css",l.onload=t,l.onerror=function(t){var a=t&&t.target&&t.target.src||s,c=new Error("Loading CSS chunk "+e+" failed.\n("+a+")");c.code="CSS_CHUNK_LOAD_FAILED",c.request=a,delete r[e],l.parentNode.removeChild(l),n(c)},l.href=s;var b=document.getElementsByTagName("head")[0];b.appendChild(l)}).then(function(){r[e]=0}));var a=s[e];if(0!==a)if(a)t.push(a[2]);else{var c=new Promise(function(t,n){a=s[e]=[t,n]});t.push(a[2]=c);var d,i=document.createElement("script");i.charset="utf-8",i.timeout=120,u.nc&&i.setAttribute("nonce",u.nc),i.src=o(e),d=function(t){i.onerror=i.onload=null,clearTimeout(f);var n=s[e];if(0!==n){if(n){var a=t&&("load"===t.type?"missing":t.type),r=t&&t.target&&t.target.src,c=new Error("Loading chunk "+e+" failed.\n("+a+": "+r+")");c.type=a,c.request=r,n[1](c)}s[e]=void 0}};var f=setTimeout(function(){d({type:"timeout",target:i})},12e4);i.onerror=i.onload=d,document.head.appendChild(i)}return Promise.all(t)},u.m=e,u.c=a,u.d=function(e,t,n){u.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},u.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},u.t=function(e,t){if(1&t&&(e=u(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(u.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var a in e)u.d(n,a,function(t){return e[t]}.bind(null,a));return n},u.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return u.d(t,"a",t),t},u.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},u.p="/",u.oe=function(e){throw console.error(e),e};var d=window["webpackJsonp"]=window["webpackJsonp"]||[],i=d.push.bind(d);d.push=t,d=d.slice();for(var f=0;f<d.length;f++)t(d[f]);var l=i;c.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},4678:function(e,t,n){var a={"./af":"2bfb","./af.js":"2bfb","./ar":"8e73","./ar-dz":"a356","./ar-dz.js":"a356","./ar-kw":"423e","./ar-kw.js":"423e","./ar-ly":"1cfd","./ar-ly.js":"1cfd","./ar-ma":"0a84","./ar-ma.js":"0a84","./ar-sa":"8230","./ar-sa.js":"8230","./ar-tn":"6d83","./ar-tn.js":"6d83","./ar.js":"8e73","./az":"485c","./az.js":"485c","./be":"1fc1","./be.js":"1fc1","./bg":"84aa","./bg.js":"84aa","./bm":"a7fa","./bm.js":"a7fa","./bn":"9043","./bn.js":"9043","./bo":"d26a","./bo.js":"d26a","./br":"6887","./br.js":"6887","./bs":"2554","./bs.js":"2554","./ca":"d716","./ca.js":"d716","./cs":"3c0d","./cs.js":"3c0d","./cv":"03ec","./cv.js":"03ec","./cy":"9797","./cy.js":"9797","./da":"0f14","./da.js":"0f14","./de":"b469","./de-at":"b3eb","./de-at.js":"b3eb","./de-ch":"bb71","./de-ch.js":"bb71","./de.js":"b469","./dv":"598a","./dv.js":"598a","./el":"8d47","./el.js":"8d47","./en-SG":"cdab","./en-SG.js":"cdab","./en-au":"0e6b","./en-au.js":"0e6b","./en-ca":"3886","./en-ca.js":"3886","./en-gb":"39a6","./en-gb.js":"39a6","./en-ie":"e1d3","./en-ie.js":"e1d3","./en-il":"7333","./en-il.js":"7333","./en-nz":"6f50","./en-nz.js":"6f50","./eo":"65db","./eo.js":"65db","./es":"898b","./es-do":"0a3c","./es-do.js":"0a3c","./es-us":"55c9","./es-us.js":"55c9","./es.js":"898b","./et":"ec18","./et.js":"ec18","./eu":"0ff2","./eu.js":"0ff2","./fa":"8df4","./fa.js":"8df4","./fi":"81e9","./fi.js":"81e9","./fo":"0721","./fo.js":"0721","./fr":"9f26","./fr-ca":"d9f8","./fr-ca.js":"d9f8","./fr-ch":"0e49","./fr-ch.js":"0e49","./fr.js":"9f26","./fy":"7118","./fy.js":"7118","./ga":"5120","./ga.js":"5120","./gd":"f6b4","./gd.js":"f6b4","./gl":"8840","./gl.js":"8840","./gom-latn":"0caa","./gom-latn.js":"0caa","./gu":"e0c5","./gu.js":"e0c5","./he":"c7aa","./he.js":"c7aa","./hi":"dc4d","./hi.js":"dc4d","./hr":"4ba9","./hr.js":"4ba9","./hu":"5b14","./hu.js":"5b14","./hy-am":"d6b6","./hy-am.js":"d6b6","./id":"5038","./id.js":"5038","./is":"0558","./is.js":"0558","./it":"6e98","./it-ch":"6f12","./it-ch.js":"6f12","./it.js":"6e98","./ja":"079e","./ja.js":"079e","./jv":"b540","./jv.js":"b540","./ka":"201b","./ka.js":"201b","./kk":"6d79","./kk.js":"6d79","./km":"e81d","./km.js":"e81d","./kn":"3e92","./kn.js":"3e92","./ko":"22f8","./ko.js":"22f8","./ku":"2421","./ku.js":"2421","./ky":"9609","./ky.js":"9609","./lb":"440c","./lb.js":"440c","./lo":"b29d","./lo.js":"b29d","./lt":"26f9","./lt.js":"26f9","./lv":"b97c","./lv.js":"b97c","./me":"293c","./me.js":"293c","./mi":"688b","./mi.js":"688b","./mk":"6909","./mk.js":"6909","./ml":"02fb","./ml.js":"02fb","./mn":"958b","./mn.js":"958b","./mr":"39bd","./mr.js":"39bd","./ms":"ebe4","./ms-my":"6403","./ms-my.js":"6403","./ms.js":"ebe4","./mt":"1b45","./mt.js":"1b45","./my":"8689","./my.js":"8689","./nb":"6ce3","./nb.js":"6ce3","./ne":"3a39","./ne.js":"3a39","./nl":"facd","./nl-be":"db29","./nl-be.js":"db29","./nl.js":"facd","./nn":"b84c","./nn.js":"b84c","./pa-in":"f3ff","./pa-in.js":"f3ff","./pl":"8d57","./pl.js":"8d57","./pt":"f260","./pt-br":"d2d4","./pt-br.js":"d2d4","./pt.js":"f260","./ro":"972c","./ro.js":"972c","./ru":"957c","./ru.js":"957c","./sd":"6784","./sd.js":"6784","./se":"ffff","./se.js":"ffff","./si":"eda5","./si.js":"eda5","./sk":"7be6","./sk.js":"7be6","./sl":"8155","./sl.js":"8155","./sq":"c8f3","./sq.js":"c8f3","./sr":"cf1e","./sr-cyrl":"13e9","./sr-cyrl.js":"13e9","./sr.js":"cf1e","./ss":"52bd","./ss.js":"52bd","./sv":"5fbd","./sv.js":"5fbd","./sw":"74dc","./sw.js":"74dc","./ta":"3de5","./ta.js":"3de5","./te":"5cbb","./te.js":"5cbb","./tet":"576c","./tet.js":"576c","./tg":"3b1b","./tg.js":"3b1b","./th":"10e8","./th.js":"10e8","./tl-ph":"0f38","./tl-ph.js":"0f38","./tlh":"cf75","./tlh.js":"cf75","./tr":"0e81","./tr.js":"0e81","./tzl":"cf51","./tzl.js":"cf51","./tzm":"c109","./tzm-latn":"b53d","./tzm-latn.js":"b53d","./tzm.js":"c109","./ug-cn":"6117","./ug-cn.js":"6117","./uk":"ada2","./uk.js":"ada2","./ur":"5294","./ur.js":"5294","./uz":"2e8c","./uz-latn":"010e","./uz-latn.js":"010e","./uz.js":"2e8c","./vi":"2921","./vi.js":"2921","./x-pseudo":"fd7e","./x-pseudo.js":"fd7e","./yo":"7f33","./yo.js":"7f33","./zh-cn":"5c3a","./zh-cn.js":"5c3a","./zh-hk":"49ab","./zh-hk.js":"49ab","./zh-tw":"90ea","./zh-tw.js":"90ea"};function r(e){var t=s(e);return n(t)}function s(e){var t=a[e];if(!(t+1)){var n=new Error("Cannot find module '"+e+"'");throw n.code="MODULE_NOT_FOUND",n}return t}r.keys=function(){return Object.keys(a)},r.resolve=s,e.exports=r,r.id="4678"},"46b8":function(e,t,n){"use strict";n.d(t,"f",function(){return a}),n.d(t,"b",function(){return s}),n.d(t,"d",function(){return r}),n.d(t,"a",function(){return c}),n.d(t,"e",function(){return o}),n.d(t,"c",function(){return u});var a=localStorage.getItem("token"),r=localStorage.getItem("permission");function s(e){a=e,localStorage.setItem("token",e)}function c(e){r=e,localStorage.setItem("permission",e)}var o={PERMISSIONMONIT:1,PERMISSIONDEV:2,PERMISSIONEXEC:4};function u(){localStorage.removeItem("token"),localStorage.removeItem("permission")}},"56d7":function(e,t,n){"use strict";n.r(t);n("cadf"),n("551c"),n("f751"),n("097d");var a=n("2b0e"),r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app"}},[n("router-view")],1)},s=[],c=(n("7faf"),n("2877")),o={},u=Object(c["a"])(o,r,s,!1,null,null,null),d=u.exports,i=n("75fc"),f=n("8c4f"),l={path:"*",name:"error_404",component:function(){return n.e("chunk-2d0ab4d9").then(n.bind(null,"1565"))},meta:{title:"404"}},b=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("a-layout",{attrs:{id:"components-layout-demo-custom-trigger"}},[n("a-layout-sider",{style:{height:"100vh",overflow:"auto",position:"fixed",left:0},attrs:{trigger:null,collapsible:""}},[n("div",{staticClass:"logo"},[n("img",{style:{width:"auto",height:"100%"}})]),n("a-menu",{attrs:{theme:"dark",mode:"vertical",defaultSelectedKeys:[e.$route.meta.Index]}},[n("a-menu-item",{key:"1"},[n("router-link",{attrs:{to:"/"}},[n("a-icon",{attrs:{type:"cloud-o"}}),n("span",[e._v("数据源")])],1)],1),n("a-menu-item",{attrs:{selectable:"false"},on:{click:function(t){e.clear(),e.logout()}}},[n("span",[e._v("登出")])])],1)],1),n("a-layout",{staticStyle:{"box-sizing":"border-box"},style:{marginLeft:"200px",minHeight:"100vh"}},[n("a-layout-content",{style:{margin:"24px 16px 0"}},[n("div",{style:{padding:"24px"}},[n("router-view")],1)]),n("a-layout-footer",{staticStyle:{textAlign:"center"}},[e._v("octopus v0.0.1")])],1)],1)},j=[],h=n("46b8"),m={name:"common",data:function(){return{clear:h["c"]}},methods:{logout:function(){this.$router.push({path:"/login"})}}},p=m,g=Object(c["a"])(p,b,j,!1,null,null,null),v=g.exports,k={path:"/opcap",name:"setting_node_opcap",component:function(){return n.e("chunk-05522dae").then(n.bind(null,"9868"))},meta:{title:"redis-opcap",Index:"1"}},y={path:"/redis_dev",name:"setting_redis_dev",component:function(){return n.e("chunk-3db9aeea").then(n.bind(null,"bb26"))},meta:{title:"数据源-运维模式",Index:"1"}},_={path:"/clusterSlots",name:"setting_redis_monit",component:function(){return n.e("chunk-7c55d574").then(n.bind(null,"4cf9"))},meta:{title:"数据源-节点列表",Index:"1"}},w={path:"/redis_monit_main",name:"setting_redis_monit_main",component:function(){return n.e("chunk-a23410ec").then(n.bind(null,"2771"))},meta:{title:"数据源-实时监控",Index:"1"}},x={path:"/redis",name:"setting_redis",component:function(){return n.e("chunk-89528e28").then(n.bind(null,"fbbd"))},meta:{title:"数据源-redis-列表",Index:"1"}},S={path:"",name:"setting",component:function(){return n.e("chunk-66716589").then(n.bind(null,"1e4b"))},meta:{title:"数据源",Index:"1"}},z=[{path:"/login",name:"login",component:function(){return n.e("chunk-12efb352").then(n.bind(null,"dc3f"))},meta:{title:"octopus"}},{path:"/",component:v,children:[S,x,_,w,y,k,l]}];a["a"].use(f["a"]);var I=new f["a"]({mode:"history",routes:Object(i["a"])(z)});I.afterEach(function(e){h["f"]||I.push({path:"/login"})});var O=I,E=n("2819"),P=n.n(E),N=n("f23d"),T=(n("202f"),n("3af9"),n("998c")),C=n.n(T),M=n("2ead"),A=n.n(M);a["a"].use(P.a),a["a"].use(C.a),a["a"].use(A.a),a["a"].config.productionTip=!1,a["a"].use(N["a"]),new a["a"]({router:O,render:function(e){return e(d)}}).$mount("#app")},"7faf":function(e,t,n){"use strict";var a=n("dc71"),r=n.n(a);r.a},dc71:function(e,t,n){}});