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