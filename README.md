# osutil

WARNING! The API of version v2.0.0 could change.

Access to operating system functionality dependent of every platform.

+ config/env: set persistent environment variables
+ config/shconf: parser and scanner for the configuration in format shell-variable
+ edi: editing of files
+ executil: executes commands in shells
+ fileutil: handles common operations at files
+ sysutil: defines operating systems, detects Linux distributions and handles different package systems
+ userutil: provides access to UNIX users database in local files
+ user/crypt: password hashing used in UNIX

[Documentation online](http://godoc.org/github.com/tredoe/osutil)

## Testing

`go test ./...`

`sudo env PATH=$PATH go test -v ./...`

'sudo' command is necessary to copy the files '/etc/{passwd,group,shadow,gshadow}' to the temporary directory, where the tests are run.
Also, it uses 'sudo' to check the package manager, at installing and removing the package 'nano'.


## License

The source files are distributed under the [Mozilla Public License, version 2.0](http://mozilla.org/MPL/2.0/),
unless otherwise noted.  
Please read the [FAQ](http://www.mozilla.org/MPL/2.0/FAQ.html)
if you have further questions regarding the license.
