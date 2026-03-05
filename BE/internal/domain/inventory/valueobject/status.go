package valueobject

// StockInStatus represents the state of a stock-in transaction.
type StockInStatus string

const (
	StockInCreated    StockInStatus = "CREATED"
	StockInInProgress StockInStatus = "IN_PROGRESS"
	StockInDone       StockInStatus = "DONE"
	StockInCancelled  StockInStatus = "CANCELLED"
)

var stockInTransitions = map[StockInStatus][]StockInStatus{
	StockInCreated:    {StockInInProgress, StockInCancelled},
	StockInInProgress: {StockInDone, StockInCancelled},
}

// CanTransitionTo returns true if the transition from s to next is valid.
func (s StockInStatus) CanTransitionTo(next StockInStatus) bool {
	allowed, ok := stockInTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == next {
			return true
		}
	}
	return false
}

// StockOutStatus represents the state of a stock-out transaction.
type StockOutStatus string

const (
	StockOutDraft      StockOutStatus = "DRAFT"
	StockOutInProgress StockOutStatus = "IN_PROGRESS"
	StockOutDone       StockOutStatus = "DONE"
	StockOutCancelled  StockOutStatus = "CANCELLED"
)

var stockOutTransitions = map[StockOutStatus][]StockOutStatus{
	StockOutDraft:      {StockOutInProgress, StockOutCancelled},
	StockOutInProgress: {StockOutDone, StockOutCancelled},
}

// CanTransitionTo returns true if the transition from s to next is valid.
func (s StockOutStatus) CanTransitionTo(next StockOutStatus) bool {
	allowed, ok := stockOutTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == next {
			return true
		}
	}
	return false
}

// ReservationStatus represents the state of a reservation.
type ReservationStatus string

const (
	ReservationActive   ReservationStatus = "ACTIVE"
	ReservationConsumed ReservationStatus = "CONSUMED"
	ReservationReleased ReservationStatus = "RELEASED"
)
