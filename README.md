# Paola GO Bot

## Available commands

Quick note here: remember to add a trailing `-a paola-go-bot` to each command, if you run it from outside this project folder.

- Check the status of process:

```bash
heroku ps
```

- Open Postgres CLI:

```bash
heroku pg:psql
```

- Check the logs:

```bash
heroku logs --tail
```

- Set the worker to run on 1 dyno (always available):

```bash
heroku ps:scale listen=1
```

- Run the birthday reminder for today:

```bash
heroku run reminder
```

- Run the birthday reminder for today:

```bash
heroku run reminder 3
```

- Run listener locally:

```bash
go build -o bin/paola-go-bot && heroku local listen=1
```
