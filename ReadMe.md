# Что это

# Как настроить

# Как запустить

## Тестовый запуск

Приложение запускается командой

```sh
cd bnhmmr
sudo go run main.go
```

Ожидаемый результат

```
2025/01/20 10:49:33 Tool to ban XRay users running.
2025/01/20 10:49:33 Update Ban list interval: 5s
2025/01/20 10:49:33 Ban time: 1m0s
2025/01/20 10:49:33 Banning xray outbound: blacklist
```
Для адресов, попадающих в blacklist outbound, будут появляться записи о попадании IP-адресов в бан.

## Запуск сервиса

Создадим сервис, который будет запускаться при старте системы:

### Сборка приложения

Соберем приложение и поместим его и его конфиг в `/opt/goipban/`

```sh
sudo go build -o /opt/goipban/goipban main.go
sudo cp config/default_config.json /opt/goipban/config.json
```

### Создание systemd-сервиса

Создадим конфиг systemd-сервиса (`sudo nano /etc/systemd/system/goipban.service`), который будет запускать наше приложение при старте системы и перезапускать при падениях, и вставляем в него следующий текст:

```ini
[Unit]
Description=GoIPBan
[Service]
Type=simple
Restart=on-failure
RestartSec=30
WorkingDirectory=/opt/goipban
ExecStart=/opt/goipban/goipban run -c /opt/goipban/config.json
[Install]
WantedBy=multi-user.target
```

Активируем:

```sh
sudo systemctl daemon-reload
sudo systemctl enable goipban
```

Запускаем, смотрим статус и логи:
```sh
sudo systemctl start goipban
sudo systemctl status goipban
sudo journalctl -u goipban
```




