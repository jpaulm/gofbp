module github.com/jpaulm/gofbp

replace github.com/jpaulm/gofbp/components/io => ./components/io

replace github.com/jpaulm/gofbp/core => ./core

replace github.com/jpaulm/gofbp/components/subnets => ./components/subnets

replace github.com/jpaulm/gofbp/components/testrtn => ./components/testrtn

go 1.17

require (
	github.com/jpaulm/gofbp/components/io v0.0.0-00010101000000-000000000000
	github.com/jpaulm/gofbp/components/subnets v0.0.0-00010101000000-000000000000
	github.com/jpaulm/gofbp/components/testrtn v0.0.0-00010101000000-000000000000
	github.com/jpaulm/gofbp/core v0.0.0-00010101000000-000000000000
)
