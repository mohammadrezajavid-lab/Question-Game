```bash
protoc --proto_path=contract/protobuf/presence --go_out=contract/goprotobuf/presence --go_opt=paths=source_relative ./contract/protobuf/presence/presence.proto

protoc --proto_path=contract/protobuf/matching --go_out=contract/goprotobuf/matching --go_opt=paths=source_relative ./contract/protobuf/matching/matching.proto

protoc --proto_path=contract/protobuf/notification --go_out=contract/goprotobuf/notification --go_opt=paths=source_relative ./contract/protobuf/notification/notification.proto

protoc --proto_path=contract/protobuf/game --go_out=contract/goprotobuf/game --go_opt=paths=source_relative ./contract/protobuf/game/created_game.proto
```

#### --proto_path=[مسیر فایل های proto]

#### --go_out=[محل ساخت فایل pb.go]

#### --go_opt=paths=source_relative [کد تولید شده را در همان مسیری ذخیره کن که فایل .proto نسبت به --proto_path داشته]

#### --go_opt=paths=import [اگر مقدار import داده شود اون وقت از مسیر option go_package برای اون استفاده میکنه]

```bash
protoc --proto_path=contract/protobuf/presence --go_out=contract/goprotobuf/presence --go_opt=paths=source_relative --go-grpc_out=contract/goprotobuf/presence --go-grpc_opt=paths=source_relative ./contract/protobuf/presence/presence.proto
```

### profiling commands
```bash
curl http://127.0.0.1:2111/debug/pprof/goroutine --output goroutine.tar.gz
```
```bash
go tool pprof -http=:1999 ./goroutine.tar.gz
```
- To collect an execution trace, use the debug/pprof/trace endpoint:
```bash
curl http://127.0.0.1:2111/debug/pprof/trace?seconds=30 --output trace.prof
```
- To analyze the collected trace, use the go tool trace command:
```bash
go tool trace ./trace.prof
```

- CPU Profiling
```bash
curl http://127.0.0.1:2111/debug/pprof/profile?seconds=30 --output cpu.prof
go tool pprof -http=:8090 cpu.prof
```

- Memory Profiling
```bash
go tool pprof -http=:8090 http://127.0.0.1:2111/debug/pprof/heap
```

- Test WebSocket
```bash
 websocat --header="Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE0MzcwMjAsIlN1YmplY3QiOiJhdCIsInVzZXJfaWQiOjIsInJvbGUiOjF9.vwc5RRE_yOhVy-o9n7aSJ5cmCG5sRc9GztByZHUC1q8 --header="Origin: http://127.0.0.1:3000" ws://172.20.10.2:8090/ws
```