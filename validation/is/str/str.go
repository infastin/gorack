package isstr

import (
	"net"
	"net/mail"
	"net/netip"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/infastin/gorack/validation"
)

var (
	ErrLowerCase      = validation.NewRuleError("is_lower_case", "must be in lower case")
	ErrUpperCase      = validation.NewRuleError("is_upper_case", "must be in upper case")
	ErrAlpha          = validation.NewRuleError("is_alpha", "must contain English letters only")
	ErrNumeric        = validation.NewRuleError("is_numeric", "must contain digits only")
	ErrAlphanumeric   = validation.NewRuleError("is_alphanumeric", "must contain English letters and digits only")
	ErrASCII          = validation.NewRuleError("is_ascii", "must contain ASCII characters only")
	ErrPrintableASCII = validation.NewRuleError("is_printable_ascii", "must contain printable ASCII characters only")
	ErrEmail          = validation.NewRuleError("is_email", "must be a valid email address")
	ErrURL            = validation.NewRuleError("is_url", "must be a valid URL")
	ErrUUID           = validation.NewRuleError("is_uuid", "must be a valid UUID")
	ErrIP             = validation.NewRuleError("is_ip", "must be a valid IP address")
	ErrIPv4           = validation.NewRuleError("is_ipv4", "must be a valid IPv4 address")
	ErrIPv6           = validation.NewRuleError("is_ipv6", "must be a valid IPv6 address")
	ErrDNSName        = validation.NewRuleError("is_dns_name", "must be a valid DNS name")
	ErrAddrPort       = validation.NewRuleError("is_addr_port", "must be a valid address with port")
	ErrDNSNamePort    = validation.NewRuleError("is_dns_name_port", "must be a valid DNS name with port")
	ErrCIDR           = validation.NewRuleError("is_cidr", "must be a valid CIDR")
	ErrHost           = validation.NewRuleError("is_host", "must be a valid IP address or DNS name")
	ErrPort           = validation.NewRuleError("is_port", "must be a valid port number")
	ErrPath           = validation.NewRuleError("is_path", "must be a valid path")
	ErrFile           = validation.NewRuleError("is_file", "must be a valid path to a file")
	ErrDirectory      = validation.NewRuleError("is_directory", "must be a valid path to a directory")
	ErrCRON           = validation.NewRuleError("is_cron", "must be a valid CRON expression")
)

var (
	rxUUID    = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	rxDNSName = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)
)

func LowerCase[T ~string](v T) error {
	for _, r := range v {
		if !unicode.IsLower(r) {
			return ErrLowerCase
		}
	}
	return nil
}

func UpperCase[T ~string](v T) error {
	for _, r := range v {
		if !unicode.IsUpper(r) {
			return ErrUpperCase
		}
	}
	return nil
}

func Alpha[T ~string](v T) error {
	for _, r := range v {
		if !unicode.IsLetter(r) {
			return ErrAlpha
		}
	}
	return nil
}

func Numeric[T ~string](v T) error {
	for _, r := range v {
		if !unicode.IsNumber(r) {
			return ErrNumeric
		}
	}
	return nil
}

func Alphanumeric[T ~string](v T) error {
	ranges := []*unicode.RangeTable{unicode.Letter, unicode.Number}
	for _, r := range v {
		if !unicode.IsOneOf(ranges, r) {
			return ErrAlphanumeric
		}
	}
	return nil
}

func ASCII[T ~string](v T) error {
	for i := 0; i < len(v); i++ {
		if v[i] >= unicode.MaxASCII {
			return ErrASCII
		}
	}
	return nil
}

func PrintableASCII[T ~string](v T) error {
	for i := 0; i < len(v); i++ {
		if v[i] >= unicode.MaxASCII || !unicode.IsPrint(rune(v[i])) {
			return ErrPrintableASCII
		}
	}
	return nil
}

func Email[T ~string](v T) error {
	if _, err := mail.ParseAddress(string(v)); err != nil {
		return ErrEmail
	}
	return nil
}

func URL[T ~string](v T) error {
	if _, err := url.Parse(string(v)); err != nil {
		return ErrURL
	}
	return nil
}

func UUID[T ~string](v T) error {
	if !rxUUID.MatchString(string(v)) {
		return ErrUUID
	}
	return nil
}

func IP[T ~string](v T) error {
	if _, err := netip.ParseAddr(string(v)); err != nil {
		return ErrIP
	}
	return nil
}

func CIDR[T ~string](v T) error {
	if _, _, err := net.ParseCIDR(string(v)); err != nil {
		return ErrCIDR
	}
	return nil
}

func IPv4[T ~string](v T) error {
	addr, err := netip.ParseAddr(string(v))
	if err != nil || !addr.Is4() {
		return ErrIPv4
	}
	return nil
}

func IPv6[T ~string](v T) error {
	addr, err := netip.ParseAddr(string(v))
	if err != nil || !addr.Is6() {
		return ErrIPv6
	}
	return nil
}

func DNSName[T ~string](v T) error {
	if v == "" || len(strings.ReplaceAll(string(v), ".", "")) > 255 {
		return ErrDNSName
	}
	if IP(v) != nil && !rxDNSName.MatchString(string(v)) {
		return ErrDNSName
	}
	return nil
}

func Host[T ~string](v T) error {
	if IP(v) != nil && DNSName(v) != nil {
		return ErrHost
	}
	return nil
}

func Port[T ~string](v T) error {
	i, err := strconv.Atoi(string(v))
	if err != nil || i <= 0 || i >= 65536 {
		return ErrPort
	}
	return nil
}

func AddrPort[T ~string](v T) error {
	if _, err := netip.ParseAddrPort(string(v)); err != nil {
		return ErrAddrPort
	}
	return nil
}

func DNSNamePort[T ~string](v T) error {
	i := strings.LastIndexByte(string(v), ':')
	if i == -1 {
		return ErrDNSNamePort
	}

	name, port := v[:i], v[i+1:]
	if DNSName(name) != nil || Port(port) != nil {
		return ErrDNSName
	}

	return nil
}

func Path[T ~string](v T) error {
	if _, err := os.Stat(string(v)); err != nil {
		return ErrPath
	}
	return nil
}

func File[T ~string](v T) error {
	if stat, err := os.Stat(string(v)); err != nil || !stat.Mode().IsRegular() {
		return ErrFile
	}
	return nil
}

func Directory[T ~string](v T) error {
	if stat, err := os.Stat(string(v)); err != nil || !stat.Mode().IsDir() {
		return ErrDirectory
	}
	return nil
}

func CRON[T ~string](v T) error {
	if err := cronValid(string(v)); err != nil {
		return ErrCRON
	}
	return nil
}
