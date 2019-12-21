`多线程并发下载器-gorc`
------------------------------------
gorc是类wget多线程下载器，支持直接从资源url并发获取资源

#### 使用说明：</br>
##### 1.手动选择模式和自动分配模式，参数：manual，默认为false/自动</br>
##### 2.指定并发线程数，参数：thread，默认为5</br>
##### 3.指定下载的url，参数：url</br>
##### 4.指定分块下载的块大小，参数：blockSize,例如，默认1代表16m,2代表32m,4代表64m，以此类推</br>
##### 5.指定分块下载失败后尝试次数，参数：attempt，默认为3</br>
##### 6.指定文件存放位置，参数：root,默认为项目的lib目录</br>
##### 7.程序使用秩序调用gorc.Download(url string)函数即可</br>

#### 功能点：</br>
##### 1.支持多线程并发下载</br>
##### 2.支持断点续传</br>
##### 3.支持进度条显示</br>
##### 4.支持手动设置临时文件大小</br>
##### 5.支持自动清理缓存文件</br>

#### 效果示例：
Windows7下</br>
自动模式</br>
![](https://github.com/V-I-C-T-O-R/gorc/blob/master/pic/windows_auto.png)</br>
自动续传模式</br>
![](https://github.com/V-I-C-T-O-R/gorc/blob/master/pic/windows_auto_comsu.png)</br>
手动模式</br>
![](https://github.com/V-I-C-T-O-R/gorc/blob/master/pic/windows_manu.png)</br>
手动续传模式</br>
![](https://github.com/V-I-C-T-O-R/gorc/blob/master/pic/windows_manu_consumer.png)
