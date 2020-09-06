echo "Compiling..."
go build
echo "Installing..." 
cp mkhomework /usr/local/bin
echo "Creating template folder at $HOME/.mkhomework/" 
cp -r templates/ $HOME/.mkhomework/