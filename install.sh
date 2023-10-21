#install chrome
set -x && \
apt update && \
apt upgrade -y && \
apt install -y wget curl gnupg libappindicator1 fonts-liberation locales fonts-noto-cjk && \
echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list && \
wget -qO - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - && \
apt update && \
apt install -y google-chrome-stable && \
rm -rf /var/lib/apt/lists/ && \
apt-get autoremove -y && \
apt-get autoclean -y

#download linux
curl -s https://api.github.com/repos/srcrs/magicbox/releases/latest | grep browser_download_url | grep linux | cut -d'"' -f4 | wget -i -