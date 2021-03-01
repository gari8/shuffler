# shuffler


## What you can do with this
- You can randomly swap the lines in the CSV file.
- You can pin any column when swapping rows.
- You can select multiple columns that can be fixed.

## How to install
```
go get github.com/gari8/shuffler
```

## How to use
```
// Outputs a guide message.
shuffler -h

// Creates a csv file sorted based on the csv in the entered path.
shuffler -p <filepath>

// Options
shuffler -p <filepath> -n <newfile's name> -c <csv's row count>
```
