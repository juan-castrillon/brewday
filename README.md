# BrewDay

![logo](web/static/images/logo_new_250_black.png)



BrewDay is a self-contained web application aimed at helping homebrewers with their brewing process. 

It is intended to be used the day of the brew, and it is designed to be used on all devices, from desktop to mobile.

The app is intended to be self-hosted and does not have multiple users in mind. It is designed to be used by a single user at a time.

The app helps the user with the following tasks:

- **Follow the recipe**. The user can import a recipe from any of the supported formats (see below), and the app will guide the user through the brewing process, step by step. 
- **Note taking**. The user can take notes during the brew, and the app will save them for future reference. Each step in the process gives the opportunity to input real data (to compare with the recipe) and notes (to keep track of the brew).
- **Timers**. The app will set timers for each step in the process, and will notify the user when the time is up. 
- **Statistics**. The app will calculate the efficiency of the brew, evaporation rate, and other useful statistics. If using the Brewing Systems feature, it will split these per system.
- **Timeline and summary**. The app will ley the users download a timeline of the brew, and a summary of the brew day, with all the relevant data. Supported summary formats are listed below.

## Supported recipe formats

The app supports the following recipe formats:
- Native (`brewday`): An example can be found under `test/recipe/brewday/All.json`
- [Maische Malz und Mehr](https://www.maischemalzundmehr.de/index.php?inhaltmitte=lr) ([JSON](https://www.maischemalzundmehr.de/rezept.json.txt))
- [Braureka](https://braureka.de/) (JSON) (This is supposed to be MMUM, but it differs in implementation of some fields that are parsed as strings instead of numbers)


## Supported summary formats

The app supports the following summary formats:
- [Markdown](https://www.markdownguide.org/basic-syntax/): Markdown summary will create a summary of the brew day in Markdown format. This is useful to copy and paste the summary in a blog post, or to share it with other people. The timeline is just a list of timestamps. 

## Supported Notification servers

The app can send notifications via these external servers:

- [Gotify](https://gotify.net/)

```yaml
notification:
  enabled: true
  type: gotify
  settings:
    gotify-url: http://localhost:8080
    gotify-username: "gotify"
    gotify-password: "gotify"
```

- Home Assistant (Companion App): For its configuration, a long lived token needs to be generated and passed (More info [here](https://developers.home-assistant.io/docs/api/rest/)). In addition, the specific device id must be given

```yaml
notification:
  enabled: true
  type: ha
  settings:
    ha-url: http://localhost:8123
    ha-token: "letters1234$_%@"
    ha-device-id: "mydevice"
```

# Installation

## Configuration

The app can be configured via a YAML file, or via environment variables. Environment variables take precedence over the YAML file and can complete or override the configuration.

The application port is required and the app will not start if it is not provided. If notifications are enabled, the different settings are required.

Brewing systems are completely optional. If at least one is set up, the import wizard will give the option to set a system for the recipe and the measurements can then be used to input height instead of volume when measuring. It will also register the system in the stats. If at least one system is declared, only `power-watts` is optional. Diameters are always internal and for the height measure  (from the inside bottom) to the highest point of the pot's rim, the spot where liquid would spill out if you kept filling it.

To pass a configuration file, the application must be run with the `--config` flag, followed by the path to the configuration file. If no configuration file is provided, the app will attempt to read the configuration from environment variables.


The following is an example of a YAML configuration file:

```yaml
app:
  port: 8080

notification:
  enabled: true
  type: ha
  settings:
    ha-url: http://localhost:8123
    ha-token: "letters1234$_%@*"
    ha-device-id: "mydevice"

store:
  type: sql
  path: "./bd.sqlite"
  
process:
  lautern-rest-time-min: 15
  refractometer-wcf: 1.00

brewing-systems:
  - name: bs1
    lower-diameter-cm: 5.1
    upper-diameter-cm: 5.2
    max-height-cm: 10
    power-watts: 2500
  - name: bs2
    lower-diameter-cm: 6.1
    upper-diameter-cm: 6.2
    max-height-cm: 11
```

Store can be `sql` or `memory` depending on the need on persistent storage.

The following is an example of the same configuration via environment variables. Note that `brewing-systems` configuration cannot be passed via env variables and **must** be defined in a YAML file:

```bash
export BREWDAY_NOTIFICATION_ENABLED=true
export BREWDAY_NOTIFICATION_TYPE="ha"
export BREWDAY_NOTIFICATION_SETTINGS_HA-DEVICE-ID='my_device'
export BREWDAY_NOTIFICATION_SETTINGS_HA-TOKEN='t0ken$#'
export BREWDAY_NOTIFICATION_SETTINGS_HA-URL="http://localhost:8123",
export BREWDAY_APP_PORT=8080
export BREWDAY_PROCESS_LAUTERN-REST-TIME-MIN=15
export BREWDAY_PROCESS_REFRACTOMETER-WCF=1.00
```

> Process variables can be skipped. The default values are shown in the example above

## Deployment

The app can be deployed as a Docker container, or as a standalone binary. In order for the notification to work, a [Gotify](https://gotify.net/) server or a Home Assistant installation must be available.

The recommended way to deploy the app is via Docker. In the `deployments` folder there is a `docker-compose.yml` file that can be used to deploy the app. The file can be used as is, or it can be used as a template to create a custom deployment.

It also includes deployment of a [Gotify](https://gotify.net/) server, which is can be used for notifications to work. The Gotify server is deployed with a volume for the data, so it can be restarted without losing the data.

In the same folder, a `docker-compose-arm.yaml` file is provided, which can be used to deploy the app on a Raspberry Pi. The file can be used as is, or it can be used as a template to create a custom deployment.
