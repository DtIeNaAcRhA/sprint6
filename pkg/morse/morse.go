package morse

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	А  = ".-"
	Б  = "-..."
	В  = ".--"
	Г  = "--."
	Д  = "-.."
	Е  = "."
	Ж  = "...-"
	З  = "--.."
	И  = ".."
	Й  = ".---"
	К  = "-.-"
	Л  = ".-.."
	М  = "--"
	Н  = "-."
	О  = "---"
	П  = ".--."
	Р  = ".-."
	С  = "..."
	Т  = "-"
	У  = "..-"
	Ф  = "..-."
	Х  = "...."
	Ц  = "-.-."
	Ч  = "---."
	Ш  = "----"
	Щ  = "--.-"
	ЪЬ = "-..-"
	Ы  = "-.--"
	Э  = "..-.."
	Ю  = "..--"
	Я  = ".-.-"

	One   = ".----"
	Two   = "..---"
	Three = "...--"
	Four  = "....-"
	Five  = "....."
	Six   = "-...."
	Seven = "--..."
	Eight = "---.."
	Nine  = "----."
	Zero  = "-----"

	Period       = "......" //.
	Comma        = ".-.-.-" //,
	Colon        = "---..." //:
	QuestionMark = "..--.." //?
	Apostrophe   = ".----." //'
	Hyphen       = "-....-" //-
	Division     = "-..-."  ///
	LeftBracket  = "-.--."  //(
	RightBracket = "-.--.-" //)
	IvertedComma = ".-..-." //“ ”
	DoubleHyphen = "-...-"  //=
	Cross        = ".-.-."  //+
	CommercialAt = ".--.-." //@

	Space = " "
)

type EncodingMap map[rune]string

// averageSize is the average size of a morse char.
const averageSize = 4.53 //Magic

var DefaultMorse = EncodingMap{
	'А': А,
	'Б': Б,
	'В': В,
	'Г': Г,
	'Д': Д,
	'Е': Е,
	'Ж': Ж,
	'З': З,
	'И': И,
	'Й': Й,
	'К': К,
	'Л': Л,
	'М': М,
	'Н': Н,
	'О': О,
	'П': П,
	'Р': Р,
	'С': С,
	'Т': Т,
	'У': У,
	'Ф': Ф,
	'Х': Х,
	'Ц': Ц,
	'Ч': Ч,
	'Ш': Ш,
	'Щ': Щ,
	'Ь': ЪЬ,
	'Ы': Ы,
	'Ъ': ЪЬ,
	'Э': Э,
	'Ю': Ю,
	'Я': Я,

	'1': One,
	'2': Two,
	'3': Three,
	'4': Four,
	'5': Five,
	'6': Six,
	'7': Seven,
	'8': Eight,
	'9': Nine,
	'0': Zero,

	'.':  Period,
	',':  Comma,
	':':  Colon,
	'?':  QuestionMark,
	'\'': Apostrophe,
	'-':  Hyphen,
	'/':  Division,
	'(':  LeftBracket,
	')':  RightBracket,
	'"':  IvertedComma,
}

var reverseDefaultMorse = reverseEncodingMap(DefaultMorse)

// ErrNoEncoding is the error used when there is no representation.
// Its primary use is inside Handlers.
type ErrNoEncoding struct{ Text string }

// Error implements the error interface.
func (e ErrNoEncoding) Error() string { return fmt.Sprintf("No encoding for: %q", e.Text) }

func RuneToMorse(r rune) string {
	r = unicode.ToUpper(r)
	return DefaultMorse[r]
}

func MorseToRune(morse string) rune {
	return reverseDefaultMorse[morse]
}

func reverseEncodingMap(encoding EncodingMap) map[string]rune {
	ret := make(map[string]rune, len(encoding))

	for k, v := range encoding {
		ret[v] = k
	}

	return ret
}

// ToText converts a morse string to his textual representation.
//
// For Example: "- . ... -" -> "TEST".
func (c Converter) ToText(morse string) string {
	out := make([]rune, 0, int(float64(len(morse))/averageSize))

	words := strings.Split(morse, c.charSeparator+Space+c.charSeparator)
	for _, word := range words {
		chars := strings.Split(word, c.charSeparator)

		for _, ch := range chars {
			text, ok := c.morseToRune[ch]
			if !ok {
				hand := []rune(c.Handling(ErrNoEncoding{string(ch)}))
				out = append(out, hand...)

				// Add a charSeparator is the len of the result is not zero
				if len(hand) != 0 {
					out = append(out, []rune(c.charSeparator)...)
				}
				continue
			}
			out = append(out, text)
		}

		out = append(out, ' ')
	}

	// Remove last charSeparator
	if !c.trailingSeparator && len(out) >= len(c.charSeparator) {
		out = out[:len(out)-len(c.charSeparator)]
	}

	return string(out)
}

// ConverterOption is a function that modifies a Converter.
// The main use of ConvertOption is inside NewConverter.
type ConverterOption func(Converter) Converter

// ErrorHandler is a function used by Converter when it encounters an unknown character.
// Returns the text to insert at the place of the unknown character.
// This may not(but can if necessary) corrupt the output inserting invalid morse character.
type ErrorHandler func(error) string

// Converter is a Morse from/to Text converter, it handles the conversion and error handling.
type Converter struct {
	runeToMorse       map[rune]string
	morseToRune       map[string]rune
	charSeparator     string
	wordSeparator     string
	convertToUpper    bool
	trailingSeparator bool

	Handling ErrorHandler
}

// IgnoreHandler ignores the error and returns nothing.
func IgnoreHandler(error) string { return "" }

// NewConverter creates a new converter with the specified configuration
// convertingMap is an EncodingMap, it contains how the characters will be translated, usually this is set to DefaultMorse
// but a custom one can be used. A nil convertingMap will panic.
func NewConverter(convertingMap EncodingMap, options ...ConverterOption) Converter {
	if convertingMap == nil {
		panic("Using a nil EncodingMap")
	}

	morseToRune := reverseEncodingMap(convertingMap)

	c := Converter{
		runeToMorse:       convertingMap,
		morseToRune:       morseToRune,
		charSeparator:     " ",
		wordSeparator:     "",
		convertToUpper:    false,
		trailingSeparator: false,

		Handling: IgnoreHandler,
	}

	for _, opt := range options {
		c = opt(c)
	}

	// Set wordSeparator as default
	if c.wordSeparator == "" {

		// Use custom space if avaible
		sp, ok := c.runeToMorse[' ']
		if !ok {
			// Fallback to the default Space
			sp = Space
		}
		c.wordSeparator = c.charSeparator + sp + c.charSeparator
	}

	return c
}

// ToMorse converts a text to his morse representation.
// Lowercase characters are automatically converted to Uppercase.
//
// For Example: "Test" -> "- . ... -".
func (c Converter) ToMorse(text string) string {
	out := make([]rune, 0, int(float64(len(text))*averageSize))

	for _, ch := range text {
		if c.convertToUpper {
			ch = unicode.ToUpper(ch)
		}

		if _, ok := c.runeToMorse[ch]; !ok {
			hand := []rune(c.Handling(ErrNoEncoding{string(ch)}))
			out = append(out, hand...)

			// Add a charSeparator is the len of the result is not zero
			if len(hand) != 0 {
				out = append(out, []rune(c.charSeparator)...)
			}
			continue
		}

		out = append(out, []rune(c.runeToMorse[ch])...)
		out = append(out, []rune(c.charSeparator)...)
	}

	// Remove last charSeparator
	if !c.trailingSeparator && len(out) >= len(c.charSeparator) {
		out = out[:len(out)-len(c.charSeparator)]
	}

	return string(out)
}

var DefaultConverter = NewConverter(
	DefaultMorse,

	WithCharSeparator(" "),
	WithWordSeparator("   "),
	WithLowercaseHandling(true),
	WithHandler(IgnoreHandler),
	WithTrailingSeparator(false),
)

// ToText converts a morse string to his textual representation, it is an alias to DefaultConverter.ToText.
func ToText(morse string) string { return DefaultConverter.ToText(morse) }

// ToMorse converts a text to his morse rrpresentation, it is an alias to DefaultConverter.ToMorse.
func ToMorse(text string) string { return DefaultConverter.ToMorse(text) }

// WithHandler sets the handler for the Converter.
func WithHandler(handler ErrorHandler) ConverterOption {
	return func(c Converter) Converter {
		c.Handling = handler
		return c
	}
}

// WithLowercaseHandling sets if the Converter may convert to uppercase before checking inside the EncodingMap.
func WithLowercaseHandling(lowercaseHandling bool) ConverterOption {
	return func(c Converter) Converter {
		c.convertToUpper = lowercaseHandling
		return c
	}
}

// WithTrailingSeparator sets if the Converter may trail the charSeparator.
func WithTrailingSeparator(trailingSpace bool) ConverterOption {
	return func(c Converter) Converter {
		c.trailingSeparator = trailingSpace
		return c
	}
}

// WithCharSeparator sets the Character Separator.
// The CharSeparator is the character used to separate two characters inside a Word.
func WithCharSeparator(charSeparator string) ConverterOption {
	return func(c Converter) Converter {
		c.charSeparator = charSeparator
		return c
	}
}

// WithWordSeparator sets the Word Separator.
// The Word Separator is used to separate two words, usually this is the Character Separator, a Space and another Character Separator.
func WithWordSeparator(wordSeparator string) ConverterOption {
	return func(c Converter) Converter {
		c.wordSeparator = wordSeparator
		return c
	}
}
