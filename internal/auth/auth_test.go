package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	id := uuid.New()
	secret := "secret"
	expiresIn := time.Duration(24 * time.Hour)

	tests := []struct {
		name          string
		id            uuid.UUID
		secret        string
		expiresIn     time.Duration
		wantErr       bool
	}{
		{
			name:          "Valid parameters",
			id:            id,
			secret:        secret,
			expiresIn:     expiresIn,
			wantErr:       false,
		},
		{
			name:          "Invalid id",
			id:            uuid.Nil,
			secret:        secret,
			expiresIn:     expiresIn,
			wantErr:       true,
		},
		{
			name:          "Invalid secret",
			id:            id,
			secret:        "",
			expiresIn:     expiresIn,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.id, tt.secret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Errorf("MakeJWT() token = %v, want non-empty string", token)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name          string
		id            uuid.UUID
		secret        string
		expiresIn     time.Duration
		wantErr       bool
		wantIdMatch   bool
	}{
		{
			name:          "Valid token",
			id:            id,
			secret:        "secret",
			expiresIn:     time.Duration(24 * time.Hour),
			wantErr:       false,
			wantIdMatch:   true,
		},
		{
			name:          "Invalid secret",
			id:            id,
			secret:        "invalidsecret",
			expiresIn:     time.Duration(24 * time.Hour),
			wantErr:       true,
			wantIdMatch:   false,
		},
		{
			name:          "Invalid expiresIn",
			id:            id,
			secret:        "secret",
			expiresIn:     time.Duration(-1 * time.Hour),
			wantErr:       true,
			wantIdMatch:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.id, "secret", tt.expiresIn)

			if err != nil {
				t.Errorf("MakeJWT() error = %v", err)
			}
			if token == "" {
				t.Errorf("MakeJWT() token = %v, want non-empty string", token)
			}	

			validatedID, err := ValidateJWT(token, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v", err)
			}

			if !tt.wantErr && !tt.wantIdMatch && validatedID != tt.id {
				t.Errorf("ValidateJWT() validatedID = %v, want %v", validatedID, tt.id)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	

	tests := []struct {
		name        	string
		headers		    func () http.Header
		expectedResult	string
		wantErr 		bool
	}{
		{
			name: "Valid token",
			headers: func() http.Header {
				headers := make(http.Header)
				headers.Set("Authorization", "Bearer mytoken")
				return headers
			},
			expectedResult: "mytoken",
			wantErr: false,
		},
		{
			name: "No header",
			headers: func() http.Header {
				headers := make(http.Header)
				return headers
			},
			expectedResult: "",
			wantErr: true,
		},
		{
			name: "Invalid format",
			headers: func() http.Header {
				headers := make(http.Header)
				headers.Set("Authorization", "Bearermytoken")
				return headers
			},
			expectedResult: "",
			wantErr: true,
		},
		{
			name: "Invalid format no token",
			headers: func() http.Header {
				headers := make(http.Header)
				headers.Set("Authorization", "Bearer ")
				return headers
			},
			expectedResult: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() err = %v", err)
			}

			if !tt.wantErr && token != tt.expectedResult {
				t.Errorf("GetBearerToken() err = %v", err)
			}
		})
	}
}