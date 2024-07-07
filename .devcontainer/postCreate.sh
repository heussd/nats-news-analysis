go get ./...
#go install github.com/cosmtrek/air@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/nats-io/natscli/nats@latest


echo "LC_ALL=en_US.UTF-8" >> /etc/environment
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
echo "LANG=en_US.UTF-8" > /etc/locale.conf

apt-get clean && apt-get update
apt-get install -y locales

locale-gen en_US.UTF-8