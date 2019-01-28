package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"testing"
	"time"
)

func TestSendTime(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	sendedTime := time.Unix(0, 0)
	sendTime(buf, sendedTime)

	bts, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Error(err)
		return
	}

	if lenbts := len(bts); lenbts > 4 || lenbts < 4 {
		t.Errorf("len of readed bytes is: %d expected: %d\n", lenbts, 4)
		return
	}

	from1970 := binary.BigEndian.Uint32(bts)
	givenTime := time.Unix(int64(from1970-epochDelta), 0)
	if givenTime != sendedTime {
		t.Errorf("given time is: %v expected: %v\n", givenTime, sendedTime)
		return
	}
}
