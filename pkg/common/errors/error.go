package errors

type CohesiveError struct {
	Message string
}

func (c CohesiveError) Error() string {
	return c.Message
}
