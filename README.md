# SplitVPN
Split Internet and VPN routing

**ATTENTION:** This is a quick and dirty implementation

```yaml ###  repeat = 1; ftoken[2] = "\"" + NewVersion + "\""; Join(ftoken, " ")
Version: 0.9.9
```

## Install with brew
```sh
# install
brew tap jurjevic/tap
brew install splitvpn
```
```sh
# start
sudo -b splitvpn
```
```sh
# update
brew upgrade splitvpn
```

## Install with Go 1.16+
```sh
# install
go install github.com/jurjevic/SplitVPN@latest
```
```sh
# start
sudo /bin/sh -c "$(go env GOPATH)/bin/SplitVPN &"
```
