# user

Provides access to the users database. It is available for Linux (by now).

[Documentation online](http://gowalker.org/github.com/p3ls/osutil/user)

## Installation

    go get github.com/p3ls/osutil/user

To run the tests, it is necessary to run them as root.  
Do not worry because the tests are done in copies of original files.

    sudo env PATH=$PATH go test -v

## Status

BSD systems and Windows unsopported.

The only backend built is to handle files (such as '/etc/passwd'), and it it not
my priority to handle other backends like LDAP or Kerberos since my goal was
to can use it in home systems.

My list of priorities are (for when I have time):

- BSD systems (included Mac OS)
- Windows

## Configuration

Some values are got from the system configuration, i.e. to get the next
available UID or GID, but every distribution of a same system can have a
different configuration system.

In the case of Linux, the research has been done in 10 different distributions:

    Arch
    CentOS
    Debian
    Fedora
    Gentoo
    Mageia (Mandriva's fork)
    OpenSUSE
    PCLinuxOS
    Slackware
    Ubuntu

## License

The source files are distributed under the [Mozilla Public License, version 2.0](http://mozilla.org/MPL/2.0/),
unless otherwise noted.  
Please read the [FAQ](http://www.mozilla.org/MPL/2.0/FAQ.html)
if you have further questions regarding the license.
