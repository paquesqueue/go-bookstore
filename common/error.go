package common

import (
	"strconv"
	"strings"
)

type Err struct {
	Code     int
	Remark   string
	Original error
}

func (e Err) Error() string {
	return strings.Join([]string{strconv.Itoa(e.Code), e.Remark}, " : ")
}
