package utils

import (
	"strings"
	"user/constants"
)

func RedisKey(suffixes ...string) string { //["unik","users"]
	prefix := []string{constants.RedisPrefix} //prefix = ["user"]
	//["user","unik","users"]
	return strings.Join(append(prefix, suffixes...), "-") //prefix = user-unik-users
}
