package flash

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	Error  = "error"
	Info   = "info"
	Sucess = "success"

	cookieName = "flash"
)

type Message struct {
	Level   string
	Content string
}

func Set(w http.ResponseWriter, msg Message) {
	encodedContent, err := encode(msg)
	if err != nil {
		slog.Warn("Unable to encode msg, no cookie was created", "error", err)
		return
	}

	c := &http.Cookie{Name: cookieName, Path: "/", Value: encodedContent}
	http.SetCookie(w, c)
}

func Get(w http.ResponseWriter, r *http.Request) Message {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return Message{}
	}

	msg, err := decode(c.Value)
	if err != nil {
		slog.Warn("Unable to decode flash cookie", "error", err)
		return Message{}
	}

	expiredCookie := &http.Cookie{Name: cookieName, Path: "/", MaxAge: -1, Expires: time.Unix(0, 0)}
	http.SetCookie(w, expiredCookie)
	return msg
}

func encode(msg Message) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return "", fmt.Errorf("unable to encode msg into buf: %w", err)
	}

	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func decode(encodedContent string) (Message, error) {
	content, err := base64.URLEncoding.DecodeString(encodedContent)
	if err != nil {
		return Message{}, fmt.Errorf("unable to decode base64 content string: %w", err)
	}

	var msg Message
	dec := gob.NewDecoder(bytes.NewReader(content))
	if err := dec.Decode(&msg); err != nil {
		return Message{}, fmt.Errorf("unable to decode content into msg: %w", err)
	}

	return msg, nil
}

func EntryCreated(w http.ResponseWriter) {
	Set(w, Message{Level: Sucess, Content: "Success! Your entry has been created."})
}

func EntryUpdated(w http.ResponseWriter) {
	Set(w, Message{Level: Sucess, Content: "Success! Your entry has been updated."})
}
