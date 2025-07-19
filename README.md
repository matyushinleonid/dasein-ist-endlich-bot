# dasein-ist-endlich-bot

The topic of death remains largely taboo in today’s society (as noted, for example, by Philippe Ariès). This project is designed to help users face the inevitability of death as a phenomenon that will (and indeed **is**) happen to them.

Users are invited to answer a series of questions about their lifestyle, after which a rough estimate of their remaining lifespan is calculated by an LLM-based model. 
Of course, the resulting number of days left is not meant to be precisely accurate. 
The real value lies in the daily countdown of remaining days that the bot will send to the user.

* The current Telegram implementation is available as [Sein zum Tode](https://t.me/DaseinIstEndlichBot)
* The corresponding ArgoCD application and related Kubernetes resources can be found [here](https://github.com/matyushinleonid/k8s.leonid.sh/tree/main/argocd/dasein-ist-endlich-bot)

## Description

At its core, this project is a wrapper around the OpenAI API, combined with functionality to send regular notifications.

The bot is implemented as a Go application consisting of two components:

* `listener` — the Telegram bot itself, which connects to the Telegram Bot API and handles user interaction.
* `notifier` — a separate service responsible for sending messages to users. In production, it runs as a CronJob.

For persistent storage of user data (e.g., the user’s estimated date of death, timestamp of the last notification sent, etc.), MongoDB is used. 
For temporary storage of user responses (answers to the questionnaire), Redis is employed. These messages are not stored persistently to ensure user privacy.

## Test/Run locally

Unit tests can be executed with:

```bash
go test -v ./...
```

For local development, we recommend using Docker Compose (see `docker-compose.yaml`), which will conveniently launch MongoDB and Redis.

Set the following environment variables in a `.env` file at the project root.
* `TELEGRAM_BOT_TOKEN`
* `OPENAI_API_KEY`

If you prefer to use a remote MongoDB instance, also specify `MONGO_PASSWORD`.

Then run:
```bash
docker compose up -d --build
```
This will start the bot using the configs/local.yaml configuration file.
