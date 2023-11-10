package static

import "time"

var (
	AuthSessionExpiration  = 24 * 7 * time.Hour
	RefreshTokenCookieName = "snip-refresh-token"
)
