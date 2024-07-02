blacklist_url=https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt
adguard_url=https://adguardteam.github.io/HostlistsRegistry/assets/filter_1.txt

custom="
m.vpon.com
"
blacklist=`curl $blacklist_url | grep -v '^#' | awk '{ print $2 }'`
adguard=`curl $adguard_url | grep -v '^!' | grep -v '^#' | sed -E 's/^@+//' | sed -E 's/^\|+//' | sed -E 's/^-//' |sed -E 's/\^\|?$//' | grep -v '*'`

echo "$custom\n$blacklist\n$adguard" | sort | uniq | grep -E '^([a-zA-Z0-9\-_]+.)+[a-zA-Z]{2,}$' > adservers.txt

git add .
git commit --allow-empty -m `date +%Y-%m-%d`
git push origin master
