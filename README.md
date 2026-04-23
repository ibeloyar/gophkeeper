# Менеджер паролей GophKeeper

## Server
Запуск сервера
```terminaloutput
gophkeeper -c config/server/config.yaml
```
Пример config.yaml
```yaml
grpc_server:
  addr: ":50022"

http_server:
  addr: ":8080"

database:
  dsn: "postgres://user:password@192.168.0.1051:5432/dbname?sslmode=disable"

security:
  user_password_cost: 3
  token_secret: "secret"
  token_lifetime: 1h
  secret_password_key: "12345678901234567890123456789012" # 32 bytes AES‑256

swagger:
  enabled: true
  json_path: "./api/apidocs.swagger.json"
```

> [!IMPORTANT]
> При передаче в конфигурационном файле http_server.addr будет запущен http-сервер.
> А так же при передаче swagger.enabled -> true, по адресу /swagger будет доступен
> Swagger UI с интерактивной документацией API

![swagger.png](docs/swagger.png)

## CLI
Собраные под разные платформы CLI-приложения, находятся в ./cmd/client
- gophkeeper-cli-darwin (MacOS)
- gophkeeper-cli-linux (Linux)
- gophkeeper-cli-win.exe (Windows)

**Запуск**
```terminaloutput
gophkeeper-cli-linux -a grpc://server_addr:50022
```
или с файлом конфигурации
```terminaloutput
gophkeeper-cli-linux -с path/to/config.yaml
```

Пример config.yaml
```yaml
grpc_server:
  addr: ":50022"
```

После запуска вы должны зарегестрироваться 

![register.png](docs/register.png)

или войти в учётную запись

![login.png](docs/login.png)

и далее взаимодействовать с приложением с помощью подсказок

![list.png](docs/list.png)

![card.png](docs/card.png)

## Общие требования
GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины,
пароли, бинарные данные и прочую приватную информацию.

### Сервер должен реализовывать следующую бизнес-логику:
- [x] регистрация, аутентификация и авторизация пользователей;
- [x] хранение приватных данных (логин/пароль, текстовые данные, бинарные данные, данные банковских карт);
- [x] синхронизация данных между несколькими авторизованными клиентами одного владельца;
- [x] передача приватных данных владельцу по запросу;
- [x] для любых данных должна быть возможность хранения произвольной текстовой метаинформации.

### Клиент должен реализовывать следующую бизнес-логику:
- [x] аутентификация и авторизация пользователей на удалённом сервере;
- [x] доступ к приватным данным по запросу.

### Дополнительные требования:
- [x] клиент должен распространяться в виде CLI-приложения с возможностью запуска на платформах Windows, Linux и Mac OS;
- [x] клиент должен давать пользователю возможность получить информацию о версии и дате сборки бинарного файла клиента.

### Тестирование и документация
- [x] Код всей системы должен быть покрыт юнит-тестами не менее чем на 70%.
- [x] Каждая экспортированная функция, тип, переменная, а также пакет системы должны содержать исчерпывающую документацию.


## Разработка и сборка
Все необходимые команды для разработки и сборке указаны в make файле.
```terminaloutput
> make help

🦫  gophkeeper

command                | description
====================================================
install-pg-tools       | install golang-migrate CLI
install-mock-tools     | install mockgen/gomock
install-proto-tools    | download protoc + plugins
install-all-tools      | install proto/mock/pg tools
proto                  | generate proto files
migrate-up             | apply DB migrations
migrate-down           | rollback migration
migrate-create         | create migration (NAME=...) 
mock                   | generate Storage/Service mocks
test                   | run tests
test_cover             | tests + coverage report
gofmt                  | format all Go files
build-cli              | cross-platform CLI builds
build                  | build server
```

### Последнее тестовое покрытие
```
go tool cover -func=coverage.filtered.out
github.com/ibeloyar/gophkeeper/cmd/server/gophkeeper.go:10:                     main                            0.0%
github.com/ibeloyar/gophkeeper/internal/app/app.go:27:                          Run                             0.0%
github.com/ibeloyar/gophkeeper/internal/app/build_grpc_server.go:17:            buildGRPCServer                 100.0%
github.com/ibeloyar/gophkeeper/internal/app/build_http_server.go:17:            buildHTTPServer                 70.0%
github.com/ibeloyar/gophkeeper/internal/client/config/config.go:23:             Read                            100.0%
github.com/ibeloyar/gophkeeper/internal/client/create_binary_secret.go:30:      updateCreateBinarySecret        47.8%
github.com/ibeloyar/gophkeeper/internal/client/create_binary_secret.go:145:     createBinarySecretView          88.9%
github.com/ibeloyar/gophkeeper/internal/client/create_card_secret.go:31:        updateCreateCardSecret          36.6%
github.com/ibeloyar/gophkeeper/internal/client/create_card_secret.go:149:       createCardSecretView            90.9%
github.com/ibeloyar/gophkeeper/internal/client/create_log_pass_secret.go:29:    updateCreateLogPassSecret       38.3%
github.com/ibeloyar/gophkeeper/internal/client/create_log_pass_secret.go:133:   createLogPassSecretView         90.0%
github.com/ibeloyar/gophkeeper/internal/client/create_text_secret.go:29:        updateCreateTextSecret          95.9%
github.com/ibeloyar/gophkeeper/internal/client/create_text_secret.go:118:       createTextSecretView            88.9%
github.com/ibeloyar/gophkeeper/internal/client/grpcclient/grpcclient.go:17:     New                             0.0%
github.com/ibeloyar/gophkeeper/internal/client/grpcclient/grpcclient.go:29:     Close                           0.0%
github.com/ibeloyar/gophkeeper/internal/client/init.go:10:                      initialModel                    100.0%
github.com/ibeloyar/gophkeeper/internal/client/login.go:26:                     updateLogin                     66.7%
github.com/ibeloyar/gophkeeper/internal/client/login.go:96:                     loginView                       100.0%
github.com/ibeloyar/gophkeeper/internal/client/main.go:58:                      Init                            0.0%
github.com/ibeloyar/gophkeeper/internal/client/main.go:66:                      Update                          0.0%
github.com/ibeloyar/gophkeeper/internal/client/main.go:94:                      View                            0.0%
github.com/ibeloyar/gophkeeper/internal/client/main.go:119:                     main                            0.0%
github.com/ibeloyar/gophkeeper/internal/client/register.go:26:                  updateRegister                  66.7%
github.com/ibeloyar/gophkeeper/internal/client/register.go:95:                  registerView                    100.0%
github.com/ibeloyar/gophkeeper/internal/client/secret_card.go:13:               updateSecretCard                100.0%
github.com/ibeloyar/gophkeeper/internal/client/secret_card.go:33:               secretCardView                  86.7%
github.com/ibeloyar/gophkeeper/internal/client/secret_delete.go:15:             updateSecretDelete              100.0%
github.com/ibeloyar/gophkeeper/internal/client/secret_delete.go:50:             secretDeleteView                100.0%
github.com/ibeloyar/gophkeeper/internal/client/secrets.go:22:                   pollSecrets                     50.0%
github.com/ibeloyar/gophkeeper/internal/client/secrets.go:31:                   updateSecrets                   90.2%
github.com/ibeloyar/gophkeeper/internal/client/secrets.go:107:                  secretsView                     100.0%
github.com/ibeloyar/gophkeeper/internal/config/config.go:43:                    Read                            91.7%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/convert.go:13:          convertSecretTypeToProto        100.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/convert.go:28:          convertSecretTypeToDTO          100.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/convert.go:46:          convertCreateSecretToDTO        80.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/convert.go:77:          convertGetSecretToProto         100.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/convert.go:103:         convertGetSecretsToProto        100.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:38:             New                             100.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:47:             HandlePanic                     0.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:57:             Register                        0.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:87:             Login                           0.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:114:            CreateSecret                    0.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:147:            GetSecret                       71.4%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:165:            GetSecrets                      0.0%
github.com/ibeloyar/gophkeeper/internal/controller/grpc/grpc.go:183:            DeleteSecret                    0.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/helpers.go:14:          readBody                        100.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/helpers.go:52:          writeJSON                       100.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:29:             New                             100.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:38:             Register                        100.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:68:             Login                           80.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:94:             GetSecret                       100.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:119:            GetSecrets                      100.0%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:134:            CreateSecret                    86.7%
github.com/ibeloyar/gophkeeper/internal/controller/http/http.go:172:            DeleteSecret                    100.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/pgstorage.go:34:   New                             0.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/pgstorage.go:71:   Shutdown                        0.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/secrets.go:14:     GetSecret                       100.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/secrets.go:35:     GetSecrets                      92.9%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/secrets.go:77:     CreateSecret                    100.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/secrets.go:99:     DeleteSecret                    90.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/users.go:14:       CreateUser                      100.0%
github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/users.go:31:       GetUserByLogin                  100.0%
github.com/ibeloyar/gophkeeper/internal/service/service.go:36:                  New                             100.0%
github.com/ibeloyar/gophkeeper/internal/service/service.go:50:                  Register                        92.9%
github.com/ibeloyar/gophkeeper/internal/service/service.go:84:                  Login                           90.9%
github.com/ibeloyar/gophkeeper/internal/service/service.go:111:                 CreateSecret                    90.0%
github.com/ibeloyar/gophkeeper/internal/service/service.go:133:                 GetSecret                       84.6%
github.com/ibeloyar/gophkeeper/internal/service/service.go:159:                 GetSecrets                      90.0%
github.com/ibeloyar/gophkeeper/internal/service/service.go:180:                 DeleteSecret                    100.0%
github.com/ibeloyar/gophkeeper/internal/service/validate.go:19:                 validateLoginDTO                100.0%
github.com/ibeloyar/gophkeeper/internal/service/validate.go:31:                 validateRegisterDTO             100.0%
github.com/ibeloyar/gophkeeper/internal/service/validate.go:43:                 validateLogin                   100.0%
github.com/ibeloyar/gophkeeper/internal/service/validate.go:51:                 validatePassword                100.0%
github.com/ibeloyar/gophkeeper/internal/service/validate.go:59:                 validateCreateSecretDTO         100.0%
github.com/ibeloyar/gophkeeper/pkg/auth/auth.go:33:                             GenerateBearerToken             80.0%
github.com/ibeloyar/gophkeeper/pkg/auth/auth.go:53:                             VerifyJWTBearerToken            80.0%
github.com/ibeloyar/gophkeeper/pkg/auth/auth.go:83:                             AuthBearerMiddlewareInit        100.0%
github.com/ibeloyar/gophkeeper/pkg/auth/auth.go:100:                            GetTokenInfo                    100.0%
github.com/ibeloyar/gophkeeper/pkg/auth/auth.go:113:                            GetTokenInfoFromContext         100.0%
github.com/ibeloyar/gophkeeper/pkg/auth/auth.go:125:                            AuthGRPCUnaryInterceptor        100.0%
github.com/ibeloyar/gophkeeper/pkg/logger/logger.go:13:                         New                             85.7%
github.com/ibeloyar/gophkeeper/pkg/logger/logger.go:28:                         LoggingMiddleware               100.0%
github.com/ibeloyar/gophkeeper/pkg/logger/logger.go:68:                         Write                           100.0%
github.com/ibeloyar/gophkeeper/pkg/logger/logger.go:76:                         WriteHeader                     100.0%
github.com/ibeloyar/gophkeeper/pkg/password/password.go:24:                     HashPassword                    100.0%
github.com/ibeloyar/gophkeeper/pkg/password/password.go:42:                     CheckPasswordHash               100.0%
github.com/ibeloyar/gophkeeper/pkg/password/password.go:50:                     EncryptPassword                 84.2%
github.com/ibeloyar/gophkeeper/pkg/password/password.go:87:                     DecryptPassword                 90.0%
total:                                                                          (statements)                    70.0%
```
