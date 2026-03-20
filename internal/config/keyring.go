package config

import "github.com/zalando/go-keyring"

const (
	keyringService = "jiru"
	keyringUser    = "api-token"
)

// keyringUserForProfile returns the profile-aware keyring user key.
func keyringUserForProfile(profile string) string {
	if profile == "" || profile == "default" {
		return keyringUser // backward compat
	}
	return keyringUser + "-" + profile
}

// deleteKeyringToken removes the API token from the OS keychain (legacy default key).
func deleteKeyringToken() {
	_ = keyring.Delete(keyringService, keyringUser)
}

// setKeyringTokenForProfile stores a token for a named profile.
func setKeyringTokenForProfile(profile, token string) error {
	return keyring.Set(keyringService, keyringUserForProfile(profile), token)
}

// getKeyringTokenForProfile retrieves a token for a named profile.
func getKeyringTokenForProfile(profile string) (string, error) {
	return keyring.Get(keyringService, keyringUserForProfile(profile))
}

// deleteKeyringTokenForProfile removes a token for a named profile.
func deleteKeyringTokenForProfile(profile string) {
	_ = keyring.Delete(keyringService, keyringUserForProfile(profile))
}
