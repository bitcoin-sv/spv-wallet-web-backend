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
	// EnvHTTPServerReadTimeout http server read timeout.
	EnvHTTPServerReadTimeout = "http.server.readTimeout"
	// EnvHTTPServerWriteTimeout http server write timeout.
	EnvHTTPServerWriteTimeout = "http.server.writeTimeout"
	// EnvHTTPServerPort http server port.
	EnvHTTPServerPort = "http.server.port"
	// EnvHTTPServerCookieDomain http server cookie domain parameter.
	EnvHTTPServerCookieDomain = "http.server.cookie.domain"
	// EnvHTTPServerCookieSecure http server cookie secure parameter.
	EnvHTTPServerCookieSecure = "http.server.cookie.secure"
	// EnvHTTPServerCorsAllowedDomains http server cors origin allowed domains.
	EnvHTTPServerCorsAllowedDomains = "http.server.cors.allowedDomains"
	// EnvHTTPServerSessionSecret gin session store secret to encrypt session data in database.
	EnvHTTPServerSessionSecret = "http.server.session.secret" //nolint: gosec
)

// Define basic spv-wallet config keys.
const (
	// EnvAdminXpriv define the admin xpriv.
	EnvAdminXpriv = "spvwallet.admin.xpriv"
	// EnvServerURL define the url of the spv-wallet (non-custodial wallet) service.
	EnvServerURL = "spvwallet.server.url"
	// EnvPaymailDomain define the paymail domain.
	EnvPaymailDomain = "spvwallet.paymail.domain"
	// EnvPaymailAvatar define the paymail avatar url.
	EnvPaymailAvatar = "spvwallet.paymail.avatar"
)

const (
	// EnvWebsocketHistoryMax max number of published events that should be hold
	// and send to client in case of restored lost connection.
	EnvWebsocketHistoryMax = "websocket.history.max"
	// EnvWebsocketHistoryTTL max minutes for which published events should be hold
	// and send to client in case of restored lost connection.
	EnvWebsocketHistoryTTL = "websocket.history.ttl"
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
	// EnvContactsPasscodePeriod define the contacts passcode validity period.
	EnvContactsPasscodePeriod = "contacts.passcode.period"
	// EnvContactsPasscodeDigits define the contacts passcode digits number.
	EnvContactsPasscodeDigits = "contacts.passcode.digits"
)

const (
	// EnvCacheSettingsTTL define the cache settings ttl used for exchange rates storage.
	EnvCacheSettingsTTL = "cache.settings.ttl"
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
