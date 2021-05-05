# SplitVPN
Split Internet and VPN routing

**ATTENTION:** This is a quick and dirty implementation

## Install with brew
```sh
# install
brew tap jurjevic/tap
brew install splitvpn

# start
sudo splitvpn
```

## Install with Go
```sh
# install
go install github.com/jurjevic/SplitVPN@latest

# start
sudo /bin/sh -c "$(go env GOPATH)/bin/SplitVPN &"
```
