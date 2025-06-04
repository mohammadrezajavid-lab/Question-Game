```bash
protoc --proto_path=contract/protobuf/presence --go_out=contract/goprotobuf/presence --go_opt=paths=source_relative ./contract/protobuf/presence/presence.proto

protoc --proto_path=contract/protobuf/matching --go_out=contract/goprotobuf/matching --go_opt=paths=source_relative ./contract/protobuf/matching/matching.proto
```

#### --proto_path=[مسیر فایل های proto]

#### --go_out=[محل ساخت فایل pb.go]

#### --go_opt=paths=source_relative [کد تولید شده را در همان مسیری ذخیره کن که فایل .proto نسبت به --proto_path داشته]

#### --go_opt=paths=import [اگر مقدار import داده شود اون وقت از مسیر option go_package برای اون استفاده میکنه]

```bash
protoc --proto_path=contract/protobuf/presence --go_out=contract/goprotobuf/presence --go_opt=paths=source_relative --go-grpc_out=contract/goprotobuf/presence --go-grpc_opt=paths=source_relative ./contract/protobuf/presence/presence.proto
```