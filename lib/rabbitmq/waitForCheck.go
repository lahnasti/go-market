package rabbitmq

import "fmt"

func (r *RabbitMQ) WaitForUserCheckResponse() (bool, error) {
	// Потребляем сообщения из очереди ответов
	messages, err := r.Channel.Consume(
		"user_check_response", // Имя очереди с ответами
		"",                    // consumer tag
		true,                  // auto-ack, автоматическое подтверждение
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // аргументы
	)
	if err != nil {
		return false, fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Ожидаем получения сообщения из очереди
	for msg := range messages {
		response := string(msg.Body)
		if response == "user exists" {
			return true, nil
		} else if response == "user not found" {
			return false, nil
		} else {
			return false, fmt.Errorf("unknown response: %s", response)
		}
	}
	return false, fmt.Errorf("no message received")
}
