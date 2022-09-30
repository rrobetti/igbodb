# igbodb
## Object oriented database

## To compile the Proto Buffers definitionn 
`protoc grpc/*.proto     --go_out=.     --go_opt=paths=source_relative     --proto_path=.`

To generate client and server grpc code
`protoc grpc/*.proto     --go_out=.     --go_opt=paths=source_relative     --go-grpc_out=.     --go-grpc_opt=paths=source_relative     --proto_path=.`

