1.生成proto
```
jupiter protoc -f ./user.proto -o ./user -g
```

2. 生成grpc服务端实现
```
jupiter protoc -f ./user.proto -o ../internal/app/grpc -p vid_user -s
```

