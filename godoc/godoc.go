package godoc

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/hhhapz/doc"
)

const (
	// Base is the base path to godocs.io
	Base = "https://godocs.io/"
	// Selectors matches methods, types, and functions on the documentation
	// page.
	Selectors = `[data-kind="function"], [data-kind="type"], [data-kind="method"]:not([class*="decl"])`
)

// GodocParser implements doc.Parser.
type GodocParser struct{}

// Parser is an implementation of godoc.Parser that retrieves documentation
// from https://godocs.io.
var Parser doc.Parser = GodocParser{}

// URL returns a url to the path to see the documentation for the provided
// module on https://godocs.io/.
func (GodocParser) URL(module string) string {
	return Base + module
}

func (p GodocParser) Parse(document *goquery.Document) (doc.Package, error) {
	s := newState(document)

	var err error
	document.Find(Selectors).EachWithBreak(func(_ int, sel *goquery.Selection) bool {
		kind := sel.AttrOr("data-kind", "")
		switch kind {
		case "function":
			err = s.function(sel)
		case "type":
			err = s.typ(sel)
		case "method":
			err = s.method(sel)
		default:
			// this should never happen
		}
		// true when err is nil
		return err == nil
	})
	s.pkg.Types[s.current.Name] = *s.current

	if err != nil {
		return doc.Package{}, err
	}
	return s.pkg, nil
}

var stdlibPackages = map[string]string{
	"tar": "archive/tar",
	"zip": "archive/zip",

	"bzip2": "compress/bzip2",
	"flate": "compress/flate",
	"gzip":  "compress/gzip",
	"lzw":   "compress/lzw",
	"zlib":  "compress/zlib",

	"heap": "container/heap",
	"list": "container/list",
	"ring": "container/ring",

	"aes":      "crypto/aes",
	"cipher":   "crypto/cipher",
	"des":      "crypto/des",
	"dsa":      "crypto/dsa",
	"ecdsa":    "crypto/ecdsa",
	"ed25519":  "crypto/ed25519",
	"elliptic": "crypto/elliptic",
	"hmac":     "crypto/hmac",
	"md5":      "crypto/md5",
	"rc4":      "crypto/rc4",
	"rsa":      "crypto/rsa",
	"sha1":     "crypto/sha1",
	"sha256":   "crypto/sha256",
	"sha512":   "crypto/sha512",
	"subtle":   "crypto/subtle",
	"tls":      "crypto/tls",
	"x509":     "crypto/x509",
	"pkix":     "crypto/x509/pkix",

	"dwarf":    "debug/dwarf",
	"elf":      "debug/elf",
	"gosym":    "debug/gosym",
	"macho":    "debug/macho",
	"pe":       "debug/pe",
	"plan9obj": "debug/plan9obj",

	"ascii85": "encoding/ascii85",
	"asn1":    "encoding/asn1",
	"base32":  "encoding/base32",
	"base64":  "encoding/base64",
	"binary":  "encoding/binary",
	"csv":     "encoding/csv",
	"gob":     "encoding/gob",
	"hex":     "encoding/hex",
	"json":    "encoding/json",
	"pem":     "encoding/pem",
	"xml":     "encoding/xml",

	"ast":           "go/ast",
	"build":         "go/build",
	"constraint":    "go/build/constraint",
	"constant":      "go/constant",
	"docformat":     "go/docformat",
	"importer":      "go/importer",
	"parserprinter": "go/parserprinter",
	"scanner":       "go/scanner",
	"token":         "go/token",
	"types":         "go/types",

	"adler32": "hash/adler32",
	"crc32":   "hash/crc32",
	"crc64":   "hash/crc64",
	"fnv":     "hash/fnv",
	"maphash": "hash/maphash",

	"color":   "image/color",
	"draw":    "image/draw",
	"gif":     "image/gif",
	"jpeg":    "image/jpeg",
	"parsing": "image/parsing",

	"suffixarray": "index/suffixarray",

	"fs":     "io/fs",
	"ioutil": "io/ioutil",

	"big":   "math/big",
	"bits":  "math/bits",
	"cmplx": "math/cmplx",

	"multipart":       "mime/multipart",
	"quotedprintable": "mime/quotedprintable",

	"http":      "net/http",
	"cgi":       "net/http/cgi",
	"cookiejar": "net/http/cookiejar",
	"fcgi":      "net/http/fcgi",
	"httptest":  "net/http/httptest",
	"httptrace": "net/http/httptrace",
	"httputil":  "net/http/httputil",
	"mail":      "net/mail",
	"rpc":       "net/rpc",
	"jsonrpc":   "net/rpc/jsonrpc",
	"smtp":      "net/smtp",
	"textproto": "net/textproto",

	"exec":   "os/exec",
	"signal": "os/signal",
	"user":   "os/user",

	"filepath": "path/filepath",

	"syntax": "regexp/syntax",

	"cgo":     "runtime/cgo",
	"metrics": "runtime/metrics",
	"msan":    "runtime/msan",
	"race":    "runtime/race",
	"trace":   "runtime/trace",

	"js": "syscall/js",

	"fstest": "testing/fstest",
	"iotest": "testing/iotest",
	"quick":  "testing/quick",

	"tabwriter": "text/tabwriter",

	"parse": "text/template/parse",

	"tzdata": "time/tzdata",

	"utf16": "unicode/utf16",
	"utf8":  "unicode/utf8",
}
