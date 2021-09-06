module github.com/tredoe/osutil

go 1.16

require (
	github.com/tredoe/fileutil v1.0.5
	github.com/tredoe/goutil v1.0.0
)

retract (
        // Breaking changes
        // curl proxy.golang.org/github.com/tredoe/osutil/@v/list|sort -V
        v1.1.1
        v1.1.2
        v1.1.3
        v1.1.4
        v1.1.5
        v1.1.6
        v1.1.7
        v1.1.8
        v1.2.0
        v1.2.1
        v1.3.0
        v1.3.1
        v1.3.2
        v1.3.3
        v1.3.4
        v1.3.5

        v1.3.6 // Retract breaking releases
)
