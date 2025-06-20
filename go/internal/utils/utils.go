package utils

import (
	"log"
	"strconv"
)

// Проверяет наличие и корректность переменных окружения, указывающих продолжительность
// JWT-токенов и в случае успеха возвращает указанное число или число по умолчанию (60) при отсутствии
func CheckJWTDuration(duration string, isAdmin bool) int {
	var user string
	if isAdmin {
		user = "ADMIN"
	} else {
		user = "USER"
	}

	if duration == "" {
		log.Println("ПРЕДУПРЕЖДЕНИЕ: в переменной окружения JWT_" + user + "_DURATION, которая указывает на " +
			"продолжительность в минутах сохранения JWT-токенов, было пропущено значение, в следствии чего будет " +
			"использоваться значение по умолчанию - 60 минут")
		return 60
	} else {
		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			log.Fatalf("ошибка окружения: необходимо указывать значение JWT_" + user + "_DURATION в формате int (в минутах)")
		}
		return durationInt
	}
}
