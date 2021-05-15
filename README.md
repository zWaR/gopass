# gopass

CLI tool for working with keepass files written in go.

## Features

* Reading single file
* Reading multiple files
* Search single or multiple files from a single console

## Usage

### CLI commands

CLI supports history and colored command hints.

Avaliable commands:
* u     - copies the username
* p     - copies the password
* find  - searches for an entry in the keepass file(s)
* cont  - continues with opening of other kdbx files
* ls    - Lists currently opened KeePass databases
* open  - Opens a KeePass database
* close - Closes a KeePass database
* help  - Shows usage help
* quit  - Quits gopass

![gopass demo](img/gopass-single-file.gif)

### Open a single file:

```bash
gopass open --kdbx myDatabase.kdbx
```

### Open multiple files

If you have not done so already, initialize the gopass config, which will store the paths to your kdbx files:

```bash
gopass initConfig
```

This will create a new config in `$HOME/.gopass/gopass.json` Edit it and enter paths to your kdbx files. gopass supports environment variable expansion, so you're free to use them in the paths.

Open the kdbx files with a single command:

```bash
gopass many
```

gopass will ask you for password for each file it tries to import. If you do not enter anything in the password prompt and just press enter, you will enter the CLI, which will allow you to search for and copy credentials. Once you want to continue opening the files, you can use the `cont` command. This feature allows you to copy the password of a kdbx file from the already opened ones.

### Searching records

If you use `find` you can search for the records. gopass case-insesitively evaluates username and title fields, but it does not evaluate the paths. `find` displays a summary of all records matching your criteria, prefixing them with IDs. The IDs can be used with `u` and `p` commands to copy username or password, respectively.

`u` and `p` commands act like `find` if there are more than one entry that matches the criteria.

### Copying passwords

After copying a passwords with `p`, clipboard will be cleared 10 seconds after copy took place, if the clipboard still holds the copied password.
