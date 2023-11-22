des_key = 12345678

serverx:
	DES_KEY=$(des_key) go run server/main.go server/exec.go server/tools.go

clientx:
	DES_KEY=$(des_key) go run client/main.go
