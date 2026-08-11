package articles

import (
	"strings"

	"shanraq.org/internal/config"
)

// applyOperator substitutes the operator-identity tokens in a legal page body
// with the configured values, localized to the page language. Values come from
// config (never from git), so the BIN and registered address are disclosed on
// the live site without being committed to the public repository.
//
// Tokens:
//
//	{{op}}              — the operator's legal name (short, inline).
//	{{operator_block}}  — a full identity sentence: name, BIN, address, contact.
func applyOperator(body string, op config.OperatorConfig, lang string) string {
	body = strings.ReplaceAll(body, "{{op}}", operatorName(op, lang))
	body = strings.ReplaceAll(body, "{{operator_block}}", operatorBlock(op, lang))
	return body
}

func operatorName(op config.OperatorConfig, lang string) string {
	if lang == LangEN && strings.TrimSpace(op.LegalNameEN) != "" {
		return op.LegalNameEN
	}
	if strings.TrimSpace(op.LegalName) != "" {
		return op.LegalName
	}
	switch lang {
	case LangKZ:
		return "Shanraq.org платформасының иесі"
	case LangEN:
		return "the owner of the Shanraq.org platform"
	default:
		return "владелец платформы Shanraq.org"
	}
}

// operatorBlock assembles the operator's disclosure sentence from whatever
// fields are set, so an unconfigured field is simply omitted rather than
// printing an empty label.
func operatorBlock(op config.OperatorConfig, lang string) string {
	var binLabel, contactLabel, defaultEmail string
	switch lang {
	case LangKZ:
		binLabel, contactLabel = "БСН", "Хабарласу үшін:"
	case LangEN:
		binLabel, contactLabel = "BIN", "Contact for enquiries:"
	default:
		binLabel, contactLabel = "БИН", "Для обращений:"
	}
	defaultEmail = "support@shanraq.org"

	parts := []string{operatorName(op, lang)}
	if b := strings.TrimSpace(op.BIN); b != "" {
		parts = append(parts, binLabel+" "+b)
	}
	address := strings.TrimSpace(op.Address)
	if lang == LangEN && strings.TrimSpace(op.AddressEN) != "" {
		address = strings.TrimSpace(op.AddressEN)
	}
	if address != "" {
		parts = append(parts, address)
	}
	out := strings.Join(parts, ", ") + "."

	email := strings.TrimSpace(op.Email)
	if email == "" {
		email = defaultEmail
	}
	contact := contactLabel + " " + email
	if p := strings.TrimSpace(op.Phone); p != "" {
		contact += ", " + p
	}
	return out + " " + contact + "."
}
