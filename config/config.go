package config

// Define basic db config keys.
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

// Define basic http server config keys.
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
	// EnvHttpServerSessionSecret gin session store secret to encrypt session data in database.
	EnvHttpServerSessionSecret = "http.server.session.secret" //nolint: gosec
)

// Define basic spv-wallet config keys.
const (
	// EnvAdminXpriv define the admin xpriv.
	EnvAdminXpriv = "spvwallet.admin.xpriv"
	// EnvServerUrl define the url of the spv-wallet (non-custodial wallet) service.
	EnvServerUrl = "spvwallet.server.url"
	// EnvPaymailDomain define the paymail domain.
	EnvPaymailDomain = "spvwallet.paymail.domain"
	// EnvPaymailAvatar define the paymail avatar url.
	EnvPaymailAvatar = "spvwallet.paymail.avatar"
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
	// EnvLoggingInstanceName define the instance name for logging.
	EnvLoggingInstanceName = "logging.instance.name"
	// EnvLoggingFormat define the logging format - console/json.
	EnvLoggingFormat = "logging.format"
	// EnvLoggingLogOrigin define whether to log origin.
	EnvLoggingLogOrigin = "logging.log.origin"
)

const (
	// EnvEndpointsExchangeRate define the exchange rate endpoint.
	EnvEndpointsExchangeRate = "endpoints.exchangeRate"
)

const (
	EnvContactsPasscodePeriod = "contacts.passcode.period"
	EnvContactsPasscodeDigits = "contacts.passcode.digits"
)

const (
	// EnvCacheSettingsTtl define the cache settings ttl used for exchange rates storage.
	EnvCacheSettingsTtl = "cache.settings.ttl"
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
