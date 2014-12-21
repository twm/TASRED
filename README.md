# TASRED

Join a Twitch IRC channel and classify messages there as "naughty" or "nice".

## Compilation

You'll need a functional Go setup per <http://golang.org/doc/install>.  1.4 is known to work; 1.3 probably will as well.

## Configuration

You'll need to create a `config.toml` file from this template:

    [irc]
    nick = "tasred"
    password = "..."
    channel = "#agdq"

You will need to generate a special password here: <http://twitchapps.com/tmi/>.  The account's normal password will *not* work.
