# 魅族相册同步脚本

采用魅族云相册接口，进行相册备份。

## 使用

1、下载编译好的程序

2、在 `Console` 中输入下面的代码获取 [云相册](https://photos.flyme.cn/photo/albums) 的token

```javascript
function getCookie(name){var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");if(arr=document.cookie.match(reg)){return unescape(arr[2])}else{return null}}getCookie("_utoken");
```

3、将token放在程序同级目录下的 `token.txt` 中

4、解压并运行程序

```shell
# Linux && MacOS
./flyme-album

# Windows
.\flyme-album.exe
```

5、等待程序响应完成

> 魅族的相册数据是放在阿里云上的，所以速度会很快。我一年多张图片用了5分钟不到就同步完成了。

## 参考资料

* [知乎-魅族云相册的照片怎么一次性全部下载？](https://www.zhihu.com/question/66221241/answer/2078686584)
* [油猴插件 - 获取魅族云空间的token](https://openuserjs.org/scripts/moreantfoxmail.com/copy-flyme-photo-token)
* [github - vue批量下载魅族云相册](https://github.com/moreant/mpcb)
* [github - php备份魅族云相册](https://github.com/dingdayu/mzstorage)