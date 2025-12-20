#Задумка какая - так называемый BFF (Backend For Frontend), который должен:
- сборка нескольких сервисов в один payload (значит готовность к взаимодействию через GRPC/брокеров)
- auth / rate limiting / caching

#Что надо сделать:
- [ ]  Fiber v3 ядро
- [ ]  аутентификация - почта, телефон (смс-коды, телеграм), oauth (gmail, Facebook, vk, Яндекс)
- [ ]  асимметричное шифрование токенов, то есть разделить секреты для создания токенов и для проверки в других сервисах
- [ ]  сессии, Http only куки выдавать на клиент. Возможно контролировать сессии только через refresh-токены
- [ ]  Refresh-токены хранить в БД для инвалидизации сессий (есть в закладках хрома пример)
- [ ]  Несколько crud-таблиц. В том числе реализовать: управление записями таблиц только своим пользователем, soft-delete
- [ ]  RBAC (role-based access control), есть в репоизитории Gorsk, реализовать 2-3 роли для данных
- [ ]  Di через интерфейсы
- [ ]  Redis для сессий, кеша отдельных простых запросов
- [ ]  Postgres (настройка work mem, shared buffers), pgx, scany, squirrel или huandu/go-sqlbuilder
- [ ]  Nats Jetstream для отправки сообщений
- [ ]  Swagger (https://github.com/gofiber/swagger)
- [ ]  Prometheus клиент
- [ ]  Docker compose как стандартный режим
- [ ]  Ci/CD - action в гитхаб, сборка контейнеров для server и migrator, тесты, lint, пуш в докер-хаб
- [ ]  Потом приделать grpc

#Возможные источники для примера:
- https://github.com/indravscode/go-fiber-boilerplate
- https://github.com/ribice/gorsk
- https://habr.com/ru/companies/ozontech/articles/976950/
