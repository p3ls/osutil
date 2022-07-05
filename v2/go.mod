module github.com/p3ls/osutil/v2

go 1.16

require (
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.62.0
)

retract (
	// Breaking changes
	// curl proxy.golang.org/github.com/p3ls/osutil/v2/@v/list |sort -V
	v2.0.0
	v2.0.0-rc
	v2.0.0-rc.1
	v2.0.0-rc.2
	v2.0.0-rc.3
	v2.0.0-rc.4
	v2.0.0-rc.6
	v2.0.0-rc.7
	v2.0.0-rc.8
	v2.0.0-rc.9
	v2.0.0-rc.10
	v2.0.0-rc.11
	v2.0.0-rc.12
	v2.0.0-rc.14
	v2.0.0-rc.15
	v2.0.0-rc.16
	v2.0.0-rc.17

	v2.0.1 // Retract breaking releases
)
