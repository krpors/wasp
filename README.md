# Wasp

A web-server/web-application based remote control for Linux 'media-devices', implemented
in the Go programming language (actually just an interface to Mplayer).

# Introduction

Small Linux devices are gaining popularity, especially since the 'invention' of 
ASUS EeePC and of course the [Raspberry Pi](http://raspberrypi.org). One of the 
many applications of such Linux devices, is functioning as a sort of media player.
The Raspberry Pi in particular, since it's not a netbook but just a headless computer.

A mediaplayer is obviously nothing much without some sort of lazy remote controlling :)
Wasp's intention is to be that remote control, albeit not in traditional manner. Instead
of using LIRC (infrared) or Bluetooth, the controlling is done with a web interface Wasp
is providing.

Once Wasp is started, it acts as a webserver to which you can browse to. That means you
can open a browser on Android, iOS, a PC, laptop or netbook to the IP address of the 
device where Wasp is running, like ``http://192.168.1.10:8080/index``. Once opened, it will
provide you with a basic interface where you can:

* Browse for media files in an intuitive manner.
* Play/pause media files.
* Change volume.
* Stop playback.
* and likely some more features added in the future.

# Getting started

The project is of course still in its infancy, but follow these basic steps to start off.

**Prerequisites:**

1. Install a [Go compiler](http://code.google.com/p/go/downloads/list).
1. Install [Mplayer](http://www.mplayerhq.hu/design7/dload.html) (or use a package manager for your distro).

**Build the sources:**

Clone the git repository:

`git clone git://github.com/krpors/wasp.git`

Change directory to the just cloned repository:

`cd ./wasp`

Set the GOPATH to the current directory:

`export GOPATH=$(pwd)`

Build the sources:

`make`

Try invoking the binary:

`./bin/wasp`

This will create the initial configuration in `${HOME}/.wasp/wasp.properties`. This configuration
sets a few properties as follows:

* Mplayer fifo pipe: `/tmp/mplayer.fifo`
* Media directory: `/` (root)
* Bind address: `:8080`. This will bind the webserver to port 8080 on the local host.
* Allowed/listable audio/video file extensions.

Once Wasp is started, a FIFO will automatically be created in the configured path (if you have create/write
permissions, of course).

Last but not least, start Mplayer in slave mode, with the input to be from the FIFO:

`while true; do mplayer -noconfig all -noconsolecontrols -quiet -idle -slave -fs -zoom -input file=/tmp/mplayer.fifo; done`

The `while true` makes sure Mplayer keeps running at all times, in case unwanted intervention has been done.
This is not a necessity though, but something to consider.

After this, you should have two processes running: mplayer and wasp. 

Try opening up a browser to the host where Wasp is running, e.g. http://localhost:8080/. You'll see
web interface as follows:

![The control page](http://krpors.github.com/wasp/img/wasp-control.png "Control page in Chrome")
![Browsing page](http://krpors.github.com/wasp/img/wasp-browse.png "Browsing in Chrome")
![Configuration page](http://krpors.github.com/wasp/img/wasp-config.png "Configuration in Chrome")

# Implementation

For those interested in some implementation details:

As mentioned previously, the actual code is done using the [Go](http://golang.org) language.
The web interface is HTML5, CSS3 and Javascript, with help of [jQuery Mobile](http://jquerymobile.com)
for UI widgets and [jQuery](http://jquery.com) for other common JS tasks.

The A/V backend is [Mplayer](http://www.mplayerhq.hu), which can be run in a 'slave mode' to
accept commands from external applications.



## Self-notes

Just some notes for development.

* [Slave mode Mplayer](http://www.mplayerhq.hu/DOCS/HTML/en/MPlayer.html#slave-mode)
* [Slave mode commands](http://www.mplayerhq.hu/DOCS/tech/slave.txt)
* ``mplayer -noconfig all -noconsolecontrols -quiet -idle -slave -fs -zoom -input file=/tmp/mplayer.fifo``
