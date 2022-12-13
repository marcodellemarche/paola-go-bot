#!/bin/bash
# Needs sudo privileges

systemctl enable docker

cp systemd/paola-go-bot-reminder.* /etc/systemd/system/

systemctl start paola-go-bot-reminder.service
systemctl enable paola-go-bot-reminder.timer

docker compose up -d

