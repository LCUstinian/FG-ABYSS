import{c as d,d as x,v as z,x as M,a as i,b as e,F as y,r as k,u as s,n as _,t as v,p as S,g as b,s as I,A as $,e as m,f as w,o as l,_ as C,k as V,l as A,h as N}from"./index-DEYOvZyo.js";import{C as B}from"./check-B1JLaK0l.js";/**
 * @license lucide-vue-next v0.400.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const L=d("ArchiveIcon",[["rect",{width:"20",height:"5",x:"2",y:"3",rx:"1",key:"1wp1u1"}],["path",{d:"M4 8v11a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8",key:"1s80jp"}],["path",{d:"M10 12h4",key:"a56b0p"}]]);/**
 * @license lucide-vue-next v0.400.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const O=d("InfoIcon",[["circle",{cx:"12",cy:"12",r:"10",key:"1mglay"}],["path",{d:"M12 16v-4",key:"1dtifu"}],["path",{d:"M12 8h.01",key:"e9boi3"}]]);/**
 * @license lucide-vue-next v0.400.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const P=d("PaletteIcon",[["circle",{cx:"13.5",cy:"6.5",r:".5",fill:"currentColor",key:"1okk4w"}],["circle",{cx:"17.5",cy:"10.5",r:".5",fill:"currentColor",key:"f64h9f"}],["circle",{cx:"8.5",cy:"7.5",r:".5",fill:"currentColor",key:"fotxhn"}],["circle",{cx:"6.5",cy:"12.5",r:".5",fill:"currentColor",key:"qy21gx"}],["path",{d:"M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z",key:"12rzf8"}]]);/**
 * @license lucide-vue-next v0.400.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const T=d("ScrollTextIcon",[["path",{d:"M15 12h-5",key:"r7krc0"}],["path",{d:"M15 8h-5",key:"1khuty"}],["path",{d:"M19 17V5a2 2 0 0 0-2-2H4",key:"zz82l3"}],["path",{d:"M8 21h12a2 2 0 0 0 2-2v-1a1 1 0 0 0-1-1H11a1 1 0 0 0-1 1v1a2 2 0 1 1-4 0V5a2 2 0 1 0-4 0v2a1 1 0 0 0 1 1h3",key:"1ph1d7"}]]);/**
 * @license lucide-vue-next v0.400.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const U=d("ShieldCheckIcon",[["path",{d:"M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z",key:"oel41y"}],["path",{d:"m9 12 2 2 4-4",key:"dzmm74"}]]);/**
 * @license lucide-vue-next v0.400.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const E=d("WifiIcon",[["path",{d:"M12 20h.01",key:"zekei9"}],["path",{d:"M2 8.82a15 15 0 0 1 20 0",key:"dnpr2z"}],["path",{d:"M5 12.859a10 10 0 0 1 14 0",key:"1x1e6c"}],["path",{d:"M8.5 16.429a5 5 0 0 1 7 0",key:"1bycff"}]]),F={class:"appearance-panel"},D={class:"settings-card"},H={class:"theme-cards"},K=["onClick"],W=["data-preview"],j={class:"theme-label"},q={class:"settings-card"},R={class:"accent-circles"},G=["aria-label","onClick"],J={class:"settings-card"},Q={class:"settings-card"},X=x({__name:"AppearancePanel",setup(g){const{locale:o}=z(),t=M(),u=[{value:"system",label:"跟随系统"},{value:"dark",label:"深色"},{value:"light",label:"浅色"}],p=[{label:"12px",value:"12px"},{label:"13px",value:"13px"},{label:"14px",value:"14px"},{label:"15px",value:"15px"}],h=[{label:"简体中文",value:"zh-CN"},{label:"English",value:"en-US"}];return(n,a)=>{const f=w("n-select");return l(),i("div",F,[e("div",D,[a[1]||(a[1]=e("h3",{class:"card-title"},"主题",-1)),e("div",H,[(l(),i(y,null,k(u,c=>e("button",{key:c.value,class:_(["theme-card",{"is-active":s(t).mode===c.value}]),onClick:r=>s(t).setMode(c.value)},[e("div",{class:"theme-preview","data-preview":c.value},null,8,W),e("span",j,v(c.label),1)],10,K)),64))])]),e("div",q,[a[2]||(a[2]=e("h3",{class:"card-title"},"强调色",-1)),e("div",R,[(l(!0),i(y,null,k(s($),(c,r)=>(l(),i("button",{key:r,class:_(["accent-circle",{"is-selected":s(t).accentKey===r}]),style:S({background:c[s(t).resolvedMode]}),"aria-label":`强调色: ${r}`,onClick:oe=>s(t).setAccent(r)},[s(t).accentKey===r?(l(),b(s(B),{key:0,size:14,color:"#fff"})):I("",!0)],14,G))),128))])]),e("div",J,[a[3]||(a[3]=e("h3",{class:"card-title"},"字体大小",-1)),m(f,{value:s(t).fontSize,options:p,size:"small",style:{width:"100px"},"onUpdate:value":s(t).setFontSize},null,8,["value","onUpdate:value"])]),e("div",Q,[a[4]||(a[4]=e("h3",{class:"card-title"},"语言",-1)),m(f,{value:s(o),options:h,size:"small",style:{width:"140px"},"onUpdate:value":a[0]||(a[0]=c=>o.value=c)},null,8,["value"])])])}}}),Y=C(X,[["__scopeId","data-v-d302b45a"]]),Z={class:"settings-view"},ee={class:"settings-nav"},te=["onClick"],ae={class:"settings-content"},se={class:"section-heading"},le={key:1,class:"placeholder-panel"},ce={style:{color:"var(--text-3)","font-size":"13px"}},ne=x({__name:"SettingsView",setup(g){const o=V("appearance"),t=[{key:"appearance",label:"外观",icon:P},{key:"connection",label:"连接",icon:E},{key:"security",label:"安全",icon:U},{key:"logs",label:"日志",icon:T},{key:"backup",label:"备份",icon:L},{key:"about",label:"关于",icon:O}],u=A(()=>{var p,h;return(h=(p=t.find(n=>n.key===o.value))==null?void 0:p.label)!=null?h:""});return(p,h)=>(l(),i("div",Z,[e("div",ee,[(l(),i(y,null,k(t,n=>e("button",{key:n.key,class:_(["settings-nav-item",{"is-active":o.value===n.key}]),onClick:a=>o.value=n.key},[(l(),b(N(n.icon),{size:16})),e("span",null,v(n.label),1)],10,te)),64))]),e("div",ae,[e("h2",se,v(u.value),1),o.value==="appearance"?(l(),b(Y,{key:0})):(l(),i("div",le,[e("span",ce,v(u.value)+" — 配置项待实现",1)]))])]))}}),de=C(ne,[["__scopeId","data-v-79111b43"]]);export{de as default};
