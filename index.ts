import fs from 'fs'

const BLACKLIST_URL = 'https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt'
const ADGUARD_URL = 'https://adguardteam.github.io/HostlistsRegistry/assets/filter_1.txt'
const CUSTOM_DOMAINS = [
  'm.vpon.com'
]

async function fetchBlacklist() {
  const response = await fetch(BLACKLIST_URL)
  const body = await response.text()
  const domains = body.split('\n')
    .filter((line) => !line.startsWith('#'))
    .map((line) => line.split(' ')[1])
    .filter((line) => line)

  return domains
}

async function fetchAdguard() {
  const response = await fetch(ADGUARD_URL)
  const body = await response.text()
  const domains = body.split('\n')
    .filter((line) => !line.startsWith('!') && !line.startsWith('#'))
    .map((line) => line
      .replace(/^@+/, '')
      .replace(/^\|+/, '')
      .replace(/^-/, '')
      .replace(/\^\|?$/, '')
      .replace(/\$important$/, '')
    )
    .filter((line) => line && !line.includes('*'))

  return domains
}

async function fetchAdservers() {
  const [blacklist, adguard] = await Promise.all([
    fetchBlacklist(),
    fetchAdguard()
  ])
  const adservers = [...new Set([ ...CUSTOM_DOMAINS, ...blacklist, ...adguard])]
    .filter((domain) => domain.match(/^([a-zA-Z0-9-_]+\.)+[a-zA-Z]{2,}$/))

  return adservers
}

const adservers = await fetchAdservers()
fs.writeFileSync('adservers.txt', adservers.join('\n'))
