# 一个简单的DNS压测工具

1、在使用多个源地址进行压测时，需要先完成本地网络地址的配置，压测工具所能使用的网络地址是参数指定的网段与本地地址的交集

2、domain参数根据给定的域随机构建子域名，该参数所构建的DNS请求通常是不合法的，以达到垃圾请求压测的目的，同时也可以使用number参数限定构建子域名的数量

3、hexfile参数指定一个文件路径，文件内容必须是DNS查询请求部分的十六进制数据的文本文件，一般通过Wireshark分析器获取DNS查询请求部分

4、domainfile参数指定一个文件路径，文件内容类似于 "www.xxx.com.cn A" 的文本文件，同时也可以使用rdvalue参数设置DNS查询请求中的rd标志位

5、pcapfile参数指定一个文件路径，文件内容是通过tcpdump命令或者Wireshark等工具捕获的网络数据文件，压测工具会分析文件内容并获取其中的DNS查询请求数据