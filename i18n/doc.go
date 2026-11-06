// Package i18n loads application-owned translation catalogs and renders localized
// messages using go-i18n. Catalogs are immutable after startup and independent of
// HTTP, so the same messages can serve views, workers, or command-line programs.
//
// Catalogs support go-i18n JSON, TOML, and YAML files named for their language,
// such as active.en.yaml. Keep whole messages with named variables in catalogs
// so translations can change word order. Message keys used only in templates
// must be maintained in the source catalog; goi18n extract scans Go message
// definitions, not HTML templates. Dates, timezones, and currency policies belong
// to the application and are not inferred from a language preference.
package i18n
