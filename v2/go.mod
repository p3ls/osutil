module github.com/tredoe/osutil/v2

go 1.16

require (
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.62.0
)

// Breaking changes
// curl proxy.golang.org/github.com/tredoe/osutil/v2/@v/list |sort -V
retract v2.0.0
