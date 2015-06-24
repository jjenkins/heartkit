package models

import "time"

type User struct {
	Id           string    `json:"id"`
	EncodedId    string    `json:"encodedId"`
	Gender       string    `json:"gender"`
	DateOfBirth  string    `json:"dateOfBirth"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
