$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

$env:GOPATH = $PWD
$env:PATH = $env:GOPATH + "/bin;C:/go/bin;" + $env:PATH

Write-Host "Installing Ginkgo"
go.exe get github.com/onsi/ginkgo/ginkgo
if ($LastExitCode -ne 0) {
    throw "Ginkgo go get process returned error code: $LastExitCode"
}

go.exe install github.com/onsi/ginkgo/ginkgo
if ($LastExitCode -ne 0) {
    throw "Ginkgo go install process returned error code: $LastExitCode"
}

cd src/code.cloudfoundry.org/guardian

go version
# TODO make vet work
go vet ./...
Write-Host "compiling test process: $(date)"

$env:GARDEN_TEST_ROOTFS = "N/A"
ginkgo -r -nodes 8 -race -keepGoing -failOnPending -skipPackage "dadoo,gqt,kawasaki,locksmith,socket2me,signals"
if ($LastExitCode -ne 0) {
    throw "Ginkgo run returned error code: $LastExitCode"
}
ginkgo -r -nodes 8 -race -keepGoing -failOnPending -randomizeSuites -randomizeAllSpecs -skipPackage "dadoo,kawasaki,locksmith" -focus "Runtime Plugin" gqt
Exit $LastExitCode
