echo "Checking if Go is installed..."
type go
if [ $? != 0 ]; then echo "Go is not installed, please install it via your package manager to compile and install zermelo-cli" && exit 1; fi
echo "Go is installed"

echo "
Fetching packages..."
go get -u -v github.com/shibukawa/configdir github.com/mgutz/ansi github.com/alexeyco/simpletable
if [ $? != 0 ]; then  echo "Some error occured..." && exit 1; fi
echo "Done!"

echo "
Compiling and installing zermelo-cli..."
go install
if [ $? != 0 ]; then echo "Some error occured..." && exit 1; fi
echo "Done! Enjoy zermelo-cli!"

echo "
Zermelo-CLI is made by Eli Saado"
