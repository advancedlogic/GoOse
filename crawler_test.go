package goose

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func ReadRawHTML(a Article) string {
	path := fmt.Sprintf("sites/%s.html", a.Domain)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("cannot read %q", path))
	}

	return string(file)
}

func ValidateArticle(expected Article, removed *[]string) error {
	g := New()
	result := g.ExtractFromRawHTML(expected.FinalURL, ReadRawHTML(expected))

	//fmt.Println("%v\n", result.CleanedText) // DEBUG
	//fmt.Println("%v\n", result.Links) // DEBUG

	if result.Title != expected.Title {
		return fmt.Errorf("article title does not match. Got %q", result.Title)
	}

	if result.MetaDescription != expected.MetaDescription {
		return fmt.Errorf("article metaDescription does not match. Got %q", result.MetaDescription)
	}

	if !strings.Contains(result.CleanedText, expected.CleanedText) {
		return fmt.Errorf("article cleanedText does not contains %q", expected.CleanedText)
	}

	// check if the specified strings where properly removed
	for _, rem := range *removed {
		if strings.Contains(result.CleanedText, rem) {
			return fmt.Errorf("article cleanedText does contains %q", rem)
		}
	}

	if result.MetaKeywords != expected.MetaKeywords {
		return fmt.Errorf("article keywords does not match. Got %q", result.MetaKeywords)
	}
	if result.CanonicalLink != expected.CanonicalLink {
		return fmt.Errorf("article CanonicalLink does not match. Got %q", result.CanonicalLink)
	}

	if result.TopImage != expected.TopImage {
		return fmt.Errorf("article topImage does not match. Got %q", result.TopImage)
	}

	if expected.Links != nil && !reflect.DeepEqual(result.Links, expected.Links) {
		return fmt.Errorf("article Links do not match")
	}

	return nil
}

func Test_ExampleParse(t *testing.T) {
	article := Article{
		Domain:          "example.com",
		Title:           "Example HTML Page TITLE",
		MetaDescription: "Example page for testing",
		CleanedText:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n\nexample 1 link content\n\nexample 2 link content\n\nDuis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n\nSed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		MetaKeywords:    "example,testing",
		CanonicalLink:   "http://www.example.com/index.html",
		TopImage:        "/example_top_image.png",
	}
	article.Links = []string{
		"http://www.example.com/page1.html",
		"http://www.example.com/page2.html",
	}

	removed := []string{
		"~HTMLComment~",
		"~div_id_hidden~",
		"~div_class_hidden~",
		"~div_name_hidden~",
		"~style_display_none~",
		"~style_visibility_hidden~",
		"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_GloboEsporteParse(t *testing.T) {
	article := Article{
		Domain:          "globoesporte.globo.com",
		Title:           "Rodrigo Caio treina até nas férias e tenta acelerar retorno aos gramados",
		MetaDescription: "Rodrigo Caio treina na esteira durante as férias em Dracena-SP (Foto: Divulgação)Rodrigo Caio quer ganhar tempo na recuperação da lesão que sofreu no joelho esquerdo. Apesar de ter sido liberado pelo departamento médico do São Paulo para as férias, o ...",
		CleanedText:     "Comissão técnica planeja volta dele para o fim de fevereiro ou início de março",
		MetaKeywords:    "notícias, notícia, presidente prudente região",
		CanonicalLink:   "http://globoesporte.globo.com/sp/presidente-prudente-regiao/noticia/2014/12/rodrigo-caio-treina-ate-nas-ferias-e-tenta-acelerar-retorno-aos-gramados.html",
		TopImage:        "http://s.glbimg.com/es/ge/f/original/2014/12/26/10863872_894379987249341_2406060334390226774_o.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_EditionCnnParse(t *testing.T) {
	article := Article{
		Domain:          "edition.cnn.com",
		Title:           "What if you could make anything you wanted?",
		MetaDescription: "Massimo Banzi's pocket-sized open-source circuit board has become a key building block in the creation of a huge variety of innovative devices.",
		CleanedText:     "In the 20th century, getting your child a toy car meant a trip to a shopping mall.",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.cnn.com/2012/07/08/opinion/banzi-ted-open-source/index.html",
		TopImage:        "http://i2.cdn.turner.com/cnn/dam/assets/120706022111-ted-cnn-ideas-massimo-banzi-00003302-story-top.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_BbcParse(t *testing.T) {
	article := Article{
		Domain:          "bbc.com",
		Title:           "Crunch talks on new Greek bailout",
		MetaDescription: "German and Greek finance ministers meet IMF and Eurogroup chiefs ahead of a crucial finance ministers' meeting on Greece's bailout request.",
		CleanedText:     "Mr Tsipras won elections in late January on a platform of rejecting the austerity measures tied to the bailout.",
		MetaKeywords:    "keywords, added, to, test, case insensitive",
		CanonicalLink:   "http://www.bbc.com/news/business-31545115",
		TopImage:        "http://news.bbcimg.co.uk/media/images/81120000/jpg/_81120901_81120501.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

// multiple og:image, according to http://ogp.me/, the first one should be preferred
func Test_LindorffParse(t *testing.T) {
	article := Article{
		Domain:          "profit.lindorff.fi",
		Title:           "Lindorff24.fi muuttaa maksujen hoidon mobiiliksi - Lindorff Profit",
		MetaDescription: "Lindorffin verkkopalvelu kuluttajille tunnetaan nyt nimellä Lindorff24.fi. Uusien ominaisuuksien lisäksi palvelu on käytettävissä tietokoneen lisäksi älypuhelimella ja tabletilla. Verkon itseasioinnin uskotaan kasvavan lähivuosina merkittävästi nykyisestä.",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "http://profit.lindorff.fi/lindorff24-fi-muuttaa-maksujen-hoidon-mobiiliksi/",
		TopImage:        "http://profit.lindorff.fi/wp-content/uploads/2015/02/Iso_Lindorff24_2_600x2501.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

// Facebook photo
func Test_FacebookParse(t *testing.T) {
	article := Article{
		Domain:          "facebook.com",
		Title:           "Facebook - Facebook's Photos",
		MetaDescription: "Stay connected with all of your groups with the new Facebook Groups app. Learn more: http://www.facebookgroups.com",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "https://www.facebook.com/facebook/photos/a.376995711728.190761.20531316728/10153398878696729/",
		TopImage:        "https://fbcdn-sphotos-g-a.akamaihd.net/hphotos-ak-xpa1/v/t1.0-9/p180x540/10408016_10153398878696729_8237363642999953356_n.png?oh=c6ae71220447f363ec41ea54c38341e1&oe=55B6D827&__gda__=1436749528_5c72e92a5105c1cc6df97163a64e72ce",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

// Relative image test
func Test_RelativeImageWithSpecialChars(t *testing.T) {
	article := Article{
		Domain:          "emeia.ey-vx.com",
		Title:           "Nordics - NO - E - IFRS9 - Bergen - Mai 2015",
		MetaDescription: "",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "https://emeia.ey-vx.com/707/43100/april-2015/nordics---no---e---ifrs9---bergen---mai-2015.asp?sid=51a92e43-8903-43bd-8cfd-8431639dfb5e",
		FinalURL:        "https://emeia.ey-vx.com/707/43100/april-2015/nordics---no---e---ifrs9---bergen---mai-2015.asp?sid=51a92e43-8903-43bd-8cfd-8431639dfb5e",
		TopImage:        "https://emeia.ey-vx.com/707/43100/_images/bergen%201%283%29.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

// Relative image test
func Test_MatchExactDescriptionMetaTag(t *testing.T) {
	article := Article{
		Domain:          "vnexpress.net",
		Title:           "Khánh Ly đến viếng mộ Trịnh Công Sơn - VnExpress Giải Trí",
		MetaDescription: "Chiều 1/5, danh ca mang theo đóa hoa hồng vàng và chai rượu đến thăm người bạn tri kỷ sau lần gặp gỡ cuối cùng vào năm 2000.  - VnExpress Giải Trí",
		CleanedText:     "",
		MetaKeywords:    "Khánh Ly đến viếng mộ Trịnh Công Sơn - VnExpress Giải Trí",
		CanonicalLink:   "http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/khanh-ly-den-vieng-mo-trinh-cong-son-2985539.html",
		FinalURL:        "http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/khanh-ly-den-vieng-mo-trinh-cong-son-2985539.html",
		TopImage:        "http://l.f11.img.vnecdn.net/2014/05/02/2-5456-1398995030_490x294.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}
