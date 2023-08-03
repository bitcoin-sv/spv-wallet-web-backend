# BUX Wallet Backend

The `bux-wallet-backend` is an HTTP server that serves as the referential backend server for a custodial web wallet for Bitcoin SV (BSV). It integrates and utilizes the `bux-server` component to handle various BSV-related operations, including the creation of transactions and listing incoming and outgoing transactions.

## Running as a Docker Image

To run the `bux-wallet-backend` as a Docker image with custom configuration, you can set the required environment variables using the `-e` flag. The environment variables will be used by the application to configure various aspects of the backend.

Here's the updated command with environment variables for the Docker container:

```bash
docker run -p 8080:8080 -e BUX_PAYMAIL_DOMAIN=bux.com ${DOCKERHUB_OWNER}/${DOCKERHUB_REPO}:latest
```

### Configuration

The `bux-wallet-backend` can be configured using environment variables. Here is a list of available environment variables and their purpose:

| Environment Variable                  | Description                                                                       | Default Value                                           |
|--------------------------------------|-----------------------------------------------------------------------------------|---------------------------------------------------------|
| `DB_HOST`                            | Database host address.                                                            | `localhost`                                             |
| `DB_PORT`                            | Database port number.                                                             | `5432`                                                  |
| `DB_USER`                            | Database username.                                                                | `postgres`                                              |
| `DB_PASSWORD`                        | Database password.                                                                | `postgres`                                              |
| `DB_NAME`                            | Database name.                                                                    | `postgres`                                              |
| `DB_SSL_MODE`                        | Database SSL mode.                                                                | `disable`                                               |
| `HTTP_SERVER_READ_TIMEOUT`           | HTTP server read timeout (in seconds).                                            | `40`                                                    |
| `HTTP_SERVER_WRITE_TIMEOUT`          | HTTP server write timeout (in seconds).                                           | `40`                                                    |
| `HTTP_SERVER_PORT`                   | HTTP server port.                                                                 | `8080`                                                  |
| `HTTP_SERVER_COOKIE_DOMAIN`          | HTTP server cookie domain parameter.                                              | `localhost`                                             |
| `HTTP_SERVER_COOKIE_SECURE`          | HTTP server cookie secure parameter.                                              | `false`                                                 |
| `HTTP_SERVER_CORS_ALLOWED_DOMAINS`   | HTTP server CORS origin allowed domains.                                          | `[]`                                                    |
| `BUX_ADMIN_XPRIV`                    | BUX admin xpriv.                                                                  | `xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK` |
| `BUX_SERVER_URL`                     | BUX server URL.                                                                   | `http://localhost:3003/v1`                              |
| `BUX_WITH_DEBUG`                     | Enable debugging for BUX connection.                                              | `true`                                                  |
| `BUX_SIGN_REQUEST`                   | Enable signing of all requests for BUX connection.                                | `true`                                                  |
| `BUX_PAYMAIL_DOMAIN`                 | BUX paymail domain.                                                               | `bux.com`                                               |
| `BUX_PAYMAIL_AVATAR`                 | BUX paymail avatar URL.                                                           | `http://localhost:3003/static/paymail/avatar.jpg`        |
| `HASH_SALT`                          | Hash salt for the application.                                                    | `bux`                                                   |
| `LOGGING_LEVEL`                      | Logging level for the running application.                                        | `Debug`                                                 |
| `ENDPOINTS_EXCHANGE_RATE`            | Exchange rate endpoint URL used in the app.                                       | `https://api.whatsonchain.com/v1/bsv/main/exchangerate` |

