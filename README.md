# TASRED

Join a Twitch IRC channel and classify messages there as "naughty" or "nice".

## Installation

You'll need a functional Go setup per [the installation instructions][goinstall].  1.4 is known to work; 1.3 probably will as well.  The version in Debian Jessie is known to be too old.

[goinstall]: http://golang.org/doc/install

Install from git:

    go get github.com/twm/TASRED

The `TASRED` binary should now be available in `$GOPATH/bin`.  You may want to put this directory on your `PATH`.

## Configuration

You'll need to create a `config.toml` file from this template:

    [irc]
    nick = "tasred"
    password = "..."
    channel = "#agdq"

You will need to generate a special password here: <http://twitchapps.com/tmi/>.  The account's normal password will *not* work.

You may also want to grab `bad-words.txt` from the repository.
