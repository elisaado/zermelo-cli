echo "Checking if Go is installed..."
which go && echo "Go is installed!" || (echo "Go is not installed, please install it via your package manager to compile and install zermelo-cli" && return 1)

echo "
Fetching packages..."
go get -u -v github.com/shibukawa/configdir github.com/mgutz/ansi github.com/alexeyco/simpletable
echo "Done!"

echo "
Compiling and installing zermelo-cli..."
go install
echo "Done! Enjoy zermelo-cli!"

echo "
Zermelo-CLI is made by Eli Saado"
