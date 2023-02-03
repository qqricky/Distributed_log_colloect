# 分布式日志收集系统

简易分布式日志收集系统。跟着B站七米老师所做。并且修改了一些老师的错误。

## 启动

依次启动（后面跟着对应的版本）

zookeeper、kafka_2.12-2.3.0、etcd-v3.3.13-windows-amd64、elasticsearch-7.2.1、kibana-7.2.1-windows-x86_64

同时电脑需要有jdk-12

往访问的etcd节点插入配置(使用json格式)

```bash
value := `[{"path":"f:/tmp/nginx.log","topic":"web_log"},{"path":"f:/xxx/redis.log","topic":"redis_log"},{"path":"f:/xxx/mysql.log","topic":"mysql_log"}]`
```

path代表所收集的主机的日志文件的路径