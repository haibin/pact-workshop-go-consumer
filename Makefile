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
	go test -v -tags=integration -count=1 github.com/haibin/pact-workshop-go-consumer/consumer/client -run 'TestClientPact'

publish: install
	@echo "--- 📝 Publishing Pacts"
	go run consumer/client/pact/publish.go
	@echo
	@echo "Pact contract publishing complete!"
	@echo
	@echo "Head over to $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL) and login with $(PACT_BROKER_USERNAME)/$(PACT_BROKER_PASSWORD)"
	@echo "to see your published contracts.	"
