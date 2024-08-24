# 1. 安装 Docker CE（即 Docker 社区版）

# 2. 配置 Docker 开机自启

```shell
sudo systemctl enable docker
sudo systemctl start docker
```

查看docker信息

```shell
docker info
```

# 3. 安装 Docker Compose

下载Docker Compose

```shell
curl -L https://get.daocloud.io/docker/compose/releases/download/1.25.4/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
```

配置执行权限

```shell
sudo chmod +x /usr/local/bin/docker-compose
```

检查是否安装成功

```shell
docker-compose -v
```
