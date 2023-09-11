package config

// Define basic db config.
const (
	// EnvDbHost define the database host.
	EnvDbHost = "db.host"
	// EnvDbPort define the database port.
	EnvDbPort = "db.port"
	// EnvDbUser define the database user.
	EnvDbUser = "db.user"
	// EnvDbPassword define the database password.
	EnvDbPassword = "db.password"
	// EnvDbName define the database name.
	EnvDbName = "db.name"
	// EnvDbSslMode define the database ssl mode.
	EnvDbSslMode = "db.sslMode"
	// EnvDbMigrationsPath define the database migrations path.
	EnvDbMigrationsPath = "db.migrationsPath"
)

// Define basic http server config.
const (
	// EnvHttpServerReadTimeout http server read timeout.
	EnvHttpServerReadTimeout = "http.server.readTimeout"
	// EnvHttpServerWriteTimeout http server write timeout.
	EnvHttpServerWriteTimeout = "http.server.writeTimeout"
	// EnvHttpServerPort http server port.
	EnvHttpServerPort = "http.server.port"
	// EnvHttpServerCookieDomain http server cookie domain parameter.
	EnvHttpServerCookieDomain = "http.server.cookie.domain"
	// EnvHttpServerCookieSecure http server cookie secure parameter.
	EnvHttpServerCookieSecure = "http.server.cookie.secure"
	// EnvHttpServerCorsAllowedDomains http server cors origin allowed domains.
	EnvHttpServerCorsAllowedDomains = "http.server.cors.allowedDomains"
)

// Define basic bux config.
const (
	// EnvBuxAdminXpriv define the bux admin xpriv.
	EnvBuxAdminXpriv = "bux.admin.xpriv"
	// EnvBuxServerUrl define the bux server url.
	EnvBuxServerUrl = "bux.server.url"
	// EnvBuxWithDebug define whether to turn debugging on.
	EnvBuxWithDebug = "bux.withDebug"
	// EnvBuxSignRequest define whether to sign all requests.
	EnvBuxSignRequest = "bux.sign.request"
	// EnvBuxPaymailDomain define the bux paymail domain.
	EnvBuxPaymailDomain = "bux.paymail.domain"
	// EnvBuxPaymailAvatar define the bux paymail avatar url.
	EnvBuxPaymailAvatar = "bux.paymail.avatar"
)

const (
	// EnvWebsocketHistoryMax max number of published events that should be hold
	// and send to client in case of restored lost connection.
	EnvWebsocketHistoryMax = "websocket.history.max"
	// EnvWebsocketHistoryTtl max minutes for which published events should be hold
	// and send to client in case of restored lost connection.
	EnvWebsocketHistoryTtl = "websocket.history.ttl"
)

// EnvHashSalt define the hash salt.
const EnvHashSalt = "hash.salt"

const (
	// EnvLoggingLevel define logging level for running application.
	EnvLoggingLevel = "logging.level"
)

const (
	// EnvEndpointsExchangeRate define the exchange rate endpoint.
	EnvEndpointsExchangeRate = "endpoints.exchangeRate"
)

// Config returns strongly typed config values.
type Config struct {
	Db *Db
}

// Db represents a database connection.
type Db struct {
	Host           string
	Port           int
	User           string
	Password       string
	Name           string
	SslMode        string
	MigrationsPath string
}
