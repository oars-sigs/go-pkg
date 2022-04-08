package constant

const (
	UserIDSessionKey      = "user_id"
	UserDetailSessionKey  = "user_detail"
	AccessTokenSessionKey = "access_token"
	SessionHeader         = "Authorization"
	ConnectorIDHeader     = "Oars-IdaaS-Connector"
	ProxyUserIDHeader     = "X-Oars-User-Id"
	ProxyUserAuthHeader   = "X-Oars-Authorization"
	ProxyUserTokenHeader  = "X-Oars-Token"
	ProxyAppIDHeader      = "X-Oars-App-Id"
	ProxySidHeader        = "X-Oars-Sid"
	ProxyAuthKindHeader   = "X-Oars-Auth-Kind"
	SessionCookieName     = "OARS-IDASS-SESSION"
	ReqURIHeader          = "XURI"
	ReqMethodHeader       = "XMETHOD"
	ReqHostHeader         = "XHOST"
	ReqAuthentication     = "XAUTHENTICATION"
	SignatureKey          = "signature"
	OarsAuthKind          = "X-Oars-Auth-Kind"
	OarsHmacSignatureKind = "Hmac"
	ExpireTimeKey         = "expireTime"
)

const (
	AppAuthKindHeader  = "app"
	UserAuthKindHeader = "user"
)
