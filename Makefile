include ./make/config.mk

install:
	@if [ ! -d pact/bin ]; then\
		echo "--- 🛠 Installing Pact CLI dependencies";\
		curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash;\
    fi

run-consumer:
	@go run consumer/client/cmd/main.go

unit:
	@echo "--- 🔨Running Unit tests "
	go test -count=1 github.com/haibin/pact-workshop-go-consumer/consumer/client -run 'TestClientUnit'

consumer: install
	@echo "--- 🔨Running Consumer Pact tests "
	go test -tags=integration -count=1 github.com/haibin/pact-workshop-go-consumer/consumer/client -run 'TestClientPact'
