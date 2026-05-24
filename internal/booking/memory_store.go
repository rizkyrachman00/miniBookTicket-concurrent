package booking

type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	// seat is taken
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	// populate
	s.bookings[b.SeatID] = b

	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var results []Booking

	for _, b := range s.bookings {
		if b.MovieID == movieID {
			results = append(results, b)
		}
	}

	return results
}
