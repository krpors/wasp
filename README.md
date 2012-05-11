# Wasp

A web-server/application based remote control for Linux 'media-devices', implemented
in the Go programming language.

# Introduction

Small Linux devices are gaining popularity, especially since the 'invention' of 
ASUS EeePC and of course the [Raspberry Pi](http://raspberrypi.org). One of the 
many applications of such Linux devices, is functioning as a sort of media player.
The Raspberry Pi in particular, since it's not a netbook but just a headless computer.

A mediaplayer is obviously nothing much without some sort of lazy remote controlling :)
Wasp's intention is to be that remote control, albeit not in traditional manner. Instead
of using LIRC (infrared) or Bluetooth, the controlling is done with a web interface.

# Plan

I'm attempting to write this using the [Go](http://www.golang.org) programming language,
because I think it's so awesome, and I want to keep it 'lightweight': a single executable 
for the actual program with a few 'helper files', like HTML template files. My idea is 
to expose certain URIs on the localhost, and when a POST or GET (I'm not sure yet) is done,
an MPlayer command is executed. Example run:

1. Wasp is started on localhost:8080 (lets say this is 192.168.1.2).
1. User points to http://192.168.1.2:8080 using a phone.
1. The web interface from Wasp will list available media files.
1. User selects the file to be played. A HTTP post is executed to for example 
http://192.168.1.2/play?v=/home/user/myvideo.mp4
1. Wasp forwards the request to the AV player to play the file.

## Self-notes

* [Slave mode Mplayer](http://www.mplayerhq.hu/DOCS/HTML/en/MPlayer.html#slave-mode)
* [Slave mode commands](http://www.mplayerhq.hu/DOCS/tech/slave.txt)
* Use a named pipe to issue commands (mkfifo)
* SQLite for storing media files etc? See [sqlite.go](http://code.google.com/p/gosqlite/) for an
interface to SQLite.
* ``mplayer -noconfig all -noconsolecontrols -quiet -idle -slave -fs -zoom -input file=/tmp/mplayer.fifo``
