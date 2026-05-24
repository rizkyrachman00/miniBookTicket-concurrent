package booking

type Service struct {
	// Dependency Injection
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) Book(b Booking) error {
	return s.store.Book(b)
}
