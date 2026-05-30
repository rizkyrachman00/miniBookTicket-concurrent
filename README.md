<!-- do test (cached)-->
go test ./... -v

<!-- do test (no-cached) -->
go test ./... -v -count=1

<!-- do test (no-cached ~ race condition detection) -->
go test ./... -v -count=1 -race


The idea is :

<!-- Current behavior -->
Current Code :
Book immediately
Seat becomes taken forever in memory

<!-- Todo -->
A TTL version would do this:
Hold seat first
Set expiration time
Release seat if expired


use redis commander for interact with Redis


<!-- --- How To Run Project --- -->

run redis container / docker compose
go run cmd/main.go

serve at http://localhost:8080/ 