# tool

日常开发用得到的工具




## 1.检查端口占用情况

```shell
./tool cp 8080 
```



## 2.kafka测试

```shell
./tool kt 127.0.0.1:9092 #主题默认为 test
```

也可以指定主题

```shell
./tool kt 127.0.0.1:9092 hello
```



## 3.向kafka发送自定义数据

```shell
./tool ks 127.0.0.1:9092 hello '{"a":"1","b":"2"}'
```



## 4.从kafka接受数据

```shell
./tool kr 127.0.0.1:9092 hello
```



## 5.运行指定秒数

```shell
./tool rs 60 succeed
```

```
./tool rs 秒数 [succeed|failed|exit|panic|loop|maxLoad]
```

