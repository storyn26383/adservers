package main

import (
  "fmt"
  "net/http"
  "io"
  "strings"
  "regexp"
  "os"
  "log"
)

const (
  BLACKLIST_URL = "https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt"
  ADGUARD_URL = "https://adguardteam.github.io/HostlistsRegistry/assets/filter_1.txt"
)

var CUSTOM_DOMAINS = []string{
  "m.vpon.com",
}

type ParallelCallback[T any] func() (T, error)
type ParallelResult[T any] struct {
  index int
  result T
  e error
}

func debug(args ...any) {
  log.New(os.Stderr, "", 0).Println(args...)
}

func merge[T any](arrays [][]T) []T {
  result := []T{}

  for _, array := range arrays {
    result = append(result, array...)
  }

  return result
}

func unique[T comparable](array []T) []T {
  result := []T{}

  set := make(map[T]bool)
  for _, item := range array {
    if _, ok := set[item]; !ok {
      result = append(result, item)
      set[item] = true
    }
  }

  return result
}

func fetch(url string) ([]byte, error) {
  response, e := http.Get(url)

  if e != nil {
    return nil, e
  }

  defer response.Body.Close()

  body, e := io.ReadAll(response.Body)

  if e != nil {
    return nil, e
  }

  return body, nil
}

func parallel[T any](callbacks []ParallelCallback[T]) ([]T, []error) {
  results := make([]T, len(callbacks))
  errors := make([]error, len(callbacks))
  channel := make(chan ParallelResult[T])

  for index, callback := range callbacks {
    go func() {
      result, e := callback()

      parallelRresult := ParallelResult[T]{}
      parallelRresult.index = index
      parallelRresult.result = result
      parallelRresult.e = e

      channel <- parallelRresult
    }()
  }

  for i := 0; i < len(callbacks); i++ {
    result := <-channel
    results[result.index] = result.result
    errors[result.index] = result.e
  }

  return results, errors
}

func fetchBlacklist() ([]string, error) {
  debug("Fetching adservers from blacklist...")

  response, e := fetch(BLACKLIST_URL)

  if e != nil {
    return nil, e
  }

  domains := []string{}

  lines := strings.Split(string(response), "\n")
  for _, line := range lines {
    if strings.HasPrefix(line, "#") {
      continue
    }

    fields := strings.Fields(line)

    if len(fields) < 2 || len(fields[1]) == 0 {
      continue
    }

    domains = append(domains, fields[1])
  }

  debug("Fetched", len(domains), "adservers from blacklist")

  return domains, nil
}

func fetchAdguard() ([]string, error) {
  debug("Fetching adservers from adguard...")

  response, e := fetch(ADGUARD_URL)

  if e != nil {
    return nil, e
  }

  domains := []string{}

  lines := strings.Split(string(response), "\n")
  for _, line := range lines {
    if strings.HasPrefix(line, "!") || strings.HasPrefix(line, "#") {
      continue
    }

    line = regexp.MustCompile(`^@+`).ReplaceAllString(line, "")
    line = regexp.MustCompile(`^\|+`).ReplaceAllString(line, "")
    line = regexp.MustCompile(`^-`).ReplaceAllString(line, "")
    line = regexp.MustCompile(`\^\|?$`).ReplaceAllString(line, "")
    line = regexp.MustCompile(`\$important$`).ReplaceAllString(line, "")

    if len(line) == 0 || strings.Contains(line, "*") {
      continue
    }

    domains = append(domains, line)
  }

  debug("Fetched", len(domains), "adservers from adguard")

  return domains, nil
}

func fetchAdservers() ([]string, error) {
  callbacks := []ParallelCallback[[]string]{fetchBlacklist, fetchAdguard}
  responses, errors := parallel(callbacks)

  for _, e := range errors {
    if e != nil {
      return nil, e
    }
  }

  adservers := []string{}

  domains := unique(merge([][]string{CUSTOM_DOMAINS, responses[0], responses[1]}))
  for _, domain := range domains {
    if regexp.MustCompile(`^([a-zA-Z0-9-_]+\.)+[a-zA-Z]{2,}$`).MatchString(domain) {
      adservers = append(adservers, domain)
    }
  }

  debug("Total", len(adservers), "adservers fetched")

  return adservers, nil
}

func main() {
  adservers, e := fetchAdservers()

  if e != nil {
    debug("Error fetching adservers:", e)
    return
  }

  fmt.Print(strings.Join(adservers, "\n"))
}
