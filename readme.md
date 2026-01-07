# Задумка какая - так называемый BFF (Backend For Frontend), который должен:

- сборка нескольких сервисов в один payload (значит готовность к взаимодействию через GRPC/брокеров)
- auth / rate limiting / caching

# Что надо сделать:

- [v] Fiber v3 ядро
- [v] аутентификация - почта/телефон, 
        - ручка POST auth/login, требуется почта/телефон, пароль
- [ ] смена пароля
        - ручка POST auth/change-password, требуется access-токен (1), старый пароль (2), новый пароль (3)
- [ ] сброс пароля
        - ручка POST auth/verify-code, отправка 6-значного кода верификации, требуется почта/телефон, причина. Проверять частоту отправки!
        - ручка POST auth/reset-password, проверка 6-значного кода верификации, требуется почта/телефон, новый пароль, код
- [ ] oauth (gmail, Facebook, vk, Яндекс)
- [ ] верификация через email и телефон при регистрации - В РАБОТЕ
        - ручка POST auth/verify-code, отправка 6-значного кода верификации, требуется почта/телефон, причина. Проверять частоту отправки!
        - для почты использовать пакет github.com/jordan-wright/email
        - ручка POST auth/register, проверка 6-значного кода верификации и регистрация (добавить в DTO)
        - коды верификации хранить в Redis, TTL 5 минут, 
        - дату верификации хранить в БД Users
- [ ] асимметричное шифрование токенов, то есть разделить секреты для создания токенов и для проверки в других сервисах
- [v] login / refresh с выдачей access и refresh токенов
- [v] контроль сессий через refresh-токены, Refresh-токены хранить в Redis для инвалидизации сессий
- [ ] Несколько crud-таблиц. В том числе реализовать: управление записями таблиц только своим пользователем, soft-delete
- [ ] создать таблицу для денег, для отработки конвертаций кастомного decimal в БД и обратно. Использовать внешний пакет (https://github.com/shopspring/decimal)
- [ ] RBAC (role-based access control), есть в репоизитории Gorsk, реализовать 2-3 роли для данных
- [ ] Di через интерфейсы
- [v] Redis для сессий
- [ ] Redis для кэширования отдельных простых запросов
- [ ] Postgres (настройка work mem, shared buffers), pgx, scany, squirrel или huandu/go-sqlbuilder
- [ ] Nats Jetstream для отправки сообщений
- [v] Swagger (https://github.com/gofiber/swagger)
- [v] Prometheus клиент
- [v] Docker compose как стандартный режим
- [ ] Dockerfile для server, seeder, migrator
- [ ] Ci/CD - action в гитхаб, сборка контейнеров для server и migrator, тесты, lint, пуш в докер-хаб
- [ ] Потом приделать grpc

# Возможные источники для примера:

- https://github.com/indravscode/go-fiber-boilerplate
- https://github.com/ribice/gorsk
- https://habr.com/ru/companies/ozontech/articles/976950/

# Контейнеры:

- server, seeder, migrator
- orders - вынесенный GRPC-сервис для примера
- postgres
- NATS
- redis
