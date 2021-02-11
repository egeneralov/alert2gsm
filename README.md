# alert2gsm

## WARNING

Project just a [MVP](https://en.wikipedia.org/wiki/Minimum_viable_product).

## Installation

- pre-check
  - requirements:
    - CPU: ~ 100 millicores
    - RAM: ~ 20 mb, based on concurent count for your alerts
    - NET:
      - public port, mapped on each instance, twilio will execute a webhook
      - access to dns resolving and api.twilio.com:443

- Twilio account
  - [register](https://www.twilio.com/try-twilio)
  - [buy a number](https://www.twilio.com/console/phone-numbers/search)
  - get `ACCOUNT SID` and `AUTH TOKEN` from (console)[https://www.twilio.com/console]
    - put in env as `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN`
  - Enshure your destination contry are in low-risk lists (example4russia)[https://www.twilio.com/console/voice/calls/geo-permissions/low-risk?countryIsoCode=rus]

- configuration
  - **external_endpoint**: "http://${IP_OR_DOMAIN}:${EXT_PORT}"
  - **http.enabled**: {en/dis}able http server
  - **http.listen**: "host:port"

  - **https.enabled**: {en/dis}able https server
  - **https.listen**: "host:port"
  - **https.ssl_certificate**: path to cert
  - **https.ssl_certificate_key**: path to key

  - **webhooks.call**: test route `GET /call/` for execute a test call

  - **webhooks.grafana.enabled**: {en/dis}able grafana webhook endpoint
  - **webhooks.grafana.username**: grafana auth username
  - **webhooks.grafana.password**: grafana auth password

  - **webhooks.twilio.enabled**: {en/dis}able twilio endpoint, must be `true`
  - **webhooks.twilio.voice**: you can choose robot voice from (polly voices)[https://docs.aws.amazon.com/polly/latest/dg/voicelist.html]
  - **webhooks.twilio.from**: your phone number in twilio
  - **webhooks.twilio.notify**: `[]string` - list of alert recivers in international format

