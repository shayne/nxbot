[blogurl]: https://medium.com/@thinking_43465/nxbot-programmable-surveillance-2d8c2de9c81e
[wikiurl]: ../../wiki

# NxBot

A Telegram bot, paired with Nx Witness, providing motion-event and on-demand camera snapshots.

This project accompanies a [blog post][blogurl] outlining the concept and implementation details of this bot.

The [Wiki][wikiurl] contains helpful guides for setting up Telegram and Nx Witness to work with NxBot.

![1_c3_3visxh25cxf6f_hzyeq](https://user-images.githubusercontent.com/47304512/52295195-0759b780-2930-11e9-9a4a-bf352fd91900.png)
    ![1_sfxaoui8bx6793hecdycfq](https://user-images.githubusercontent.com/47304512/52295192-04f75d80-2930-11e9-8c94-4894aadad133.gif)


# Docker

NxBot is made available as a Docker container. It's packaged as a standalone binary, see [Dockerfile](Dockerfile) for build the process.


## Usage

```
docker create \
--name=nxbot \
-e NX_IP_PORT=<nx-server-ip>:<port> \
-e NX_USER=<USER> \
-e NX_PASS=<PASS> \
-e HTTP_IP_PORT=<http-ip-port> \
--expose <http-port> \
-e TG_TOKEN=<telegram-bot-token> \
-e TG_USER_WHITELIST=<whitelist-ids> \
-e TG_GROUP_WHITELIST=<whitelist-ids> \
-e TG_MOTION_RECIPIENTS=<recipient-ids> \
jacknx/nxbot
```

## Parameters

All parameters are required, except for `TG_USER_WHITELIST`, `TG_GROUP_WHITE_LIST` and `TG_MOTION_RECIPIENTS`. It is, however, highly suggested that you use these lists, as refusal to do so leaves your bot open to the public.

* `-e NX_IP_PORT=` - The IP and port of the Nx Server. *The port is typically 7001, on a standard install*
* `-e NX_USER=` - The user account to the Nx Server. *Any account with desired bot accessible camera access*
* `-e NX_PASS=` - The password for the above user account.
* `-e HTTP_IP_PORT=` - The IP and port the container listens on for motion events. *Typically the IP will be 0.0.0.0 or left blank, for all interfaces (e.g. ":8012")*
* `--expose <http-port>` - It's required to expose the above in typical Docker systems. *Alternatively you can use the `-p <port>:<port>` syntax if desired*
* `-e TG_TOKEN=` - The Telegram Bot token *See the [Wiki][wikiurl] for help setting up a Telegram Bot*
* `-e TG_USER_WHITELIST=` - A list of Telegram users IDs allowed access to use the bot *See the [Wiki][wikiurl] for more help*
* `-e TG_GROUP_WHITELIST=` - A list of Telegram group IDs allowed access to use the bot *See the [Wiki][wikiurl] for more help*
* `-e TG_MOTION_RECIPIENTS=` - A list of Telegram users or group IDs to receive motion event messages *See the [Wiki][wikiurl] for more help*


## Running the bot

Starting NxBot:

    docker start nxbot

You can see status and any potential errors with the bot by accessing its log:

    docker logs nxbot

You may find errors or warnings pointing out issues with the above environment variables or invalid credentials to Nx and/or Telegram.

When successfully started, the bot will report:

    Starting Nx Telegram Bot and motion HTTP server

## Info

* To monitor the logs of the container in realtime: `docker logs -f nxbot`
* There's a dedicated [blog post][blogurl] over on Medium about this project
* Detailed information about NxBot and for help setting up Telegram can be found on the [Wiki][wikiurl]