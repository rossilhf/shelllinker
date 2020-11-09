# 远程设备管理工具--shelllinker“谢凌客”

linux shell version of teamviewer

2020.11.9

本工具可以进行大批量的智能设备(目前仅支持linux，部分支持windows)的远程管理和远程控制，主要用于在设备故障时，远程控制调试。可以看作是命令行模式的teamviewer。

本工具组成结构如下图所示：

![image](https://github.com/rossilhf/shelllinker/blob/master/shelllinker_stru.png)

本工具使用流程如下图所示：

![image](https://github.com/rossilhf/shelllinker/blob/master/shelllinker_howtouse.png)

分成三个独立部分：


# 设备端：

设备端代码每次启动时，都会主动连接指定mqtt服务器，执行三类工作：

1.平时定时上报自身mac地址/产品序列号等信息；

2.收听mqtt消息。当服务器端发送来执行消息，设备端都将消息当作系统命令进行执行，并返回控制台打印信息给服务器端。因此可以执行pwd，ls，cat，ifconfig等简单命令；

3.收听mqtt消息。当服务器端发送来升级消息，设备端去指定web服务器上下载最新版本的客户端程序，进行自我更新。

设备端必须配置：
1.安装sshpass
2./etc/ssh/ssh_config 中，StrictHostKeyChecking 设置为no


# 服务器(控制)端：

服务器端有三个独立进程。

进程1：该进程永久运行。监听指定mqtt服务器的消息队列，每当有设备上报自身信息时，将设备信息写入本地数据库；

进程2：该进程永久运行。定期检查当前程序最新版本号，检查当前设备数据库，查看所有活跃设备的当前版本号。给需要升级版本的设备发送升级指令。

进程3：该进程手动开启和关闭，用于账号登录、设备列表查看、指定设备交互控制。进程启动后，首先进入本地命令状态，支持若干私有命令，如“.help”,".list",".echo"；使用.echo xxxx命令进入远程命令控制状态，即可以与设备端进行命令行交互(但有不少限制，如不能sudo，top等)。使用私有命令".help"可以查看如何解除限制。使用私有命令".exit"退出远程命令控制状态，回到本地命令控制状态。


# 服务器(web资源)端：

开启http服务，提供最新版本客户端程序的下载服务


# To Do:

1.完善设备端查找本机mac和ip功能的鲁棒性。用ip addr 补充 ifconfig

2.服务端需完善对大量设备的查看功能，如按某项排序等

3.根据系统开放的端口，选择合适的通信协议方式

4.完善对其他系统的支持，如windows，android，ios
