package helpers

import "net/mail"

var popularEmailDomains []string = []string{
	"gmail.com", "yahoo.com", "hotmail.com", "aol.com", "hotmail.co.uk", "msn.com", "yahoo.co.uk", "live.com", "ymail.com", "outlook.com", "hotmail.it", "verizon.net", "googlemail.com", "rocketmail.com", "yahoo.ca", "sky.com", "me.com", "mail.com", "live.ca", "aim.com",
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func EmailDomainCheck(emaildn string) bool {
	for _, dn := range popularEmailDomains {
		if dn == emaildn {
			return true
		}
	}
	return false
}
