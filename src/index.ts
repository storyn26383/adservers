import fs from 'fs'

const DOMAIN_LIST_URLS = [
  'https://adguardteam.github.io/HostlistsRegistry/assets/filter_1.txt',
  'https://adguardteam.github.io/HostlistsRegistry/assets/filter_50.txt',
  'https://filters.adtidy.org/extension/ublock/filters/2_optimized.txt',
  'https://filters.adtidy.org/extension/ublock/filters/11_optimized.txt',
  'https://filters.adtidy.org/extension/ublock/filters/224_optimized.txt',
]
const CUSTOM_DOMAINS = [
  'm.vpon.com'
]

async function fetchDomainList(url: string) {
  const response = await fetch(url)
  const body = await response.text()
  const domains = body.split('\n')
    .filter((line) => !line.startsWith('!') && !line.startsWith('#'))
    .map((line) => line
      .trim()
      .replace(/^\|\|/, '.')
      .replace(/^:\/\//, '')
      .replace(/\^$/, '')
      .replace(/\^\$important$/, '')
    )
    .filter((line) => line)

  return domains
}

async function fetchAdservers() {
  const domains = await Promise.all(DOMAIN_LIST_URLS.map((url) => fetchDomainList(url)))
  const adservers = [...new Set([ ...CUSTOM_DOMAINS, ...domains.flat()])]
    .filter((domain) => domain.match(/^\.?([a-zA-Z0-9-_*]+\.)+[a-zA-Z]{2,}$/))

  return adservers
}

const adservers = await fetchAdservers()
fs.writeFileSync('adservers.txt', adservers.join('\n'))
