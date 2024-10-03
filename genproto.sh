git submodule update --remote proto-definitions

protoc -I=proto-definitions/ \
    --go_out=./proto --go_opt=paths=source_relative --go_opt=Mgame.proto=gitlab.com/hackathon-rainbow-2024/go-client/proto \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mgame.proto=gitlab.com/hackathon-rainbow-2024/go-client/proto \
    proto-definitions/game.proto