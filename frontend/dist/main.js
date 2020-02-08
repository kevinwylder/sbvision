!function(e){var t={};function n(r){if(t[r])return t[r].exports;var o=t[r]={i:r,l:!1,exports:{}};return e[r].call(o.exports,o,o.exports,n),o.l=!0,o.exports}n.m=e,n.c=t,n.d=function(e,t,r){n.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:r})},n.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},n.t=function(e,t){if(1&t&&(e=n(e)),8&t)return e;if(4&t&&"object"==typeof e&&e&&e.__esModule)return e;var r=Object.create(null);if(n.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var o in e)n.d(r,o,function(t){return e[t]}.bind(null,o));return r},n.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return n.d(t,"a",t),t},n.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},n.p="",n(n.s=1)}([function(e,t){e.exports=React},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0)),i=r(n(2)),a=n(12),s=n(3);function l(){let[e,t]=o.useState();return o.createElement("div",{className:"main-layout"},e?o.createElement(s.VideoDisplay,{video:e}):o.createElement(a.Listing,{onVideoSelected:e=>{t(e)}}))}t.PanelLayout=l,i.render(o.createElement(l,null),document.getElementById("root"))},function(e,t){e.exports=ReactDOM},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0));n(4);const i=n(9),a=n(8);t.VideoDisplay=function(e){let t=o.createRef(),n=o.createRef(),[r,s]=o.useState([100,100]);o.useEffect(()=>{n.current&&(n.current.onresize=({target:e})=>{let t=e;s([t.videoWidth,t.videoHeight])},n.current.load())},[n.current,t.current]);let l=o.createRef(),[c,u]=o.useState({top:0,left:0,width:0,height:0});o.useEffect(()=>{if(!l.current)return;let[e,t]=r,n=e/t,{width:o,height:i}=l.current.getBoundingClientRect(),a=Math.max(o-n*i,0)/2,s=o-2*a,c=s/n,d=Math.max((i-c)/2,0);u({top:d,left:a,width:s,height:c})},[l.current,r]);let[d,f]=o.useState(!1),[h,p]=o.useState(0),[m,v]=o.useState(0),[b,g]=o.useState(!1);return o.useEffect(()=>{m&&n.current&&(n.current.currentTime+=m/e.video.fps,v(0)),b&&n.current&&(n.current.play(),g(!1))},[m,b,n.current]),o.createElement("div",{ref:l,className:"video-container"},o.createElement("video",{ref:n,onClick:e=>{var t;return null===(t=n.current)||void 0===t?void 0:t.pause()},style:Object.assign({position:"absolute"},c)},o.createElement("source",{src:`/video?type=${e.video.type}&id=${e.video.id}`,type:e.video.format})),o.createElement(i.VideoBox,{layout:c,video:n,setHasPlayed:f,videoWidth:r[0],videoHeight:r[1],onPause:e=>{},onSubmit:e=>{v(1)},onRefuse:()=>{g(!0)}}),o.createElement(a.VideoScrubber,{video:n,fps:e.video.fps,onFrame:p,bounds:{top:c.top+c.height-15,left:c.left,width:c.width}}),o.createElement("img",{style:Object.assign(Object.assign({},c),{position:"absolute",objectFit:"contain",display:d?"none":"block"}),src:`/images/${e.video.thumbnail}`,onClick:()=>n.current&&n.current.play()}))}},function(e,t,n){var r=n(5),o=n(6);"string"==typeof(o=o.__esModule?o.default:o)&&(o=[[e.i,o,""]]);var i={insert:"head",singleton:!1},a=(r(o,i),o.locals?o.locals:{});e.exports=a},function(e,t,n){"use strict";var r,o=function(){return void 0===r&&(r=Boolean(window&&document&&document.all&&!window.atob)),r},i=function(){var e={};return function(t){if(void 0===e[t]){var n=document.querySelector(t);if(window.HTMLIFrameElement&&n instanceof window.HTMLIFrameElement)try{n=n.contentDocument.head}catch(e){n=null}e[t]=n}return e[t]}}(),a=[];function s(e){for(var t=-1,n=0;n<a.length;n++)if(a[n].identifier===e){t=n;break}return t}function l(e,t){for(var n={},r=[],o=0;o<e.length;o++){var i=e[o],l=t.base?i[0]+t.base:i[0],c=n[l]||0,u="".concat(l," ").concat(c);n[l]=c+1;var d=s(u),f={css:i[1],media:i[2],sourceMap:i[3]};-1!==d?(a[d].references++,a[d].updater(f)):a.push({identifier:u,updater:v(f,t),references:1}),r.push(u)}return r}function c(e){var t=document.createElement("style"),r=e.attributes||{};if(void 0===r.nonce){var o=n.nc;o&&(r.nonce=o)}if(Object.keys(r).forEach((function(e){t.setAttribute(e,r[e])})),"function"==typeof e.insert)e.insert(t);else{var a=i(e.insert||"head");if(!a)throw new Error("Couldn't find a style target. This probably means that the value for the 'insert' parameter is invalid.");a.appendChild(t)}return t}var u,d=(u=[],function(e,t){return u[e]=t,u.filter(Boolean).join("\n")});function f(e,t,n,r){var o=n?"":r.media?"@media ".concat(r.media," {").concat(r.css,"}"):r.css;if(e.styleSheet)e.styleSheet.cssText=d(t,o);else{var i=document.createTextNode(o),a=e.childNodes;a[t]&&e.removeChild(a[t]),a.length?e.insertBefore(i,a[t]):e.appendChild(i)}}function h(e,t,n){var r=n.css,o=n.media,i=n.sourceMap;if(o?e.setAttribute("media",o):e.removeAttribute("media"),i&&btoa&&(r+="\n/*# sourceMappingURL=data:application/json;base64,".concat(btoa(unescape(encodeURIComponent(JSON.stringify(i))))," */")),e.styleSheet)e.styleSheet.cssText=r;else{for(;e.firstChild;)e.removeChild(e.firstChild);e.appendChild(document.createTextNode(r))}}var p=null,m=0;function v(e,t){var n,r,o;if(t.singleton){var i=m++;n=p||(p=c(t)),r=f.bind(null,n,i,!1),o=f.bind(null,n,i,!0)}else n=c(t),r=h.bind(null,n,t),o=function(){!function(e){if(null===e.parentNode)return!1;e.parentNode.removeChild(e)}(n)};return r(e),function(t){if(t){if(t.css===e.css&&t.media===e.media&&t.sourceMap===e.sourceMap)return;r(e=t)}else o()}}e.exports=function(e,t){(t=t||{}).singleton||"boolean"==typeof t.singleton||(t.singleton=o());var n=l(e=e||[],t);return function(e){if(e=e||[],"[object Array]"===Object.prototype.toString.call(e)){for(var r=0;r<n.length;r++){var o=s(n[r]);a[o].references--}for(var i=l(e,t),c=0;c<n.length;c++){var u=s(n[c]);0===a[u].references&&(a[u].updater(),a.splice(u,1))}n=i}}}},function(e,t,n){(t=n(7)(!1)).push([e.i,'body{\n    margin: 0px;\n    padding: 0px;\n}\n\n.video-container{\n    overflow: hidden;\n    width: 100%;\n    height: 100%;\n    background-color: black;\n    position: relative;\n}\n\n.video-control-bar{\n    background-color: "#FF0000";\n    color: white;\n    position: absolute;\n}',""]),e.exports=t},function(e,t,n){"use strict";e.exports=function(e){var t=[];return t.toString=function(){return this.map((function(t){var n=function(e,t){var n=e[1]||"",r=e[3];if(!r)return n;if(t&&"function"==typeof btoa){var o=(a=r,s=btoa(unescape(encodeURIComponent(JSON.stringify(a)))),l="sourceMappingURL=data:application/json;charset=utf-8;base64,".concat(s),"/*# ".concat(l," */")),i=r.sources.map((function(e){return"/*# sourceURL=".concat(r.sourceRoot||"").concat(e," */")}));return[n].concat(i).concat([o]).join("\n")}var a,s,l;return[n].join("\n")}(t,e);return t[2]?"@media ".concat(t[2]," {").concat(n,"}"):n})).join("")},t.i=function(e,n,r){"string"==typeof e&&(e=[[null,e,""]]);var o={};if(r)for(var i=0;i<this.length;i++){var a=this[i][0];null!=a&&(o[a]=!0)}for(var s=0;s<e.length;s++){var l=[].concat(e[s]);r&&o[l[0]]||(n&&(l[2]?l[2]="".concat(n," and ").concat(l[2]):l[2]=n),t.push(l))}},t}},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0));t.VideoScrubber=function({bounds:e,video:t,fps:n,onFrame:r}){let[i,a]=o.useState(0),[s,l]=o.useState(),[c,u]=o.useState(0),[d,f]=o.useState(!1),[h,p]=o.useState(new Date),m=o.createRef();o.useEffect(()=>{b()},[i,s]),o.useEffect(()=>{b();let e=window.setTimeout(()=>b(),2e3);return()=>window.clearTimeout(e)},[h]),o.useEffect(()=>{t.current&&(t.current.ontimeupdate=({target:e})=>{let t=e;a(t.currentTime),l(t.buffered),u(t.duration),r(Math.round(t.currentTime*n))})},[t.current]);const v=n=>{if(b(),!m.current||!t.current)return;let{x:r}=m.current.getBoundingClientRect(),o=t.current.duration*(n.clientX-r)/e.width;t.current.currentTime=o,a(o)},b=()=>{var t;let n=null===(t=m.current)||void 0===t?void 0:t.getContext("2d");if(!n)return;let r=(new Date).getTime()-h.getTime()>2e3?5:30,o=(30-r)/2,a=e.width/c;if(n.clearRect(0,0,e.width,30),s){n.fillStyle="grey";for(let e=0;e<s.length;e++)n.fillRect(s.start(e)*a,o,s.end(e)*a,r)}n.fillStyle="#33b5e5",n.fillRect(0,o,i*a,r)};return o.createElement("canvas",{ref:m,width:e.width,onMouseDown:e=>f(!0),onMouseMove:e=>{p(new Date),d&&v(e)},onMouseUp:e=>f(!1),onMouseLeave:e=>f(!1),onMouseOut:e=>f(!1),onTouchStart:e=>p(new Date),onTouchMove:e=>{p(new Date),v(e.targetTouches[0])},height:30,style:Object.assign(Object.assign({},e),{position:"absolute",height:30})})}},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0)),i=n(10);t.VideoBox=function(e){let t=o.createRef(),[n,r]=o.useState(!1);return o.useEffect(()=>{e.video.current&&(e.video.current.onpause=({target:t})=>{let n=function(e){let t=document.createElement("canvas");t.width=e.videoWidth,t.height=e.videoHeight;let n=t.getContext("2d");if(!n)return;return n.drawImage(e,0,0),t.toDataURL()}(t);n?e.onPause(n):e.onRefuse(),r(!1)},e.video.current.onplay=()=>{e.setHasPlayed(!0),r(!0)})},[e.video.current,t.current]),o.useEffect(()=>{t.current&&(t.current.style.display=n?"none":"block")},[t.current,n]),o.useEffect(()=>{if(!t.current)return;let n=t.current.getContext("2d");if(!n)return;let r=new i.Box(e.videoWidth,e.videoHeight,10,n);const o=({clientX:t,clientY:n})=>{let{top:r,left:o,width:i,height:a}=e.layout;return[e.videoWidth*(t-o)/i,e.videoHeight*(n-r)/a]};t.current.onmousedown=e=>{e.preventDefault(),r.grab(o(e),"click")},t.current.onmousemove=e=>{e.preventDefault(),r.drag(o(e))},t.current.onmouseup=t=>{t.preventDefault(),r.release(e.onSubmit,e.onRefuse)},t.current.onmouseleave=e=>{e.preventDefault(),r.release(()=>{},()=>{})},t.current.onmouseout=e=>{e.preventDefault(),r.release(()=>{},()=>{})},t.current.ontouchstart=e=>{e.preventDefault(),r.grab(o(e.targetTouches[0]),"tap")},t.current.ontouchmove=e=>{e.preventDefault(),r.drag(o(e.targetTouches[0]))},t.current.ontouchend=t=>{t.preventDefault(),r.release(e.onSubmit,e.onRefuse)}},[t.current,e.layout]),o.createElement("canvas",{ref:t,width:e.videoWidth,height:e.videoHeight,style:Object.assign(Object.assign({},e.layout),{position:"absolute",display:"none"})}," ")}},function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0});t.Box=class{constructor(e,t,n,r){this.areaWidth=e,this.areaHeight=t,this.border=n,this.ctx=r,this.dragStartTime=0,this.helpTextInterval=0,this.dragDistance=0,this.lastPosition=[0,0],this.grabbed=[],this.lastCoordinates=[0,0,0,0],this.coordinates=[0,0,0,0],this.type="click/tap",this.dragging=!1}describe(){let[e,t,n,r]=this.coordinates;return{top:Math.min(t,r),bottom:Math.max(t,r),left:Math.min(e,n),right:Math.max(e,n)}}bounds(){let[e,t,n,r]=this.coordinates;return{left:Math.min(e,n),top:Math.min(t,r),width:Math.abs(e-n),height:Math.abs(t-r)}}draw(){this.ctx.clearRect(0,0,this.areaWidth,this.areaHeight);let{top:e,bottom:t,left:n,right:r}=this.describe();this.ctx.fillStyle="rgba(0, 0, 0, .7)",this.ctx.fillRect(0,0,n,t),this.ctx.fillRect(n,0,this.areaWidth,e),this.ctx.fillRect(r,e,this.areaWidth,this.areaHeight),this.ctx.fillRect(0,t,r,this.areaHeight),this.ctx.lineCap="round",this.ctx.lineWidth=this.border/2,this.ctx.strokeStyle="#33b5e5",this.ctx.beginPath(),this.ctx.moveTo(n,e),this.ctx.lineTo(n,t),this.ctx.lineTo(r,t),this.ctx.lineTo(r,e),this.ctx.closePath(),this.ctx.stroke()}grab([e,t],n){this.helpTextInterval&&window.clearInterval(this.helpTextInterval),this.type=n,this.dragging=!0,this.lastCoordinates=this.coordinates,this.dragStartTime=(new Date).getTime(),this.dragDistance=0,this.lastPosition=[e,t];let{top:r,bottom:o,left:i,right:a}=this.describe(),s=i-this.border<=e&&e<=a+this.border,l=r-this.border<=t&&t<=o+this.border,[c,u,d,f]=this.coordinates;if(this.grabbed=[],l){let t=Math.abs(e-c),n=Math.abs(e-d);t<Math.min(this.border,n)?this.grabbed.push(0):n<this.border&&this.grabbed.push(2)}if(s){let e=Math.abs(t-u),n=Math.abs(t-f);e<Math.min(this.border,n)?this.grabbed.push(1):n<this.border&&this.grabbed.push(3)}return 0==this.grabbed.length&&(this.coordinates=[e,t,e,t],this.grabbed=[0,1]),s&&l}drag([e,t]){if(!this.dragging)return;e=Math.max(0,Math.min(e,this.areaWidth)),t=Math.max(0,Math.min(t,this.areaHeight));let[n,r]=this.lastPosition;this.dragDistance+=Math.sqrt((e-n)*(e-n)+(t-r)*(t-r)),this.lastPosition=[e,t],-1!=this.grabbed.indexOf(0)&&(this.coordinates[0]=e),-1!=this.grabbed.indexOf(1)&&(this.coordinates[1]=t),-1!=this.grabbed.indexOf(2)&&(this.coordinates[2]=e),-1!=this.grabbed.indexOf(3)&&(this.coordinates[3]=t),this.draw()}release(e,t){if(this.dragging=!1,this.dragDistance<5&&(new Date).getTime()-this.dragStartTime<200){this.coordinates=this.lastCoordinates;const[n,r]=this.lastPosition,{top:o,bottom:i,left:a,right:s}=this.describe();a-this.border<n&&n<s+this.border&&o-this.border<r&&r<i+this.border?e(this.bounds()):t(),this.draw()}else this.helpTextInterval=window.setTimeout(()=>{})}wasInside(){}}},function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var r=n(16);t.addVideo=r.addVideo,t.getVideos=r.getVideos},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0));n(13);const i=n(15),a=n(11),s=n(18),l=n(19);t.Listing=function(e){let[t,n]=o.useState(0),[r,c]=o.useState(0),[u,d]=o.useState(),[f,h]=o.useState([]);return o.useEffect(()=>{a.getVideos(t,7).then(({total:e,videos:t})=>{h(t),c(e)}).catch(e=>{console.log(e)})},[t]),o.useEffect(()=>u&&e.onVideoSelected(u),[u]),o.createElement("div",{className:"listing"},o.createElement(s.ListSidebar,{onVideoAdded:d}),o.createElement("div",{className:"list-container"},f.map((e,t)=>o.createElement(i.ListRow,{selected:!!u&&e.id==u.id,onSelect:d,key:t,video:e})),o.createElement(l.ListPagenation,{pageSize:7,maxTotal:r,start:t,end:t+f.length,onPageSelected:n})))}},function(e,t,n){var r=n(5),o=n(14);"string"==typeof(o=o.__esModule?o.default:o)&&(o=[[e.i,o,""]]);var i={insert:"head",singleton:!1},a=(r(o,i),o.locals?o.locals:{});e.exports=a},function(e,t,n){(t=n(7)(!1)).push([e.i,"\n* {\n    box-sizing: border-box;\n}\n  \nhtml {\n    font-size: 62.5%;\n}\n\n.list-row{\n    border: 1px solid #DDD;\n    border-radius: 10px;\n    margin: 5px;\n    height: 90px;\n    vertical-align: top;\n    margin: 5px;\n    display: grid;\n    grid-template-columns: 90px 1fr;\n    box-shadow: 1px 1px 3px rgba(0,0,0,0.05);\n}\n\n.list-row-image{\n    width: 90px;\n    height: 90px;\n    object-fit: contain;\n}\n\n.list-row-text{\n    font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif;\n    padding: 15px;\n}\n\n.list-row-tile{\n    font-size: 2rem;\n    margin: 0;\n}\n\n.list-row-stats{\n    display: flex;\n    justify-content: space-between;\n    font-size: .6rem;\n    line-height: 2rem;\n}\n\n.listing{\n    margin: 14px;\n    display: flex;\n    flex-wrap: wrap;\n}\n\n.list-add{\n    font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif;\n    top: 0px;\n    width: 200px;\n    z-index: 10;\n    background-color: white;\n    position: fixed;\n}\n\n.list-add-space{\n    width: 200px;\n    height: 85px;\n}\n\n.list-add-always{\n    margin: auto;\n}\n\n.list-add-always input {\n    height: 30px;\n    width: 160px;\n    border: 1px solid #cdcdcd;\n    border-top-left-radius: 10px;\n    border-bottom-left-radius: 10px;\n    margin: 0px;\n    vertical-align: middle;\n    font-size: 16px;\n    box-shadow: 1px 1px 3px rgba(0,0,0,0.05);\n}\n\n.list-add-always button {\n    height: 30px;\n    width: 40px;\n    border: 1px solid #cdcdcd;\n    border-top-right-radius: 10px;\n    border-bottom-right-radius: 10px;\n    margin: 0px;\n    vertical-align: middle;\n    font-size: 16px;\n    box-shadow: 1px 1px 3px rgba(0,0,0,0.05);\n}\n\n.list-add-error {\n    color: red;\n}\n\n.list-pagenation {\n    display: flex;\n    justify-content: space-between;\n    font-weight: bold;\n    font-size: larger;\n}\n\n@media screen and (max-width: 499px) {\n    .list-add {\n        text-align: center;\n        width: 100%;\n        height: 85px;\n        box-shadow: 1px 1px 3px rgba(0,0,0,0.05);\n    }\n\n    .list-add-extra {\n        display: none;\n    }\n}\n\n@media screen and (max-width: 750px) {\n    .list-container{\n        max-width: 350px;\n        margin: auto;\n    }\n}\n\n@media screen and (min-width: 751px) {\n    .list-container{\n        width: 470px;\n        margin: auto;\n    }\n}",""]),e.exports=t},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0));t.ListRow=function({video:e,selected:t,onSelect:n}){return o.createElement("div",{className:"list-row",onClick:()=>n(e)},o.createElement("img",{className:"list-row-image",src:`/image/${e.thumbnail}`}),o.createElement("div",{className:"list-row-text"},o.createElement("h3",{className:"list-row-title",style:{color:t?"red":"black"}}," ",e.title," "),o.createElement("div",{className:"list-row-stats"},o.createElement("div",null," ",(r=e.duration)>60?`${Math.floor(r/60)}:${r%60<10?"0":""}${r%60}`:r+"s"," "),o.createElement("div",null," Analyzed ",e.clips," frames "),o.createElement("div",null," Youtube "))));var r}},function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0});const r=n(17);t.getVideos=function(e,t){return Promise.resolve({videos:[{clips:0,duration:90,format:"video/mp4",fps:22,id:2,thumbnail:"thumbnail/1U-cgn3cEGA.jpg",title:"This is a test video so I can code without the Internet!",type:3}],total:1})},t.addVideo=function(e,t){return fetch("/videos",{method:"POST",body:JSON.stringify({type:t,url:e}),headers:{Session:r.session}}).then(e=>200!=e.status?e.text().then(e=>Promise.reject(e)):e.json())}},function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),function e(){console.log("Getting session header"),fetch("/session").then(e=>200==e.status?e.text():Promise.reject(e.text())).then(e=>{t.session=e}).catch(t=>{console.log("Could not get session header:",t),setTimeout(e,1e4)})}()},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(0)),i=n(11);t.ListSidebar=function(e){let t=o.createRef(),[n,r]=o.useState("");return o.createElement(o.Fragment,null,o.createElement("div",{className:"list-add-space"}),o.createElement("div",{className:"list-add"},o.createElement("div",{className:"list-add-always"},o.createElement("h3",null,"Add Video"),o.createElement("input",{type:"url",ref:t,placeholder:"https://www.youtube.com/watch?v=",className:"list-url-textbox"}),o.createElement("button",{onClick:()=>{t.current&&(t.current.value?t.current.value.startsWith("https://www.youtube.com")?t.current.value.startsWith("https://www.youtube.com/watch?v=")?i.addVideo(t.current.value,1).then(t=>e.onVideoAdded(t)).catch(e=>r("Server Error - "+e.toString())):r('Error - Please use the "/watch?v= url" format'):r("Error - Please use a full youtube URL"):r("Error - Please put a URL in the box"))}}," Go"),o.createElement("div",{className:"list-add-error"},n)),o.createElement("div",{className:"list-add-extra"},o.createElement("p",null,"This section is reserved for future additions. Some features that might go here include"),o.createElement("ol",null,o.createElement("li",null,"User Login"),o.createElement("li",null,"Clip Review"),o.createElement("li",null,"Total Dataset Visualization")))))}},function(e,t,n){"use strict";var r=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)Object.hasOwnProperty.call(e,n)&&(t[n]=e[n]);return t.default=e,t};Object.defineProperty(t,"__esModule",{value:!0});const o=r(n(20));t.ListPagenation=function(e){let t=e.pageSize<=1?e.start+1+"":e.start+1+" to "+e.end;const n=()=>e.start>0,r=()=>e.end<e.maxTotal;return o.createElement("div",{className:"list-pagenation"},o.createElement("div",{style:{color:n()?"black":"#9a9a9a"},onClick:()=>{n()&&e.onPageSelected(Math.max(0,e.start-e.pageSize))}},"Back"),o.createElement("div",null," ",t," "),o.createElement("div",{style:{color:r()?"black":"#9a9a9a"},onClick:()=>{r()&&e.onPageSelected(Math.min(e.start+e.pageSize,e.maxTotal-e.pageSize))}},"Forward"))}},function(e,t,n){"use strict";e.exports=n(21)},function(e,t,n){"use strict";
/** @license React v16.12.0
 * react.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var r=n(22),o="function"==typeof Symbol&&Symbol.for,i=o?Symbol.for("react.element"):60103,a=o?Symbol.for("react.portal"):60106,s=o?Symbol.for("react.fragment"):60107,l=o?Symbol.for("react.strict_mode"):60108,c=o?Symbol.for("react.profiler"):60114,u=o?Symbol.for("react.provider"):60109,d=o?Symbol.for("react.context"):60110,f=o?Symbol.for("react.forward_ref"):60112,h=o?Symbol.for("react.suspense"):60113;o&&Symbol.for("react.suspense_list");var p=o?Symbol.for("react.memo"):60115,m=o?Symbol.for("react.lazy"):60116;o&&Symbol.for("react.fundamental"),o&&Symbol.for("react.responder"),o&&Symbol.for("react.scope");var v="function"==typeof Symbol&&Symbol.iterator;function b(e){for(var t="https://reactjs.org/docs/error-decoder.html?invariant="+e,n=1;n<arguments.length;n++)t+="&args[]="+encodeURIComponent(arguments[n]);return"Minified React error #"+e+"; visit "+t+" for the full message or use the non-minified dev environment for full errors and additional helpful warnings."}var g={isMounted:function(){return!1},enqueueForceUpdate:function(){},enqueueReplaceState:function(){},enqueueSetState:function(){}},y={};function x(e,t,n){this.props=e,this.context=t,this.refs=y,this.updater=n||g}function w(){}function S(e,t,n){this.props=e,this.context=t,this.refs=y,this.updater=n||g}x.prototype.isReactComponent={},x.prototype.setState=function(e,t){if("object"!=typeof e&&"function"!=typeof e&&null!=e)throw Error(b(85));this.updater.enqueueSetState(this,e,t,"setState")},x.prototype.forceUpdate=function(e){this.updater.enqueueForceUpdate(this,e,"forceUpdate")},w.prototype=x.prototype;var _=S.prototype=new w;_.constructor=S,r(_,x.prototype),_.isPureReactComponent=!0;var E={current:null},O={current:null},j=Object.prototype.hasOwnProperty,M={key:!0,ref:!0,__self:!0,__source:!0};function P(e,t,n){var r,o={},a=null,s=null;if(null!=t)for(r in void 0!==t.ref&&(s=t.ref),void 0!==t.key&&(a=""+t.key),t)j.call(t,r)&&!M.hasOwnProperty(r)&&(o[r]=t[r]);var l=arguments.length-2;if(1===l)o.children=n;else if(1<l){for(var c=Array(l),u=0;u<l;u++)c[u]=arguments[u+2];o.children=c}if(e&&e.defaultProps)for(r in l=e.defaultProps)void 0===o[r]&&(o[r]=l[r]);return{$$typeof:i,type:e,key:a,ref:s,props:o,_owner:O.current}}function R(e){return"object"==typeof e&&null!==e&&e.$$typeof===i}var k=/\/+/g,C=[];function T(e,t,n,r){if(C.length){var o=C.pop();return o.result=e,o.keyPrefix=t,o.func=n,o.context=r,o.count=0,o}return{result:e,keyPrefix:t,func:n,context:r,count:0}}function D(e){e.result=null,e.keyPrefix=null,e.func=null,e.context=null,e.count=0,10>C.length&&C.push(e)}function N(e,t,n){return null==e?0:function e(t,n,r,o){var s=typeof t;"undefined"!==s&&"boolean"!==s||(t=null);var l=!1;if(null===t)l=!0;else switch(s){case"string":case"number":l=!0;break;case"object":switch(t.$$typeof){case i:case a:l=!0}}if(l)return r(o,t,""===n?"."+$(t,0):n),1;if(l=0,n=""===n?".":n+":",Array.isArray(t))for(var c=0;c<t.length;c++){var u=n+$(s=t[c],c);l+=e(s,u,r,o)}else if(null===t||"object"!=typeof t?u=null:u="function"==typeof(u=v&&t[v]||t["@@iterator"])?u:null,"function"==typeof u)for(t=u.call(t),c=0;!(s=t.next()).done;)l+=e(s=s.value,u=n+$(s,c++),r,o);else if("object"===s)throw r=""+t,Error(b(31,"[object Object]"===r?"object with keys {"+Object.keys(t).join(", ")+"}":r,""));return l}(e,"",t,n)}function $(e,t){return"object"==typeof e&&null!==e&&null!=e.key?function(e){var t={"=":"=0",":":"=2"};return"$"+(""+e).replace(/[=:]/g,(function(e){return t[e]}))}(e.key):t.toString(36)}function L(e,t){e.func.call(e.context,t,e.count++)}function V(e,t,n){var r=e.result,o=e.keyPrefix;e=e.func.call(e.context,t,e.count++),Array.isArray(e)?I(e,r,n,(function(e){return e})):null!=e&&(R(e)&&(e=function(e,t){return{$$typeof:i,type:e.type,key:t,ref:e.ref,props:e.props,_owner:e._owner}}(e,o+(!e.key||t&&t.key===e.key?"":(""+e.key).replace(k,"$&/")+"/")+n)),r.push(e))}function I(e,t,n,r,o){var i="";null!=n&&(i=(""+n).replace(k,"$&/")+"/"),N(e,V,t=T(t,i,r,o)),D(t)}function H(){var e=E.current;if(null===e)throw Error(b(321));return e}var A={Children:{map:function(e,t,n){if(null==e)return e;var r=[];return I(e,r,null,t,n),r},forEach:function(e,t,n){if(null==e)return e;N(e,L,t=T(null,null,t,n)),D(t)},count:function(e){return N(e,(function(){return null}),null)},toArray:function(e){var t=[];return I(e,t,null,(function(e){return e})),t},only:function(e){if(!R(e))throw Error(b(143));return e}},createRef:function(){return{current:null}},Component:x,PureComponent:S,createContext:function(e,t){return void 0===t&&(t=null),(e={$$typeof:d,_calculateChangedBits:t,_currentValue:e,_currentValue2:e,_threadCount:0,Provider:null,Consumer:null}).Provider={$$typeof:u,_context:e},e.Consumer=e},forwardRef:function(e){return{$$typeof:f,render:e}},lazy:function(e){return{$$typeof:m,_ctor:e,_status:-1,_result:null}},memo:function(e,t){return{$$typeof:p,type:e,compare:void 0===t?null:t}},useCallback:function(e,t){return H().useCallback(e,t)},useContext:function(e,t){return H().useContext(e,t)},useEffect:function(e,t){return H().useEffect(e,t)},useImperativeHandle:function(e,t,n){return H().useImperativeHandle(e,t,n)},useDebugValue:function(){},useLayoutEffect:function(e,t){return H().useLayoutEffect(e,t)},useMemo:function(e,t){return H().useMemo(e,t)},useReducer:function(e,t,n){return H().useReducer(e,t,n)},useRef:function(e){return H().useRef(e)},useState:function(e){return H().useState(e)},Fragment:s,Profiler:c,StrictMode:l,Suspense:h,createElement:P,cloneElement:function(e,t,n){if(null==e)throw Error(b(267,e));var o=r({},e.props),a=e.key,s=e.ref,l=e._owner;if(null!=t){if(void 0!==t.ref&&(s=t.ref,l=O.current),void 0!==t.key&&(a=""+t.key),e.type&&e.type.defaultProps)var c=e.type.defaultProps;for(u in t)j.call(t,u)&&!M.hasOwnProperty(u)&&(o[u]=void 0===t[u]&&void 0!==c?c[u]:t[u])}var u=arguments.length-2;if(1===u)o.children=n;else if(1<u){c=Array(u);for(var d=0;d<u;d++)c[d]=arguments[d+2];o.children=c}return{$$typeof:i,type:e.type,key:a,ref:s,props:o,_owner:l}},createFactory:function(e){var t=P.bind(null,e);return t.type=e,t},isValidElement:R,version:"16.12.0",__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED:{ReactCurrentDispatcher:E,ReactCurrentBatchConfig:{suspense:null},ReactCurrentOwner:O,IsSomeRendererActing:{current:!1},assign:r}},z={default:A},U=z&&A||z;e.exports=U.default||U},function(e,t,n){"use strict";
/*
object-assign
(c) Sindre Sorhus
@license MIT
*/var r=Object.getOwnPropertySymbols,o=Object.prototype.hasOwnProperty,i=Object.prototype.propertyIsEnumerable;function a(e){if(null==e)throw new TypeError("Object.assign cannot be called with null or undefined");return Object(e)}e.exports=function(){try{if(!Object.assign)return!1;var e=new String("abc");if(e[5]="de","5"===Object.getOwnPropertyNames(e)[0])return!1;for(var t={},n=0;n<10;n++)t["_"+String.fromCharCode(n)]=n;if("0123456789"!==Object.getOwnPropertyNames(t).map((function(e){return t[e]})).join(""))return!1;var r={};return"abcdefghijklmnopqrst".split("").forEach((function(e){r[e]=e})),"abcdefghijklmnopqrst"===Object.keys(Object.assign({},r)).join("")}catch(e){return!1}}()?Object.assign:function(e,t){for(var n,s,l=a(e),c=1;c<arguments.length;c++){for(var u in n=Object(arguments[c]))o.call(n,u)&&(l[u]=n[u]);if(r){s=r(n);for(var d=0;d<s.length;d++)i.call(n,s[d])&&(l[s[d]]=n[s[d]])}}return l}}]);