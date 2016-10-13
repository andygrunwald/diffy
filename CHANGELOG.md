# Changelog

This is a high level log of changes, bugfixes, enhancements, etc
that have taken place between releases. Later versions are shown
first. For more complete details see
[the releases on GitHub.](https://github.com/andygrunwald/go-gerrit/releases)

## Versions

### 0.1.1 (not yet released)

* Minor fix to SubmitChange to use the `http.StatusConflict` constant
  instead of a hard coded value when comparing response codes.
* Updated AccountInfo.AccountID to be omitted of empty (such as when 
  used in ApprovalInfo).
* + and : in url parameters for queries are no longer escaped. This was
  causing `400 Bad Request` to be returned when the + symbol was
  included as part of the query. To match behavior with Gerrit's search
  handling, the : symbol was also excluded.

### 0.1.0

* The first official release
* Implemented digest auth and several fixes for it.
* Ensured Content-Type is included in all requests
* Fixed several internal bugs as well as a few documentation issues
