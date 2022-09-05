APP:=bolter
APP_ENTRY_POINT:=./cmd/fire.go

fire:
	MallocNanoZone=0 go run -race $(APP_ENTRY_POINT) fire



