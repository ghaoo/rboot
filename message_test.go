package rboot

import (
	"testing"
	"bytes"
	"io/ioutil"
)

var msgBody = `From: John Doe <jdoe@machine.example>
To: Mary Smith <mary@example.net>
Subject: Saying Hello
Date: Fri, 21 Nov 1997 09:55:06 -0600
Message-ID: <1234@local.machine.example>

This is a message just to say hello.
So, "Hello".`

func TestReadMessage(t *testing.T) {
	msg, err := ReadMessage(bytes.NewBuffer([]byte(msgBody)))
	if err != nil {
		t.Errorf("Failed read message: %v", err)
	}

	body, err := ioutil.ReadAll(msg.Body)

	if err != nil {
		t.Errorf("Failed reading body: %v", err)
	}

	t.Logf(`Body: %s`, string(body))

}

