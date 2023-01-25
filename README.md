# MiniDouyin

## 抖音项目服务端

### 配置

version：

**MySQL8.0**

**Go 1.19**


```shell
mysql -uroot -p"password" < example.sql
```

### 启动

工程无其他依赖，直接编译运行即可

```shell
go build && ./MiniDouyin
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试

### 目录说明

`repository` 使用MySQL数据库，定义各个数据结构

`service` 处理逻辑

`controller` 接收客户端信息


