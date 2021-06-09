version=0.1.0

rm ./sender

echo "🟩 Downloading agent"
wget https://github.com/senderchan/go-sender/releases/download/v${version}/sender_${version}_Linux_x86_64.tar.gz
tar -zxvf sender_${version}_Linux_x86_64.tar.gz
rm -f sender_${version}_Linux_x86_64.tar.gz
chmod +x ./sender
mv -f ./sendersender_${version}_Linux_x86_64.tar.gz /usr/bin

echo "\n\n\n\n"

echo "🟩 Agent version"
sender-agent version
echo "====================================="
echo "Sender agent cli install finish"
echo 'run "sender" will create a config file in current dir'
echo 'run "sender service" will create systemctl config file "sender.service" in "/etc/systemd/system/"'
echo 'run "systemctl start sender" for start'
echo "====================================="