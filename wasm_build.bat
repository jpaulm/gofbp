set GOOS=js
set GOARCH=wasm
rem go build -o wasm/merge.wasm merge.go
go build -o wasm/components/testrtn/kick.wasm components/testrtn/kick.go
pause