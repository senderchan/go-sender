# Sender Agent

信达Agent

## 获取

### 安装脚本

运行前请确认已经安装`wget`和`curl`。

```shell
sudo sh -c "$(curl -sSL https://raw.githubusercontent.com/zhshch2002/sender-agent/master/install.sh)"
```

脚本自动创建了一个配置文件，需要进行配置。配置方式参考[https://sender.xzhsh.ch/docs](https://sender.xzhsh.ch/docs)。

```shell
sudo nano /etc/sender/agent/config.yaml
```

随后，启动服务即可。

```shell
sudo systemctl start sender
```

### 二进制程序

在 [Release](https://github.com/zhshch2002/sender-agent/releases) 页面下载Agent的二进制程序。

解压压缩文件后，在压缩文件同一目录下创建`config.yaml`。配置方式参考[https://sender.xzhsh.ch/docs](https://sender.xzhsh.ch/docs)。

```shell
nano ./config.yaml
```

随后，启动程序即可。

```shell
./sender-agent
# nohup ./sender-agent >> ./output.log 2>&1 &
```
