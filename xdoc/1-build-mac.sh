rm -rf build
mkdir build

go-bindata -o=res/res.go -pkg=res res/ res/en res/zh

CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o build/zd-mac src/zd.go