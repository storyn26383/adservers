blacklist_url=https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt
adguard_url=https://adguardteam.github.io/HostlistsRegistry/assets/filter_1.txt

blacklist=`curl $blacklist_url | grep -v '^#' | awk '{ print $2 }'`
adguard=`curl $adguard_url | grep -v '^!' | grep -v '^#' | sed -E 's/^@+//' | sed -E 's/^\|+//' | sed -E 's/^-//' |sed -E 's/\^\|?$//' | grep -v '*' | grep -v '^$'`

echo "$blacklist\n$adguard" | sort | uniq > adservers.txt

git add .
git commit --allow-empty -m `date +%Y-%m-%d`
git push origin master
