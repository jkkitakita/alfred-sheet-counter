# alfred-sheet-counter

## Required

- Go >= 1.14

## Usage

```bash
$ ./alfred-sheet-counter -h
Usage:
  alfred-sheet-counter [OPTIONS]

Application Options:
  -c, --column= a column to update

Help Options:
  -h, --help    Show this help message

2020/05/16 01:29:46 Unable to parse flag: Usage:
  alfred-sheet-counter [OPTIONS]

Application Options:
  -c, --column= a column to update

Help Options:
  -h, --help    Show this help message
```

## Getting Started

1. Enable the Google Sheets API and download credentials.json
<https://developers.google.com/sheets/api/quickstart/go#step_1_turn_on_the>
2. Place credentials.json on root dir
3. Modify const in main.go
   1. `spreadsheetID`
   2. `sheetName`
   3. `readRange`
4. build
```bash
GO111MODULE=on go build
```
5. Since the URL for oauth authorization is output to the standard output, follow the instructions and execute authorization with the Google account that holds the target spreadsheet.
6. Execute alfred-sheet-counter
```bash
# e.g. Update B Column
./alfred-sheet-counter -c=B
```
