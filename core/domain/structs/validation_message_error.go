package structs

type ValidationMessagesError struct {
	Messages []string
}

func (e *ValidationMessagesError) Error() string {
	if len(e.Messages) > 0 {
		return e.Messages[0]
	}
	return "Erro de validação"
}
