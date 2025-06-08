# Scripts for dream project

### Requirements
- `jq` (for JSON encoding)
- `curl`

Install `jq` on macOS:
```sh
brew install jq
```
Install `jq` on Ubuntu/Debian:
```sh
sudo apt-get install jq
```

## send_data.sh

Sends process data to the upload endpoint. 
The script generates a unique user id (email) upon every request.
Now requires two arguments: the data file and the OS type (`windows` or `linux`).

### Usage

```sh
./send_data.sh <filename> <os_type>
```
- `<filename>`: Path to the process output file (e.g., `../testdata/windows_tasklist.txt`)
- `<os_type>`: Either `windows` or `linux` (sets the appropriate fields in the payload)

#### Examples
```sh
./send_data.sh ../testdata/windows_tasklist.txt windows
./send_data.sh ../testdata/linux_ps.txt linux
```
---

## clean_db.sh
Cleans all tables in the PostgreSQL database.

## check_db.sh
Shows the current state of the database tables.

---

### Making scripts executable
```sh
chmod +x *.sh
``` 