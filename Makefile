all:
	go build -o go-assumerole main.go
	cp -a go-assumerole ../org/go-assumerole
