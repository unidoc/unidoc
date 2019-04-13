node {
    // Install the desired Go version
    def root = tool name: 'go 1.11.5', type: 'go'

    env.GOROOT="${root}"
    env.GOPATH="${WORKSPACE}/gopath"
    env.PATH="${root}/bin:${env.GOPATH}/bin:${env.PATH}"
    env.GOCACHE="off"
    env.UNIDOC_EXTRACT_FORCETEST="1"
    env.UNIDOC_E2E_FORCE_TESTS="1"
    env.UNIDOC_EXTRACT_TESTDATA="/home/jenkins/corpus/unidoc-extractor-testdata"
    env.UNIDOC_RENDERTEST_BASELINE_PATH="/home/jenkins/corpus/unidoc-creator-render-testdata"
    env.UNIDOC_PASSTHROUGH_TESTDATA="/home/jenkins/corpus/unidoc-e2e-testdata"
    env.UNIDOC_ALLOBJECTS_TESTDATA="/home/jenkins/corpus/unidoc-e2e-testdata"
    env.UNIDOC_SPLIT_TESTDATA="/home/jenkins/corpus/unidoc-e2e-split-testdata"
    env.UNIDOC_GS_BIN_PATH="/usr/bin/gs"
    // Hack for 1.11.5 testing work.
    env.CGO_ENABLED="0"

    env.TMPDIR="${WORKSPACE}/temp"
    sh "mkdir -p ${env.TMPDIR}"

    dir("${GOPATH}/src/github.com/unidoc/unidoc") {
        sh 'go version'

        stage('Checkout') {
            echo "Pulling unidoc on branch ${env.BRANCH_NAME}"
            checkout scm
        }

        stage('Prepare') {
            // Get linter and other build tools.
            sh 'go get -u golang.org/x/lint/golint'
            sh 'go get github.com/tebeka/go2xunit'
            sh 'go get github.com/t-yuki/gocover-cobertura'

            // Get all dependencies (for tests also).
            sh 'go get -t ./...'
        }

        stage('Linting') {
            // Go vet - List issues
            sh '(go vet ./... >govet.txt 2>&1) || true'

            // Go lint - List issues
            sh '(golint ./... >golint.txt 2>&1) || true'
        }

        stage('Testing') {
            // Go test - No tolerance.
            sh "rm -f ${env.TMPDIR}/*.pdf"
            sh '2>&1 go test -v ./... | tee gotest.txt'
        }

        stage('Check generated PDFs') {
            // Check the created output pdf files.
            sh "find ${env.TMPDIR} -maxdepth 1 -name \"*.pdf\" -print0 | xargs -t -n 1 -0 gs -dNOPAUSE -dBATCH -sDEVICE=nullpage -sPDFPassword=password -dPDFSTOPONERROR -dPDFSTOPONWARNING"
        }

        stage('Test coverage') {
            sh 'go test -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./...'
            sh '/home/jenkins/codecov.sh'
            sh 'gocover-cobertura < coverage.out > coverage.xml'
            step([$class: 'CoberturaPublisher', coberturaReportFile: 'coverage.xml'])
        }

        stage('Post') {
            // Assemble vet and lint info.
            warnings parserConfigurations: [
                [pattern: 'govet.txt', parserName: 'Go Vet'],
                [pattern: 'golint.txt', parserName: 'Go Lint']
            ]

            sh 'go2xunit -fail -input gotest.txt -output gotest.xml'
            junit "gotest.xml"
        }
    }

    dir("${GOPATH}/src/github.com/unidoc/unidoc-examples") {
        stage('Build examples') {
            // Output environment variables (useful for debugging).
            sh("printenv")

            // Pull unidoc-examples from connected branch, or master otherwise.
            def examplesBranch = "master"
            if (env.BRANCH_NAME.take(3) == "PR-") {
                // Pull requests (PR) require separate handling.
                if (env.CHANGE_TARGET.take(2) == "v3") {
                    examplesBranch = "v3"
                }
            } else {
                if (env.BRANCH_NAME.take(2) == "v3") {
                    examplesBranch = "v3"
                }
                // Special cases.
                switch("${env.BRANCH_NAME}") {
                    case "v3-reduce-fonts-exports":
                        examplesBranch = "v3-reduce-fonts-exports"
                        break
                }
            }

            // Check if connected branch is defined explicitly.
            def safeName = env.BRANCH_NAME.replaceAll(/[\/\.]/, '')
            def fpath = "/home/jenkins/exbranch/" + safeName
            if (fileExists(fpath)) {
                examplesBranch = readFile(fpath).trim()
            }

            echo "Pulling unidoc-examples on branch ${examplesBranch}"
            git url: 'https://github.com/unidoc/unidoc-examples.git', branch: examplesBranch
            
            // Dependencies for examples.
            sh 'go get github.com/wcharczuk/go-chart'

            // Build all examples.
            sh 'find . -name "*.go" -print0 | xargs -0 -n1 go build'
        }

        stage('Passthrough benchmark pdfdb_small') {
            sh './pdf_passthrough_bench /home/jenkins/corpus/pdfdb_small/* | grep -v "Testing " | grep -v "copy of" | grep -v "To get " | grep -v " - pass"'
        }
    }
}
