package helpers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	kUrlTemplate = "https://etherscan.io/tokens?p=%d"
)

type Token struct {
	Address string
	Url     string
	Name    string
	Token   string
}

func (p *Token) String() string {
	return p.Url + " " + p.Name + " " + p.Token
}

func getContentByUrl(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func containsUrl(tokens []*Token, url string) bool {
	for _, t := range tokens {
		if strings.Contains(t.Url, url) {
			return true
		}
	}
	return false
}

func splitToken(str string) (string, string) {
	str2 := strings.TrimSpace(str)
	strs := strings.Split(str2, "(")
	if len(strs) == 2 {
		str1 := strings.TrimSpace(strs[0])
		str2 := strings.TrimSpace(strings.TrimSuffix(strs[1], ")"))
		return str1, str2
	}

	return str, ""
}

func parseHtmlPage(htmlReader io.ReadCloser) ([]*Token, error) {
	t := html.NewTokenizer(htmlReader)
	if t == nil {
		return nil, errors.New("invalid html content")
	}
	defer htmlReader.Close()
	result := []*Token{}

	for {
		switch t.Next() {
		case html.StartTagToken:
			token := t.Token()
			if token.Data == "a" {
				tempToken := &Token{}
				for _, a := range token.Attr {
					if a.Key == "href" {
						//fmt.Println("Found href:", a.Val)
						if strings.HasPrefix(a.Val, "/token/0") {
							tempToken.Url = "https://etherscan.io" + a.Val
							tempToken.Address = strings.TrimLeft(a.Val, "/token/")
							//result = append(result, a.Val)
							t.Next()
							token2 := t.Token()
							if strings.Contains(token2.String(), "img") {
								break
							} //endif
							if containsUrl(result, tempToken.Url) {
								break
							}
							tempToken.Name, tempToken.Token = splitToken(token2.String())
							result = append(result, tempToken)
						}
						break
					}
				}
			}
		case html.ErrorToken:
			return result, nil
		}
	}
}

var gTokens []*Token

func GetAllERC20Tokens2() ([]*Token, error) {
	if gTokens == nil {
		var e error
		gTokens, e = GetAllERC20Tokens()
		return gTokens, e
	}
	return gTokens, nil
}

func GetAllERC20Tokens() ([]*Token, error) {
	tokens := []*Token{}
	i := 1
	for {
		rc, e := getContentByUrl(fmt.Sprintf(kUrlTemplate, i))
		if e != nil {
			return tokens, e
		}
		tks, e2 := parseHtmlPage(rc)
		if e2 != nil {
			return tokens, e2
		}

		if len(tks) == 0 {
			return tokens, nil
		}
		for _, s := range tks {
			if s.Token != "" {
				tokens = append(tokens, s)
			}
		}

		i = i + 1
		if i > 23 {
			fmt.Println("xxxxxxxxxxxxxxxxxxxxx")
			break
		}
	}

	return tokens, nil
}
