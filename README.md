# NakshaDB Server
This has worked on Ubuntu 16.04 with golang 1.10. This works in conjunction
with nakshadb-tiler.

## Dependencies
* sudo apt install zip
* sudo apt install postgresql-9.5 postgresql-9.5-postgis-2.2
* sudo apt install supervisor
* sudo apt install apache2
* sudo apt install gdal-bin
* If you are installing from source,
download [golang 1.10](https://dl.google.com/go/go1.10.linux-amd64.tar.gz)

## Configure PostgreSQL Server
Edit /etc/postgresql/9.5/main/pg_hba.conf
Search for line(s) that look like:  
<pre>
# "local" is for Unix domain socket connections only
local   all             all                                     peer
</pre>

Replace 'peer' with 'md5'. Save and exit.
Restart postgresql service with `sudo service postgresql restart`

## Installation
Assumes that the directories containing the naksha server and naksha tiler code
are under the same parent directory. Something like naksha server is at
~/naksha/server/ and the naksha tiler is at ~/naksha/tiler/
To install from source, you can either clone this repository or download the
latest tagged source code release.
To install from binaries, download the latest tagged binary release.
If you are using binary release, move on to the [set up](#set-up) section.
* `cd /path/to/naksha/server`
* `export GOPATH=/path/to/naksha/server`  

Replace
> /path/to/naksha/server

with full path to the directory where the source code of naksha server is stored.
* `go install github.com/gorilla/sessions`
* `go install github.com/lib/pq`
* `go install golang.org/x/crypto/bcrypt`
* `go build -o bin/adduser cmd/adduser/adduser.go`
* `go build -o bin/delete_old_exports cmd/delete_old_exports/delete_old_exports.go`
* `go build -o bin/genconfig cmd/genconfig/genconfig.go`
* `go build -o bin/importdb cmd/importdb/importdb.go`
* `go build -o bin/server cmd/server/server.go`

## Set Up
Assumes that the go code/binaries are at ~/naksha/server/ and the tiler is at
~/naksha/tiler/. You have to generate the config files and database script. A
tool has been provided for that. Keep the following values at hand.

### Generate config files
* Domain name (or IP address) from which application will be accessed
* Name of the database in which nakshadb data will be stored
* Need three username/password pairs for database access:
  * Admin user and password
  * Application user and password
  * Api user and password

Database and database users will be created. So, make sure that they **do not**
exist. To generate the config files and database script, run the following commands.
* `cd ~/naksha/server`
* `bin/genconfig`  

You will be prompted for the above mentioned information. Enter the values as asked
for. Following files will be generated and stored in current working directory.
* db.sql
* config.xml
* supervisor.naksha.conf
* apache.vhost.conf

Move config.xml to config directory. `mv config.xml config/`

### Prepare Database
#### Create database and database users
* `cd ~/naksha/server`
* `sudo su postgres -c 'psql -f db.sql -b'`
#### Import naksha database
* `bin/importdb`

### Tiler
See the installation instructions on the
[nakshadb-tiler](https://github.com/sri-soham/nakshadb-tiler) repository. It is
important that you install the naksha tiler before moving on to the next steps.
To repeat, assumption is that naksha server is at ~/naksha/server and tiler is
at ~/naksha/tiler.

### Supervisor
Copy the config file generated in the **Generate Config Files** step to the supervisor
config files directory
* `sudo cp supervisor.naksha.conf /etc/supervisor/conf.d/naksha.conf`  

Restart supervisor with `sudo service supervisor restart`

### Apache
Copy the vhost file generated in the **Generate Config Files** step to the apache
vhosts directory
* `sudo cp apache.vhost.conf /etc/apache2/sites-available/naksha.conf`

Run following commands to enable the required apache modules and to enable the
virtual host for naksha server
* `sudo a2enmod proxy`
* `sudo a2enmod proxy_http`
* `sudo a2enomd proxy_fcgi`
* `sudo a2ensite naksha.conf`
* `sudo service apache2 restart`

### Create application user
* `cd ~/naksha/server`
* `bin/adduser`

You will be asked for your name, desired username and password. Application user
account with desired username and password will be generated. Please make note of
the username and the password.

### Ready
Type in the domain name (or the ip address) that you used while generating the config
file into the address bar of your browser. Type in the username and password used
while creating the **Create application user** step to login.

## Clean Up
Delete the db.sql, supervisor.naksha.conf and apache.vhost.conf in the ~/naksha/server
directory.
