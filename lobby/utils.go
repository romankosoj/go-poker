package lobby

import (
	"math/rand"
	"strings"
)

const idLength = 5

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateLobbyID() string {
	sb := strings.Builder{}
	sb.Grow(idLength)
	for i := 0; i < idLength; i++ {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}
	return sb.String()
}
