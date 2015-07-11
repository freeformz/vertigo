// This file contains bunch of miscful helper functions.
// The functions here are either too rare to be assiociated to some known file
// or are met more or less everywhere across the code.
package misc

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
	"time"

	. "github.com/9uuso/vertigo/settings"

	"github.com/go-martini/martini"
	"github.com/kennygrant/sanitize"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// NotFound is a shorthand JSON response for HTTP 404 errors.
func NotFound() map[string]interface{} {
	return map[string]interface{}{"error": "Not found"}
}

// TimeOffset returns timezone offset of loc in seconds from UTC.
// Loc should be valid IANA timezone location.
func TimeOffset(loc string) (int, error) {
	var timeOffset int
	l, err := time.LoadLocation(loc)
	if err != nil {
		return timeOffset, err
	}
	now := time.Now().In(l)
	_, timeOffset = now.Zone()
	return timeOffset, nil
}

// Excerpt generates 15 word excerpt from input.
// Used to make shorter summaries from blog posts.
func Excerpt(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	count := 0
	var excerpt bytes.Buffer
	for scanner.Scan() && count < 15 {
		count++
		excerpt.WriteString(scanner.Text() + " ")
	}
	return sanitize.HTML(strings.TrimSpace(excerpt.String()))
}

func Sessionchecker() martini.Handler {
	return func(session sessions.Session) {
		data := session.Get("user")
		_, exists := data.(int64)
		if exists {
			return
		}
		session.Set("user", -1)
		return
	}
}

// sessionIsAlive checks that session cookie with label "user" exists and is valid.
func sessionIsAlive(session sessions.Session) bool {
	data := session.Get("user")
	_, exists := data.(int64)
	if exists {
		return true
	}
	return false
}

// SessionRedirect in addition to sessionIsAlive makes HTTP redirection to user home.
// SessionRedirect is useful for redirecting from pages which are only visible when logged out,
// for example login and register pages.
func SessionRedirect(res http.ResponseWriter, req *http.Request, session sessions.Session) {
	if sessionIsAlive(session) {
		http.Redirect(res, req, "/user", 302)
	}
}

// ProtectedPage makes sure that the user is logged in. Use on pages which need authentication
// or which have to deal with user structure later on.
func ProtectedPage(req *http.Request, session sessions.Session, render render.Render) {
	if !sessionIsAlive(session) {
		session.Delete("user")
		render.JSON(401, map[string]interface{}{"error": "Unauthorized"})
	}
}

// root returns HTTP request "root".
// For example, calling it with http.Request which has URL of /api/user/5348482a2142dfb84ca41085
// would return "api". This function is used to route both JSON API and frontend requests in the same function.
func Root(r *http.Request) string {
	u := strings.TrimPrefix(r.URL.String(), Settings.URL.Path)
	return strings.Split(u[1:], "/")[0]
}

type Page struct {
	Data interface{}
}

type abbr struct {
	STD string
	DST string
	LOC string
}

var Abbrs = map[string]abbr{
	"Egypt Standard Time":             {"EET", "EET", "Africa/Cairo"},
	"Morocco Standard Time":           {"WET", "WEST", "Africa/Casablanca"},
	"South Africa Standard Time":      {"SAST", "SAST", "Africa/Johannesburg"},
	"W. Central Africa Standard Time": {"WAT", "WAT", "Africa/Lagos"},
	"E. Africa Standard Time":         {"EAT", "EAT", "Africa/Nairobi"},
	"Libya Standard Time":             {"EET", "EET", "Africa/Tripoli"},
	"Namibia Standard Time":           {"WAT", "WAST", "Africa/Windhoek"},
	"Alaskan Standard Time":           {"AKST", "AKDT", "America/Anchorage"},
	"Paraguay Standard Time":          {"PYT", "PYST", "America/Asuncion"},
	"Bahia Standard Time":             {"BRT", "BRST", "America/Bahia"},
	"SA Pacific Standard Time":        {"COT", "COT", "America/Bogota"},
	"Argentina Standard Time":         {"ART", "ART", "America/Buenos_Aires"},
	"Venezuela Standard Time":         {"VET", "VET", "America/Caracas"},
	"SA Eastern Standard Time":        {"GFT", "GFT", "America/Cayenne"},
	"Central Standard Time":           {"CST", "CDT", "America/Chicago"},
	"Mountain Standard Time (Mexico)": {"MST", "MDT", "America/Chihuahua"},
	"Central Brazilian Standard Time": {"AMT", "AMST", "America/Cuiaba"},
	"Mountain Standard Time":          {"MST", "MDT", "America/Denver"},
	"Greenland Standard Time":         {"WGT", "WGST", "America/Godthab"},
	"Central America Standard Time":   {"CST", "CST", "America/Guatemala"},
	"Atlantic Standard Time":          {"AST", "ADT", "America/Halifax"},
	"US Eastern Standard Time":        {"EST", "EDT", "America/Indianapolis"},
	"SA Western Standard Time":        {"BOT", "BOT", "America/La_Paz"},
	"Pacific Standard Time":           {"PST", "PDT", "America/Los_Angeles"},
	"Central Standard Time (Mexico)":  {"CST", "CDT", "America/Mexico_City"},
	"Montevideo Standard Time":        {"UYT", "UYST", "America/Montevideo"},
	"Eastern Standard Time":           {"EST", "EDT", "America/New_York"},
	"US Mountain Standard Time":       {"MST", "MST", "America/Phoenix"},
	"Canada Central Standard Time":    {"CST", "CST", "America/Regina"},
	"Pacific Standard Time (Mexico)":  {"PST", "PDT", "America/Santa_Isabel"},
	"Pacific SA Standard Time":        {"CLT", "CLST", "America/Santiago"},
	"E. South America Standard Time":  {"BRT", "BRST", "America/Sao_Paulo"},
	"Newfoundland Standard Time":      {"NST", "NDT", "America/St_Johns"},
	"Central Asia Standard Time":      {"ALMT", "ALMT", "Asia/Almaty"},
	"Jordan Standard Time":            {"EET", "EEST", "Asia/Amman"},
	"Arabic Standard Time":            {"AST", "AST", "Asia/Baghdad"},
	"Azerbaijan Standard Time":        {"AZT", "AZST", "Asia/Baku"},
	"SE Asia Standard Time":           {"ICT", "ICT", "Asia/Bangkok"},
	"Middle East Standard Time":       {"EET", "EEST", "Asia/Beirut"},
	"India Standard Time":             {"IST", "IST", "Asia/Calcutta"},
	"Sri Lanka Standard Time":         {"IST", "IST", "Asia/Colombo"},
	"Syria Standard Time":             {"EET", "EEST", "Asia/Damascus"},
	"Bangladesh Standard Time":        {"BDT", "BDT", "Asia/Dhaka"},
	"Arabian Standard Time":           {"GST", "GST", "Asia/Dubai"},
	"North Asia East Standard Time":   {"IRKT", "IRKT", "Asia/Irkutsk"},
	"Israel Standard Time":            {"IST", "IDT", "Asia/Jerusalem"},
	"Afghanistan Standard Time":       {"AFT", "AFT", "Asia/Kabul"},
	"Pakistan Standard Time":          {"PKT", "PKT", "Asia/Karachi"},
	"Nepal Standard Time":             {"NPT", "NPT", "Asia/Katmandu"},
	"North Asia Standard Time":        {"KRAT", "KRAT", "Asia/Krasnoyarsk"},
	"Magadan Standard Time":           {"MAGT", "MAGT", "Asia/Magadan"},
	"N. Central Asia Standard Time":   {"NOVT", "NOVT", "Asia/Novosibirsk"},
	"Myanmar Standard Time":           {"MMT", "MMT", "Asia/Rangoon"},
	"Arab Standard Time":              {"AST", "AST", "Asia/Riyadh"},
	"Korea Standard Time":             {"KST", "KST", "Asia/Seoul"},
	"China Standard Time":             {"CST", "CST", "Asia/Shanghai"},
	"Singapore Standard Time":         {"SGT", "SGT", "Asia/Singapore"},
	"Taipei Standard Time":            {"CST", "CST", "Asia/Taipei"},
	"West Asia Standard Time":         {"UZT", "UZT", "Asia/Tashkent"},
	"Georgian Standard Time":          {"GET", "GET", "Asia/Tbilisi"},
	"Iran Standard Time":              {"IRST", "IRDT", "Asia/Tehran"},
	"Tokyo Standard Time":             {"JST", "JST", "Asia/Tokyo"},
	"Ulaanbaatar Standard Time":       {"ULAT", "ULAT", "Asia/Ulaanbaatar"},
	"Vladivostok Standard Time":       {"VLAT", "VLAT", "Asia/Vladivostok"},
	"Yakutsk Standard Time":           {"YAKT", "YAKT", "Asia/Yakutsk"},
	"Ekaterinburg Standard Time":      {"YEKT", "YEKT", "Asia/Yekaterinburg"},
	"Caucasus Standard Time":          {"AMT", "AMT", "Asia/Yerevan"},
	"Azores Standard Time":            {"AZOT", "AZOST", "Atlantic/Azores"},
	"Cape Verde Standard Time":        {"CVT", "CVT", "Atlantic/Cape_Verde"},
	"Greenwich Standard Time":         {"GMT", "GMT", "Atlantic/Reykjavik"},
	"Cen. Australia Standard Time":    {"CST", "CST", "Australia/Adelaide"},
	"E. Australia Standard Time":      {"EST", "EST", "Australia/Brisbane"},
	"AUS Central Standard Time":       {"CST", "CST", "Australia/Darwin"},
	"Tasmania Standard Time":          {"EST", "EST", "Australia/Hobart"},
	"W. Australia Standard Time":      {"WST", "WST", "Australia/Perth"},
	"AUS Eastern Standard Time":       {"EST", "EST", "Australia/Sydney"},
	"UTC":                            {"GMT", "GMT", "Etc/GMT"},
	"UTC-11":                         {"GMT+11", "GMT+11", "Etc/GMT+11"},
	"Dateline Standard Time":         {"GMT+12", "GMT+12", "Etc/GMT+12"},
	"UTC-02":                         {"GMT+2", "GMT+2", "Etc/GMT+2"},
	"UTC+12":                         {"GMT-12", "GMT-12", "Etc/GMT-12"},
	"W. Europe Standard Time":        {"CET", "CEST", "Europe/Berlin"},
	"GTB Standard Time":              {"EET", "EEST", "Europe/Bucharest"},
	"Central Europe Standard Time":   {"CET", "CEST", "Europe/Budapest"},
	"Turkey Standard Time":           {"EET", "EEST", "Europe/Istanbul"},
	"Kaliningrad Standard Time":      {"FET", "FET", "Europe/Kaliningrad"},
	"FLE Standard Time":              {"EET", "EEST", "Europe/Kiev"},
	"GMT Standard Time":              {"GMT", "BST", "Europe/London"},
	"Russian Standard Time":          {"MSK", "MSK", "Europe/Moscow"},
	"Romance Standard Time":          {"CET", "CEST", "Europe/Paris"},
	"Central European Standard Time": {"CET", "CEST", "Europe/Warsaw"},
	"Mauritius Standard Time":        {"MUT", "MUT", "Indian/Mauritius"},
	"Samoa Standard Time":            {"WST", "WST", "Pacific/Apia"},
	"New Zealand Standard Time":      {"NZST", "NZDT", "Pacific/Auckland"},
	"Fiji Standard Time":             {"FJT", "FJT", "Pacific/Fiji"},
	"Central Pacific Standard Time":  {"SBT", "SBT", "Pacific/Guadalcanal"},
	"Hawaiian Standard Time":         {"HST", "HST", "Pacific/Honolulu"},
	"Line Islands Standard Time":     {"LINT", "LINT", "Pacific/Kiritimati"},
	"West Pacific Standard Time":     {"PGT", "PGT", "Pacific/Port_Moresby"},
	"Tonga Standard Time":            {"TOT", "TOT", "Pacific/Tongatapu"},
}
