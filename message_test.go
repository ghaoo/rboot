package rboot

import (
	"testing"
)

var msgBody = `From: John Doe <jdoe@machine.example>
To: Mary Smith <mary@example.net>
Subject: Saying Hello
Date: Fri, 21 Nov 1997 09:55:06 -0600
Message-ID: <1234@local.machine.example>

Hello.`

func TestNewStringMessage(t *testing.T) {
	msg := NewStringMessage(msgBody)

	content := msg.Content

	if content != `Hello.` {
		t.Logf(`new string message error`)
	}

	header := msg.Header

	if header.Get(`From`) != `John Doe <jdoe@machine.example>` {
		t.Logf(`new string message error, header not fond`)
	}

}
