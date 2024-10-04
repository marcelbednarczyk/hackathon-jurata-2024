git submodule update --remote proto-definitions

protoc -I=proto-definitions/ \
    --go_out=./proto --go_opt=paths=source_relative --go_opt=Mgame.proto=github.com/marcelbednarczyk/hackathon-jurata-2024/proto \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mgame.proto=github.com/marcelbednarczyk/hackathon-jurata-2024/proto \
    proto-definitions/game.proto