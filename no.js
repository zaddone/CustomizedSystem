function uigs_pv(d){
try{
if(uigs_para&&"undefined"!=typeof uigs_para){
var b=["undefined"!=typeof httpsUtil?httpsUtil.getPingbackHost():"http://pb.sogou.com","/pv.gif?uigs_t=",(new Date).getTime()],a;for(a in uigs_para)uigs_para.hasOwnProperty(a)&&b.push("&"+a+"="+encodeURIComponent(uigs_para[a]));b.push("&uigs_refer=");b.push(encodeURIComponent(document.referrer));d&&(b.push("&"),b.push(d));(new Image).src=b.join("")}}catch(c){}}function uigs_cl(d,b){try{if(uigs_para&&"undefined"!=typeof uigs_para){var a=["undefined"!=typeof httpsUtil?httpsUtil.getPingbackHost():"http://pb.sogou.com","/cl.gif?uigs_t=",(new Date).getTime()],c;for(c in uigs_para)uigs_para.hasOwnProperty(c)&&a.push("&"+c+"="+encodeURIComponent(uigs_para[c]));a.push("&uigs_cl=");a.push(d);a.push("&uigs_refer=");a.push(encodeURIComponent(document.referrer));b&&(a.push("&"),a.push(b));(new Image).src=a.join("")}}catch(e){}};
