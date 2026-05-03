<p align="center">
  <img src="https://github.com/user-attachments/assets/fb9a3ed6-8435-4317-b1d0-d67477bbaadb" width="150" height="150" alt="icon">
  <h2 align="center">Housekeepy</h2>

</p>

<p align="center">
  Never neglect your household responsibilities again!
</p>

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Release](https://img.shields.io/github/v/release/ashiven/housekeepy)](https://github.com/ashiven/housekeepy/releases)
[![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/ashiven/housekeepy)](https://github.com/ashiven/housekeepy/issues)
[![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/ashiven/housekeepy)](https://github.com/ashiven/housekeepy/pulls)

</div>

**Housekeepy** is a simple application that organizes your household and prevents fights about who should do what and when and how often.
Given a list of members and regular household tasks, it automatically assigns these tasks each day and sends all members a notification on their chosen messaging service.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features

- :clipboard: Automatic household task assignments
- :alarm_clock: Daily reminders via SMS or WhatsApp
- :wrench: Browser interface for configuration

## Getting Started

### Prerequisites

1. Register for an account on [Twilio](https://www.twilio.com/en-us) and create a [Messaging Service](https://console.twilio.com/us1/develop/sms/services).
2. Have [Docker](https://docs.docker.com/get-started/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
3. If you want to host this service, prepare your domain name and install [Certbot](https://www.digitalocean.com/community/tutorials/how-to-use-certbot-standalone-mode-to-retrieve-let-s-encrypt-ssl-certificates-on-ubuntu-20-04).

### Setup

1. Clone the repository to your machine.

```bash
git clone https://github.com/ashiven/housekeepy.git
```

2. Fill out the `.env` file with the required variables.

```
TWILIO_ACCOUNT_SID=AC55...
TWILIO_AUTH_TOKEN=25fe...
WHATSAPP_SENDER=+155...
SMS_SENDER=+174...
TEMPLATE_SID=HX8b...
SERVICE_SID=MGf7...
ADMIN_PASSWORD=myAdminPassword
COOKIE_SECRET=jhkwe234kjhajhfa234
DOMAIN_NAME=myDomainName.com
```

3. Add tasks and household members to `config.ini`.

```
[Members]
Thomas : +4924352425243
Peter : +4923452354235

[Daily Tasks]
Throw out the trash :
Wipe the floor :

[Weekly Tasks]
Buy groceries :
Clean the living room :

[Monthly Tasks]
Clean the fridge :
```

4. Create an SSL certificate.

```bash
certbot certonly --standalone -d <domain-name>
```

5. Start the services.

```bash
docker compose up --detach --build
```

## Configuration

- Navigate to your domain name in the browser to open up the configuration interface.
- Any edit has to be confirmed with the `ADMIN_PASSWORD` configured in the `.env` file.
- Add members that should be notified about household tasks with phone numbers in the [E.164](https://en.wikipedia.org/wiki/E.164) format.
- Add household tasks that need to be completed every day/week/month.
- Once a day at 12:00 PM, every member will be reminded of their tasks via SMS or WhatsApp.

## Contributing

Please feel free to submit a [pull request](https://github.com/ashiven/housekeepy/pulls) or open an [issue](https://github.com/ashiven/housekeepy/issues).

1. Fork the repository
2. Create a new branch: `git checkout -b feature-name`.
3. Make your changes
4. Push your branch: `git push origin feature-name`.
5. Submit a PR

## License

This project is licensed under the [MIT License](./LICENSE).

---

> GitHub [@ashiven](https://github.com/Ashiven) &nbsp;&middot;&nbsp;
> Twitter [ashiven\_](https://twitter.com/ashiven_)
