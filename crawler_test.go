package goose

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

// ReadRawHTML reads the specified HTML file (article.domain) and return the content
func ReadRawHTML(a Article) string {
	path := fmt.Sprintf("sites/%s.html", a.Domain)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("cannot read %q", path))
	}

	return string(file)
}

// ValidateArticle validates (test) the specified article
func ValidateArticle(expected Article, removed *[]string) error {
	g := New()
	result, err := g.ExtractFromRawHTML(expected.FinalURL, ReadRawHTML(expected))
	if err != nil {
		return err
	}

	// DEBUG
	//fmt.Printf("article := Article{\n\tDomain:          %q,\n\tTitle:           %q,\n\tMetaDescription: %q,\n\tCleanedText:     %q,\n\tMetaKeywords:    %q,\n\tCanonicalLink:   %q,\n\tTopImage:        %q,\n}\n\n", expected.Domain, result.Title, result.MetaDescription, result.CleanedText, result.MetaKeywords, result.CanonicalLink, result.TopImage)
	//fmt.Printf("%#v\n", result.Links)

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

func Test_AbcNewsGoCom(t *testing.T) {
	article := Article{
		Domain:          "abcnews.go.com",
		Title:           "NHL Owner Apologizes for Landing Helicopter at Kids' Soccer Game",
		MetaDescription: "A co-owner of the NHL's New Jersey Devils said today that he's \"truly sorry\" after landing in a helicopter in the middle of a kids' soccer game in Newark.",
		CleanedText:     "A co-owner of the NHL's New Jersey Devils said today that he's \"truly sorry\" after landing in a helicopter in the middle of a kids' soccer game in Newark.\n\nDevils co-owner Joshua Harris said in a statement that he unexpectedly arrived in a chopper in the middle of Saint Benedict Preparatory School's soccer field Sunday night, causing many parents and kids \"frustration\" because the game ended up having to be canceled.\n\n\"I sincerely apologize to the kids and their coaches and families for the cancellation of their soccer game in Newark on Sunday night,\" said Harris, who also owns the NBA's Philadelphia 76ers. \"As a dad, who has spent hundreds of hours watching my kids play sports, I can understand the frustration, and for that, I am truly sorry.\"\n\nHelicopter 'Sounded a Little Funny' Before Crashing Into Florida Home\n\nMan Drives Car Into Ocean to Escape Police During Chase, Helicopter Video Shows\n\nNYPD Chopper Ride-Along: Here’s What Can Happen If You Fly Your Drone Near Aircraft\n\nHarris had been attending a Devils game and was indeed scheduled to land at St. Benedict's soccer field, which is regularly used as a helipad, according to an agreement with the school, a Prudential Center spokesman told ABC station WABC-TV in New York.\n\nBut the problem arose when the Devils game unexpectedly went into overtime and went into the kids' scheduled soccer game.\n\n\"Working with St. Benedict's, we have fixed the process to prevent any future issues,\" Harris said in the statement. \"While I can't take back what happened, I hope the coaches, the teams and their families would be open to being my guest at an upcoming Devils game, and I will be extending an invitation.\"\n\nThe Associated Press contributed to this report.",
		MetaKeywords:    "nj devils owner lands helicopter kids soccer game, helicopter youth soccer game, newark, new jersey, nj nj devils, nhl, josh harris, helicopter cancels soccer game, st benedict preparatory school, sta u13, youth soccer, us news, national news, local news",
		CanonicalLink:   "http://abcnews.go.com/US/nj-devils-owner-apologizes-landing-helicopter-middle-kids/story?id=35155591",
		TopImage:        "http://a.abcnews.go.com/images/US/ht_devils_helicopter_landing_hb_151112_16x9_992.jpg",
	}
	article.Links = []string{
		"http://abcnews.go.com/topics/sports/nhl.htm",
		"http://abcnews.go.com/topics/sports/hockey/new-jersey-devils.htm",
		"http://abcnews.go.com/topics/sports/nba.htm",
		"http://abcnews.go.com/topics/sports/basketball/philadelphia-76ers.htm",
		"http://abcnews.go.com/US/helicopter-sounded-funny-crashing-florida-home/story?id=29836015",
		"http://abcnews.go.com/US/nypd-chopper-ride-heres-happen-fly-drone-aircraft/story?id=33394237",
		"http://abcnews.go.com/US/nypd-chopper-ride-heres-happen-fly-drone-aircraft/story?id=33394237",
		"http://abc7ny.com/sports/devils-co-owners-helicopter-on-newark-field-prompts-cancelation-of-youth-soccer-game/1079546/",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_BbcCom(t *testing.T) {
	article := Article{
		Domain:          "bbc.com",
		Title:           "Crunch talks on new Greek bailout",
		MetaDescription: "German and Greek finance ministers meet IMF and Eurogroup chiefs ahead of a crucial finance ministers' meeting on Greece's bailout request.",
		CleanedText:     "Greek bailout\n\nGreece bailout talks - in 60 secs\n\nEuro's existential threat\n\nNothing left to lose?\n\nWhat we know\n\nThe German and Greek finance ministers are holding talks with IMF and Eurogroup chiefs ahead of a meeting of eurozone finance ministers on Friday.\n\nThe talks are aimed at striking a deal on the request made on Thursday by Greece for a new six-month bailout.\n\nGermany rejected the request despite it being welcomed by the European Commission.\n\nThe existing bailout deal expires at the end of the month and Greece could run out of money without a new accord.\n\nGermany's Wolfgang Schaeuble and Greece's Yanis Varoufakis are meeting in Brussels with IMF managing director Christine Lagarde and Jeroen Dijsselbloem, the Dutch finance minister who heads the Eurogroup.\n\nDuring a break in the talks, Mr Dijsselbloem said the situation was quite complicated: \"I am talking to the main players trying to find a solution. It will take some time, but there is still reason for some optimism, but it is still very difficult. I hope to tell you the outcome in a couple of hours time.\"\n\nThe unscheduled negotiations have delayed the start of the finance ministers' meeting, which was due to commence at 1400 GMT.\n\nArriving for the Eurogroup meeting, Mr Varoufakis said he hoped there would be a deal struck on Friday.\n\n\"The Greek government has not just gone the extra mile, but the extra 10 miles, and now we are expecting our partners not to meet us halfway, but a fifth of the way... Hopefully at the end of this, we come out with some white smoke,\" he said.\n\nMeanwhile, French President Francois Hollande reiterated that Greece belonged in the eurozone and there were no plans for it leaving, following talks in Paris with German Chancellor Angela Merkel.\n\n\"Greece is in the eurozone and it must remain in the eurozone,\" he told a joint news conference with Mrs Merkel.\n\nMrs Merkel said German politicians were \"very much geared towards Greece remaining in the euro\", adding that the Greek people had \"made a lot of sacrifices\" to do so.\n\nHowever, she said there was a need for \"significant improvements in the substance\" of the Greek request ahead of a vote in the German parliament next week.\n\nEarlier on Friday the German government's stance appeared to soften after a spokeswoman for Mrs Merkel said Greece's request for a loan extension from its eurozone partners provided \"a starting point\" for more talks.\n\n\"From the German government's point of view, [the request] is still not sufficient,\" said Christiane Wirtz. But \"it certainly offers a starting point for further talks.\"\n\nOne Greek government official described the phone call as \"constructive\", adding: \"The conversation was held in a positive climate, geared towards finding a mutually beneficial solution for Greece and the eurozone.\"\n\nGermany stands to lose up to €80bn if Greece were to leave the eurozone.\n\nAnalysis: Andrew Walker, economics correspondent\n\nGreece has certainly shifted its position. The letter from the Finance Minister, Yanis Varoufakis, to the Eurogroup asked for a six-month master financial assistance facility agreement.\n\nPayments under that agreement require Greece to comply with the measures set out in another document, the memorandum of understanding.\n\nThat is the hated economic policy programme agreed with the equally hated bailout lenders.\n\nIn the meantime, Mr Varoufakis was offering to refrain from unilateral actions that that would undermine the fiscal targets, economic recovery and financial stability and to ensure any new measures were fully funded.\n\nThose certainly look like concessions to Germany and others.\n\nWhat Berlin doesn't like is the manifest desire of the Greek government to use the proposed extension to revise the programme.\n\nGerman press 'fed up' with Greece\n\nMr Tsipras won elections in late January on a platform of rejecting the austerity measures tied to the bailout.\n\nA Greek government source said on Thursday the Eurogroup had \"just two choices: to accept or reject the Greek request. We will now discover who wants to find a solution, and who does not\".\n\nGreece formally requested a six-month extension to its eurozone loan agreement on Thursday, offering major concessions as it raced to avoid running out of cash within weeks.",
		MetaKeywords:    "keywords, added, to, test, case insensitive",
		CanonicalLink:   "http://www.bbc.com/news/business-31545115",
		TopImage:        "http://news.bbcimg.co.uk/media/images/81120000/jpg/_81120901_81120501.jpg",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_BbcCoUk(t *testing.T) {
	article := Article{
		Domain:          "bbc.co.uk",
		Title:           "Homeopathy 'could be blacklisted'",
		MetaDescription: "Ministers are considering whether homeopathy should be put on an NHS blacklist of banned treatments, the BBC learns.",
		CleanedText:     "Ministers are considering whether homeopathy should be put on a blacklist of treatments GPs in England are banned from prescribing, the BBC has learned.The controversial practice is based on the principle that \"like cures like\", but critics say patients are being given useless sugar pills. The Faculty of Homeopathy said patients supported the therapy. A consultation is expected to take place in 2016.The total NHS bill for homeopathy, including homeopathic hospitals and GP prescriptions, is thought to be about £4m.\n\nHomeopathy is based on the concept that diluting a version of a substance that causes illness has healing properties. So pollen or grass could be used to create a homeopathic hay-fever remedy.One part of the substance is mixed with 99 parts of water or alcohol, and this is repeated six times in a \"6c\" formulation or 30 times in a \"30c\" formulation. The end result is combined with a lactose (sugar) tablet.Homeopaths say the more diluted it is, the greater the effect. Critics say patients are getting nothing but sugar. Common homeopathic treatments are for asthma, ear infections, hay-fever, depression, stress, anxiety, allergy and arthritis.But the NHS itself says: \"There is no good-quality evidence that homeopathy is effective as a treatment for any health condition.\" What do you think about homeopathic treatments? Join our Facebook Q&A on Friday 13th November from 3pm, on the BBC News Facebook page, with the BBC website's health editor, James Gallagher. The Good Thinking Society has been campaigning for homeopathy to be added to the NHS blacklist - known formally as Schedule 1 - of drugs that cannot be prescribed by GPs. Drugs can be blacklisted if there are cheaper alternatives or if the medicine is not effective. After the Good Thinking Society threatened to take their case to the courts, Department of Health legal advisers replied in emails that ministers had \"decided to conduct a consultation\".Officials have now confirmed this will take place in 2016.Simon Singh, the founder of the Good Thinking Society, said: \"Given the finite resources of the NHS, any spending on homeopathy is utterly unjustifiable.\"The money spent on these disproven remedies can be far better spent on treatments that offer real benefits to patients.\"But Dr Helen Beaumont, a GP and the president of the Faculty of Homeopathy, said other drugs such as SSRIs (selective serotonin reuptake inhibitors) for depression would be a better target for saving money, as homeopathic pills had a \"profound effect\" on patients.She told the BBC News website: \"Patient choice is important; homeopathy works, it's widely used by doctors in Europe, and patients who are treated by homeopathy are really convinced of its benefits, as am I.\"The result of the consultation would affect GP prescribing, but not homeopathic hospitals which account for the bulk of the NHS money spent on homeopathy. Estimates suggest GP prescriptions account for about £110,000 per year. And any decision would not affect people buying the treatments over the counter or privately.Health Secretary Jeremy Hunt was criticised for supporting a parliamentary motion on homeopathy, but in an interview last year argued \"when resources are tight we have to follow the evidence\".Minister for Life Sciences, George Freeman, told the BBC: \"With rising health demands, we have a duty to make sure we spend NHS funds on the most effective treatments. \"We are currently considering whether or not homeopathic products should continue to be available through NHS prescriptions.\"We expect to consult on proposals in due course.\"",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.bbc.co.uk/news/health-34744858",
		TopImage:        "http://ichef.bbci.co.uk/news/1024/cpsprodpb/B4FE/production/_86643364_m7410098-homeopathic_pills-spl.jpg",
	}
	article.Links = []string{
		"http://www.britishhomeopathic.org/how-are-homeopathic-medicines-made/",
		"http://www.nhs.uk/Conditions/homeopathy/Pages/Introduction.aspx#when-used",
		"https://www.facebook.com/bbcnews/",
		"http://www.legislation.gov.uk/uksi/2004/629/schedule/1/made",
		"https://www.newscientist.com/article/dn22241-hail-jeremy-hunt-the-new-minister-for-magic/",
		"http://www.lbc.co.uk/watch-jeremy-hunt-live-on-lbc-from-7pm-96835",
		"https://twitter.com/JamesTGallagher",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_BizJournalsCom(t *testing.T) {
	article := Article{
		Domain:          "bizjournals.com",
		Title:           "Activist investor sells off $1 billion worth of Microsoft stock",
		MetaDescription: "ValueAct will still retain a 0.7 percent stake in Microsoft after the sale.",
		CleanedText:     "The San Francisco-based activist investing firm that helped pushed Steve Ballmer out of Microsoft’s top job announced Thursday it will sell some of its shares of the company.\n\nValueAct Capital bought a $2 billion stake in Microsoft 2013 and then gained a spot on the Microsoft company’s board. The firm was part of the group that forced Ballmer into retiring ahead of schedule, and ushering in a new era for the company under Satya Nadella.\n\nThat move has paid off well for shareholders, including Ballmer, who remains the company's largest shareholder. Microsoft (Nasdaq: MSFT) share prices have climbed more than 90 percent since Nadella took office. Microsoft stock is now selling for more than $53 a share, almost up to the highs the company hit ahead of the dot com crisis.\n\nNow, ValueAct will sell about a quarter of its stake or nearly 18.7 million Microsoft shares worth just under $1 billion. The firm will retain a 0.7 percent stake in the company.\n\nValueAct Capital's President G. Mason Morfit said Microsoft shares represent more than 20 percent of the firm’s overall portfolio and will sell some to diversify and buy stock in another company.\n\nMicrosoft's stock price increase is partially why the company has come to represent so much of ValueAct's portfolio. But ValueAct also was a major shareholder in Valeant Pharmaceuticals, whose stock has dropped 70 percent over the last three months, according to Forbes.\n\nMorfit will run for re-election to Microsoft’s board of directors at a shareholder’s meeting later this year. He says Microsoft will remain one of the firm’s top positions.",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.bizjournals.com/seattle/blog/techflash/2015/11/activist-investor-sells-off-1-billion-worth-of.html",
		TopImage:        "http://media.bizj.us/view/img/2167041/mason-morfit*400xx306-307-0-25.jpg",
	}
	article.Links = []string{
		"http://www.bizjournals.com/profiles/company/us/wa/redmond/microsoft_corporation/1087001",
		"http://www.bizjournals.com/profiles/company/us/ca/san_francisco/valueact_capital_partners_lp/13646",
		"http://www.bizjournals.com/seattle/blog/techflash/2014/03/microsoft-adds-activist-investor-to-board.html",
		"http://www.bizjournals.com/seattle/print-edition/2013/08/30/not-just-a-new-ceo-steve-ballmers.html",
		"http://www.bizjournals.com/profiles/company/us/ca/aliso_viejo/valeant_pharmaceuticals_international/20416",
		"http://www.forbes.com/sites/antoinegara/2015/11/12/hedge-fund-valueact-hurt-by-valeant-sells-1-billion-of-surging-microsoft-stock/?utm_campaign=yahootix&partner=yahootix",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_BlogSpotCoUK(t *testing.T) {
	article := Article{
		Domain:          "blogspot.co.uk",
		Title:           "Five ways to grow your business this Small Business Week",
		MetaDescription: "",
		CleanedText:     "Susan Brown, owner of Los Angeles gardening store Potted, recently updated her business listing on Google. Susan says, “Putting your business on Google lets people find you easily. Your directions are right there, your hours are right there, what you sell is right there.”\n\nThanks to her decision, Susan has seen more customers walk through her door: “So many of the customers that come in here find us on Google. As a small business, you want to use every opportunity to help your business grow.”\n\nNational Small Business Week is one of those opportunities. So from May 4-8, instead of three cheers, we’re giving you five—five simple ways to get your small business online and growing.\n\nCelebrating National Small Business Week with Google\n\nA handful of bright ideas and quick-fixes, all five ways are doable in a week or less and will help you throw a digital spotlight on your business all year round.\n\n1. SHOW UP ON GOOGLE\n\nCheck to see how your business shows up on Google. Then, claim your listing so that customers can find the right info about your business on Google Search and Maps. When you claim your listing this week: You could be one of 100 randomly selected businesses to get a 360° virtual tour photoshoot—a $255 value.\n\n2. LEARN FROM PROS & PEERS\n\nGet business advice from experts and colleagues in the Google Small Business Community. They're ready to chat! When you visit or join this week: Share your tips for summertime business success and we'll feature your tip in front of an audience of 400K members.\n\nWith professional email, calendars, and docs that you can access anywhere, Google Apps for Work makes it easy for your team to create and collaborate. When you sign up this week you’ll receive 25% off Google Apps for Work for one year.\n\n4. CLAIM YOUR DOMAIN:\n\nWith a custom domain name and website, Google Domains helps you create a place for your business on the web. When you sign up and purchase a .co, .com or .company domain this week you could be one of 1,500 randomly selected businesses to get reimbursed for the first year of registration.\n\n5. GET ADVICE FROM AN ADVERTISING PRO:\n\nLearn how you can promote your business online and work with a local digital marketing expert to craft a strategy that’s right for your business goals. When you RSVP this week you’ll get help from an expert who knows businesses like yours.\n\nWhile these resources are available year-round, there’s no better time to embark on a digital reboot.\n\nFor more information, visit google.com/smallbusinessweek.\n\nWishing everyone a happy and productive Small Business Week!\n\nPS: To join the conversation, use #5Days5Ways and #SBW15 on G+, Facebook or Twitter.",
		MetaKeywords:    "",
		CanonicalLink:   "http://googlewebmastercentral.blogspot.com/2015/05/five-ways-to-grow-your-business-this.html",
		TopImage:        "http://3.bp.blogspot.com/-6SCcCupadL0/VUnQdhs_98I/AAAAAAAAA7Q/wCdIXm6v9Sg/s72-c/Screen%2BShot%2B2015-05-06%2Bat%2B10.22.08%2BAM.png",
	}
	article.Links = []string{
		"http://gybo.com/resources",
		"https://www.gybo.com/ca/mountain-view/resources#way1",
		"http://gybo.com/resources#way2",
		"http://gybo.com/resources#way3",
		"http://gybo.com/resources#way4",
		"http://gybo.com/resources#way5",
		"https://www.gybo.com/ca/mountain-view/resources",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_BloombergCom(t *testing.T) {
	article := Article{
		Domain:          "bloomberg.com",
		Title:           "U.K. Needs 'Urgent Action' to Keep Banks in London, BBA Says",
		MetaDescription: "British lawmakers need to take “urgent action” to ensure the U.K. maintains its position as the leading global financial center or risk the departure of banks to cities such as Singapore and Hong Kong, according to the British Bankers’ Association.",
		CleanedText:     "British lawmakers need to take “urgent action” to ensure the U.K. maintains its position as the leading global financial center or risk the departure of banks to cities such as Singapore and Hong Kong, according to the British Bankers’ Association.New regulations, taxes and depressed economic activity in Europe have resulted in an 8 percent drop in British banking jobs, with two-thirds of BBA members saying they’ve moved business elsewhere since 2010, the lobby group said in a report Friday. The BBA recommends a softening of the law separating retail operations from investment banking, further tax cuts and a reworking of visa limits to make it easier to hire from abroad.“We have now reached a watershed moment in Britain’s competitiveness as an international banking center” and “many international banks have been moving jobs overseas or deciding not to invest in the U.K.,” BBA Chief Executive Officer Anthony Browne said in the report. “Wholesale banking is an internationally mobile industry and there is a real risk this decline could accelerate.”Chancellor of the Exchequer George Osborne, 44, outlined a “new settlement” for the City of London in a speech in June, pledging to curtail huge fines and amend regulations to “get the balance right.” As memories fade of the 1 trillion pounds ($1.5 trillion) of U.K. taxpayer support given to banks amid the 2008 crisis, this year the government has backed down on some issues after lobbying from the BBA, while HSBC Holdings Plc has said it may leave London.Osborne diluted a levy on U.K. banks and pushed out the regulator’s chief misconduct enforcer, Martin Wheatley, and most recently u-turned on a plan to assume senior bank managers are guilty until proven innocent, which lenders blamed for hindering recruitment of top foreign executives.“We recognize the change of tone in conduct regulation, important developments in the senior managers regime, the proposed reduction in the bank levy, greater certainty over tax for international banks,” the BBA said.Nevertheless, London’s financial sector continues to shrink while its rivals grow, according to the report. Compared with 35,000 jobs losses and a 12 percent fall in U.K. banking assets in the past four years, assets in the U.S. have grown by the same percentage, while in Singapore and Hong Kong they have climbed by 24 percent and 34 percent respectively.European firms are also losing market share to U.S. rivals in wholesale banking, which is the part of banks that cater to large corporates and other financial institutions. From 2010 to 2014, the wholesale market share of the top five European banks fell to 24 percent from 26 percent, whereas the share of the top five U.S. banks has risen to 48 percent from 44 percent, the BBA said.London is also losing market share in lending and initial public offerings, the BBA said. Wholesale banking’s global return-on-equity, a measure of profitability, is expected to fall to an average of 6.5 percent by 2017, about a third of the 18 percent-average between 2000 and 2006, according to the report, co-authored by consulting firm Oliver Wyman.Osborne’s overtures to the industry were counterbalanced by the high cost of ring-fencing -- a law that requires splitting off retail units to protect them from investment banking losses, the BBA said. “Uncertainty arising from the rapidly changing tax regime and European Union referendum are inhibiting business planning and discouraging investment,” according to the report.The BBA’s wishlist includes a demand the Chancellor cut the bank levy faster. Under current plans the tax will be reduced over six years and then limited to domestic balance sheets until 2021. The lobby group also wants an 8 percent surcharge on bank profits to be phased out over time.Financial services is the U.K.’s biggest export industry selling 62 billion pounds abroad every year, and employing more than 405,000 people, the BBA said. Before it's here, it's on the Bloomberg Terminal.",
		MetaKeywords:    "Jobs,Banking,London",
		CanonicalLink:   "http://www.bloomberg.com/news/articles/2015-11-13/u-k-needs-urgent-action-to-keep-banks-in-london-bba-says",
		TopImage:        "http://assets.bwbx.io/images/ifXjLu6rC3Tg/v1/-1x-1.jpg",
	}
	article.Links = []string{
		"http://bloom.bg/dg-ws-core-bcom-a1",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_businessInsiderCom(t *testing.T) {
	article := Article{
		Domain:          "businessinsider.com",
		Title:           "This little change could make Black Friday even more miserable this year",
		MetaDescription: "A change to how retailers process payments could make Americans stand in line longer this Black Friday.",
		CleanedText:     "This little change could make Black Friday even more miserable this year\n\nJust when you thought there couldn't be another way to make Black Friday any more miserable for shoppers and retail employees, the credit-card industry came up with one.\n\nCredit-card companies last month began to mandate new technology that uses chips instead of magnetic stripes. It's a change made for a very good reason: card security.\n\nThe credit-card industry self-imposed October 1 as the deadline for the new card readers, though many consumers had received chip-enabled credit and debit cards — which will still work on the old \"swipe\" card processors — long before that.\n\nThe timing of this wider rollout, however, has retail and payments experts warning that this will slow things down at the checkouts on the November 27 shopping day.\n\n\"Any time you introduce a major change like this, there's going to be confusion,\" said Matt Schulz, senior industry analyst with CreditCards.com. \"There's no question this is going to cause some slowdown on Black Friday.\"\n\nThe change itself is simple: Instead of swiping the card through the magnetic-strip reader, shoppers now have to insert it — chip side up — into a slot on the bottom of the device.\n\nBut here's where the delays come in. People who are unfamiliar with the process will swipe as they always have, then be told it didn't work because they have a new chip-enabled card. Then they must be shown how to insert it, and leave it in, so the payment can be processed.\n\nNow multiply that by thousands, and add in the fact that people have been in line since the crack of dawn, elbowed their way to that bargain bin, and then had to wait again just to get to the register, and you can see why even a small delay will test patience. It's called the EMV chip, and it just might wreak havoc on holiday shopping.\n\nA 'rude awakening' for retailers\n\n\"There is going to be a rude awakening\" for retailers, said Jared Drieling, business intelligence manager for The Strawhecker Group, an Omaha, Nebraska-based advisory firm focusing on payments. \"The industry is still bickering over how long an EMV transaction takes.\"\n\nAs many as 47% of US merchants will have new technology in place by the end of 2015, according to a survey conducted earlier this year by the Payments Security Task Force, an industry-backed group of financial services firms and leading retailers. Already, 40% of Americans have been issued new chip-enabled cards.\n\nOf course the nightmare scenario that Drieling is warning about is dependent on a lot of factors. Some customers have been using the chip technology for weeks, and some retailers don't have the readers yet. There is a wide disparity in how individual retailers have gotten ready for the switch.\n\nBest Buy, Macy's, and Walmart stores have been fully outfitted with new card readers, representatives for those companies said. Macy's and Walmart have also reissued store-branded credit cards with new EMV chips embedded in them. Sears, on the other hand, says it is \"continuously working to further enhance the security of our systems,\" according to a spokesman — but declined to provide specifics for Black Friday.\n\nJ. Craig Shearman, a spokesman for the National Retail Federation, said the new card readers would be at \"most major retailers and large national chains.\" The progress of smaller shops\u00a0in\u00a0adapting the chips is not as clear, but those shops\u00a0are less likely to\u00a0be open the day after Thanksgiving anyway.\n\nShearman didn't argue with the notion that things could slow down, but he said it was not clear how much longer it would take to process each transaction.\n\nFor retailers, Black Friday and the ensuing weekend is crucial to performance. Americans packed malls and stores last year after Thanksgiving, driving more than $50 billion in revenue to retailers, the National Retail Federation reported in 2014.\n\nOf course, there are lots of ways to avoid even having to find out. Stay home. Turkey and stuffing is better on day two anyway.\n\nNOW WATCH: JAMES ALTUCHER: 'Warren Buffett is a f-----g liar'\n\nLatest Deals\n\nLatest Research\n\nRead Business Insider On The Go\n\nAvailable for iPhone, iPad, and Android\n\nFind A Job\n\nTech Jobs\n\nC-Level Jobs\n\nMedia Jobs\n\nDesign Jobs\n\nFinance Jobs\n\nSales Jobs\n\nSee All Jobs »\n\nThanks to our partners",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.businessinsider.com/credit-card-chips-could-slow-black-friday-lines-2015-11",
		TopImage:        "http://static5.businessinsider.com/image/56410a64bd86ef18008c8901/this-little-change-could-make-black-friday-even-more-miserable-this-year.jpg",
	}
	article.Links = []string{
		"http://www.businesswire.com/news/home/20150504005631/en/Issuers-Forecast-U.S.-Shift-Chip-Cards-Complete",
		"http://www.usatoday.com/story/money/business/2015/10/01/chip-credit-debit-card-readers-october-1/73140516/",
		"http://www.businessinsider.com/james-altucher-warren-buffett-rant-holding-period-2015-10",
		"http://www.amazon.com/b/ref=as_li_ss_tl?_encoding=UTF8&adid=0KADBBMAY2FXXYKHFED8&camp=1789&creative=390957&linkCode=ur2&node=7762829011&tag=bi_rightrail_promo_thanksgiving_store-20&linkId=4L4I3IMGNBYS6OZ2",
		"http://www.amazon.com/Black-Friday/b/ref=as_li_ss_tl?_encoding=UTF8&camp=1789&creative=390957&linkCode=ur2&node=384082011&tag=bi_rightrail_promo_countdown_to_black_friday-20&linkId=JWBPWCJVHG2FALYN",
		"http://www.amazon.com/b/ref=as_li_ss_tl?_encoding=UTF8&camp=1789&creative=390957&linkCode=ur2&linkId=ALE2HQMFFTHY7RQQ&node=10161501011&tag=bi_rightrail_promo_trendsetter_gift_guide-20&linkId=O2DSMRD62IEXW36C",
		"http://www.amazon.com/b/ref=as_li_ss_tl?_encoding=UTF8&camp=1789&creative=390957&linkCode=ur2&node=12745394011&pf_rd_i=10044414011&pf_rd_m=ATVPDKIKX0DER&pf_rd_p=2259643402&pf_rd_r=092K3C4J2XEG5K6B9GPA&pf_rd_s=merchandised-search-top-1&pf_rd_t=101&tag=bi_rightrail_promo_holiday_central-20&linkId=F4ME27F5X6VUWHCF",
		"http://www.businessinsider.com/intelligence/mobile-payments-free-report?IR=T&utm_source=businessinsider&utm_medium=banner&utm_term=RR&utm_content=Mobile_Free_Report&utm_campaign=Mobile_Free_Report_RR",
		"http://www.businessinsider.com/intelligence/research-store?IR=T&utm_source=House&utm_term=RR-AffiliateMark&utm_campaign=RR#!/THE-AFFILIATE-MARKETING-REPORT/p/56467038/category=11987293",
		"http://www.businessinsider.com/intelligence/research-store?IR=T&utm_source=House&utm_term=RR-IoT2015&utm_campaign=RR#!/The-Internet-of-Things-Report/p/46301489/category=11987294",
		"http://www.businessinsider.com/intelligence/research-store?IR=T&utm_source=House&utm_term=RR-AheadoftheCurve&utm_campaign=RR#!/Ahead-of-the-Curve-The-Digital-Disruption-of-Retail-Banking/p/55853103/category=11987295",
		"http://www.businessinsider.com/about/mobile",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=technology&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin020&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-tech-applys",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=c-level&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin023&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-c-level-applys",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=media&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin021&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-media-applys",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=design&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin024&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-design-applys",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=finance&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin022&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-finance-applys",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=sales&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin025&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-sales-applys",
		"http://www.careerbuilder.com/Jobseeker/Jobs/JobResults.aspx?IPath=QH&qb=1&s_rawwords=&s_freeloc=&s_jobtypes=ALL&lr=cbpar_busiin&siteid=cbpar_busiin026&utm_source=businessinsider&utm_medium=partner&utm_campaign=businessinsider-all-jobs-applys",
		"http://www.catchpoint.com/",
		"http://www.sailthru.com",
		"http://www.ooyala.com/?utm_source=BusinessInsider&utm_medium=Sponsor&utm_campaign=Rebranding",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_CnnCom(t *testing.T) {
	article := Article{
		Domain:          "cnn.com",
		Title:           "Exhausted F1 star Lewis Hamilton crashes car",
		MetaDescription: "After running away with the Formula One world championship, Lewis Hamilton explains he has run out of gas after crashing his car in Monaco.",
		CleanedText:     "Story highlightsLewis Hamilton reveals Monaco car accident on eve of Brazilian GPF1 world champion says he was exhausted and had a feverHamilton organized surprise party for his Mum after Mexico GPThe Mercedes driver revealed he crashed his car in Monaco after \"heavy partying\" last weekend. He turned up for this weekend's Brazilian Grand Prix a day late after taking time off to recover.\"I've not been well with a fever but I also had a road accident in Monaco on Monday night,\" Hamilton explained on his Instagram account.\"Nobody was hurt, which is the most important thing. I made very light contact with a stationary vehicle.\"Talking with the team and my doctor, we decided together that it was best for me to rest at home and leave a day later.\" Dear TeamLH, just wanted to let you know why things have been quiet on social media the past few days. I've not been well with a fever but I also had a road accident in Monaco on Monday night. Whilst ultimately, it is nobody's business, there are people knowing my position that will try to take advantage of the situation and make a quick buck. NO problem. Nobody was hurt, which is the most important thing. But the car was obviously damaged and I made very light contact with a stationary vehicle. Talking with the team and my doctor, we decided together that it was best for me to rest at home and leave a day later. But i am feeling better and am currently boarding the plane to Brazil. However, I am informing you because I feel we all must take responsibility for our actions. Mistakes happen to us all but what's important is that we learn from them and grow. Can't wait for the weekend Brazil🙌🏾 Bless Lewis A photo posted by Lewis Hamilton (@lewishamilton) on Nov 11, 2015 at 2:50pm PST\n\nHamilton posted the news to his fans, who he refers to as \"Team LH,\" but he also added: \"Ultimately, it is nobody's business, there are people knowing my position that will try to take advantage of the situation and make a quick buck.\"After arriving in Sao Paulo for the penultimate race of the 2015 season, the three-time world champion inevitably faced questions from the assembled media. JUST WATCHEDReplayMore Videos ...MUST WATCH Both Hamilton and his Mercedes teammate Nico Rosberg always speak to reporters on the Thursday before a race weekend, while the British driver also has obligations with the UK press.Hamilton explained that his busy schedule since the last race in Mexico 12 days ago had included throwing a surprise 60th birthday party for his mother Carmen in London last Sunday, the night before his Monaco prang.\"\"It was a result of heavy partying and not much rest for 10 days. I am a bit run down,\" Hamilton, who spent four more days in Mexico after the race, said in his BBC Sport column.\"When I got back to the UK, I was trying to organize my Mum's 60th birthday. The party turned out great but by the end of it I was exhausted. I had been busy for two solid weeks and I basically collapsed.\"JUST WATCHEDReplayMore Videos ...MUST WATCH Although an element of mystery still surrounds Hamilton's Monaco car crash, it's not the first time the 30-year-old has been involved in driving drama off the track.At the 2010 Australian Grand Prix, Hamilton was fined for dangerous driving after deliberately spinning his wheels and skidding on his way out of the Albert Park circuit. In 2007, when he was an F1 rookie, his car was impounded in France after he was caught speeding.Hamilton, who wrapped up the 2015 world title at the U.S. Grand Prix in Austin, Texas with three races to spare, is now focused on getting back to business in Brazil.\"I feel good, I'm on an up slope, so a lot closer to 100%\" Hamilton told reporters at the Interlagos track. \"I'm excited to be here. I'm definitely cherishing the moments I'm in the car.\"Tell us what you think of Hamilton's crash on CNN Sport's Facebook page",
		MetaKeywords:    "f1, lewis hamilton, brazilian grand prix, monaco, mercedes, motorsport, Exhausted F1 star Lewis Hamilton crashes car - CNN.com",
		CanonicalLink:   "http://edition.cnn.com/2015/11/13/motorsport/formula-one-lewis-hamilton-crashes-car-news/index.html",
		TopImage:        "http://i2.cdn.turner.com/cnnnext/dam/assets/151113115049-lewis-hamilon-media-brazil-large-169.jpg",
	}
	article.Links = []string{
		"https://instagram.com/p/99kB_8L00w/",
		"https://instagram.com/p/99kB_8L00w/",
		"http://www.bbc.co.uk/sport/features/34783569",
		"http://edition.cnn.com/2010/SPORT/motorsport/08/24/motorsport.f1.hamilton.fine.melbourne/",
		"http://edition.cnn.com/2015/10/25/motorsport/motorsport-usgp-hamilton-vettel-rosberg/",
		"https://www.facebook.com/cnnsport",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_DailyMailCoUk(t *testing.T) {
	article := Article{
		Domain:          "dailymail.co.uk",
		Title:           "High street stores start charging customers for paper bags",
		MetaDescription: "Major high street stores including Debenhams and House of Fraser have started charging up to 10p for paper carrier bags – despite them being exempt from the new laws brought in last month.",
		CleanedText:     "Major high street stores have been accused of ripping off shoppers by charging up to 10p for paper carrier bags – despite them being exempt from the new laws brought in last month.Outraged shoppers have hit out at Debenhams and House of Fraser claiming they are 'cashing in' by charging for paper bags when other high street shops offer them for free.House of Fraser has said the charge for paper bags had been introduced for 'ethical and moral' reasons, and that all proceeds would be donated to charity.However, shoppers have taken to Twitter to express their anger at the charge.\n\nHouse of Fraser has said the paper bag charge has been brought in at stores for 'ethical and moral' reasons\n\nPaper bags are being handed out to shoppers at London branch of Tesco weeks after 5p charge introducedTwitter user Jimmy said: 'Absolutely disgusted! Just spent £180 on shoes and you have the audacity to make me pay 5p for a 'cardboard' bag #shocking'Anthony Bongos added: 'I can't understand why you are charging for paper carrier bags. This isn't the law, is it you cashing in on the law?'A spokesperson for House of Fraser said: 'We have made the ethical and moral decision to support the introduction of a 5p charge on all plastic and paper bags.'Shoppers in Debenhams have also reported being charged to paper bags, with some saying they have been made to pay up to 10p.Suzanne Foley said: '£162 for a suit no suit bags and then get charged 10p for a large bag, what's that all about debenhams!' (sic)And Martena David added: '£162 for a suit no suit bags and then get charged 10p for a large bag, what's that all about debenhams!' (sic)Elsewhere, some Tesco stores have started giving customers free paper bags just weeks after the 5p charge for plastic bags caused chaos around the country.\n\nTwitter uses have expressed their outrage after being made to pay for paper bags at House of Fraser\n\nThe rules are being rolled out by the Government's Department for Environment, Food & Rural Affairs. It claims the change will save £60m in litter clean-up costs and £13m in carbon savings.The levy for supermarkets and big shops employing more than 250 staff will raise more than £70m a year for 'good causes'. Shops can also take a 'reasonable costs' cut. The Government will pocket the VAT raising an estimated £19m a year.Yes. If you have bought food such as fish, uncooked meat or prescription medicines then the retailer should still offer bags for nothing.But problems occur if you buy anything else at the same time. For example, if the bag shares space with a packet of cornflakes it will cost you 5p. You should not be charged if a shop uses paper bags.\n\nA London store has been handing out recyclable small bags as an alternative to shoppers just picking up a handful of groceries.The bags feature the phrase 'love food hate waste'.The new law does not prevent shops from handing out free paper bags, a source from the Department for Food, Environment and Rural Affairs told the Evening Standard.'The key thing is encouraging people to reuse bags,' they said.'The best thing to do is to have a plastic bag in your pocket.'But clearly paper bags can be recycled and do degrade better than plastic bags, and they won't end up strangling a turtle.'England was the last place in the UK to introduce the 5p bag charge.Some supermarkets around the UK where forced to put security tags on baskets and trolleys after shoppers began taking them home to carry their groceries.MailOnline has contacted Debenhams and House of Fraser for comment.\n\nDebenhams has been accused of ripping off customers across the UK by charging up to 10p for paper bags\n\nMOST WATCHED NEWS VIDEOS\n\nPrevious\n\n1\n\n2\n\n3\n\nNext\n\nMOST READ NEWS\n\nPrevious\n\nNext\n\n●\n\n●\n\n●",
		MetaKeywords:    "Debenhams,House,Fraser,charge,PAPER,BAGS,Tesco,started,giving,away,free",
		CanonicalLink:   "http://www.dailymail.co.uk/news/article-3316789/Debenhams-House-Fraser-charge-PAPER-BAGS-Tesco-started-giving-away-free.html",
		TopImage:        "http://i.dailymail.co.uk/i/pix/2015/11/13/10/2E6847FA00000578-0-image-a-9_1447409694956.jpg",
	}
	article.Links = []string{
		"http://www.standard.co.uk/news/uk/tesco-is-giving-out-paper-bags-to-dodge-the-5p-carrier-bag-charge-a3112131.html",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

// Relative image test
func Test_MatchExactDescriptionMetaTag(t *testing.T) {
	article := Article{
		Domain:          "vnexpress.net",
		Title:           "Khánh Ly đến viếng mộ Trịnh Công Sơn",
		MetaDescription: "Chiều 1/5, danh ca mang theo đóa hoa hồng vàng và chai rượu đến thăm người bạn tri kỷ sau lần gặp gỡ cuối cùng vào năm 2000.  - VnExpress Giải Trí",
		CleanedText:     "",
		MetaKeywords:    "Khánh Ly đến viếng mộ Trịnh Công Sơn - VnExpress Giải Trí",
		CanonicalLink:   "http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/khanh-ly-den-vieng-mo-trinh-cong-son-2985539.html",
		FinalURL:        "http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/khanh-ly-den-vieng-mo-trinh-cong-son-2985539.html",
		TopImage:        "http://l.f11.img.vnecdn.net/2014/05/02/2-5456-1398995030_490x294.jpg",
	}
	article.Links = []string{
		"http://giaitri.vnexpress.net/photo/trong-nuoc/ngoc-diem-khoe-con-gai-5-tuoi-3294807.html",
		"http://giaitri.vnexpress.net/photo/trong-nuoc/con-trai-truong-quynh-anh-do-danh-con-gai-xuan-lan-3294397.html",
		"http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/huong-ly-toi-khong-ngac-nhien-khi-chien-thang-next-top-3294195.html",
		"http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/dam-vinh-hung-hat-o-le-cuoi-cua-40-doi-vo-chong-khuyet-tat-3294598.html",
		"http://giaitri.vnexpress.net/photo/trong-nuoc/vo-chong-tuan-hung-du-dam-cuoi-vu-duy-khanh-3294280.html",
		"http://giaitri.vnexpress.net/photo/trong-nuoc/diem-my-9x-khoe-hinh-the-khi-tap-vo-3293926.html",
		"http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/luong-viet-quang-toi-that-bai-vi-qua-tu-tin-vao-giong-hat-3292227.html",
		"http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/sao-viet-buc-xuc-vi-bi-su-dung-hinh-anh-trai-phep-3293246.html",
		"http://giaitri.vnexpress.net/photo/trong-nuoc/ha-tran-om-con-nhun-nhay-theo-nhac-duoi-mua-3293618.html",
		"http://giaitri.vnexpress.net/tin-tuc/gioi-sao/trong-nuoc/cuoc-song-sau-bao-benh-cua-chu-van-quenh-3291824.html",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_EditionCnnCom(t *testing.T) {
	article := Article{
		Domain:          "edition.cnn.com",
		Title:           "What if you could make anything you wanted?",
		MetaDescription: "Massimo Banzi's pocket-sized open-source circuit board has become a key building block in the creation of a huge variety of innovative devices.",
		CleanedText:     "In the 20th century, getting your child a toy car meant a trip to a shopping mall.",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.cnn.com/2012/07/08/opinion/banzi-ted-open-source/index.html",
		TopImage:        "http://i2.cdn.turner.com/cnn/dam/assets/120706022111-ted-cnn-ideas-massimo-banzi-00003302-story-top.jpg",
	}
	article.Links = []string{
		"http://blog.ted.com/2012/06/26/open-source-your-projects-and-upload-them-to-space-massimo-banzi-at-tedglobal-2012/",
		"http://www.cnn.com/video/#/video/us/2012/07/06/ted-massimo-banzi-arduino.ted",
		"http://gizmodo.com/5822319/a-chilean-teen-tweets-about-earthquakes-better-than-his-whole-government",
		"http://mattrichardson.com/blog/2011/08/17/the-enough-already/",
		"http://www.botanicalls.com/",
		"http://code.google.com/p/arducopter/wiki/ArduCopter",
		"http://www.ted.com/talks/boaz_almog_levitates_a_superconductor.html ",
		"http://www.ted.com",
		"http://dontapscott.com/",
		"http://www.ted.com/talks/don_tapscott_four_principles_for_the_open_world_1.html",
		"http://www.youtube.com/watch?v=yNAGkSbt1xI",
		"http://genspace.org/person/Ellen%20D./Jorgensen,%20Ph.D.",
		"http://www.marcgoodman.net/",
		"http://www.nyls.edu/faculty/faculty_profiles/beth_simone_noveck",
		"http://itp.tisch.nyu.edu/object/ShirkyC.html",
		"http://www.ted.com/talks/clay_shirky_how_cognitive_surplus_will_change_the_world.html",
		"http://edition.cnn.com/2012/06/15/world/europe/uk-school-dinner-blog/index.html",
		"http://www.twitter.com/CNNOpinion",
		"http://www.facebook.com/CNNOpinion",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_EntrepreneurCom(t *testing.T) {
	article := Article{
		Domain:          "entrepreneur.com",
		Title:           "6 Thoughts on Why Facing Your Fears Could Help You Achieve Massive Success",
		MetaDescription: "Fire-breathing dragons are a good reason to be afraid. Business fears, maybe not so much.",
		CleanedText:     "Everyone has fears. They’re important, and they’ve helped keep us alive throughout our evolution. Think about the fears\u00a0characters understandably\u00a0feel at certain points in\u00a0Game of Thrones, the hugely successful HBO dramatic series which\u00a0combines elements of medieval times and fantasy. We're talking outrageously murderous kings here, plus scheming\u00a0lords and ladies. Large men with even larger swords. Even fire-breathing dragons.Related:\u00a0Why Fear Is the Entrepreneur's Best FriendIn Season One of GOT, a great line\u00a0illustrates the point about fears perfectly. The speaker is\u00a0Robb Stark, eldest son of the lord of Winterfell and generally a good guy, who\u00a0decides to declare war and march south to Kings Landing, the capital of the Seven (usually warring) Kingdoms and home to\u00a0a lot more of those men with swords . Theon\u00a0Greyjoy, the son of another royal house,\u00a0asks Stark if he’s afraid. And Stark, his hands trembling, replies,\u00a0“I guess I must be.” To which\u00a0Greyjoy’s response is perfect:\u00a0“Good, that means you’re not stupid.”It certainly was appropriate for the denizens of GOT's medieval era to be afraid, but does the same apply to you? For, while fear was an important factor in our hereditary past, in our modern day and age, our fears today\u00a0are often based more in psychology\u00a0than\u00a0actual physical threats. Drawing on some of the books I've enjoyed, I offer\u00a0six thoughts on why facing your fears will assist you in creating massive success.I've had a lot of worries in my life, most of which never happened.\u00a0- -Mark TwainWhen you take the time to actually define your fears, you\u00a0learn to separate fact from fiction. This is an important distinction. Some things you’re afraid of will be valid, but many will be mental worst-case scenarios that have simply spiraled further in your mind than they ever will or would in reality.What about the fears on your list that you’ve defined that are actually valid, like losing a client or\u00a0employee, gettng backlash from a layoff\u00a0or encountering some other tangible fear?\u00a0Easy. When you face fears that have merit -- now that you’ve defined them --\u00a0you can come up with an action plan of responses to mitigate the damages.Think of this list as your \"fear emergency\u00a0plan.\" You know what you’d do in the case of a fire or earthquake, so why not enact a plan of appropriate responses you could take against some of your more valid business\u00a0fears?Related:\u00a07 Ways to Think Differently About FearBran thought about it. \"Can a man still be brave if he's afraid?\" \"That is the only time a man can be brave,\"\u00a0his father told him.” -- George R.R. Martin, series author, A Song of Ice and Fire, on which HBO's GOT series is based.Perhaps I’m just missing Game of Thrones in the offseason, but this quote really struck me and is an important facet of facing your fears. You don’t develop bravery and courage in the good times, you develop them when you actually confront fears. If you were once afraid of starting your own business, but did it anyway, you know the terror, but also the reward, that comes from facing fears head on. Your courage grows with each fear you face.There is wisdom that comes from the experience of working through fears. Some of your fears may have even come true. If you are a business owner and have seen your business falter or fail, perhaps you’ve already lived through adversity. The silver lining of these experiences is that you learn from them. Wisdom comes from all of life’s experiences, but the fearful or bad ones in particular teach us great lessons. Wisdom is always the by-product of facing your fears, and that’s an important quality to develop.Dealing with fears helps your develop compassion. When you yourself have been afraid,\u00a0you’re more likely to have patience and feel compassion toward others experiencing similar situations. After all, we all want a good life. When you push hard for what you want, and experience the joys and failures of success, you learn compassion you can use to help others push through their early fears.You can put yourself in\u00a0the shoes of someone who is just starting out, and that empathy can help guide that person to have deeper courage.“Life doesn't get easier or more forgiving;\u00a0we get stronger and more resilient.” -- Steve Maraboli, Life, the Truth, and Being FreeResilience comes from facing your fears. You become better than your surroundings and transform yourself above the fear and into bigger and bigger success. Resiliience starts with you, and it begins in your mind. Face your fears and learn to rise to face whatever is in front of you.Related:\u00a0What Companies Can Learn From 'Game of Thrones' When Hiring Their Next Chief Information Officer",
		MetaKeywords:    "Growth Strategies,Fear,Success Stories,Courage",
		CanonicalLink:   "http://www.entrepreneur.com/article/252739",
		TopImage:        "https://assets.entrepreneur.com/content/3x2/822/20151112203147-fire-breathing-dragon.jpeg",
	}
	article.Links = []string{
		"http://www.entrepreneur.com/article/239581",
		"https://www.youtube.com/watch?v=fNxvFgysbvU",
		"http://www.entrepreneur.com/article/244277",
		"http://www.entrepreneur.com/article/247456",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_ExampleCom(t *testing.T) {
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

// Facebook photo
func Test_FacebookCom(t *testing.T) {
	article := Article{
		Domain:          "facebook.com",
		Title:           "Facebook - Facebook's Photos",
		MetaDescription: "Stay connected with all of your groups with the new Facebook Groups app. Learn more: http://www.facebookgroups.com",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "https://www.facebook.com/facebook/photos/a.376995711728.190761.20531316728/10153398878696729/",
		TopImage:        "https://fbcdn-sphotos-g-a.akamaihd.net/hphotos-ak-xpa1/v/t1.0-9/p180x540/10408016_10153398878696729_8237363642999953356_n.png?oh=c6ae71220447f363ec41ea54c38341e1&oe=55B6D827&__gda__=1436749528_5c72e92a5105c1cc6df97163a64e72ce",
	}
	article.Links = []string{
		"https://www.facebook.com/facebook?fref=photo",
		"http://l.facebook.com/l.php?u=http%3A%2F%2Fwww.facebookgroups.com%2F&h=gAQEbndf0&enc=AZNwbqa7wrhRCkIAQcDAt9ivI6lNENnpagDgNd4WzF4di3sKJDzKaxBVXeEChFPdrgWkyEHV0H7Kj9a3Y2PWgHbuGr2k_yamwC5KvANw_2Mq5X8ySXJaGXXj22haJvHJhrw-5IFcBmFwRJnUG1t9DHx9&s=1",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_ForbesCom(t *testing.T) {
	article := Article{
		Domain:          "forbes.com",
		Title:           "The World's Most Expensive Passports [Infographic]",
		MetaDescription: "Passports are valuable and expensive items, with the price of applying for one varying tremendously by nationality. The U.S. passport may seem expensive with a $110 application fee and a $25 acceptance fee adding up to $135 in total. According to a report by Go Euro, however, American travellers actually [...]",
		CleanedText:     "",
		MetaKeywords:    "Lifestyle,Lists,On The Move,Travel",
		CanonicalLink:   "http://www.forbes.com/sites/niallmccarthy/2015/11/13/the-worlds-most-expensive-passports-infographic/",
		TopImage:        "http://blogs-images.forbes.com/niallmccarthy/files/2015/11/20151109_Passports_Fo.jpg",
	}
	//article.Links = []string{""}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_FoxNewsCom(t *testing.T) {
	article := Article{
		Domain:          "foxnews.com",
		Title:           "Party insiders give Clinton early, commanding delegate edge",
		MetaDescription: "Hillary Rodham Clinton has locked up public support from half of the Democratic insiders who cast ballots at the party's national convention, giving her a commanding advantage over her rivals for the party's presidential nomination.",
		CleanedText:     "Published November 13, 2015Associated Press\n\n0\n\n0\n\nHillary Rodham Clinton has locked up public support from half of the Democratic insiders who cast ballots at the party's national convention, giving her a commanding advantage over her rivals for the party's presidential nomination.\n\nClinton's margin over Vermont Sen. Bernie Sanders and former Maryland Gov. Martin O'Malley is striking. Not only is it big, but it comes more than two months before primary voters head to the polls -- an early point in the race for so many of the people known as superdelegates to publicly back a candidate.\n\n\"She has the experience necessary not only to lead this country, she has experience politically that I think will help her through a tough campaign,\" said Unzell Kelley, a county commissioner from Alabama.\n\n\"I think she's learned from her previous campaign,\" he said. \"She's learned what to do, what to say, what not to say -- which just adds to her electability.\"\n\nThe Associated Press contacted all 712 superdelegates in the past two weeks, and heard back from more than 80 percent. They were asked which candidate they plan to support at the convention next summer.\n\nThe 712 superdelegates make up about 30 percent of the 2,382 delegates needed to clinch the Democratic nomination. That means that more than two months before voting starts, Clinton already has 15 percent of the delegates she needs.\n\nThat sizable lead reflects Clinton's advantage among the Democratic Party establishment, an edge that has helped the 2016 front-runner build a massive campaign organization, hire top staff and win coveted local endorsements.\n\nSuperdelegates are convention delegates who can support the candidate of their choice, regardless of who voters choose in the primaries and caucuses. They are members of Congress and other elected officials, party leaders and members of the Democratic National Committee.\n\nClinton is leading most preference polls in the race for the Democratic nomination, most by a wide margin. Sanders has made some inroads in New Hampshire, which holds the first presidential primary, and continues to attract huge crowds with his populist message about income inequality.\n\nBut Sanders has only recently started saying he's a Democrat after a decades-long career in politics as an independent. While he's met with and usually voted with Democrats in the Senate, he calls himself a democratic socialist.\n\n\"We recognize Secretary Clinton has enormous support based on many years working with and on behalf of many party leaders in the Democratic Party,\" said Tad Devine, a senior adviser to the Sanders campaign. \"But Sen. Sanders will prove to be the strongest candidate, with his ability to coalesce and bring young people to the polls the way that Barack Obama did.\"\n\n\"The best way to win support from superdelegates is to win support from voters,\" added Devine, a longtime expert on the Democrats' nominating process.\n\nThe Clinton campaign has been working for months to secure endorsements from superdelegates, part of a strategy to avoid repeating the mistakes that cost her the Democratic nomination eight years ago.\n\nIn 2008, Clinton hinged her campaign on an early knockout blow on Super Tuesday, while Obama's staff had devised a strategy to accumulate delegates well into the spring.\n\nThis time around, Clinton has hired Obama's top delegate strategist from 2008, a lawyer named Jeff Berman, an expert on the party's arcane rules for nominating a candidate for president.\n\nClinton's increased focus on winning delegates has paid off, putting her way ahead of where she was at this time eight years ago. In December 2007, Clinton had public endorsements from 169 superdelegates, according to an AP survey. At the time, Obama had 63 and a handful of other candidates had commitments as well from the smaller fraction of superdelegates willing to commit to a candidate.\n\n\"Our campaign is working hard to earn the support of every caucus goer, primary voter and grassroots and grasstop leaders,\" said Clinton campaign spokesman Jesse Ferguson. \"Since day one we have not taken this nomination for granted and that will not change.\"\n\nSome superdelegates supporting Clinton said they don't think Sanders is electable, especially because of his embrace of socialism. But few openly criticized Sanders and a handful endorsed him.\n\n\"I've heard him talk about many subjects and I can't say there is anything I disagree with,\" said Chad Nodland, a DNC member from North Dakota who is backing Sanders.\n\nHowever, Nodland added, if Clinton is the party's nominee, \"I will knock on doors for her. There are just more issues I agree with Bernie.\"\n\nSome superdelegates said they were unwilling to publicly commit to candidates before voters have a say, out of concern that they will be seen as undemocratic. A few said they have concerns about Clinton, who has been dogged about her use of a private email account and server while serving as secretary of state.\n\n\"If it boils down to anything I'm not sure about the trust factor,\" said Danica Oparnica, a DNC member from Arizona. \"She has been known to tell some outright lies and I can't tolerate that.\"\n\nStill others said they were won over by Clinton's 11 hours of testimony before a GOP-led committee investigating the attack on a U.S. consulate in Benghazi, Libya. Clinton's testimony won widespread praise as House Republicans struggled to trip her up.\n\n\"I don't think that there's any candidate right now, Democrat or Republican, that could actually face up to that and come out with people shaking their heads and saying, `That is one bright, intelligent person,\"' said California Democratic Rep. Tony Cardenas.",
		MetaKeywords:    "Democratic National Committee,Hillary Rodham Clinton,Barack Obama,presidential primary,primary voters,superdelegates",
		CanonicalLink:   "http://www.foxnews.com/politics/2015/11/13/party-insiders-give-clinton-early-commanding-delegate-edge/",
		TopImage:        "http://a57.foxnews.com/global.fncstatic.com/static/managed/img/fn2/video/0/0/111215_otr_clinton_1280.jpg",
	}
	article.Links = []string{
		"http://www.ap.org/",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_GloboesporteGloboCom(t *testing.T) {
	article := Article{
		Domain:          "globoesporte.globo.com",
		Title:           "Rodrigo Caio treina até nas férias e tenta acelerar retorno aos gramados",
		MetaDescription: "Rodrigo Caio treina na esteira durante as férias em Dracena-SP (Foto: Divulgação)Rodrigo Caio quer ganhar tempo na recuperação da lesão que sofreu no joelho esquerdo. Apesar de ter sido liberado pelo departamento médico do São Paulo para as férias, o ...",
		CleanedText:     "Rodrigo Caio treina até nas férias e tenta acelerar retorno aos gramadosJogador segue programação de exercícios em Dracena, interior de São Paulo. Comissão técnica planeja volta dele para o fim de fevereiro ou início de março Rodrigo Caio treina na esteira durante as férias em Dracena-SP (Foto: Divulgação)Rodrigo Caio quer ganhar tempo na recuperação da lesão que\n\nsofreu no joelho esquerdo. Apesar de ter sido liberado pelo departamento médico\n\ndo São Paulo para as férias, o jogador vem treinando diariamente para acelerar\n\na recuperação após ser submetido a uma cirurgia.O zagueiro e volante passa férias com a família em Dracena, interior\n\nde São Paulo, e alterna os períodos de descanso com uma rotina de\n\nexercícios. Ele vem realizando trabalhos de reforço muscular e corridas na\n\nesteira.\u00a0O jogador lesionou o joelho esquerdo no dia 2 de agosto,\n\ncontra o Criciúma, no Morumbi, pelo Campeonato Brasileiro, e precisou passar por\n\numa cirurgia. O defensor vinha sendo um dos destaques do São Paulo na\n\ntemporada.\u00a0Na avaliação do departamento médico, Rodrigo Caio deve\n\nser liberado para treinos com o elenco e jogos entre fevereiro e março. Com\n\nisso, é provável que seja inscrito pelo técnico Muricy Ramalho para disputar a\n\nfase de grupos da Taça Libertadores.\n\nsobre\n\nSão Paulo\n\n+\n\nAnterior\n\n30\n\nDez\n\n18:27\n\nBLOG: Corinthians corre risco de perder Dudu para o São Paulo\n\n16:30\n\nBLOG: RETROSPECTIVA 2014: Entre variações táticas, o ano foi da intensidade, 3 zagueiros e contragolpe\n\n12:01\n\nEm último teste antes da Copinha, São Paulo empata com Botafogo-SP\n\n11:06\n\nSão Paulo tenta fazer acordo para se livrar do 'mico' Clemente Rodríguez\n\n10:00\n\nVolante Hudson vê São Paulo pronto para conquistar títulos em 2015\n\n07:25\n\nTricolor recebe Botafogo em Cotia e faz último amistoso antes da Copinha\n\n29\n\nDez\n\n19:05\n\nApós ano no Dragão, Caramelo pode ser cedido pelo São Paulo à Chape\n\n16:28\n\nConmebol divulga tabela detalhada da Taça Libertadores de 2015; veja\n\n09:15\n\nAidar acredita em brilho de Pato, mas cobra: \"Ainda não mostrou a que veio\"\n\n08:00\n\nTimes paulistas tentam manter hegemonia recente na Copinha\n\nProximo\n\n+\n\nAnterior\n\n24\n\nDez\n\n08:10\n\nSem dor, Rodrigo Caio vence etapas e já pensa na volta aos gramados\n\n19\n\nDez\n\n21:44\n\nVice do São Paulo diz que Alvaro fica, revela parceria e quer comprar Pato\n\n13:52\n\nEm recuperação, Rodrigo Caio segue rotina no CT e ganha apoio de Ganso\n\n09\n\nDez\n\n15:04\n\nRodrigo Caio vence nova etapa de recuperação e inicia corrida na esteira\n\n25\n\nSet\n\n17:54\n\nSaudade! Toloi e Rodrigo Caio observam treino dos reservas no CT\n\n11\n\nSet\n\n16:02\n\nRodrigo Caio inicia nova etapa de recuperação e festeja evolução\n\n14\n\nAgo\n\n17:13\n\nRodrigo Caio começa fisioterapia após cirurgia no joelho esquerdo\n\n12\n\nAgo\n\n18:23\n\nCom a família por perto, Rodrigo Caio comenta dificuldades após a operação\n\n07\n\nAgo\n\n16:23\n\nRodrigo Caio passa por cirurgia e inicia fisioterapia na próxima semana\n\n06\n\nAgo\n\n08:05\n\nLesão de Rodrigo Caio trará reflexos dentro e fora de campo no São Paulo\n\nProximo",
		MetaKeywords:    "notícias, notícia, presidente prudente região",
		CanonicalLink:   "http://globoesporte.globo.com/sp/presidente-prudente-regiao/noticia/2014/12/rodrigo-caio-treina-ate-nas-ferias-e-tenta-acelerar-retorno-aos-gramados.html",
		TopImage:        "http://s.glbimg.com/es/ge/f/original/2014/12/26/10863872_894379987249341_2406060334390226774_o.jpg",
	}
	article.Links = []string{
		"http://globoesporte.globo.com/atleta/rodrigo-caio.html",
		"http://globoesporte.globo.com/sp/ribeirao-preto-e-regiao/noticia/2014/12/em-ultimo-teste-antes-da-copinha-sao-paulo-empata-com-botafogo-sp.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/sao-paulo-tenta-fazer-acordo-para-se-livrar-do-mico-clemente-rodriguez.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/mg/zona-da-mata-centro-oeste/noticia/2014/12/volante-hudson-ve-sao-paulo-pronto-para-conquistar-titulos-em-2015.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/tricolor-recebe-botafogo-em-cotia-e-faz-ultimo-amistoso-antes-da-copinha.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/noticia/2014/12/apos-ano-no-atletico-go-caramelo-pode-ser-cedido-pelo-sao-paulo-chape.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/libertadores/noticia/2014/12/conmebol-divulga-tabela-e-timao-x-sao-paulo-pode-abrir-fase-de-grupos.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/aidar-acredita-em-brilho-de-pato-mas-cobra-ainda-nao-mostrou-que-veio.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/Copa-SP-de-futebol-junior/noticia/2014/12/times-paulistas-tentam-manter-hegemonia-recente-na-copinha.html#equipe-sao-paulo",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/sem-dor-rodrigo-caio-vence-etapas-e-ja-pensa-na-volta-aos-gramados.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/vice-do-sao-paulo-diz-que-alvaro-fica-revela-parceria-e-quer-comprar-pato.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/fora-ha-cinco-meses-rodrigo-caio-dispensa-ferias-e-tem-papo-com-ganso.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/12/rodrigo-caio-vence-nova-etapa-de-recuperacao-e-inicia-corrida-na-esteira.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/09/saudade-toloi-e-rodrigo-caio-observam-treino-dos-reservas-no-ct.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/09/rodrigo-caio-inicia-nova-etapa-de-recuperacao-e-festeja-evolucao.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/08/rodrigo-caio-comeca-fisioterapia-apos-cirurgia-no-joelho-esquerdo.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/sp/presidente-prudente-regiao/noticia/2014/08/com-familia-por-perto-rodrigo-caio-comenta-dificuldades-apos-operacao.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/08/rodrigo-caio-opera-joelho-e-iniciara-fisioterapia-no-ct-na-proxima-semana.html#atleta-rodrigo-caio",
		"http://globoesporte.globo.com/futebol/times/sao-paulo/noticia/2014/08/lesao-de-rodrigo-caio-trara-reflexos-dentro-e-fora-de-campo-no-sao-paulo.html#atleta-rodrigo-caio",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_HbrOrg(t *testing.T) {
	article := Article{
		Domain:          "hbr.org",
		Title:           "Email Is the Best Way to Reach Millennials",
		MetaDescription: "It’s still the channel with the highest ROI.",
		CleanedText:     "With $200 billion in annual buying power by 2017, Millennials have become every brand’s coveted customer. But what’s the best way to reach them?\n\nThe answer is email.\n\nFor all the talk of email being dead — Too much noise! Too much spam! Too many distractions! Snapchat! — email remains\u00a0the standard for digital communication. In fact, Millennials check email more than any other age group, and nearly half can’t even use the bathroom without checking it, according to a\u00a0recent Adobe study.\n\nThat same study\u00a0found nearly 98% of Millennials check their personal email at least every few hours at work, while almost 87% of Millennials check their work email outside of work.\n\nEmail is not only relevant for Millennials, it also happens to remain the channel where direct marketers get the highest ROI ($39 for every dollar spent, according to the Direct Marketing Association). But that doesn’t mean the same old email marketing will work on Millennials. Instead, marketers need to adjust, or run the risk of that dreaded swipe to the trash bin. Consider these ideas the next time you’re planning an email campaign and Millennials are a key part of the audience:\n\nMobile is a must. Millennials are more likely than any other age group to check email on smartphones, with 88% reporting that they regularly using a smartphone to check email. If you’re not mobile first, you’re not putting your Millennial customers first. Responsive design has been a mantra for some time, but if you’re not employing it, you’re alienating an important generation of consumers who live, breathe, and sleep with their mobile devices.\n\nTiming is everything. Looking at opens and clicks won’t get you anywhere without analyzing the day of week and time of day those emails are opened and clicked. For example, we found that Millennials are more likely than any other age group to check email while in bed (45.2%). Why not experiment with sending emails first thing in the morning or late in the evening with content relevant to that time of day?\n\nPictures are worth a thousand words. They’re also an important mechanism for Millennials to filter messages. Why send an email survey asking for written feedback when all you need to do is provide a choice between a smiley face and a frown? Images are an integral part of Millennial language, even in the workplace. A third of Millennials believe it is appropriate to use an emoji when communicating with a direct manager or senior executive, so it’s a safe bet they’re even more comfortable when it comes to emoji from brands. Millennials are thinking and communicating in images, so marketers need to optimize emails for images and allow for quick feedback through emoji.\n\nLess is more. Email marketing to Millennials isn’t about sending more of the same. Many Millennials want to see fewer emails (39%) and fewer repetitive emails from brands (32%). Marketers take note — stop spamming your lists and start marketing to individuals by understanding who they are first.\n\nNot every Millennial communicates the same way, of course. And digital communication is constantly evolving. Nonetheless, for now it seems safe to say that email is here to stay and will remain a critical channel even for reaching mobile customers. Just don’t expect the same old email tactics to work.",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "https://hbr.org/resources/images/article_assets/2015/11/nov15-12-169799513-horz.jpg",
	}
	article.Links = []string{
		"https://blogs.adobe.com/conversations/2015/08/email.html",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_HuffingtonPostCoUk(t *testing.T) {
	article := Article{
		Domain:          "huffingtonpost.co.uk",
		Title:           "How We Are Controlling The Future Of TV Scheduling",
		MetaDescription: "var isMobile = {\n    Android: function() {\n        return navigator.userAgent.match(/Android/i);\n    },\n    BlackBerry: function() {\n        return navigator.userAgent.match(/BlackBerry/i);",
		CleanedText:     "Since its inception, television has been a unifying social force, bringing family, friends and different groups of people together. Even watching television on your own connects you to the multitudes of others watching the same thing across the globe.TV has come a long way: from black-and-white to colour, from a rare treat accessible to few to a household staple for everyone, from standard definition to tomorrow's ultra-HD screens.Perceptions of TV audiences have also changed over time. While theorists once believed TV viewers were passive, zombie-like figures transfixed in front of their televisions, numerous studies have proven that TV audiences are engaged, active and critical of the programmes they watch.In the last several years, we've seen a dramatic shift that's placed viewers in control of their own scheduling. There's also more choice than ever before when it comes to accessing favourite programmes and watching them when and where they like.\"There are two simultaneous trends emerging when it comes to our TV watching habits, and they're two opposite trends, which is interesting,\" says Professor Sonia Livingstone OBE, a full professor in the Department of Media and Communications at the London School of Economics.\"One: we're watching TV on our laptops, tablets and phones, wherever and on whatever.\"And two, somewhat paradoxically, we're seeing a growth in the size of the screen in the living room. People talk about how everyone is watching TV on a 'small screen', but there's also a new viewing growing up around this enormous screen, as well as the more individualised viewing.\"Now, we watch shows wherever we want, whether it's relaxing in the bath with Corrie characters, catching up with a favourite drama on our phone during a morning commute or settling down in the sitting room every week to enjoy GBBO, gathered around the biggest 'and best' screen in the house. Equally, thanks to the latest in wearable tech, our most beloved television content has become a coveted accessory, accessible with a swipe on our watch.Subscription-free services like Freeview Play have also given us more options than ever before, with over 60 TV channels, 12 HD channels and over 25 radio stations a remote click away, plus the freedom to catch up on shows from the BBC, ITV, Channel 4 and Channel 5. Other services like Netflix and Amazon Prime also give us the opportunity to watch shows we missed the first time around - in one sitting, if we so desire! - while simultaneously introducing us to new and original programming.\"We keep fearing that people won't talk to each other anymore,\" says Professor Livingstone. \"There's the choice to watch separately and the choice to come together, whether it's binge viewing or the greater choice of programmes than ever before.\"All of this choice has had a positive impact on TV consumers, according to Professor Livingstone.\"Most of the evidence is that people are feeling empowered and delighted. There's been an enormous welcome from people about the joys of having so much control and more choice than ever before.\"People are also prepared to pay to improve their television watching experience, whether that's spending on bigger HD screens or subscription services.While scheduling is fairly unimportant for younger generations, the middle-aged and young elderly population that remembers how television used to be is growing, so scheduling continues to play an important role for them.For those younger generations, the definition of whether TV is 'a five minute clip of a beauty vlogger's latest haul on YouTube or a critically respected docudrama' calls into question what TV viewing really means these days.\"People have been saying for a while that scheduling is dead, but there's no getting rid of schedule for the 40s or 50-pluses who absolutely adhere to traditions of what to watch and when,\" says Professor Livingstone.Rapidly emerging trends, like the increase in individual TV consumption across new tech and the importance of the living room big screen as the centrepoint of family life, ensure that the landscape of television scheduling is in constant flux and the future of television remains uncertain.One thing we know? We'll still be watching.",
		MetaKeywords:    "changing, channels:, how, we, are, controlling, the, future, of, tv, scheduling, uk, entertainment",
		CanonicalLink:   "http://www.huffingtonpost.co.uk/2015/10/29/how-we-are-changing-the-future-of-tv-scheduling_n_8303736.html",
		TopImage:        "http://i.huffpost.com/gen/3507100/images/o-TELEVISION-REMOTE-CONTROL-facebook.jpg",
	}
	//article.Links = []string{""}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_IncCom(t *testing.T) {
	article := Article{
		Domain:          "inc.com",
		Title:           "How Rent the Runway Plans to Own Your Closet",
		MetaDescription: "A new business model, brick-and-mortar stores, and $70 million in venture capital funding. Here's how this business lit up runways (and sidewalks) in 2015.",
		CleanedText:     "Forgot Password?\n\nNew member? Sign up now.\n\nSign in if you're already registered.\n\nMark Cuban: What I Would Do If I Were President\n\nSamuel Adams Creator Jim Koch on Scaling up, One Barrel at a Time\n\nWhy America Needs a CEO in the White House\n\n2 Traits That Give Veterans an Entrepreneurial Advantage\n\n3 Key Traits Shared by the Most Successful Business Leaders\n\nWhy Startups Need to Be Able to Survive Without Their Founders\n\nRussell Simmons: Why It's Important to Do What You Love\n\nMark Cuban: How You'll Know You're Ready to Launch\n\nThe 4 Mentors Every Entrepreneur Needs\n\nDaymond John: 5 Traits That Make a Good Business Leader\n\nHow Marcus Lemonis Knows If You're Making Good Money\n\nArianna Huffington: The Wake-Up Call That Helped Arianna Huffington Learn to Thrive\n\nThe Making of Inc.'s Jessica Alba Cover Story\n\nSecrets of Wealth and Success From Tony Robbins\n\nBarbara Corcoran's 8 Lessons for Entrepreneurs\n\nMint Founder: How to Learn From Your Early Mistakes\n\nOne Nightly Productivity Tip to Get the Most out of Your Day\n\nHow to Keep the Fear of Failure From Stalling Personal Growth\n\nWhy Entrepreneurship Is a 24/7 Lifestyle\n\nWhy the Only Guaranteed Path to Success Is Through Hard Work and Hustle",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.inc.com/zoe-henry/rent-the-runway-2015-company-of-the-year-nominee.html",
		TopImage:        "http://www.inc.com/uploaded_files/image/970x450/OUT63313304-web_70674.jpg",
	}
	article.Links = []string{
		"http://www.inc.com/",
		"https://magazine.inc.com/servlet/ConvertibleGateway?cds_mag_code=ICM&cds_page_id=136768&cds_response_key=XB5KNNGF1",
		"https://www.facebook.com/Inc",
		"https://twitter.com/inc",
		"https://www.linkedin.com/company/inc--magazine",
		"https://plus.google.com/+incmagazine",
		"https://www.pinterest.com/incmagazine/",
		"http://www.youtube.com/user/incmagazine?sub_confirmation=1",
		"https://instagram.com/incmagazine",
		"https://flipboard.com/@incmagazine",
		"http://www.inc.com/mark-cuban/what-i-would-do-if-i-were-president.html",
		"http://www.inc.com/jim-koch/samuel-adams-creator-on-scaling-up-one-barrel-at-a-time.html",
		"http://www.inc.com/donny-deutsch/why-america-needs-a-ceo-in-the-white-house.html",
		"http://www.inc.com/norm-brodsky/2-traits-that-give-veterans-a-leg-up-as-entrepreneurs.html",
		"http://www.inc.com/donny-deutsch/3-key-traits-shared-by-the-most-successful-business-leaders.html",
		"http://www.inc.com/gary-vaynerchuk/askgaryvee-episode-84-surviving-without-a-founder.html",
		"http://www.inc.com/russell-simmons/why-its-important-to-do-what-you-love.html",
		"http://www.inc.com/mark-cuban/how-youll-know-youre-ready-to-launch.html",
		"http://www.inc.com/kim-kaupe/4-mentors-that-every-entrepreneur-needs.html",
		"http://www.inc.com/daymond-john/5-traits-that-make-a-good-business-leader.html",
		"http://www.inc.com/marcus-lemonis-bees-knees-spicy-honey.html",
		"http://www.inc.com/arianna-huffington/founders-forum-how-huffington-learned-to-thrive.html",
		"http://www.inc.com/jessica-alba/the-making-of-inc-jessica-alba-cover-story.html",
		"http://www.inc.com/tony-robbins/tony-robbins-reveals-his-secrets-on-wealth-success-and-financial-freedom.html",
		"http://www.inc.com/barbara-corcoran/eight-lessons-for-entrepreneurs.html",
		"http://www.inc.com/aaron-patzer/how-to-learn-from-early-mistakes.html",
		"http://www.inc.com/adam-miller/one-nightly-productivity-tip-to-get-the-most-out-of-your-day.html",
		"http://www.inc.com/jen-groover/how-to-keep-failure-away-from-personal-growth.html",
		"http://www.inc.com/ravin-gandhi/why-entrepreneurship-is-a-lifestyle.html",
		"http://www.inc.com/gary-vaynerchuk/askgaryvee-episode-86-hard-work-and-hustle.html",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_LinkedinCom(t *testing.T) {
	article := Article{
		Domain:          "linkedin.com",
		Title:           "An Unexpected Way to Achieve a Better Work-Life Balance",
		MetaDescription: "Work-life balance. Everyone talks about it. And everyone struggles to achieve it. Yet finding a reasonable work-life balance is easier than you think",
		CleanedText:     "Work-life balance. Everyone talks about it. And everyone struggles to achieve it.\n\nYet finding a reasonable work-life balance is easier than you think. While it's true the equilibrium point is constantly shifting, most of the same attitudes, perspectives, and skills apply to both \"work\" and \"life.\"\n\nSo why not take advantage of that fact? Pick the right \"life\" pursuits and they inform and enhance your professional skills -- and add a healthy dose of perspective and humility along the way.\n\nIn my case I like to take on extremely difficult (at least for me) physical goals. (Granted my approach to goal achievement in general is a little unconventional. Just like\u00a0Fight Club,\u00a0the first rule of achieving a goal is\u00a0you don't talk about achieving that goal. And achieving a goal has a lot less to do with the goal itself and\u00a0a lot more to do with the routine you develop\u00a0to support that goal.)\n\nSo a few years ago, after just four months of training, I rode the\u00a0Alpine Loop Gran Fondo, a 92-mile, four-mountain ride that included 11,000 feet of climbing. (Those four months felt like a lifetime, though, since pro mountain biker Jeremiah Bishop trained me. But then again I never could have been ready without him.)\n\nAfter a few years of cycling I got tired of being cycling skinny -- 6' tall, 150 lbs is not a particularly good look -- and decided to see if I could pull off some semblance of the\u00a0\"movie star becomes an action hero\"\u00a0physical transformation. I gained over 20 pounds, lost a few percentage points of body fat, and got a lot stronger. (That training sucked too, since\u00a0Jeffrey Del Favero\u00a0of\u00a0Bodybuilding.com\u00a0created my program, but then again I never could have done it without him.)\n\nSo why do I do take on (feel free to insert your own adjective) personal challenges? And how does that help me professionally? It's all about the habits, skills, and perspectives gained. Here are some reasons.\n\nSuccess is ultimately based on numbers. Sure, you can try to \"hack\" a goal. Sure, you can look for shortcuts. (People have\u00a0built entire careers\u00a0off the premise.) But eventually achieving a huge goal is all about volume and repetition.\n\nWant to eventually ride a tough gran fondo? You'll have to ride hundreds of miles along the way. Want to go from only being able to do three pull-ups to eventually being able to do four sets of twenty? You'll have to lift a ton of weight along the way.\n\nThe same is true for professional success; it's largely based on doing the work. Want twenty new customers? Expect to cold call two or three hundred prospects. Want to hire a superstar? Expect to screen dozens and then interview ten or fifteen people.\n\nThe surest path to success is to do an incredible amount of work. If you're willing to do the work, you can succeed at almost anything.\n\nThe armor that protects us eventually destroys us. We all wear armor. That armor protects us but also, over time, wears us down.\n\nOur armor is primarily forged by success. Every accomplishment adds an additional layer of protection from vulnerability. In fact, when we feel particularly insecure we unconsciously strap on more armor so we feel less vulnerable:\n\nArmor protects when we're unsure, tentative, or at a perceived disadvantage. Our armor says, \"That's okay; I may not be good at this... but I'm really good at\u00a0that.\"\n\nOver time armor also encourages us to narrow our focus to our strengths. That way we stay safe. The more armor we put on the more we can hide our weaknesses and failings--from others and from ourselves.\n\nWe use our armor all the time. I use my armor all the time--I feel sure more than you. But I get really tired of wearing it.\n\nWhen I ride a bike the guy who passes me doesn't care if I've ghostwritten bestsellers or drive a fancy car or live in a nice neighborhood. At the gym, the guy who lifts more than me also doesn't care about any of that stuff. He's stronger and fitter than me. Period.\n\nIn those situations no amount of armor, real or imagined, can protect me. I'm just a guy on a bike. I'm just a guy at the gym. I'm just me.\n\nBeing just me is pretty scary.\n\nBut being who you really are is something we all need to do more often. It keeps things in perspective. It reminds us that we can always be better. It reminds us that no matter how good we think we are at something there is always someone who is a lot better.\n\nAnd that's not depressing -- that's motivating.\n\nGrace is an awesome feeling -- one we can never experience enough. Outstanding athletes exist in a state of grace, a place where calculation and strategy and movement happen almost unconsciously. Great athletes can focus in a way that, to us, is unrecognizable because through skill, training, and experience their ability to focus is nearly effortless.\n\nWe've all felt a sense of grace, if only for a few precious moments, when we performed better than we ever imagined possible... and realized what we assumed to be limits weren't really limits at all.\n\nThose moments don't happen by accident, though. Grace is never given; grace must be earned through discipline and training and sacrifice.\n\nI want to ride up a mountain and experience the feeling that I can climb and climb and climb and I don't have to think about anything because I can just\u00a0go....\n\nI want to struggle with a weight and experience the feeling that I can do a few more reps because I know, without a doubt, I always have a little more in me...\n\nAnd I want to sometimes write almost effortlessly and without thinking because years of effort and practice have brought me to a place where occasionally I am the writer I would like to be...\n\nAll those are moments of grace. They're awesome. They're amazing.\n\nAnd they feed off each other because the confidence you build after experiencing a moment of grace in one pursuit helps you keep pushing when the going gets tough in other pursuits.\n\nWith work, \"then\" is always better than \"now.\"\u00a0 \"Now\" and \"then\" are wonderful words when they appear in the same sentence.\n\nWhen you work to improve at something -- especially in the beginning stages -- \"now\" is often a terrible place. At one point my \"now\" was riding like an asthmatic hippo. At one point my \"now\" was doing four dips and feeling like I was tearing my chest apart.\n\nBut with time and effort my \"now\" was transformed. I could ride\u00a0with more speed, power, and confidence. I could do\u00a0sets of ten, then twenty, then thirty dips. I was able to look back with satisfaction at a \"now\" I had transformed into a vastly inferior \"then.\"\n\nThink about something you wanted to do. Then think about where you would be\u00a0now\u00a0if you had actually gotten started on it\u00a0then.\n\nWhen you do the work, then always pales in comparison to now: family, business, and every aspect of your life. When you don't do the work, now is just like then -- except now you also get to live with regret.\n\nQuitting is a habit anyone can learn to break. We're all busy. Each of us face multiple, ongoing demands. Every day we are forced a number of times to say, \"That's not perfect, but it works... and I need to move on to something else.\"\n\nStopping short of excellence is something we are not just forced to do but are also\u00a0trained\u00a0to do. Most of the time we have no choice so we get really good at \"quitting.\"\n\nI'm really good at quitting. I raised wonderful kids and did a good job... but I know I could have done more. I've built a decent business... but I know I could have done more. I've tackled challenges before and tried really hard... but I know I could have done more.\n\nWhere physical challenges are concerned there are hundreds if not thousands of times I want to quit. Training is hard and only gets harder. Balancing family and work and everything else is hard and only gets harder.\n\nAt weak moments, struggle shatters our resolve and make us want to quit.\n\nIt's hard not to stop, by choice or otherwise, at \"good enough.\" But sometimes, if the goal is big enough, we have to be\u00a0great: not great compared to other people... but great compared to ourselves.\n\nThat comparison is the only comparison that really matters and is the best reason of all to try to accomplish more than you -- or anyone around you -- ever thought possible.\n\nWhen you succeed, you become something you were not. And then you get to do it again, and become\u00a0something else you once were not -- but definitely are now.\n\nCheck out my book of personal and professional advice,\u00a0TransForm: Dramatically Improve Your Career, Business, Relationships, and Life -- One Simple Step At a Time. (PDF version here,\u00a0Kindle version here,\u00a0Nook version here.)\n\nIf after 10 minutes you don't find at least 5 things you can do to make your life better I'll refund your money.\n\nThat way you have nothing to lose... and everything to gain.",
		MetaKeywords:    "",
		CanonicalLink:   "https://www.linkedin.com/pulse/unexpected-way-achieve-better-work-life-balance-jeff-haden",
		TopImage:        "http://m.c.lnkd.licdn.com/mpr/mpr/AAEAAQAAAAAAAATuAAAAJGRiODU4MjBjLTFlZTEtNGQ3NS05ZDk1LTZiNjVkYjE5NWZlNA.jpg",
	}
	article.Links = []string{
		"http://www.inc.com/jeff-haden/silence-the-surprising-way-to-achieve-a-goal.html",
		"http://www.inc.com/jeff-haden/an-nearly-foolproof-way-to-achieve-every-goal-you-set-wed.html",
		"http://www.alpineloopgranfondo.com/",
		"http://www.huffingtonpost.com/2014/12/01/jake-gyllenhaal-southpaw_n_6251010.html",
		"https://www.linkedin.com/pub/jeffrey-del-favero/23/b5a/a15",
		"http://www.bodybuilding.com/",
		"http://fourhourworkweek.com/blog/",
		"http://www.inc.com/author/jeff-haden",
		"https://gumroad.com/l/YHadh",
		"https://gumroad.com/l/YHadh",
		"http://amzn.to/1EiaVXV",
		"http://www.barnesandnoble.com/w/books/1121702502?ean=2940151263917",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_NyTimesCom(t *testing.T) {
	article := Article{
		Domain:          "nytimes.com",
		Title:           "How Gun Traffickers Get Around State Gun Laws",
		MetaDescription: "The effect of state gun control laws is diluted by a thriving underground market for firearms brought from states with few restrictions.",
		CleanedText:     "Wash.\n\nWhere guns used in crimes came from\n\nMe.\n\nArrow sizes show the number of guns traced to other states in 2014\n\nMont.\n\nN.D.\n\nMinn.\n\nVt.\n\nOre.\n\nN.H.\n\nIdaho\n\nN.Y.\n\nWis.\n\nS.D.\n\nMass.\n\nMich.\n\nR.I.\n\nWyo.\n\nConn.\n\nPa.\n\nIowa\n\nN.J.\n\nNeb.\n\nNev.\n\nMore than two-thirds of guns connected to crimes in New York and New Jersey were brought in from other states, mostly from the South.\n\nMd.\n\nOhio\n\nDel.\n\nUtah\n\nIll.\n\nW.Va.\n\nColo.\n\nD.C.\n\nInd.\n\nVa.\n\nKan.\n\nMo.\n\nCalif.\n\nKy.\n\nN.C.\n\nTenn.\n\nOkla.\n\n1,184 guns\n\nfrom arizona\n\nN.M.\n\nArk.\n\nS.C.\n\nAriz.\n\nGa.\n\nAla.\n\nMiss.\n\nCriminals in California used about 6,000 guns from other states, mainly from those with few gun-buying restrictions like Arizona and Nevada.\n\nLa.\n\nTexas\n\nCrime rings smuggle guns from Orlando, Fla., to Puerto Rico.\n\nFla.\n\n349 guns\n\nfrom florida\n\nState gun control laws\n\nLENIENT\n\nstrict\n\nPuerto Rico\n\nWhere guns used in crimes came from\n\nArrow sizes show the number of guns traced to other states in 2014\n\nOre.\n\nN.Y.\n\nN.J.\n\nNev.\n\nIllinois\n\nIndiana\n\nCalif.\n\nN.C.\n\nAriz.\n\nS.C.\n\n1,184 guns\n\nfrom arizona\n\nGa.\n\nTexas\n\nFlorida\n\nState gun control laws\n\n349 guns\n\nfrom florida\n\nLENIENT\n\nstrict\n\nPuerto Rico\n\nWhere guns used in crimes came from\n\nWash.\n\nArrow sizes show the number of guns traced to other states in 2014\n\nMe.\n\nMont.\n\nN.D.\n\nMinn.\n\nVt.\n\nOre.\n\nN.H.\n\nIdaho\n\nN.Y.\n\nWis.\n\nS.D.\n\nMass.\n\nMich.\n\nWyo.\n\nPa.\n\nIowa\n\nN.J.\n\nNeb.\n\nNev.\n\nMd.\n\nOhio\n\nUtah\n\nMost guns connected to crimes in New York and New Jersey were brought in from other states, mostly from the South.\n\nIll.\n\nW.Va.\n\nColo.\n\nD.C.\n\nInd.\n\nVa.\n\nKan.\n\nMo.\n\nCalif.\n\nKy.\n\nN.C.\n\nTenn.\n\n1,184\n\nguns from arizona\n\nOkla.\n\nN.M.\n\nArk.\n\nS.C.\n\nAriz.\n\nGa.\n\nAla.\n\nMiss.\n\nLa.\n\nTexas\n\nCriminals in California used about 6,000 guns from other states, mainly from those with few gun-buying restrictions like Arizona and Nevada.\n\n349 guns\n\nfroM\n\nflorida\n\nFla.\n\nState gun control laws\n\nCrime rings smuggle guns from Orlando, Fla., to Puerto Rico.\n\nPuerto Rico\n\nLENIENT\n\nstrict\n\nWhere guns used in crimes came from\n\nWashington\n\nArrow sizes show the number of guns traced to other states in 2014\n\nMe.\n\nMontana\n\nNorth Dakota\n\nMinnesota\n\nVt.\n\nOregon\n\nN.H.\n\nIdaho\n\nNew York\n\nWisconsin\n\nSouth Dakota\n\nMass.\n\nMichigan\n\nR.I.\n\nWyoming\n\nConn.\n\nPa.\n\nIowa\n\nNew Jersey\n\nNeb.\n\nNevada\n\nMd.\n\nOhio\n\nMore than two-thirds of guns connected to crimes in New York and New Jersey were brought in from other states, mostly from the South.\n\nDel.\n\nUtah\n\nIllinois\n\nW.Va.\n\nColorado\n\nD.C.\n\nIndiana\n\nVa.\n\nKansas\n\nMo.\n\nCalifornia\n\nKy.\n\nN.C.\n\nTenn.\n\nOklahoma\n\n1,184 guns\n\nfrom arizona\n\nNew Mexico\n\nArkansas\n\nS.C.\n\nArizona\n\nGeorgia\n\nAlabama\n\nMiss.\n\nCriminals in California used about 6,000 guns from other states, mainly from those with few gun-buying restrictions like Arizona and Nevada.\n\nTexas\n\nCrime rings smuggle guns from Orlando, Fla., to Puerto Rico.\n\nLouisiana\n\nFlorida\n\nState gun control laws\n\n349 guns\n\nfrom florida\n\nLENIENT\n\nstrict\n\nPuerto Rico\n\nIn California, some gun smugglers use FedEx. In Chicago, smugglers drive just across the state line into Indiana, buy a gun and drive back. In Orlando, Fla., smugglers have been known to fill a $500 car with guns and send it on a ship to crime rings in Puerto Rico.\n\nIn response to mass shootings in the last few years, more than 20 states, including some of the nation’s biggest, have passed new laws restricting how people can buy and carry guns. Yet the effect of those laws has been significantly diluted by a thriving underground market for firearms brought from states with few restrictions.\n\nAbout 50,000 guns are found to be diverted to criminals across state lines every year, federal data shows, and many more are likely to cross state lines undetected.\n\nIn New York and New Jersey, which have some of the strictest laws in the country, more than two-thirds of guns tied to criminal activity were traced to out-of-state purchases in 2014. Many were brought in via the so-called Iron Pipeline, made up of Interstate 95 and its tributary highways, from Southern states with weaker gun laws, like Virginia, Georgia and Florida.\n\nNew York\n\nThe Iron Pipeline\n\nPa.\n\nGuns used in recent shootings of New York City police officers were traced to pawn shops in Georgia.\n\nJONESBORO\n\nVa.\n\nNew Jersey\n\n386 guns\n\nN.C.\n\nPERRY\n\nGa.\n\nS.C.\n\nMany guns used in crimes are brought to New York and New Jersey along Interstate 95. In recent years, more guns have started coming from Pennsylvania gun shows, a federal official said.\n\n292 guns\n\nFla.\n\nNew York\n\nPa.\n\nNew\n\nJersey\n\nGuns used in recent shootings of New York City police officers were traced to pawn shops in Georgia.\n\nVa.\n\nN.C.\n\nS.C.\n\nJONESBORO\n\n386\n\nguns\n\nThe Iron Pipeline\n\nPERRY\n\nMany guns used in crimes are brought to New York and New Jersey along Interstate 95. In recent years, more guns have started coming from Pennsylvania gun shows, a federal official said.\n\nGa.\n\n292\n\nguns\n\nFla.\n\nA handgun used in the killing of two Brooklyn officers last year was traced to a pawnshop just south of Atlanta. A revolver used in a fatal shooting of an officer in Queens in May was traced to a roadside pawnshop, also in Georgia, about 100 miles from Atlanta. And a handgun used to kill an officer in East Harlem last month was traced to South Carolina.\n\n“We’re trying to deal with it, but we have a spigot that’s wide open down there and we don’t have a national or local ability to shut that spigot down at the moment,” said the New York City police commissioner, William J. Bratton, as he announced an indictment against gun traffickers last week.\n\nThe economics are straightforward: A low-quality handgun that sells for $100 in an Atlanta store might sell for $500 or $600 in New York City, researchers say — and it can be transported cheaply. By contrast, the majority of guns used in crimes in Texas, Georgia and other states with more lenient gun laws are purchased in-state.\n\nThe New York Times examined gun trafficking by analyzing nine years of data compiled by the Bureau of Alcohol, Tobacco, Firearms and Explosives, as well as an index of state gun laws developed by researchers at Johns Hopkins University.\n\nLaw enforcement officials express frequent frustration that they are not able to track every gun that crosses state lines, which means the estimates here are conservative. When the police do recover a gun tied to criminal activity, typically after an arrest, they can trace the gun to where it was last sold through a federally licensed dealer.\n\nChicago offers perhaps the starkest example of trafficking. There are no retail gun dealers within city limits, because Chicago has some of the tightest municipal gun regulations. Yet bringing a gun into Chicago can be as simple as driving less than an hour to a gun show in Indiana, where private sales are not recorded and do not require a background check.\n\n“If you’re in the city of Chicago on the South Side, you may be closer to Indiana than you are to the Magnificent Mile,” said Roseanna Ander, executive director of the University of Chicago Crime Lab, referring to a well-known part of Chicago’s downtown.\n\nThe Route Into Chicago\n\nWisconsin\n\nMost guns used in crimes in Illinois were recovered in the Chicago area.\n\nMichigan\n\nIowa\n\nCHICAGO\n\n1,041 guns\n\nIllinois\n\nGun shows in Indiana are a frequent source for guns used in crimes in Illinois.\n\nIndiana\n\nMissouri\n\nMany people in Illinois have family ties to Mississippi, the second most common source for crime guns.\n\nThe Route Into Chicago\n\nMost guns used in crimes in Illinois were recovered in the Chicago area.\n\nWisconsin\n\nIowa\n\nCHICAGO\n\n1,041\n\nguns\n\nIllinois\n\nIndiana\n\nGun shows in Indiana are a frequent source for guns used in crimes in Illinois.\n\nMany people in Illinois have family ties to Mississippi, the second most common source for crime guns.\n\nMissouri\n\nMany guns follow a complex path from the original sale to the underground market. Most guns are originally bought from retail stores, but people who can’t pass a background check typically obtain guns from friends, family or illegal dealers.\n\nAccording to an anonymous survey of inmates in Cook County, Ill., covering 135 guns they had access to, only two had been purchased directly from a gun store. Many inmates reported obtaining guns from friends who had bought them legally and then reported them stolen, or from locals who had brought the guns from out of state.\n\nOne inmate said, “Some people get on a train and bring them back, can be up to five or six guns, depending on how much risk they want to take.”\n\nSome larger traffickers use more elaborate techniques. Buying a gun in Puerto Rico requires an expensive permit and a lengthy application process, but Florida has no such restrictions. Traffickers in Orlando tied to organized gangs in Puerto Rico send guns in the mail, through FedEx, or even encased in cars that travel by ship to the island.\n\n“They’ll buy a $500 car and stuff it with as many guns as possible,” said Carlos Gonzalez, an agent with the Miami division of the Bureau of Alcohol, Tobacco, Firearms and Explosives.\n\nGuns by Mail\n\nOrlando, which has a large Puerto Rican population, is the source for many guns trafficked to Puerto Rico.\n\nORLANDO\n\nFlorida\n\nMIAMI\n\n349 guns\n\nIn 2014, more guns used in crimes in Puerto Rico were traced to purchases in Florida than on the island itself.\n\nCuba\n\nPuerto\n\nRico\n\nHaiti\n\nDom.\n\nRep.\n\nGuns by Mail\n\nOrlando, which has a large Puerto Rican population, is the source for many guns trafficked to Puerto Rico.\n\nORLANDO\n\nFlorida\n\nMIAMI\n\n349\n\nguns\n\nCuba\n\nHaiti\n\nDom.\n\nRep.\n\nIn 2014, more guns used in crimes in Puerto Rico were traced to purchases in Florida than on the island itself.\n\nPuerto\n\nRico\n\nFederal agents and postal inspectors have caught some traffickers, leading to modified techniques, such as shipping guns in newer, more expensive cars or mailing guns from Jacksonville, Fla., instead of Orlando. Stopping such smuggling is logistically hard. “If the U.S. Postal Service were to screen every single package that entered into Puerto Rico, it would bring the economy to a halt,” Mr. Gonzalez said.\n\nMost gun trafficking patterns have remained remarkably constant over time. But some researchers point to a significant shift in Missouri as evidence that changes to one state’s laws can have broad implications.\n\nBefore 2007, Missouri required gun buyers to get a state permit and to undergo background checks on private sales, two restrictions strongly associated with states that provide fewer guns to interstate traffickers, according to research by Daniel Webster, director of the Johns Hopkins Center for Gun Policy and Research. At the time, nearly half of the guns used in crimes and recovered in Missouri were traced to other states, largely from neighboring Kansas and Illinois.\n\nBut when Missouri relaxed its gun control laws in 2007, the flow started to change. The number of guns traced to other states decreased, while the number of guns from within Missouri increased to nearly three-quarters.\n\nSource of guns used in crimes in Missouri\n\n80 percent\n\nGuns from Missouri\n\n74%\n\n60\n\nMore criminals used guns from Missouri after guns became easier to purchase.\n\n40\n\nMissouri repealed strict gun control laws\n\nin August 2007.\n\n26%\n\n20\n\nGuns imported from other states\n\n2014\n\n’12\n\n’10\n\n’08\n\n’06\n\n’04\n\n2002\n\nSource of guns used in crimes in Missouri\n\n80 percent\n\nGuns from Missouri\n\n74%\n\n60\n\nMore criminals used guns from Missouri after guns became easier to purchase.\n\n40\n\nMissouri repealed strict gun control laws in August 2007.\n\n26%\n\nGuns imported from other states\n\n20\n\n2014\n\n’12\n\n’10\n\n’08\n\n’06\n\n’04\n\n’02",
		MetaKeywords:    "Gun Control,Attacks on Police,Firearms,Bureau of Alcohol  Tobacco and Firearms",
		CanonicalLink:   "http://www.nytimes.com/interactive/2015/11/12/us/gun-traffickers-smuggling-state-gun-laws.html",
		TopImage:        "http://static01.nyt.com/images/2015/11/12/us/gun-traffickers-smuggling-state-gun-laws-1447372488027/gun-traffickers-smuggling-state-gun-laws-1447372488027-articleLarge-v4.png",
	}
	article.Links = []string{
		"http://www.motherjones.com/politics/2013/12/state-gun-laws-after-newtown",
		"https://www.atf.gov/resource-center/data-statistics",
		"http://www.nytimes.com/2014/12/25/nyregion/tracing-the-gun-used-to-kill-2-new-york-city-police-officers.html",
		"http://www.nytimes.com/2015/05/06/nyregion/guns-from-georgia-are-linked-to-another-new-york-officers-death.html",
		"http://www.nytimes.com/2015/10/27/nyregion/gun-fished-from-harlem-river-is-linked-to-officers-killing.html",
		"https://crimelab.uchicago.edu/",
		"http://www.sciencedirect.com/science/article/pii/S0091743515001486",
		"http://www.jhsph.edu/research/centers-and-institutes/johns-hopkins-center-for-gun-policy-and-research/",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_PostFacebookCom(t *testing.T) {
	article := Article{
		Domain:          "post.facebook.com",
		Title:           "Science - Spewings from Earth’s deep mantle reveal clues...",
		MetaDescription: "Spewings from Earth’s deep mantle reveal clues into the origin of our planet’s water. These findings serve as evidence for primordial water on Earth and...",
		CleanedText:     "Cookies help us to provide, protect and improve Facebook's services. By continuing to use our site, you agree to our cookie policy.",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "",
	}
	article.Links = []string{
		"https://www.facebook.com/help/cookies?fref=cub",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

// multiple og:image, according to http://ogp.me/, the first one should be preferred
func Test_ProfitLindorffFi(t *testing.T) {
	article := Article{
		Domain:          "profit.lindorff.fi",
		Title:           "Lindorff24.fi muuttaa maksujen hoidon mobiiliksi",
		MetaDescription: "Lindorffin verkkopalvelu kuluttajille tunnetaan nyt nimellä Lindorff24.fi. Uusien ominaisuuksien lisäksi palvelu on käytettävissä tietokoneen lisäksi älypuhelimella ja tabletilla. Verkon itseasioinnin uskotaan kasvavan lähivuosina merkittävästi nykyisestä.",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "http://profit.lindorff.fi/lindorff24-fi-muuttaa-maksujen-hoidon-mobiiliksi/",
		TopImage:        "http://profit.lindorff.fi/wp-content/uploads/2015/02/Iso_Lindorff24_2_600x2501.jpg",
	}
	article.Links = []string{
		"http://profit.lindorff.fi/teemat/uudistaja/",
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
	article.Links = []string{
		"http://www.ey.com/NO/no/Newsroom/PR-activities/PR-programs/Kurs-Bergen-2015-05-20-IFRS9",
		"http://www.ey.com/NO/no/Newsroom/PR-activities/PR-programs/Kurs-Bergen-2015-05-12-IFRSnyheter",
		"http://www.ey.no",
		"http://www.ey.no/publikasjoner",
		"http://www.ey.no/nyhetsbrev",
		"https://emeia.ey-vx.com/2180/26360/home/topics-of-interest.aspx?sid=51a92e43-8903-43bd-8cfd-8431639dfb5e",
		"https://emeia.ey-vx.com/2180/26360/home/contact-information.aspx?sid=51a92e43-8903-43bd-8cfd-8431639dfb5e",
	}

	err := ValidateArticle(article, &[]string{"~~~REMOVED~~~"})
	if err != nil {
		t.Error(err)
	}
}

func Test_TechCrunchCom(t *testing.T) {
	article := Article{
		Domain:          "techcrunch.com",
		Title:           "Gmail Will Soon Warn Users When Emails Arrive Over Unencrypted\u00a0Connections",
		MetaDescription: "Soon, you may see a warning in Gmail that tells you that an email has arrived over an unencrypted connection. Gmail already defaults to using HTTPS for the..",
		CleanedText:     "Soon, you may see a warning in Gmail that tells you that an email has arrived over an unencrypted connection.\n\nGmail already defaults to using HTTPS for the connections between your browser and its servers, but for the longest time, the standard practice for sending email between providers was to leave them unencrypted. If somebody managed to intercept those messages, it was pretty trivial to snoop on them.\n\nOver the last few years (and especially after the Snowden leaks), Google and other email providers started to change this and today, 57 percent of messages that users on other email providers send to Gmail are encrypted (and 81 percent of outgoing messages from Gmail are, too). Gmail-to-Gmail traffic is always encrypted.\n\nWhy does all of this matter? Unencrypted email makes for a great target. The good news is that email security is getting better. A joint research project\u00a0between Google, the University of Michigan, and the University of Illinois found that 94 percent of inbound messages to Gmail\u00a0can now be authenticated, which makes life harder for phishers. But at the same time, these researchers also found that there are “regions of the Internet actively preventing message encryption by tampering with requests to initiate SSL connections.”\n\nThe team also saw a number of malicious DNS servers that tried to intercept traffic. “These nefarious servers are like telephone directories that intentionally list misleading phone numbers for a given name,” the researchers write. “While this type of attack is rare, it’s very concerning as it could allow attackers to censor or alter messages before they are relayed to the email recipient.”\n\nGiven that there are still plenty of email servers that don’t support encryption, chances are you’ll see one or two of these warning labels in the next few months.",
		MetaKeywords:    "",
		CanonicalLink:   "http://techcrunch.com/2015/11/12/gmail-will-soon-warn-users-when-emails-arrive-over-unencrypted-connections/",
		TopImage:        "https://tctechcrunch2011.files.wordpress.com/2015/02/gmail-autocomplete.png?w=764&h=400&crop=1",
	}
	article.Links = []string{
		"https://googleonlinesecurity.blogspot.com/2015/11/new-research-encouraging-trends-and.html",
		"http://www.google.com/transparencyreport/saferemail/?hl=en",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_TheGuardianCom(t *testing.T) {
	article := Article{
		Domain:          "theguardian.com",
		Title:           "Thousands without power as storm Abigail forces school closures",
		MetaDescription: "Scotland is worst hit, with winds of up to 84mph, but rest of Britain can expect heavy, thundery showers on Friday",
		CleanedText:     "The Met Office has amber “be prepared” warnings in place for rain and wind in the north-west of Scotland, while yellow “be aware” warnings cover much of the rest of Scotland.\n\nThousands of homes lost power, dozens of schools were shut and bridges were closed to high-sided vehicles as storm Abigail brought gale-force winds of up to 84mph to northern Britain overnight.\n\nSorry, your browser is unable to play this video.\n\nA golfer struggles in the high winds and driving rain of Storm Abigail\n\nScotland has been worst hit by the strong gusts, which prompted a number of Met Office amber warnings, but the rest of the UK can expect heavy, thundery showers throughout the day as Britain’s first named storm sweeps its way down the country.\n\nThe Met Office, which said surface water and gusts could cause problems during rush hour, issued amber weather warnings for the Highlands, Orkney Islands and Shetland Islands; a yellow warning covered most of Scotland and part of the south-west of England and Wales.\n\nThe storm reached its peak in the early hours of this morning and while it was expected to ease throughout the day, it would be a slow process, forecaster Simon Partridge said. “It’s going to be a blustery, wet day for most parts and feel much cooler than it has done in recent weeks. Temperatures will be much closer to the average for this time of year and in Scotland it might even drop to a ‘feels-like’ temperature of around 1C (33.8F).”\n\nA number of Caledonian MacBrayne ferry sailings were cancelled before the storm and commuters on the trains and roads faced disruption. Western Isles council said every school and nursery in its area would be closed to pupils on Friday; schools would be open for teaching staff from 10am.\n\nShetland Islands council also announced that its schools would be shut to pupils due to the forecast of strong winds and lightning. Orkney Islands council said any decision on such closures would be taken on Friday morning.\n\nT he Met Office warned of likely gusts of 70-80mph, potentially reaching up to 90mph across exposed locations in the north-west of Scotland. The storm, which was expected to reach its height overnight, had already brought strengthening winds and heavy rain to many parts of Scotland.\n\nEmma Sharples, a Met Office meteorologist, said: “The main centre of the low pressure system around which all the winds are going to be strongest is moving from the Atlantic towards the north-west parts of Scotland at the moment. That’s going to continue to edge towards us.\n\n“There’s obviously rain already setting in and winds strengthening across the country and that will continue to be the case through the rest of this evening, with the band of rain spreading eastwards across Scotland and then the wind turning from a south westerly to more of a westerly as we go through towards midnight.”\n\nSharples said the Western Isles had experienced gusts of more than 55mph by mid-afternoon on Thursday. By 5pm, CalMac said 24 of its 26 ferry routes were disrupted. The company urged travellers to think carefully if they were planning to visit the west coast.\n\nScotRail said there was minor disruption on its routes from Glasgow to Carlisle/Newcastle, Glasgow to Ardrossan/Ayr/Largs and Kilmarnock to Ayr. The Forth road bridge has been closed to high-sided vehicles, cars with trailers, caravans, motorcycles, bicycles and pedestrians.\n\nHigh wind warnings are in place for key crossings, including the Erskine and Kessock bridges, and warnings of surface water have been issued for key commuter routes the M90 and M74.\n\nDublin airport said it was experiencing some minor disruption to flight schedules due to strong winds.\n\nMeanwhile, Dumfries and Galloway police said there are a number of trees down across the region. Traffic Scotland said a fallen tree on the A82 is partially blocking the road and affecting traffic in both directions. The Scottish Environment Protection Agency (Sepa) has flood alerts and warnings in place for Dumfries and Galloway, Argyll and Bute, Ayrshire and Arran, Skye and Lochaber, and Speyside.\n\nMembers of the public have been asked to secure any loose debris, while builders have been advised to secure scaffolding and any loose items on building sites. People are also being asked to look out for the elderly and vulnerable.\n\nThe Scottish Fire and Rescue service has urged people to take extra care if they are using candles during any power cuts. Scottish Hydro Electric Power Distribution said it had moved to yellow alert and had more than 500 workers in place in advance of the storm hitting.\n\nThe storm is the first such weather system affecting the country to merit a name as part of a Met Office project that invited the public to suggest names. Officials hope the initiative will help raise awareness of severe weather and ensure greater public safety.",
		MetaKeywords:    "Weather,Met Office,Scotland,UK news",
		CanonicalLink:   "http://www.theguardian.com/uk-news/2015/nov/12/storm-abigail-forces-school-closures-in-scotland",
		TopImage:        "https://i.guim.co.uk/img/static/sys-images/Guardian/Pix/pictures/2015/11/12/1447354267832/1a8862ca-64a6-492f-af0e-49a44523b360-2060x1236.jpeg?w=1200&q=85&auto=format&sharp=10&s=ded68c9fe6a7099fe1f2faf30130396f",
	}
	article.Links = []string{
		"http://www.theguardian.com/uk/met-office",
		"http://www.theguardian.com/world/video/2015/nov/13/storm-abigail-hits-irish-golf-course-video",
		"http://www.theguardian.com/uk/scotland",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_TwitterCom(t *testing.T) {
	article := Article{
		Domain:          "twitter.com",
		Title:           "Bloomberg Business on Twitter",
		MetaDescription: "",
		CleanedText:     "Language:\n\nHave an account?\n\nBloomberg Business\n\n@\n\nBloomberg Business\n\n@\n\nThe first word in business news.\n\nDisney is counting on 300 million tourists to flock to its new Shanghai theme park\n\nLoading seems to be taking a while.\n\nTwitter may be over capacity or experiencing a momentary hiccup. Try again or visit Twitter Status for more information.\n\nThis has already been marked as containing sensitive content.\n\nFlag this as containing potentially illegal content.\n\nList name\n\nDescription\n\nPublic · Anyone can follow this list\n\nPrivate · Only you can access this list\n\nThe URL of this tweet is below. Copy it to easily share with friends.\n\nAdd this Tweet to your website by copying the code below. Learn more\n\nAdd this video to your website by copying the code below. Learn more\n\nInclude parent Tweet\n\nInclude media\n\nForgot password?\n\nNot on Twitter? Sign up, tune into the things you care about, and get updates as they happen.\n\nSign up\n\nCountry\n\nCode\n\nFor customers of\n\nUnited States\n\n40404\n\n(any)\n\nCanada\n\n21212\n\n(any)\n\nUnited Kingdom\n\n86444\n\nVodafone, Orange, 3, O2\n\nBrazil\n\n40404\n\nNextel, TIM\n\nHaiti\n\n40404\n\nDigicel, Voila\n\nIreland\n\n51210\n\nVodafone, O2\n\nIndia\n\n53000\n\nBharti Airtel, Videocon, Reliance\n\nIndonesia\n\n89887\n\nAXIS, 3, Telkomsel, Indosat, XL Axiata\n\nItaly\n\n4880804\n\nWind\n\n3424486444\n\nVodafone\n\n» See SMS short codes for other countries\n\nHmm... Something went wrong. Please try again.",
		MetaKeywords:    "",
		CanonicalLink:   "https://twitter.com/business/status/665179987645964290",
		TopImage:        "https://pbs.twimg.com/media/CTsxrHVUkAAeOS2.jpg:large",
	}
	article.Links = []string{
		"https://t.co/pYwmmcd0bL",
		"http://status.twitter.com",
		"https://twitter.com/signup",
		"http://support.twitter.com/articles/14226-how-to-find-your-twitter-short-code-or-long-code",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_UsaTodayCom(t *testing.T) {
	article := Article{
		Domain:          "usatoday.com",
		Title:           "Social Security, Medicare changes are coming with new budget law",
		MetaDescription: "President Obama signed into law a bipartisan budget bill last week (Monday) that, among other things, changes — for better and worse — Social Security and Medicare laws. Here's a wrap-up.",
		CleanedText:     "Share This Story!Let friends in your social network know what you are reading aboutTwitterGoogle+LinkedInPinterestPosted!A link has been posted to your Facebook feed. Social Security, Medicare changes are coming with new budget lawPresident Obama signed into law a bipartisan budget bill last week that, among other things, changes —\u00a0for better and worse — Social Security and Medicare laws. Here's a wrap-up:•\u00a0File and suspend.\u00a0Currently,\u00a0a married person — typically the higher wage earner in a couple — who's\u00a0at least full retirement age could file for his or her own Social Security benefits and then immediately suspend those benefits while the spouse could file\u00a0for spousal benefits. By doing this, the higher wage earner’s benefits would grow 8% per year. In the meantime, the couple still get\u00a0a Social Security check, and down the road the surviving spouse could get a higher benefit.That option is ending for new filers starting May 1, 2016, so if you're\u00a0interested, now's the time to apply. People already using\u00a0this strategy will be grandfathered in until age 70.USA TODAYFull retirement age is a magic number for Social Security benefits•\u00a0Restricted application.\u00a0\u00a0This is also being phased out.\u00a0Currently, individuals\u00a0eligible for both a spousal benefit based their spouse's work record and a retirement benefit based on his or her own work record could choose to elect only a spousal benefit at full retirement age, according to Social Security Timing. That would let them collect a higher benefit later on.Under the new law, however, only those born Jan. 1, 1954, or earlier can use this option. Anyone younger will\u00a0just automatically get the larger of the two benefits,\u00a0according to Social Security Timing.•\u00a0Social Security Disability.\u00a0\u00a0The Social Security Disability trust was on pace\u00a0to run out money next year and, as a result, millions of Americans were going to receive an automatic 19% reduction in their disability benefits in the fourth quarter of 2016. The new law fixes that\u00a0by shifting payroll tax revenue from one Social Security trust fund —\u00a0the Old-Age and Survivors Insurance Trust fund —\u00a0to another,\u00a0the Disability Insurance Trust fund.USA TODAYRetirement: When you should take Social Security•\u00a0Medicare Part B.\u00a0Some 30% of Medicare beneficiaries were expecting a 52% increase in their Medicare Part B medical insurance premiums and deductible\u00a0in 2016.\u00a0Under the new law, those beneficiaries —\u00a0an estimated 17 million Americans —\u00a0will pay about $119\u00a0per month, instead of $159.30, for Part B. (Some 70% of Medicare beneficiaries will continue to pay the same premium in 2016 as they did in 2015, $104.90.)Beneficiaries, however, will also have to pay an extra $3 per month to help pay down a loan the government gave to Medicare to offset lost revenue.\u00a0 Plus, all Part B beneficiaries will see their annual deductible increase by 15% to about $166\u00a0in 2016.Robert Powell is editor of Retirement Weekly, contributes regularly to MarketWatch, The Wall Street Journal, USA TODAY, and teaches at Boston University.",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.usatoday.com/story/money/columnist/powell/2015/11/12/social-security-medicare-changes-budget-law-retirement/75164246/",
		TopImage:        "http://www.gannett-cdn.com/-mm-/eba3ab7ada1c4fcc1a671898ecfb68274260e9c9/c=0-48-508-335&r=x633&c=1200x630/local/-/media/2015/02/24/USATODAY/USATODAY/635603784536631512-177533853.jpg",
	}
	article.Links = []string{
		"https://twitter.com/intent/tweet?url=http%3A//usat.ly/1iXh7zM&text=Social%20Security%2C%20Medicare%20changes%20are%20coming%20with%20new%20budget%20law&via=usatoday",
		"http://www.linkedin.com/shareArticle?url=http%3A//usat.ly/1iXh7zM&mini=true",
		"https://twitter.com/intent/tweet?url=http%3A//usat.ly/1iXh7zM&text=Social%20Security%2C%20Medicare%20changes%20are%20coming%20with%20new%20budget%20law&via=usatoday",
		"http://www.linkedin.com/shareArticle?url=http%3A//usat.ly/1iXh7zM&mini=true",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_WashingtonpostCom(t *testing.T) {
	article := Article{
		Domain:          "washingtonpost.com",
		Title:           "The 7 big things on President Obama’s to-do list, with one year to go",
		MetaDescription: "A checklist of the political battles ahead -- only some of which involved Congress.",
		CleanedText:     "For President Obama, it's legacy time. With\u00a0less than a year before his successor is elected and he officially becomes a\u00a0lame-duck president, time is running short. Obama has moved the ball forward on a number of legacy items already this year. Some have solidified; others remain in limbo. His 2010 health care reform law will already be mentioned at the top of the 44th president's Wikipedia page. But the Obama White House is moving quickly\u00a0on a number of issues that could be\u00a0listed\u00a0in the\u00a0first few paragraphs, too. \"You do get a sense they are aware of the legacy, and there is a kind of a presidential scorecard being filled out,\" says Gil Troy, a visiting fellow at the Brookings Institution. Obama's ambitions are high. They start with one last shot at the seemingly impossible task of closing the prison in Guantanamo Bay, Cuba, and cover everything from signing an international climate change deal to finalizing\u00a0one of the world's largest free-trade agreements in a generation. It's notable that most\u00a0of Obama's goals are abroad; that's because a Republican-controlled Congress has less authority to intervene.\u00a0But that doesn't mean crossing things off his final to-do list is going to be easy. The scope of what he wants to do means finishing it will take a\u00a0lot of late nights for Obama and his staff in their\u00a0final year, said Jacob Stokes, an associate fellow at the bipartisan Center for a New American Security. But if the stars align for the president — as they seem to have done this summer — Stokes thinks Obama can get most of them done.\u00a0\"The president and the administration have a relatively large amount of agency to get these things done,\" he said. \"If they really focus on it.\" Here are seven things on Obama's final to-do list. Shuttering Guantanamo is less of a legacy issue and more of a moral one for the president, Stokes said. Since the first days of his presidency, Obama has maintained the prison, where men can be held indefinitely, is a propaganda tool for terrorists. But congressional Republicans say closing it will create more risk than it's worth, and they — and the realities of what to do with existing prisoners there — have successfully blocked the president for six years from doing anything about it. The clock's ticking for Obama to fulfill one of his oldest\u00a0campaign promises. He's\u00a0planning a final standoff with Congress, by dropping a\u00a0plan as soon as this week to close Guantanamo without Congress's help. (On Tuesday, Congress passed its annual defense spending bill that, per usual, restricts the president from transferring detainees to the United States.) But Obama must decide how badly he wants Guantanamo closed. Trying to transfer the remaining 112 prisoners by himself could start a much broader fight with Republicans over the president's constitutional power. Sen. John McCain (R-Ariz.) is threatening to sue the president if he acts over Congress's wishes. Obama already scored a major legislative victory this summer when he persuaded\u00a0enough congressional Democrats — yes, he was working against much of his own party on this one — to give him authority to negotiate the Trans-Pacific Partnership without congressional say-so on\u00a0every little detail. His job is only half done, though. After the United States and 11 other nations came to an agreement on the deal in October, Obama now needs to sway enough lawmakers on both sides to approve the whole package. Lawmakers will soon review the sprawling deal and could vote on it this spring at the earliest. Getting it passed is going to be an uphill battle for Obama, reports The Washington Post's David Nakamura. Liberal Democrats are concerned about the trade deal's environmental impacts and potential drag on U.S. manufacturing jobs, while some Republicans worry the deal isn't strong enough. The stakes are also higher for Obama than simply completing\u00a0the biggest free-trade deal in modern history. This trade agreement is a major economic cornerstone of Obama's pivot to Asia, Stokes said. Without TPP, he'll lose one of his most concrete examples\u00a0of a shift to Asia that has struggled to take shape. Here's one place Obama might\u00a0not need to do battle with Congress. Whatever comes out of a major United Nations summit on climate change held in Paris at the end of this month will likely not have to ratified by the Senate. That's a good thing, because Obama didn't make any friends on Capitol Hill on Friday when he announced he won't approve an extension to the Keystone XL oil pipeline from Canada. In bucking a Republican priority and taking the environmentalists' side, Obama indicated he's ready to make some serious changes to U.S.\u00a0pollution levels, which rank among the top in the world. Rejecting Keystone demonstrates to the rest of the world that Obama is \"willing to pull out all the stops on climate change,\" writes The Washington Post's Chris Mooney. Plus, a 2009 climate change meeting of world leaders in Copenhagen was kind of a bust, so Paris could be Obama's last chance to effect any meaningful change on the world stage. \"If they miss putting something together in Paris, it's going to be very tough to do anything beyond that,\" Stokes said. Obama is leaning in on Syria in his final months in office. In addition to stepping up airstrikes on the Islamic State there, he announced in October he's putting 50 Special Ops forces\u00a0on the ground, appearing to go back on his past statements\u00a0he wouldn't commit ground troops to Syria. This is happening as\u00a0Russia jumped into the\u00a0Syrian conflict, aiming to help Syrian President Bashar al-Assad keep control of his crumbling country. That's a major problem for the White House, which is facing an already no-win situation in an increasingly violent Middle East. But Russia's sudden prioritization of the region could be just the crisis the world and Obama need to find a solution, Stokes said. Already, it\u00a0forced diplomats with stakes in the region to a hastily organized conference in Vienna last week to talk about what to do, he noted. \"I think there's a sense that a political agreement is not imminent by any stretch of the imagination,\" Stokes said, \"but that in the next year you may get parties into a region where they can start thinking more broadly.\" Either way, Obama would really like to leave office without a civil war in Syria — something which\u00a0is fueling the Islamic State's movement — still raging. Yes, Obama announced a historic nuclear agreement with Iran in June, and yes, he managed to avoid a reluctant Congress from blocking it in September. But the deal is still mostly on-paper, which means several GOP presidential candidates' campaign promises to \"rip it to shreds\" — as Sen. Ted Cruz (Tex.) likes to say — are legitimate threats. That is, unless Obama can spend the next year or so setting\u00a0key elements of the deal in motion. That would make\u00a0it much tougher for another president to come along and undo one of his biggest foreign policy achievements, Stokes said. These days, Obama and Republicans celebrate when they can agree on a budget just\u00a0to keep the government running. So there's little hope they'll come to an agreement on the president's other major domestic policy goals, like immigration reform. One bright spot is criminal justice reform. A bipartisan bill to change federal sentencing mandates is moving quickly in the Senate and has the potential for bipartisan support in the House of Representatives, too. Obama has made reforming sentencing laws a priority recently.\u00a0In 2014, then-Attorney General Eric Holder announced the department would stop charging nonviolent drug offenders with crimes that require judges to enact so-called mandatory minimum sentences. And the Justice Department recently released 6,000 federal prisoners, the largest one-time release ever, who were sentenced for non-violent drug crimes. Not just any court case, mind you. After 26 states challenged his executive actions on immigration, Obama is betting it all on the Supreme Court. A federal court upheld the states' challenge on Monday, and by Tuesday, the White House confirmed it would ask the Supreme Court to rule next year on whether he stayed within his constitutionally limited powers by deferring deportations for millions of young immigrants and some of their parents. If the\u00a0Supreme Court takes up the case, it could rule by June, leaving just months for the administration to start enrolling immigrants and create a buffer for whoever comes into the White House next — and whatever vision of Obama's they might\u00a0try to undo.",
		MetaKeywords:    "Obama; Syria; TPP; immigration; Iran; Guantanamo",
		CanonicalLink:   "https://www.washingtonpost.com/news/the-fix/wp/2015/11/12/the-7-big-things-on-president-obamas-to-do-list-with-one-year-to-go/",
		TopImage:        "http://img.washingtonpost.com/rf/image_908w/2010-2019/WashingtonPost/2015/11/11/National-Politics/Images/05020997.jpg",
	}
	article.Links = []string{
		"https://www.washingtonpost.com/news/the-fix/wp/2015/11/06/by-nixing-the-keystone-pipeline-obama-finalizes-the-third-facet-of-his-legacy/",
		"https://www.washingtonpost.com/news/the-fix/wp/2015/07/09/the-role-of-congress-or-lack-thereof-in-the-iran-deal-explained/",
		"https://www.washingtonpost.com/news/the-fix/wp/2015/11/12/just-8-percent-of-the-gop-likes-the-gop-controlled-congress-thats-bad-for-paul-ryan-and-great-for-trump-and-carson/",
		"http://www.politico.com/story/2015/11/guantanamo-gitmo-john-mccain-barack-obama-constitution-executive-orders-215779",
		"https://www.washingtonpost.com/business/economy/deal-reached-on-pacific-rim-trade-pact/2015/10/05/7c567f00-6b56-11e5-b31c-d80d62b53e28_story.html",
		"https://www.washingtonpost.com/politics/obama-aims-to-reinvigorate-asia-strategy/2014/04/16/4a46ed5e-c4bf-11e3-bcec-b71ee10e9bc3_story.html",
		"https://www.washingtonpost.com/politics/obama-aims-to-reinvigorate-asia-strategy/2014/04/16/4a46ed5e-c4bf-11e3-bcec-b71ee10e9bc3_story.html",
		"https://www.washingtonpost.com/news/post-politics/wp/2015/11/06/obama-set-to-reject-keystone-xl-project-citing-climate-concerns/",
		"https://www.washingtonpost.com/news/energy-environment/wp/2015/11/06/how-obamas-keystone-xl-rejection-gives-him-momentum-for-the-paris-climate-talks/",
		"https://www.washingtonpost.com/politics/obame-decides-on-small-special-operations-force-for-syria/2015/10/30/a8f69c0e-7f13-11e5-afce-2afd1d3eb896_story.html",
		"https://www.washingtonpost.com/news/the-fix/wp/2015/10/30/5-times-president-obama-said-there-would-be-no-ground-troops-or-no-combat-mission-in-syria/",
		"https://www.washingtonpost.com/news/the-fix/wp/2015/09/10/senate-democrats-just-pinned-a-bow-on-obamas-iran-deal/",
		"https://www.washingtonpost.com/news/the-fix/wp/2015/07/31/why-the-iran-deal-is-huge-for-obamas-legacy/",
		"https://www.washingtonpost.com/world/national-security/justice-department-about-to-free-6000-prisoners-largest-one-time-release/2015/10/06/961f4c9a-6ba2-11e5-aa5b-f78a98956699_story.html",
		"https://www.washingtonpost.com/politics/obama-administration-seeks-supreme-court-involvement-in-immigration-case/2015/11/10/ce13d802-87bb-11e5-9a07-453018f9a0ec_story.html",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_WordpressCom(t *testing.T) {
	article := Article{
		Domain:          "wordpress.com",
		Title:           "Raj's Lab",
		MetaDescription: "Strategy and Entrepreneurship (by Raj Shankar)",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "https://secure.gravatar.com/blavatar/749e313b7d7ba65e9f0d0fabb2b5fd36?s=200&ts=1447426008",
	}
	//article.Links = []string{""}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_WsjCom(t *testing.T) {
	article := Article{
		Domain:          "wsj.com",
		Title:           "Big Obama Donors Stay on Sidelines in 2016 Race",
		MetaDescription: "President Obama’s biggest campaign donors are mostly sitting on the sidelines of the 2016 Democratic presidential primary so far, not opening their wallets in support of Hillary Clinton or Bernie Sanders.",
		CleanedText:     "WASHINGTON—President Barack Obama’s biggest campaign donors are mostly sitting on the sidelines of the 2016 Democratic presidential primary so far, not opening their wallets in support of Hillary Clinton or Bernie Sanders. Almost four-fifths of the people who gave the 2012 maximum $5,000 to the president’s re-election committee hadn’t donated to a presidential candidate by Oct. 1, a Wall Street Journal analysis of federal campaign finance records found. In interviews ahead of this Saturday’s Democratic debate in Iowa, donors said Mrs. Clinton, the party’s front-runner, hadn’t motivated them to give the way Mr. Obama and previous Democratic candidates had. Still others said they are put off by the larger role of super PACs and that their donations to candidates, which are limited in this election cycle to $5,400 for the eventual nominee, just don’t matter much anymore.\n\nRelated Some Candidates, Super PACs Draw Closer (Oct. 25) The 2016 Money Race (Oct. 15)\n\n“I’m just not ready for Hillary yet,” said Robert Finnell, a Rome, Ga., lawyer who gave the maximum allowed contribution to Mr. Obama’s 2008 and 2012 campaigns and gave significant sums to 2008 hopeful John Edwards and 2004 Democratic nominee John Kerry. “It’s not that I don’t think she’s competent—she is competent, she’s just hard to like.” The donors’ reluctance could be a troubling trend for Mrs. Clinton. They are some of the easiest prospective contributors to identify, given that their names are on Mr. Obama’s campaign disclosure reports, and that they’ve already made a habit of cutting checks to politicians. Julianna Smoot, finance director on President Obama’s 2008 campaign and deputy campaign manager of his re-election effort, said: “Most Democrats will be behind Hillary if she’s the nominee. Once that becomes clear, the rest of that money should be easy for her to get. I do think these folks will be there.” Mrs. Clinton has outpaced Mr. Obama’s fundraising in the first two quarters of his initial presidential campaign. The former secretary of state has raised $77.5 million for those six months through October. By July 2007, when Mr. Obama had been in the race a comparable length of time, he had raised $58.9 million. In 2012, roughly 4,000 individuals donated the maximum to Mr. Obama’s campaign committee, delivering $20 million to his account, according to disclosure reports filed with the Federal Election Commission. Of them, about 830 can be identified as having donated to a candidate in the 2016 presidential race. Mrs. Clinton is the largest recipient of their money at $1.8 million. The big Obama donors gave about $109,000 to Mr. Sanders, about $94,000 to former Maryland Gov. Martin O’Malley, and about $70,000 to Republican Jeb Bush.\n\nFor the analysis, The Wall Street Journal cross-referenced a list of individuals who had donated the maximum amount to Mr. Obama in 2012 with those who have given to candidates in the current presidential race. The maximum donation total is based on rules that allow a donor to give a candidate up to $2,700 each for the 2016 primary and general election. Michael Briggs, a spokesman for the campaign of Mr. Sanders, said he expected to be outspent and that “he’s taking on the establishment and does not expect the establishment to support that.” The Clinton campaign’s spokesman, Josh Schwerin, said: “Thanks to the support of hundreds of thousands of people, we have been able to raise a record amount for a nonincumbent during our first two quarters in the race.” Some people inclined to support Mrs. Clinton note that it is still early in the race and the Republican field remains unsettled. “I don’t think she needs the money right now,” said Jeff Choney, a retired high-school teacher in Wellesley, Mass., who gave $5,000 to Mr. Obama in 2012 and said he may contribute to Mrs. Clinton’s later in the cycle. “I like Bernie Sanders—he speaks the truth on a lot of things. But I don’t think he has a chance of beating her, so I’m not so worried about her campaign.” Still, the Clinton campaign is building a national campaign apparatus that will be expensive to maintain through the general election, should she win the party nomination. Mr. Obama built a similar operation incrementally during the extended 2007 Democratic primary contest. As of Oct. 1, the last records available to the public, Hillary for America had spent $44.5 million, compared with $14.3 million for Mr. Sanders and $20.1 million for retired neurosurgeon Ben Carson, the current fundraising leader in the Republican field. Mrs. Clinton is also relying on support from super PACs, which can raise and spend unlimited sums as long as they don’t coordinate with her campaign. One of the largest of those, Priorities USA Action, raised $15.7 million as of July. The limitations of super PACs have been on display in the GOP primary. Former Texas Gov. Rick Perry and Wisconsin Gov. Scott Walker dropped out of the race after struggling to pay campaign expenses that can’t be covered by an outside group. Though a super PAC backing former Florida Gov. Jeb Bush raised $103.2 million as of July, his campaign in October still had to cut staff salaries and trim the head count at its Miami headquarters. The super PAC restrictions put pressure on Mrs. Clinton—and all candidates—to raise as much money as they can for their own campaign accounts. Meanwhile, some donors are demoralized witnessing the big checks pouring into super PACs. “Even though I gave the maximum [in 2012], it’s nothing compared with what these PACs do. I certainly don’t see my contribution as significant,” said Marilynn Duker, president of a Baltimore residential development and property management company. She gave to Mr. Obama’s 2008 campaign and Mr. Kerry in 2004, but has yet to donate to a White House hopeful in this cycle. “It has no real meaning relative to the gazillions of dollars that the PACs contribute to the races these days. I just don’t feel like the individual really makes a difference,” she said. Doug Curling, an Atlanta-area executive, gave significant sums to Mr. Obama’s campaigns and gave Mrs. Clinton the maximum donation in 2007. But he said he and his wife now plan to contribute to groups advocating for structural change in the political system. “Nobody needs our money,” Mr. Curling said. “I wouldn’t misinterpret it as we’re disenfranchised from our party, it’s more we’re disenfranchised from the system.” Peter Maroney, who was the Democratic National Committee’s national finance co-chairman in the 2004 presidential campaign, said many Democratic donors had been waiting on Vice President Joe Biden. “Now that the vice president has made his decision, this is an opportunity for candidates like Mrs. Clinton to proactively go after these donors and make them feel that they have a seat at her table,” he said. Write to Daniel Nasaw at daniel.nasaw@wsj.com",
		MetaKeywords:    "democratic primary,donors,election 2016,fundraising,maximum donation,super pacs,political,general news,politics,international relations,domestic politics,elections,national,presidential elections",
		CanonicalLink:   "http://www.wsj.com/articles/big-obama-donors-stay-on-sidelines-in-2016-race-1447375429",
		TopImage:        "http://si.wsj.net/public/resources/images/BN-LF842_OBADON_D_20151112192304.jpg",
	}
	article.Links = []string{
		"http://topics.wsj.com/person/O/Barack-Obama/4328",
		"http://topics.wsj.com/person/C/Hillary-Clinton/6344",
		"http://graphics.wsj.com/elections/2016/campaign-finance/",
		"http://www.wsj.com/articles/some-candidates-super-pacs-draw-closer-1445809990",
		"http://www.wsj.com/articles/some-candidates-super-pacs-draw-closer-1445809990",
		"http://graphics.wsj.com/elections/2016/campaign-finance/",
		"http://topics.wsj.com/person/E/John-Edwards/6600",
		"http://topics.wsj.com/person/K/John-Kerry/7196",
		"http://topics.wsj.com/person/B/Jeb-Bush/8217",
		"http://topics.wsj.com/person/P/Rick-Perry/5983",
		"http://topics.wsj.com/person/B/Joe-Biden/6352",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_YahooCom(t *testing.T) {
	article := Article{
		Domain:          "yahoo.com",
		Title:           "El Nino sends rare tropical visitors to California waters",
		MetaDescription: "From Yahoo News: By Alex Dobuzinskis LOS ANGELES (Reuters) - El Nino's warm currents have brought fish in an unexpected spectrum of shapes and colors from Mexican waters to the ocean off California's coast, thrilling scientists with the sight of bright tropical species and giving anglers the chance of a once-in-a-lifetime big catch. Creatures that have made a splash by venturing north in the past several weeks range from a whale shark, a gentle plankton-eating giant that ranks as the world's largest fish and was seen off Southern California, to two palm-sized pufferfish, a species with large and endearing eyes, that washed ashore on the state's central coast. Scientists say El Nino, a periodic warming of ocean surface temperatures in the eastern and central Pacific, has sent warm waves to California's coastal waters that make them more hospitable to fish from the tropics.",
		CleanedText:     ".\n\nView photo\n\nLOS ANGELES (Reuters) - El Nino's warm currents have brought fish in an unexpected spectrum of shapes and colors from Mexican waters to the ocean off California's coast, thrilling scientists with the sight of bright tropical species and giving anglers the chance of a once-in-a-lifetime big catch. Creatures that have made a splash by venturing north in the past several weeks range from a whale shark, a gentle plankton-eating giant that ranks as the world's largest fish and was seen off Southern California, to two palm-sized pufferfish, a species with large and endearing eyes, that washed ashore on the state's central coast. Scientists say El Nino, a periodic warming of ocean surface temperatures in the eastern and central Pacific, has sent warm waves to California's coastal waters that make them more hospitable to fish from the tropics. El Nino is also expected to bring some relief to the state's devastating four-year drought by triggering heavy rains onshore. But so far precipitation has been modest, and researchers say the northern migration of fish in the Pacific Ocean has been one of the most dynamic, albeit temporary, effects of the climate phenomenon. Even as marine biologists up and down the coast gleefully alert one another to each new, rare sighting, the arrival of large numbers of big fish such as wahoo and yellowtail has also invigorated California's saltwater sport fishing industry, which generates an estimated $1.8 billion a year. \"Every tropical fish seems to have punched their ticket for Southern California,\" said Milton Love, a marine science researcher at the University of California, Santa Barbara. Some fish made the journey north as larva, drifting on ocean currents, before they grew up, researchers said. The first ever sighting off California's coast of a largemouth blenny fish was made over the summer near San Diego, said Phil Hastings, a curator of marine vertebrates at the Scripps Institution of Oceanography. That species had previously only been seen further south, he said, off Mexico's Baja California. Small, colorful cardinalfish were also spotted this year off San Diego, while spotfin burrfish, a rounded and spiny species, were sighted off the coast of Los Angeles, said Rick Feeney, a fish expert at the Natural History Museum of Los Angeles County. Those tropical species are hardly ever found in Californian waters, he said. 'NEVER SEEN IT LIKE THIS' Some small tropical fish could remain in the state's waters over the coming months, researchers said, as El Nino is expected to last until early next year. \"As soon as the water gets cold, or as soon as they get eaten by something else, we'll never see them again,\" Love said. For sports fishers, it was so-called pelagic zone fish like wahoo, that live neither close to the bottom nor near the shore, which made this year special. Before the El Nino, California anglers only saw wahoo, a fish with a beak-like snout and a slim body that often measures more than 5 feet (1.5 meters) in length, when they made boat trips south to Mexican waters. This year, there were 256 recorded catches of wahoo by sport fishing party boats from Southern California, with almost all of those being taken on the U.S. side of the border, said Chad Woods, founder of the tracking company Sportfishingreport.com. Last year, he said, the same boats made 42 wahoo catches. Michael Franklin, 56, a dock master for Marina Del Rey Sportfishing near Los Angeles in the Santa Monica Bay, said this was the best year he can remember, with plentiful catches of yellowtail and marlin. \"I've been fishing this bay all my life since I was old enough to fish, and I've never seen it like this,\" he said. Many hammerhead sharks also cruised into Californian waters because of El Nino, experts say. Sport fisherman Rick DeVoe, 46, said he took a group of children out in his boat off the Southern California coast this September. A hammerhead followed them, chomping in half any tuna they tried to reel in. \"The kids were freaking out because the shark's going around our boat like 'Jaws',\" DeVoe said. (Reporting by Alex Dobuzinskis; Editing by Daniel Wallis and Andrew Hay)",
		MetaKeywords:    "",
		CanonicalLink:   "http://news.yahoo.com/el-nino-sends-rare-tropical-visitors-california-waters-110532667.html",
		TopImage:        "https://s1.yimg.com/bt/api/res/1.2/q3ifY_wb94kl1PKV7QJ8UQ--/YXBwaWQ9eW5ld3NfbGVnbztpbD1wbGFuZTtxPTc1O3c9NjAw/http://media.zenfs.com/en_us/News/Reuters/2015-11-13T110532Z_1_LYNXNPEBAC0IO_RTROPTP_2_CALIFORNIA-ELNINO-FISH.JPG",
	}
	//article.Links = []string{""}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_YouTubeCom(t *testing.T) {
	article := Article{
		Domain:          "youtube.com",
		Title:           "WTF (Where They From) ft. Pharrell Williams [Official Video]",
		MetaDescription: "Missy Elliott's new single \"WTF (Where They From)\" ft. Pharrell Williams available now! Download: http://smarturl.it/WTFdownload Stream: http://smarturl.it/W...",
		CleanedText:     "",
		MetaKeywords:    "music, official, music video, Missy Elliott (Musical Artist), Pharrell Williams (Celebrity), Hip Hop Music (Musical Genre), WTF, Where They From, Dave Meyers...",
		CanonicalLink:   "https://www.youtube.com/watch?v=KO_3Qgib6RQ",
		TopImage:        "https://i.ytimg.com/vi/KO_3Qgib6RQ/hqdefault.jpg",
	}
	//article.Links = []string{""}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}
