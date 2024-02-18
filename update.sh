curl https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt | grep -v '^#' | awk '{ print $2 }' > adservers.txt
git add .
git commit --allow-empty -m `date +%Y-%m-%d`
git push origin master
