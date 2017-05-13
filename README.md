# McDowell

[![Build Status](https://travis-ci.org/willmadison/mcdowell.svg?branch=master)](https://travis-ci.org/willmadison/mcdowell)

![Cleo McDowell](http://mcdowells.mortenjonassen.dk/img/staff/cleo1.jpg)

This is a Slack bot for the Atlanta Black Tech Slack.

You can get an invite from [here](https://atlblacktech-slack-invite.herokuapp.com/)

## Building

This bot requires Go 1.7+ and can be built as follows:

```
    go build ./cmd/mcdowell
```

## Running

To run this you need to set the the following environment variables:
- ` ABT_SLACK_BOT_TOKEN ` - the Slack bot token
- ` ABT_SLACK_BOT_DEV_MODE ` - boolean, set the bot in development mode

```
    ABT_SLACK_BOT_TOKEN=<TOKEN_HERE> ./mcdowell
```
