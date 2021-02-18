# Grafana Webhook to WeCom

> A Webhook implementation of WeCom alert channel for grafana

## Quickstart

run with default port `80`:

`./grafana-webhook-wecom YOUR_WECOM_WEBHOOK_API_ADDR`

run with specified port:

`./grafana-webhook-wecom YOUR_WECOM_WEBHOOK_API_ADDR 6000`

API endpoint

`/api/webhook/wecom`

## How to use it in Grafana

- Create a webhook channel with `POST` method
- Put `http://YOUR_ADDRESS/api/webhook/wecom` in Url
