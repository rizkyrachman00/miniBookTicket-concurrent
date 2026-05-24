<!-- do test (cached)-->
go test ./... -v

<!-- do test (no-cached) -->
go test ./... -v -count=1

<!-- do test (no-cached ~ race condition detection) -->
go test ./... -v -count=1 -race