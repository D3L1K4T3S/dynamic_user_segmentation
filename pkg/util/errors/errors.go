package errors

import "fmt"

// Wrap оборачивает ошибку в сообщение переданное функции
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapIfErr оборачивает ошибку в сообщение только в том случаем, если она есть
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}
