package Email

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"github.com/ThePiachu/Go/Log"
	"golang.org/x/net/context"
	"google.golang.org/appengine/v2/mail"
)

func SendHTMLEmail(c context.Context, subject string, to []string, sender string, mailBody string) error {
	msg := &mail.Message{
		Sender:   sender,
		To:       to,
		Subject:  subject,
		HTMLBody: mailBody,
	}
	err := mail.Send(c, msg)
	if err != nil {
		Log.Errorf(c, "SendHTMLEmail - %v", err)
		return err
	}
	return nil
}

func SendEmail(c context.Context, subject string, to []string, sender string, mailBody string) error {
	msg := &mail.Message{
		Sender:  sender,
		To:      to,
		Subject: subject,
		Body:    mailBody,
	}
	err := mail.Send(c, msg)
	if err != nil {
		Log.Errorf(c, "Email - SendHTMLEmail - %v", err)
		return err
	}
	return nil
}
