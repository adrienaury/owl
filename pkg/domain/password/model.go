package password

// Domain ...
type Domain interface {
	Len() uint
	AsSlice() []rune
	Rune(uint) rune
	MergeWith(Domain) Domain
}

type domain struct {
	chars []rune
}

// NewDomain ...
func NewDomain(chars []rune) Domain {
	return domain{chars}
}

func (d domain) Len() uint                 { return uint(len(d.chars)) }
func (d domain) AsSlice() []rune           { return d.chars }
func (d domain) Rune(idx uint) rune        { return d.chars[idx] }
func (d domain) MergeWith(o Domain) Domain { return NewDomain(append(d.chars, o.AsSlice()...)) }

// Most common domains
var (
	LowerCaseLetters = NewDomain([]rune("abcdefghijklmnopqrstuvwxyz"))
	UpperCaseLetters = NewDomain([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	Numbers          = NewDomain([]rune("0123456789"))
	ASCIISpecials    = NewDomain([]rune(`!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`))

	Letters           = LowerCaseLetters.MergeWith(UpperCaseLetters)
	LettersAndNumbers = Letters.MergeWith(Numbers)

	// LettersAndNumbers will be more common than specials
	Standard = LettersAndNumbers.MergeWith(LettersAndNumbers).MergeWith(ASCIISpecials)
)
