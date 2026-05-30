package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

// RedisStore implements session-based seat booking backed by Redis.
//
// Key design:
//
//	seat:{movieID}:{seatID}   → session JSON (TTL = held, no TTL = confirmed)
//	session:{sessionID}       → seat key     (reverse lookup)
type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

// sessionKey builds the reverse-lookup key for a session.
func sessionKey(id string) string {
	return fmt.Sprintf("Session: %s", id)
}

func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)

	if err != nil {
		return Booking{}, err
	}

	log.Printf("Session booked %v", session)

	return session, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	pattern := fmt.Sprintf("seat:%s:*", movieID) // Build Redis key pattern to find all seat bookings for this movie ID
	var sessions []Booking

	ctx := context.Background()

	// Create an iterator to loop through Redis keys matching the pattern
	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result() // Get the value stored at the current Redis key

		if err != nil {
			continue
		}

		session, err := parseSession(val) // Convert Redis JSON string into Booking struct
		if err != nil {
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	b.ID = id
	val, _ := json.Marshal(b)
	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX", // only set if key does not already exist
		TTL:  defaultHoldTTL,
	})

	oke := res.Val() == "OK"
	if !oke {
		return Booking{}, ErrSeatAlreadyBooked
	}

	s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}

func parseSession(val string) (Booking, error) {
	var data Booking

	// Convert string to []byte, then decode JSON into data
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}

	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}
