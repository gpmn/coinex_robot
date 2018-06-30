GOOS=windows GOARCH=amd64 go build robot.go
mv robot.exe robot_win.exe

go build robot.go
mv robot robot_linux

GOOS=darwin GOARCH=amd64 go build robot.go
mv robot robot_darwin_amd64

GOOS=darwin GOARCH=386 go build robot.go
mv robot robot_darwin_386
