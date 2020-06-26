# Removes all binaries.
rm ./gopkgs*

# Linux amd64
GOOS=linux GOARCH=amd64 go build -tags sqlite3 -o gopkgs-Linux-amd64-sqlit3
GOOS=linux GOARCH=amd64 go build -tags mysql -o gopkgs-Linux-amd64-mysql
GOOS=linux GOARCH=amd64 go build -tags postgres -o gopkgs-Linux-amd64-postgres
GOOS=linux GOARCH=amd64 go build -tags sqlserver -o gopkgs-Linux-amd64-sqlserver

# Windows amd64
GOOS=windows GOARCH=amd64 go build -tags sqlite3 -o gopkgs-Windows-amd64-sqlit3
GOOS=windows GOARCH=amd64 go build -tags mysql -o gopkgs-Windows-amd64-mysql
GOOS=windows GOARCH=amd64 go build -tags postgres -o gopkgs-Windows-amd64-postgres
GOOS=windows GOARCH=amd64 go build -tags sqlserver -o gopkgs-Windows-amd64-sqlserver
