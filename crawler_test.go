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
	//g.config.debug = true
	result, err := g.ExtractFromRawHTML(expected.FinalURL, ReadRawHTML(expected))
	if nil != err {
		return err
	}

	// DEBUG
	//fmt.Printf("article := Article{\n\tDomain:          %q,\n\tTitle:           %q,\n\tMetaDescription: %q,\n\tCleanedText:     %q,\n\tMetaKeywords:    %q,\n\tCanonicalLink:   %q,\n\tTopImage:        %q,\n}\n\n", expected.Domain, result.Title, result.MetaDescription, result.CleanedText, result.MetaKeywords, result.CanonicalLink, result.TopImage)
	//fmt.Printf("%#v\n", result.Links)

	if result.Title != expected.Title {
		return fmt.Errorf("article title does not match. Got '%q', Expected '%q'", result.Title, expected.Title)
	}

	if result.MetaDescription != expected.MetaDescription {
		return fmt.Errorf("article metaDescription does not match. Got '%q', Expected '%q'", result.MetaDescription, expected.MetaDescription)
	}

	if !strings.Contains(result.CleanedText, expected.CleanedText) {
		fmt.Printf("EXPECTED:       %s \n\n\n\nACTUAL:    %s\n\n", expected.CleanedText, result.CleanedText)
		return fmt.Errorf("article cleanedText does not contain %q", expected.CleanedText)
	}

	// check if the specified strings where properly removed
	for _, rem := range *removed {
		if strings.Contains(result.CleanedText, rem) {
			return fmt.Errorf("article cleanedText contains %q", rem)
		}
	}

	if result.MetaKeywords != expected.MetaKeywords {
		return fmt.Errorf("article keywords does not match. Got %q\n Expected: %q", result.MetaKeywords, expected.MetaKeywords)
	}
	if result.CanonicalLink != expected.CanonicalLink {
		return fmt.Errorf("article CanonicalLink does not match. Got %q, Expected '%q'", result.CanonicalLink, expected.CanonicalLink)
	}

	if result.TopImage != expected.TopImage {
		return fmt.Errorf("article topImage does not match. Got %q, Expected %q", result.TopImage, expected.TopImage)
	}

	if expected.Links != nil && !reflect.DeepEqual(result.Links, expected.Links) {
		return fmt.Errorf("article Links do not match. Got %#v, Expected %#v", result.Links, expected.Links)
	}

	return nil
}

func Test_AbcNewsGoCom(t *testing.T) {
	article := Article{
		Domain:          "abcnews.go.com",
		Title:           "New Jersey Devils Owner Apologizes After Landing Helicopter in Middle of Kids' Soccer Game Forces Cancellation",
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
		Title:           "Crunch talks on new Greek bailout under way",
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
		CleanedText:     "Ministers are considering whether homeopathy should be put on a blacklist of treatments GPs in England are banned from prescribing, the BBC has learned.\n\nThe controversial practice is based on the principle that \"like cures like\", but critics say patients are being given useless sugar pills.\n\nThe Faculty of Homeopathy said patients supported the therapy.\n\nA consultation is expected to take place in 2016.\n\nThe total NHS bill for homeopathy, including homeopathic hospitals and GP prescriptions, is thought to be about £4m.\n\nHomeopathy is based on the concept that diluting a version of a substance that causes illness has healing properties.\n\nSo pollen or grass could be used to create a homeopathic hay-fever remedy.\n\nOne part of the substance is mixed with 99 parts of water or alcohol, and this is repeated six times in a \"6c\" formulation or 30 times in a \"30c\" formulation.\n\nThe end result is combined with a lactose (sugar) tablet.\n\nHomeopaths say the more diluted it is, the greater the effect. Critics say patients are getting nothing but sugar.\n\nCommon homeopathic treatments are for asthma, ear infections, hay-fever, depression, stress, anxiety, allergy and arthritis.\n\nSource: British Homeopathic Association\n\nBut the NHS itself says: \"There is no good-quality evidence that homeopathy is effective as a treatment for any health condition.\"\n\nWhat do you think about homeopathic treatments? Join our Facebook Q&A on Friday 13th November from 3pm, on the BBC News Facebook page, with the BBC website's health editor, James Gallagher.\n\nThe Good Thinking Society has been campaigning for homeopathy to be added to the NHS blacklist - known formally as Schedule 1 - of drugs that cannot be prescribed by GPs.\n\nDrugs can be blacklisted if there are cheaper alternatives or if the medicine is not effective.\n\nAfter the Good Thinking Society threatened to take their case to the courts, Department of Health legal advisers replied in emails that ministers had \"decided to conduct a consultation\".\n\nOfficials have now confirmed this will take place in 2016.\n\nSimon Singh, the founder of the Good Thinking Society, said: \"Given the finite resources of the NHS, any spending on homeopathy is utterly unjustifiable.\n\n\"The money spent on these disproven remedies can be far better spent on treatments that offer real benefits to patients.\"\n\nBut Dr Helen Beaumont, a GP and the president of the Faculty of Homeopathy, said other drugs such as SSRIs (selective serotonin reuptake inhibitors) for depression would be a better target for saving money, as homeopathic pills had a \"profound effect\" on patients.\n\nShe told the BBC News website: \"Patient choice is important; homeopathy works, it's widely used by doctors in Europe, and patients who are treated by homeopathy are really convinced of its benefits, as am I.\"\n\nThe result of the consultation would affect GP prescribing, but not homeopathic hospitals which account for the bulk of the NHS money spent on homeopathy.\n\nEstimates suggest GP prescriptions account for about £110,000 per year.\n\nAnd any decision would not affect people buying the treatments over the counter or privately.\n\nHealth Secretary Jeremy Hunt was criticised for supporting a parliamentary motion on homeopathy, but in an interview last year argued \"when resources are tight we have to follow the evidence\".\n\nMinister for Life Sciences, George Freeman, told the BBC: \"With rising health demands, we have a duty to make sure we spend NHS funds on the most effective treatments.\n\n\"We are currently considering whether or not homeopathic products should continue to be available through NHS prescriptions.\n\n\"We expect to consult on proposals in due course.\"",
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
		CleanedText:     "Susan Brown, owner of Los Angeles gardening store Potted, recently updated her business listing on Google. Susan says, “Putting your business on Google lets people find you easily. Your directions are right there, your hours are right there, what you sell is right there.”\n\nThanks to her decision, Susan has seen more customers walk through her door: “So many of the customers that come in here find us on Google. As a small business, you want to use every opportunity to help your business grow.”\n\nNational Small Business Week is one of those opportunities. So from May 4-8, instead of three cheers, we’re giving you five—five simple ways to get your small business online and growing.\n\nCelebrating National Small Business Week with Google\n\nA handful of bright ideas and quick-fixes, all five ways are doable in a week or less and will help you throw a digital spotlight on your business all year round.\n\n1. SHOW UP ON GOOGLE\n\nCheck to see how your business shows up on Google. Then, claim your listing so that customers can find the right info about your business on Google Search and Maps. When you claim your listing this week: You could be one of 100 randomly selected businesses to get a 360° virtual tour photoshoot—a $255 value.\n\n2. LEARN FROM PROS & PEERS\n\nGet business advice from experts and colleagues in the Google Small Business Community. They're ready to chat! When you visit or join this week: Share your tips for summertime business success and we'll feature your tip in front of an audience of 400K members.\n\n3. WORK BETTER, TOGETHER:\n\nWith professional email, calendars, and docs that you can access anywhere, Google Apps for Work makes it easy for your team to create and collaborate. When you sign up this week you’ll receive 25% off Google Apps for Work for one year.\n\n4. CLAIM YOUR DOMAIN:\n\nWith a custom domain name and website, Google Domains helps you create a place for your business on the web. When you sign up and purchase a .co, .com or .company domain this week you could be one of 1,500 randomly selected businesses to get reimbursed for the first year of registration.\n\n5. GET ADVICE FROM AN ADVERTISING PRO:\n\nLearn how you can promote your business online and work with a local digital marketing expert to craft a strategy that’s right for your business goals. When you RSVP this week you’ll get help from an expert who knows businesses like yours.\n\nWhile these resources are available year-round, there’s no better time to embark on a digital reboot.\n\nFor more information, visit google.com/smallbusinessweek.\n\nWishing everyone a happy and productive Small Business Week!\n\nPS: To join the conversation, use #5Days5Ways and #SBW15 on G+, Facebook or Twitter.",
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
		CleanedText:     "British lawmakers need to take “urgent action” to ensure the U.K. maintains its position as the leading global financial center or risk the departure of banks to cities such as Singapore and Hong Kong, according to the British Bankers’ Association.\n\nNew regulations, taxes and depressed economic activity in Europe have resulted in an 8 percent drop in British banking jobs, with two-thirds of BBA members saying they’ve moved business elsewhere since 2010, the lobby group said in a report Friday. The BBA recommends a softening of the law separating retail operations from investment banking, further tax cuts and a reworking of visa limits to make it easier to hire from abroad.\n\n“We have now reached a watershed moment in Britain’s competitiveness as an international banking center” and “many international banks have been moving jobs overseas or deciding not to invest in the U.K.,” BBA Chief Executive Officer Anthony Browne said in the report. “Wholesale banking is an internationally mobile industry and there is a real risk this decline could accelerate.”\n\nChancellor of the Exchequer George Osborne, 44, outlined a “new settlement” for the City of London in a speech in June, pledging to curtail huge fines and amend regulations to “get the balance right.” As memories fade of the 1 trillion pounds ($1.5 trillion) of U.K. taxpayer support given to banks amid the 2008 crisis, this year the government has backed down on some issues after lobbying from the BBA, while HSBC Holdings Plc has said it may leave London.\n\nOsborne diluted a levy on U.K. banks and pushed out the regulator’s chief misconduct enforcer, Martin Wheatley, and most recently u-turned on a plan to assume senior bank managers are guilty until proven innocent, which lenders blamed for hindering recruitment of top foreign executives.\n\n“We recognize the change of tone in conduct regulation, important developments in the senior managers regime, the proposed reduction in the bank levy, greater certainty over tax for international banks,” the BBA said.\n\nNevertheless, London’s financial sector continues to shrink while its rivals grow, according to the report. Compared with 35,000 jobs losses and a 12 percent fall in U.K. banking assets in the past four years, assets in the U.S. have grown by the same percentage, while in Singapore and Hong Kong they have climbed by 24 percent and 34 percent respectively.\n\nEuropean firms are also losing market share to U.S. rivals in wholesale banking, which is the part of banks that cater to large corporates and other financial institutions. From 2010 to 2014, the wholesale market share of the top five European banks fell to 24 percent from 26 percent, whereas the share of the top five U.S. banks has risen to 48 percent from 44 percent, the BBA said.\n\nLondon is also losing market share in lending and initial public offerings, the BBA said. Wholesale banking’s global return-on-equity, a measure of profitability, is expected to fall to an average of 6.5 percent by 2017, about a third of the 18 percent-average between 2000 and 2006, according to the report, co-authored by consulting firm Oliver Wyman.\n\nOsborne’s overtures to the industry were counterbalanced by the high cost of ring-fencing -- a law that requires splitting off retail units to protect them from investment banking losses, the BBA said. “Uncertainty arising from the rapidly changing tax regime and European Union referendum are inhibiting business planning and discouraging investment,” according to the report.\n\nThe BBA’s wishlist includes a demand the Chancellor cut the bank levy faster. Under current plans the tax will be reduced over six years and then limited to domestic balance sheets until 2021. The lobby group also wants an 8 percent surcharge on bank profits to be phased out over time.\n\nFinancial services is the U.K.’s biggest export industry selling 62 billion pounds abroad every year, and employing more than 405,000 people, the BBA said.\n\nBefore it's here, it's on the Bloomberg Terminal.",
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
		Title:           "Credit-card chips could slow Black Friday lines",
		MetaDescription: "A change to how retailers process payments could make Americans stand in line longer this Black Friday.",
		CleanedText:     "Just when you thought there couldn't be another way to make Black Friday any more miserable for shoppers and retail employees, the credit-card industry came up with one.\n\nCredit-card companies last month began to mandate new technology that uses chips instead of magnetic stripes. It's a change made for a very good reason: card security.\n\nThe credit-card industry self-imposed October 1 as the deadline for the new card readers, though many consumers had received chip-enabled credit and debit cards — which will still work on the old \"swipe\" card processors — long before that.\n\nThe timing of this wider rollout, however, has retail and payments experts warning that this will slow things down at the checkouts on the November 27 shopping day.\n\n\"Any time you introduce a major change like this, there's going to be confusion,\" said Matt Schulz, senior industry analyst with CreditCards.com. \" There's no question this is going to cause some slowdown on Black Friday.\"\n\nThe change itself is simple: Instead of swiping the card through the magnetic-strip reader, shoppers now have to insert it — chip side up — into a slot on the bottom of the device.\n\nBut here's where the delays come in. People who are unfamiliar with the process will swipe as they always have, then be told it didn't work because they have a new chip-enabled card. Then they must be shown how to insert it, and leave it in, so the payment can be processed.\n\nNow multiply that by thousands, and add in the fact that people have been in line since the crack of dawn, elbowed their way to that bargain bin, and then had to wait again just to get to the register, and you can see why even a small delay will test patience. It's called the EMV chip, and it just might wreak havoc on holiday shopping.\n\n\"There is going to be a rude awakening\" for retailers, said Jared Drieling, business intelligence manager for The Strawhecker Group, an Omaha, Nebraska-based advisory firm focusing on payments. \"The industry is still bickering over how long an EMV transaction takes.\"\n\nAs many as 47% of US merchants will have new technology in place by the end of 2015, according to a survey conducted earlier this year by the Payments Security Task Force, an industry-backed group of financial services firms and leading retailers. Already, 40% of Americans have been issued new chip-enabled cards.\n\nOf course the nightmare scenario that Drieling is warning about is dependent on a lot of factors. Some customers have been using the chip technology for weeks, and some retailers don't have the readers yet. There is a wide disparity in how individual retailers have gotten ready for the switch.\n\nBest Buy, Macy's, and Walmart stores have been fully outfitted with new card readers, representatives for those companies said. Macy's and Walmart have also reissued store-branded credit cards with new EMV chips embedded in them. Sears, on the other hand, says it is \"continuously working to further enhance the security of our systems,\" according to a spokesman — but declined to provide specifics for Black Friday.\n\nJ. Craig Shearman, a spokesman for the National Retail Federation, said the new card readers would be at \"most major retailers and large national chains.\" The progress of smaller shops\u00a0in\u00a0adapting the chips is not as clear, but those shops\u00a0are less likely to\u00a0be open the day after Thanksgiving anyway.\n\nShearman didn't argue with the notion that things could slow down, but he said it was not clear how much longer it would take to process each transaction.\n\nFor retailers, Black Friday and the ensuing weekend is crucial to performance. Americans packed malls and stores last year after Thanksgiving, driving more than $50 billion in revenue to retailers, the National Retail Federation reported in 2014.\n\nOf course, there are lots of ways to avoid even having to find out. Stay home. Turkey and stuffing is better on day two anyway.\n\nNOW WATCH: JAMES ALTUCHER: 'Warren Buffett is a f-----g liar'",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.businessinsider.com/credit-card-chips-could-slow-black-friday-lines-2015-11",
		TopImage:        "http://static5.businessinsider.com/image/56410a64bd86ef18008c8901/this-little-change-could-make-black-friday-even-more-miserable-this-year.jpg",
	}
	article.Links = []string{
		"http://www.businesswire.com/news/home/20150504005631/en/Issuers-Forecast-U.S.-Shift-Chip-Cards-Complete",
		"http://www.usatoday.com/story/money/business/2015/10/01/chip-credit-debit-card-readers-october-1/73140516/",
		"http://www.businessinsider.com/james-altucher-warren-buffett-rant-holding-period-2015-10",
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
		CleanedText:     "Story highlights Lewis Hamilton reveals Monaco car accident on eve of Brazilian GP\n\nF1 world champion says he was exhausted and had a fever\n\nHamilton organized surprise party for his Mum after Mexico GP\n\nThe Mercedes driver revealed he crashed his car in Monaco after \"heavy partying\" last weekend. He turned up for this weekend's Brazilian Grand Prix a day late after taking time off to recover.\n\n\"I've not been well with a fever but I also had a road accident in Monaco on Monday night,\" Hamilton explained on his Instagram account.\n\n\"Nobody was hurt, which is the most important thing. I made very light contact with a stationary vehicle.\n\n\"Talking with the team and my doctor, we decided together that it was best for me to rest at home and leave a day later.\"\n\nDear TeamLH, just wanted to let you know why things have been quiet on social media the past few days. I've not been well with a fever but I also had a road accident in Monaco on Monday night. Whilst ultimately, it is nobody's business, there are people knowing my position that will try to take advantage of the situation and make a quick buck. NO problem. Nobody was hurt, which is the most important thing. But the car was obviously damaged and I made very light contact with a stationary vehicle. Talking with the team and my doctor, we decided together that it was best for me to rest at home and leave a day later. But i am feeling better and am currently boarding the plane to Brazil. However, I am informing you because I feel we all must take responsibility for our actions. Mistakes happen to us all but what's important is that we learn from them and grow. Can't wait for the weekend Brazil🙌🏾 Bless Lewis\n\nA photo posted by Lewis Hamilton (@lewishamilton) on Nov 11, 2015 at 2:50pm PST\n\nHamilton posted the news to his fans, who he refers to as \"Team LH,\" but he also added: \"Ultimately, it is nobody's business, there are people knowing my position that will try to take advantage of the situation and make a quick buck.\"\n\nAfter arriving in Sao Paulo for the penultimate race of the 2015 season, the three-time world champion inevitably faced questions from the assembled media.\n\nJUST WATCHED Replay More Videos ... MUST WATCH\n\nBoth Hamilton and his Mercedes teammate Nico Rosberg always speak to reporters on the Thursday before a race weekend, while the British driver also has obligations with the UK press.\n\nHamilton explained that his busy schedule since the last race in Mexico 12 days ago had included throwing a surprise 60th birthday party for his mother Carmen in London last Sunday, the night before his Monaco prang.\n\n\"\"It was a result of heavy partying and not much rest for 10 days. I am a bit run down,\" Hamilton, who spent four more days in Mexico after the race, said in his BBC Sport column.\n\n\"When I got back to the UK, I was trying to organize my Mum's 60th birthday. The party turned out great but by the end of it I was exhausted. I had been busy for two solid weeks and I basically collapsed.\"\n\nJUST WATCHED Replay More Videos ... MUST WATCH\n\nAlthough an element of mystery still surrounds Hamilton's Monaco car crash, it's not the first time the 30-year-old has been involved in driving drama off the track.\n\nAt the 2010 Australian Grand Prix, Hamilton was fined for dangerous driving after deliberately spinning his wheels and skidding on his way out of the Albert Park circuit. In 2007, when he was an F1 rookie, his car was impounded in France after he was caught speeding.\n\nHamilton, who wrapped up the 2015 world title at the U.S. Grand Prix in Austin, Texas with three races to spare, is now focused on getting back to business in Brazil.\n\n\"I feel good, I'm on an up slope, so a lot closer to 100%\" Hamilton told reporters at the Interlagos track. \"I'm excited to be here. I'm definitely cherishing the moments I'm in the car.\"\n\nTell us what you think of Hamilton's crash on CNN Sport's Facebook page",
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
		Title:           "Debenhams and House of Fraser charge for PAPER BAGS as Tesco give them away",
		MetaDescription: "Major high street stores including Debenhams and House of Fraser have started charging up to 10p for paper carrier bags – despite them being exempt from the new laws brought in last month.",
		CleanedText:     "Major high street stores have been accused of ripping off shoppers by charging up to 10p for paper carrier bags – despite them being exempt from the new laws brought in last month.\n\nOutraged shoppers have hit out at Debenhams and House of Fraser claiming they are 'cashing in' by charging for paper bags when other high street shops offer them for free.\n\nHouse of Fraser has said the charge for paper bags had been introduced for 'ethical and moral' reasons, and that all proceeds would be donated to charity.\n\nHowever, shoppers have taken to Twitter to express their anger at the charge.\n\nHouse of Fraser has said the paper bag charge has been brought in at stores for 'ethical and moral' reasons\n\nPaper bags are being handed out to shoppers at London branch of Tesco weeks after 5p charge introduced\n\nTwitter user Jimmy said: 'Absolutely disgusted! Just spent £180 on shoes and you have the audacity to make me pay 5p for a 'cardboard' bag #shocking'\n\nAnthony Bongos added: 'I can't understand why you are charging for paper carrier bags. This isn't the law, is it you cashing in on the law?'\n\nA spokesperson for House of Fraser said: 'We have made the ethical and moral decision to support the introduction of a 5p charge on all plastic and paper bags.'\n\nShoppers in Debenhams have also reported being charged to paper bags, with some saying they have been made to pay up to 10p.\n\nSuzanne Foley said: '£162 for a suit no suit bags and then get charged 10p for a large bag, what's that all about debenhams!' (sic)\n\nAnd Martena David added: '£162 for a suit no suit bags and then get charged 10p for a large bag, what's that all about debenhams!' (sic)\n\nElsewhere, some Tesco stores have started giving customers free paper bags just weeks after the 5p charge for plastic bags caused chaos around the country.\n\nTwitter uses have expressed their outrage after being made to pay for paper bags at House of Fraser\n\nThe rules are being rolled out by the Government's Department for Environment, Food & Rural Affairs. It claims the change will save £60m in litter clean-up costs and £13m in carbon savings.\n\nThe levy for supermarkets and big shops employing more than 250 staff will raise more than £70m a year for 'good causes'. Shops can also take a 'reasonable costs' cut. The Government will pocket the VAT raising an estimated £19m a year.\n\nYes. If you have bought food such as fish, uncooked meat or prescription medicines then the retailer should still offer bags for nothing.\n\nBut problems occur if you buy anything else at the same time. For example, if the bag shares space with a packet of cornflakes it will cost you 5p. You should not be charged if a shop uses paper bags.\n\nA London store has been handing out recyclable small bags as an alternative to shoppers just picking up a handful of groceries.\n\nThe bags feature the phrase 'love food hate waste'.\n\nThe new law does not prevent shops from handing out free paper bags, a source from the Department for Food, Environment and Rural Affairs told the Evening Standard.\n\n'The key thing is encouraging people to reuse bags,' they said.\n\n'The best thing to do is to have a plastic bag in your pocket.\n\n'But clearly paper bags can be recycled and do degrade better than plastic bags, and they won't end up strangling a turtle.'\n\nEngland was the last place in the UK to introduce the 5p bag charge.\n\nSome supermarkets around the UK where forced to put security tags on baskets and trolleys after shoppers began taking them home to carry their groceries.\n\nMailOnline has contacted Debenhams and House of Fraser for comment.\n\nDebenhams has been accused of ripping off customers across the UK by charging up to 10p for paper bags\n\nMOST WATCHED NEWS VIDEOS\n\nPrevious\n\n1\n\n2\n\n3\n\nNext\n\nMOST READ NEWS\n\nPrevious\n\nNext\n\n●\n\n●\n\n●",
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

func Test_EconomistCom(t *testing.T) {
	article := Article{
		Domain:          "economist.com",
		Title:           "Renting hotel rooms by the hour: A quick in and out",
		MetaDescription: "A Spanish hotel-reservation platform that allows customers to book rooms in three hour slots is looking to expand into Britain.",
		CleanedText:     "BYHOURS, a Spanish hotel-reservation platform that allows customers to rent rooms in three-hour slots, is looking to expand into Britain. Travelmole  that the website aims to sign up 25 hotels in the country by the end of the month, although so far only six have taken the plunge.\n\nMany people, when bringing to mind short-stay hotel rooms, will no doubt picture businessmen with their cinq-à-septs or, perhaps, company a little more transactional than that. Banish such grubby thoughts from your minds; having the option of booking a bedroom for three hours is a great and practical idea.\n\nIt is no coincidence that several of the establishments that have signed up with ByHours are close to airports and train stations. How often have you had several hours to kill at an airport and longed for a place to shower and snooze? And Gulliver has written before about that horrible dead time when, having checked out of a hotel in the morning, with your flight not until late in the evening, you have ages to kill wandering around a strange town dragging a wheely-bag. Then there are those day trips when you fly in to town at some ungodly early hour and are scheduled to fly out at an equally uncivilised late one; how much more pleasant if you could pop your head down for a few hours in the afternoon? In fact you needn’t even be a visitor. Back when Gulliver's daughter was a sleep-averse baby, he would have paid handsomely for the chance to close his eyes for an hour in a short-stay hotel during his lunch break.\n\nIt is also easy to see why it would appeal to hotels, which could sweat their assets more, filling gaps between guests checking out and in. According to Travelmole, in Spain last year more than 150,000 bookings were made through ByHours at more than 1,500 hotels. However, for the consumer the big drawback would appear to be pricing. Prices for a three-hour stay in London tomorrow start at €50 and quickly hit the hundreds. That is understandable. By its nature it is often likely to be a last-minute purchase, and hotels will obviously price very short reservations at a premium. But the more hotels that sign up, the easier it will be to find something more budget friendly.",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.economist.com/blogs/gulliver/2015/04/renting-hotel-rooms-hour",
		TopImage:        "https://www.economist.com/sites/default/files/images/guliver.png",
	}
	article.Links = []string{
		"http://www.travelmole.com/news_feature.php?news_id=2016292",
		"http://content.time.com/time/magazine/article/0,9171,843018,00.html",
		"http://www.economist.com/blogs/gulliver/2013/04/surreptitious-snoozing",
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
		CleanedText:     "Everyone has fears. They’re important, and they’ve helped keep us alive throughout our evolution. Think about the fears\u00a0characters understandably\u00a0feel at certain points in\u00a0Game of Thrones, the hugely successful HBO dramatic series which\u00a0combines elements of medieval times and fantasy. We're talking outrageously murderous kings here, plus scheming\u00a0lords and ladies. Large men with even larger swords. Even fire-breathing dragons.\n\nRelated:\u00a0Why Fear Is the Entrepreneur's Best Friend\n\nIn Season One of GOT, a great line\u00a0illustrates the point about fears perfectly. The speaker is\u00a0Robb Stark, eldest son of the lord of Winterfell and generally a good guy, who\u00a0decides to declare war and march south to Kings Landing, the capital of the Seven (usually warring) Kingdoms and home to\u00a0a lot more of those men with swords . Theon\u00a0Greyjoy, the son of another royal house,\u00a0asks Stark if he’s afraid. And Stark, his hands trembling, replies,\u00a0“I guess I must be.” To which\u00a0Greyjoy’s response is perfect:\u00a0“Good, that means you’re not stupid.”\n\nIt certainly was appropriate for the denizens of GOT's medieval era to be afraid, but does the same apply to you? For, while fear was an important factor in our hereditary past, in our modern day and age, our fears today\u00a0are often based more in psychology\u00a0than\u00a0actual physical threats. Drawing on some of the books I've enjoyed, I offer\u00a0six thoughts on why facing your fears will assist you in creating massive success.\n\nI've had a lot of worries in my life, most of which never happened.\u00a0- -Mark Twain\n\nWhen you take the time to actually define your fears, you\u00a0learn to separate fact from fiction. This is an important distinction. Some things you’re afraid of will be valid, but many will be mental worst-case scenarios that have simply spiraled further in your mind than they ever will or would in reality.\n\nWhat about the fears on your list that you’ve defined that are actually valid, like losing a client or\u00a0employee, gettng backlash from a layoff\u00a0or encountering some other tangible fear?\u00a0Easy. When you face fears that have merit -- now that you’ve defined them --\u00a0you can come up with an action plan of responses to mitigate the damages.\n\nThink of this list as your \"fear emergency\u00a0plan.\" You know what you’d do in the case of a fire or earthquake, so why not enact a plan of appropriate responses you could take against some of your more valid business\u00a0fears?\n\nRelated:\u00a07 Ways to Think Differently About Fear\n\nBran thought about it. \"Can a man still be brave if he's afraid?\" \"That is the only time a man can be brave,\"\u00a0his father told him.” -- George R.R. Martin, series author, A Song of Ice and Fire, on which HBO's GOT series is based.\n\nPerhaps I’m just missing Game of Thrones in the offseason, but this quote really struck me and is an important facet of facing your fears. You don’t develop bravery and courage in the good times, you develop them when you actually confront fears. If you were once afraid of starting your own business, but did it anyway, you know the terror, but also the reward, that comes from facing fears head on. Your courage grows with each fear you face.\n\nThere is wisdom that comes from the experience of working through fears. Some of your fears may have even come true. If you are a business owner and have seen your business falter or fail, perhaps you’ve already lived through adversity. The silver lining of these experiences is that you learn from them. Wisdom comes from all of life’s experiences, but the fearful or bad ones in particular teach us great lessons. Wisdom is always the by-product of facing your fears, and that’s an important quality to develop.\n\nDealing with fears helps your develop compassion. When you yourself have been afraid,\u00a0you’re more likely to have patience and feel compassion toward others experiencing similar situations. After all, we all want a good life. When you push hard for what you want, and experience the joys and failures of success, you learn compassion you can use to help others push through their early fears.\n\nYou can put yourself in\u00a0the shoes of someone who is just starting out, and that empathy can help guide that person to have deeper courage.\n\n“Life doesn't get easier or more forgiving;\u00a0we get stronger and more resilient.” -- Steve Maraboli, Life, the Truth, and Being Free\n\nResilience comes from facing your fears. You become better than your surroundings and transform yourself above the fear and into bigger and bigger success. Resiliience starts with you, and it begins in your mind. Face your fears and learn to rise to face whatever is in front of you.\n\nRelated:\u00a0What Companies Can Learn From 'Game of Thrones' When Hiring Their Next Chief Information Officer",
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
		CleanedText:     "Published November 13, 2015 Associated Press\n\nHillary Rodham Clinton has locked up public support from half of the Democratic insiders who cast ballots at the party's national convention, giving her a commanding advantage over her rivals for the party's presidential nomination.\n\nClinton's margin over Vermont Sen. Bernie Sanders and former Maryland Gov. Martin O'Malley is striking. Not only is it big, but it comes more than two months before primary voters head to the polls -- an early point in the race for so many of the people known as superdelegates to publicly back a candidate.\n\n\"She has the experience necessary not only to lead this country, she has experience politically that I think will help her through a tough campaign,\" said Unzell Kelley, a county commissioner from Alabama.\n\n\"I think she's learned from her previous campaign,\" he said. \"She's learned what to do, what to say, what not to say -- which just adds to her electability.\"\n\nThe Associated Press contacted all 712 superdelegates in the past two weeks, and heard back from more than 80 percent. They were asked which candidate they plan to support at the convention next summer.\n\nThe 712 superdelegates make up about 30 percent of the 2,382 delegates needed to clinch the Democratic nomination. That means that more than two months before voting starts, Clinton already has 15 percent of the delegates she needs.\n\nThat sizable lead reflects Clinton's advantage among the Democratic Party establishment, an edge that has helped the 2016 front-runner build a massive campaign organization, hire top staff and win coveted local endorsements.\n\nSuperdelegates are convention delegates who can support the candidate of their choice, regardless of who voters choose in the primaries and caucuses. They are members of Congress and other elected officials, party leaders and members of the Democratic National Committee.\n\nClinton is leading most preference polls in the race for the Democratic nomination, most by a wide margin. Sanders has made some inroads in New Hampshire, which holds the first presidential primary, and continues to attract huge crowds with his populist message about income inequality.\n\nBut Sanders has only recently started saying he's a Democrat after a decades-long career in politics as an independent. While he's met with and usually voted with Democrats in the Senate, he calls himself a democratic socialist.\n\n\"We recognize Secretary Clinton has enormous support based on many years working with and on behalf of many party leaders in the Democratic Party,\" said Tad Devine, a senior adviser to the Sanders campaign. \"But Sen. Sanders will prove to be the strongest candidate, with his ability to coalesce and bring young people to the polls the way that Barack Obama did.\"\n\n\"The best way to win support from superdelegates is to win support from voters,\" added Devine, a longtime expert on the Democrats' nominating process.\n\nThe Clinton campaign has been working for months to secure endorsements from superdelegates, part of a strategy to avoid repeating the mistakes that cost her the Democratic nomination eight years ago.\n\nIn 2008, Clinton hinged her campaign on an early knockout blow on Super Tuesday, while Obama's staff had devised a strategy to accumulate delegates well into the spring.\n\nThis time around, Clinton has hired Obama's top delegate strategist from 2008, a lawyer named Jeff Berman, an expert on the party's arcane rules for nominating a candidate for president.\n\nClinton's increased focus on winning delegates has paid off, putting her way ahead of where she was at this time eight years ago. In December 2007, Clinton had public endorsements from 169 superdelegates, according to an AP survey. At the time, Obama had 63 and a handful of other candidates had commitments as well from the smaller fraction of superdelegates willing to commit to a candidate.\n\n\"Our campaign is working hard to earn the support of every caucus goer, primary voter and grassroots and grasstop leaders,\" said Clinton campaign spokesman Jesse Ferguson. \"Since day one we have not taken this nomination for granted and that will not change.\"\n\nSome superdelegates supporting Clinton said they don't think Sanders is electable, especially because of his embrace of socialism. But few openly criticized Sanders and a handful endorsed him.\n\n\"I've heard him talk about many subjects and I can't say there is anything I disagree with,\" said Chad Nodland, a DNC member from North Dakota who is backing Sanders.\n\nHowever, Nodland added, if Clinton is the party's nominee, \"I will knock on doors for her. There are just more issues I agree with Bernie.\"\n\nSome superdelegates said they were unwilling to publicly commit to candidates before voters have a say, out of concern that they will be seen as undemocratic. A few said they have concerns about Clinton, who has been dogged about her use of a private email account and server while serving as secretary of state.\n\n\"If it boils down to anything I'm not sure about the trust factor,\" said Danica Oparnica, a DNC member from Arizona. \"She has been known to tell some outright lies and I can't tolerate that.\"\n\nStill others said they were won over by Clinton's 11 hours of testimony before a GOP-led committee investigating the attack on a U.S. consulate in Benghazi, Libya. Clinton's testimony won widespread praise as House Republicans struggled to trip her up.\n\n\"I don't think that there's any candidate right now, Democrat or Republican, that could actually face up to that and come out with people shaking their heads and saying, `That is one bright, intelligent person,\"' said California Democratic Rep. Tony Cardenas.",
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
		CleanedText:     "Rodrigo Caio treina até nas férias e tenta acelerar retorno aos gramados Jogador segue programação de exercícios em Dracena, interior de São Paulo. Comissão técnica planeja volta dele para o fim de fevereiro ou início de março\n\nRodrigo Caio treina na esteira durante as férias em Dracena-SP (Foto: Divulgação) Rodrigo Caio quer ganhar tempo na recuperação da lesão que\n\nsofreu no joelho esquerdo. Apesar de ter sido liberado pelo departamento médico\n\ndo São Paulo para as férias, o jogador vem treinando diariamente para acelerar\n\na recuperação após ser submetido a uma cirurgia.\n\nO zagueiro e volante passa férias com a família em Dracena, interior\n\nde São Paulo, e alterna os períodos de descanso com uma rotina de\n\nexercícios. Ele vem realizando trabalhos de reforço muscular e corridas na\n\nesteira.\n\nO jogador lesionou o joelho esquerdo no dia 2 de agosto,\n\ncontra o Criciúma, no Morumbi, pelo Campeonato Brasileiro, e precisou passar por\n\numa cirurgia. O defensor vinha sendo um dos destaques do São Paulo na\n\ntemporada.\n\nNa avaliação do departamento médico, Rodrigo Caio deve\n\nser liberado para treinos com o elenco e jogos entre fevereiro e março. Com\n\nisso, é provável que seja inscrito pelo técnico Muricy Ramalho para disputar a\n\nfase de grupos da Taça Libertadores.\n\nsobre\n\nSão Paulo\n\n+\n\nAnterior\n\n30\n\nDez\n\n18:27\n\nBLOG: Corinthians corre risco de perder Dudu para o São Paulo\n\n16:30\n\nBLOG: RETROSPECTIVA 2014: Entre variações táticas, o ano foi da intensidade, 3 zagueiros e contragolpe\n\n12:01\n\nEm último teste antes da Copinha, São Paulo empata com Botafogo-SP\n\n11:06\n\nSão Paulo tenta fazer acordo para se livrar do 'mico' Clemente Rodríguez\n\n10:00\n\nVolante Hudson vê São Paulo pronto para conquistar títulos em 2015\n\n07:25\n\nTricolor recebe Botafogo em Cotia e faz último amistoso antes da Copinha\n\n29\n\nDez\n\n19:05\n\nApós ano no Dragão, Caramelo pode ser cedido pelo São Paulo à Chape\n\n16:28\n\nConmebol divulga tabela detalhada da Taça Libertadores de 2015; veja\n\n09:15\n\nAidar acredita em brilho de Pato, mas cobra: \"Ainda não mostrou a que veio\"\n\n08:00\n\nTimes paulistas tentam manter hegemonia recente na Copinha\n\nProximo\n\n+\n\nAnterior\n\n24\n\nDez\n\n08:10\n\nSem dor, Rodrigo Caio vence etapas e já pensa na volta aos gramados\n\n19\n\nDez\n\n21:44\n\nVice do São Paulo diz que Alvaro fica, revela parceria e quer comprar Pato\n\n13:52\n\nEm recuperação, Rodrigo Caio segue rotina no CT e ganha apoio de Ganso\n\n09\n\nDez\n\n15:04\n\nRodrigo Caio vence nova etapa de recuperação e inicia corrida na esteira\n\n25\n\nSet\n\n17:54\n\nSaudade! Toloi e Rodrigo Caio observam treino dos reservas no CT\n\n11\n\nSet\n\n16:02\n\nRodrigo Caio inicia nova etapa de recuperação e festeja evolução\n\n14\n\nAgo\n\n17:13\n\nRodrigo Caio começa fisioterapia após cirurgia no joelho esquerdo\n\n12\n\nAgo\n\n18:23\n\nCom a família por perto, Rodrigo Caio comenta dificuldades após a operação\n\n07\n\nAgo\n\n16:23\n\nRodrigo Caio passa por cirurgia e inicia fisioterapia na próxima semana\n\n06\n\nAgo\n\n08:05\n\nLesão de Rodrigo Caio trará reflexos dentro e fora de campo no São Paulo\n\nProximo",
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
		CleanedText:     "With $200 billion in annual buying power by 2017, Millennials have become every brand’s coveted customer. But what’s the best way to reach them?\n\nThe answer is email.\n\nFor all the talk of email being dead — Too much noise! Too much spam! Too many distractions! Snapchat! — email remains\u00a0the standard for digital communication. In fact, Millennials check email more than any other age group, and nearly half can’t even use the bathroom without checking it, according to a\u00a0recent Adobe study.\n\nThat same study\u00a0found nearly 98% of Millennials check their personal email at least every few hours at work, while almost 87% of Millennials check their work email outside of work.\n\nEmail is not only relevant for Millennials, it also happens to remain the channel where direct marketers get the highest ROI ($39 for every dollar spent, according to the Direct Marketing Association). But that doesn’t mean the same old email marketing will work on Millennials. Instead, marketers need to adjust, or run the risk of that dreaded swipe to the trash bin. Consider these ideas the next time you’re planning an email campaign and Millennials are a key part of the audience:\n\nMobile is a must. Millennials are more likely than any other age group to check email on smartphones, with 88% reporting that they regularly using a smartphone to check email. If you’re not mobile first, you’re not putting your Millennial customers first. Responsive design has been a mantra for some time, but if you’re not employing it, you’re alienating an important generation of consumers who live, breathe, and sleep with their mobile devices.\n\nTiming is everything. Looking at opens and clicks won’t get you anywhere without analyzing the day of week and time of day those emails are opened and clicked. For example, we found that Millennials are more likely than any other age group to check email while in bed (45.2%). Why not experiment with sending emails first thing in the morning or late in the evening with content relevant to that time of day?\n\nPictures are worth a thousand words. They’re also an important mechanism for Millennials to filter messages. Why send an email survey asking for written feedback when all you need to do is provide a choice between a smiley face and a frown? Images are an integral part of Millennial language, even in the workplace. A third of Millennials believe it is appropriate to use an emoji when communicating with a direct manager or senior executive, so it’s a safe bet they’re even more comfortable when it comes to emoji from brands. Millennials are thinking and communicating in images, so marketers need to optimize emails for images and allow for quick feedback through emoji.\n\nLess is more . Email marketing to Millennials isn’t about sending more of the same. Many Millennials want to see fewer emails (39%) and fewer repetitive emails from brands (32%). Marketers take note — stop spamming your lists and start marketing to individuals by understanding who they are first.\n\nNot every Millennial communicates the same way, of course. And digital communication is constantly evolving. Nonetheless, for now it seems safe to say that email is here to stay and will remain a critical channel even for reaching mobile customers. Just don’t expect the same old email tactics to work.",
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
		CleanedText:     "Since its inception, television has been a unifying social force, bringing family, friends and different groups of people together. Even watching television on your own connects you to the multitudes of others watching the same thing across the globe.\n\nTV has come a long way: from black-and-white to colour, from a rare treat accessible to few to a household staple for everyone, from standard definition to tomorrow's ultra-HD screens.\n\nPerceptions of TV audiences have also changed over time. While theorists once believed TV viewers were passive, zombie-like figures transfixed in front of their televisions, numerous studies have proven that TV audiences are engaged, active and critical of the programmes they watch.\n\nIn the last several years, we've seen a dramatic shift that's placed viewers in control of their own scheduling. There's also more choice than ever before when it comes to accessing favourite programmes and watching them when and where they like.\n\n\"There are two simultaneous trends emerging when it comes to our TV watching habits, and they're two opposite trends, which is interesting,\" says Professor Sonia Livingstone OBE, a full professor in the Department of Media and Communications at the London School of Economics.\n\n\"One: we're watching TV on our laptops, tablets and phones, wherever and on whatever.\"\n\nAnd two, somewhat paradoxically, we're seeing a growth in the size of the screen in the living room. People talk about how everyone is watching TV on a 'small screen', but there's also a new viewing growing up around this enormous screen, as well as the more individualised viewing.\"\n\nNow, we watch shows wherever we want, whether it's relaxing in the bath with Corrie characters, catching up with a favourite drama on our phone during a morning commute or settling down in the sitting room every week to enjoy GBBO, gathered around the biggest 'and best' screen in the house. Equally, thanks to the latest in wearable tech, our most beloved television content has become a coveted accessory, accessible with a swipe on our watch.\n\nSubscription-free services like Freeview Play have also given us more options than ever before, with over 60 TV channels, 12 HD channels and over 25 radio stations a remote click away, plus the freedom to catch up on shows from the BBC, ITV, Channel 4 and Channel 5. Other services like Netflix and Amazon Prime also give us the opportunity to watch shows we missed the first time around - in one sitting, if we so desire! - while simultaneously introducing us to new and original programming.\n\n\"We keep fearing that people won't talk to each other anymore,\" says Professor Livingstone. \"There's the choice to watch separately and the choice to come together, whether it's binge viewing or the greater choice of programmes than ever before.\"\n\nAll of this choice has had a positive impact on TV consumers, according to Professor Livingstone.\n\n\"Most of the evidence is that people are feeling empowered and delighted. There's been an enormous welcome from people about the joys of having so much control and more choice than ever before.\"\n\nPeople are also prepared to pay to improve their television watching experience, whether that's spending on bigger HD screens or subscription services.\n\nWhile scheduling is fairly unimportant for younger generations, the middle-aged and young elderly population that remembers how television used to be is growing, so scheduling continues to play an important role for them.\n\nFor those younger generations, the definition of whether TV is 'a five minute clip of a beauty vlogger's latest haul on YouTube or a critically respected docudrama' calls into question what TV viewing really means these days.\n\n\"People have been saying for a while that scheduling is dead, but there's no getting rid of schedule for the 40s or 50-pluses who absolutely adhere to traditions of what to watch and when,\" says Professor Livingstone.\n\nRapidly emerging trends, like the increase in individual TV consumption across new tech and the importance of the living room big screen as the centrepoint of family life, ensure that the landscape of television scheduling is in constant flux and the future of television remains uncertain.\n\nOne thing we know? We'll still be watching.",
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

func Test_HuffingtonPostJp(t *testing.T) {
	article := Article{
		Domain:          "huffingtonpost.jp",
		Title:           "クロマグロ残り2匹　葛西臨海水族園の大量死は未だに原因不明",
		MetaDescription: "クロマグロやカツオ類が大量死した問題で、葛西臨海水族園（東京都江戸川区）は３日、病理検査の結果、海の養殖魚を大量死させることで知られる２種類のウイルスが原因ではないことが確認されたと発表した。",
		CleanedText:     "",
		MetaKeywords:    "クロマグロ残り2匹　葛西臨海水族園の大量死は未だに原因不明, japan",
		CanonicalLink:   "http://www.huffingtonpost.jp/2015/03/03/tuna-death_n_6796602.html",
		TopImage:        "http://i.huffpost.com/gen/2678692/images/o-TUNA-DEATH-facebook.jpg",
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
		Title:           "Why 2015 Was Rent the Runway's Biggest Year So Far",
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
		CleanedText:     "Work-life balance. Everyone talks about it. And everyone struggles to achieve it.\n\nYet finding a reasonable work-life balance is easier than you think. While it's true the equilibrium point is constantly shifting, most of the same attitudes, perspectives, and skills apply to both \"work\" and \"life.\"\n\nSo why not take advantage of that fact? Pick the right \"life\" pursuits and they inform and enhance your professional skills -- and add a healthy dose of perspective and humility along the way.\n\nIn my case I like to take on extremely difficult (at least for me) physical goals. (Granted my approach to goal achievement in general is a little unconventional. Just like\u00a0Fight Club,\u00a0the first rule of achieving a goal is\u00a0you don't talk about achieving that goal. And achieving a goal has a lot less to do with the goal itself and\u00a0a lot more to do with the routine you develop\u00a0to support that goal.)\n\nSo a few years ago, after just four months of training, I rode the\u00a0Alpine Loop Gran Fondo, a 92-mile, four-mountain ride that included 11,000 feet of climbing. (Those four months felt like a lifetime, though, since pro mountain biker Jeremiah Bishop trained me. But then again I never could have been ready without him.)\n\nAfter a few years of cycling I got tired of being cycling skinny -- 6' tall, 150 lbs is not a particularly good look -- and decided to see if I could pull off some semblance of the\u00a0\"movie star becomes an action hero\"\u00a0physical transformation. I gained over 20 pounds, lost a few percentage points of body fat, and got a lot stronger. (That training sucked too, since\u00a0Jeffrey Del Favero\u00a0of\u00a0Bodybuilding.com\u00a0created my program, but then again I never could have done it without him.)\n\nSo why do I do take on (feel free to insert your own adjective) personal challenges? And how does that help me professionally? It's all about the habits, skills, and perspectives gained. Here are some reasons.\n\nSuccess is ultimately based on numbers. Sure, you can try to \"hack\" a goal. Sure, you can look for shortcuts. (People have\u00a0built entire careers\u00a0off the premise.) But eventually achieving a huge goal is all about volume and repetition.\n\nWant to eventually ride a tough gran fondo? You'll have to ride hundreds of miles along the way. Want to go from only being able to do three pull-ups to eventually being able to do four sets of twenty? You'll have to lift a ton of weight along the way.\n\nThe same is true for professional success; it's largely based on doing the work. Want twenty new customers? Expect to cold call two or three hundred prospects. Want to hire a superstar? Expect to screen dozens and then interview ten or fifteen people.\n\nThe surest path to success is to do an incredible amount of work. If you're willing to do the work, you can succeed at almost anything.\n\nThe armor that protects us eventually destroys us. We all wear armor. That armor protects us but also, over time, wears us down.\n\nOur armor is primarily forged by success. Every accomplishment adds an additional layer of protection from vulnerability. In fact, when we feel particularly insecure we unconsciously strap on more armor so we feel less vulnerable:\n\nArmor protects when we're unsure, tentative, or at a perceived disadvantage. Our armor says, \"That's okay; I may not be good at this... but I'm really good at\u00a0that.\"\n\nOver time armor also encourages us to narrow our focus to our strengths. That way we stay safe. The more armor we put on the more we can hide our weaknesses and failings--from others and from ourselves.\n\nWe use our armor all the time. I use my armor all the time--I feel sure more than you. But I get really tired of wearing it.\n\nWhen I ride a bike the guy who passes me doesn't care if I've ghostwritten bestsellers or drive a fancy car or live in a nice neighborhood. At the gym, the guy who lifts more than me also doesn't care about any of that stuff. He's stronger and fitter than me. Period.\n\nIn those situations no amount of armor, real or imagined, can protect me. I'm just a guy on a bike. I'm just a guy at the gym. I'm just me.\n\nBeing just me is pretty scary.\n\nBut being who you really are is something we all need to do more often. It keeps things in perspective. It reminds us that we can always be better. It reminds us that no matter how good we think we are at something there is always someone who is a lot better.\n\nAnd that's not depressing -- that's motivating.\n\nGrace is an awesome feeling -- one we can never experience enough. Outstanding athletes exist in a state of grace, a place where calculation and strategy and movement happen almost unconsciously. Great athletes can focus in a way that, to us, is unrecognizable because through skill, training, and experience their ability to focus is nearly effortless.\n\nWe've all felt a sense of grace, if only for a few precious moments, when we performed better than we ever imagined possible... and realized what we assumed to be limits weren't really limits at all.\n\nThose moments don't happen by accident, though. Grace is never given; grace must be earned through discipline and training and sacrifice.\n\nI want to ride up a mountain and experience the feeling that I can climb and climb and climb and I don't have to think about anything because I can just\u00a0go....\n\nI want to struggle with a weight and experience the feeling that I can do a few more reps because I know, without a doubt, I always have a little more in me...\n\nAnd I want to sometimes write almost effortlessly and without thinking because years of effort and practice have brought me to a place where occasionally I am the writer I would like to be...\n\nAll those are moments of grace. They're awesome. They're amazing.\n\nAnd they feed off each other because the confidence you build after experiencing a moment of grace in one pursuit helps you keep pushing when the going gets tough in other pursuits.\n\nWith work, \"then\" is always better than \"now.\"\u00a0 \"Now\" and \"then\" are wonderful words when they appear in the same sentence.\n\nWhen you work to improve at something -- especially in the beginning stages -- \"now\" is often a terrible place. At one point my \"now\" was riding like an asthmatic hippo. At one point my \"now\" was doing four dips and feeling like I was tearing my chest apart.\n\nBut with time and effort my \"now\" was transformed. I could ride\u00a0with more speed, power, and confidence. I could do\u00a0sets of ten, then twenty, then thirty dips. I was able to look back with satisfaction at a \"now\" I had transformed into a vastly inferior \"then.\"\n\nThink about something you wanted to do. Then think about where you would be\u00a0now\u00a0if you had actually gotten started on it\u00a0then.\n\nWhen you do the work, then always pales in comparison to now: family, business, and every aspect of your life. When you don't do the work, now is just like then -- except now you also get to live with regret.\n\nQuitting is a habit anyone can learn to break. We're all busy. Each of us face multiple, ongoing demands. Every day we are forced a number of times to say, \"That's not perfect, but it works... and I need to move on to something else.\"\n\nStopping short of excellence is something we are not just forced to do but are also\u00a0trained\u00a0to do. Most of the time we have no choice so we get really good at \"quitting.\"\n\nI'm really good at quitting. I raised wonderful kids and did a good job... but I know I could have done more. I've built a decent business... but I know I could have done more. I've tackled challenges before and tried really hard... but I know I could have done more.\n\nWhere physical challenges are concerned there are hundreds if not thousands of times I want to quit. Training is hard and only gets harder. Balancing family and work and everything else is hard and only gets harder.\n\nAt weak moments, struggle shatters our resolve and make us want to quit.\n\nIt's hard not to stop, by choice or otherwise, at \"good enough.\" But sometimes, if the goal is big enough, we have to be\u00a0great: not great compared to other people... but great compared to ourselves.\n\nThat comparison is the only comparison that really matters and is the best reason of all to try to accomplish more than you -- or anyone around you -- ever thought possible.\n\nWhen you succeed, you become something you were not. And then you get to do it again, and become\u00a0something else you once were not -- but definitely are now.\n\nI also write for Inc.com:\n\nCheck out my book of personal and professional advice,\u00a0TransForm: Dramatically Improve Your Career, Business, Relationships, and Life -- One Simple Step At a Time. (PDF version here,\u00a0Kindle version here,\u00a0Nook version here.)\n\nIf after 10 minutes you don't find at least 5 things you can do to make your life better I'll refund your money.\n\nThat way you have nothing to lose... and everything to gain.",
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

func Test_PrnewswireCom(t *testing.T) {
	article := Article{
		Domain:          "prnewswire.com",
		Title:           "Atlantic Merchant Capital makes lead investment in Social Quant -- TAMPA, Fla., April 29, 2015 /PRNewswire/ --",
		MetaDescription: "TAMPA, Fla., April 29, 2015 /PRNewswire/ -- Atlantic Merchant Capital makes lead investment in Social Quant.",
		CleanedText:     "TAMPA, Fla. , April 29, 2015 /PRNewswire/ -- Atlantic Merchant Capital Investors announced today that it has made a lead investment in conjunction with Fleur De Lis Partners, into Social Quant, LLC.  Social Quant is a Tampa -based technology start-up, formed by Dr. Morten Middelfart.  It provides clients with big data algorithmic support tools to increase the size and quality of their Twitter audiences.  Dr. Middelfart is a serial entrepreneur and globally-respected author and thinker on the uses of big data and artificial intelligence.\n\n\"We are very excited to have the opportunity to back Dr. Middelfart in this exciting new social media venture.  We think the company's IT engine can be expanded for uses with other social media tools.  As a result, we think Dr. Middelfart's company and services will gain greater scale and relevance over the next several years.  He has a history of success with projects like this and we are excited to be a part of this venture.  We think the value of his IP is substantial already, but over time we think this can become an extraordinarily valuable IP asset,\" said Allan Martin , CEO of Atlantic.\n\n\"The team at Social Quant is enthusiastic about our new capital partnership with Atlantic,\" said Dr. Morten Middelfart , founder of Social Quant.  \"The capital provided by Atlantic will allow us to execute marketing initiatives we believe are necessary to scale our offerings.  We are also excited about the opportunity to work with the Atlantic principals in the ongoing development of our strategy.  They are proving to add much more than capital to our business.\"\n\nAtlantic Merchant Capital Investors is a Tampa -based private equity firm founded in 2009 with investments in middle-market companies throughout the United States.  Its primary investment focus is financial services with emphasis on insurance and banking.  However, it has a growing portfolio of investments in scalable start-ups in an effort to support entrepreneurship in the Tampa Bay area.\n\nFleur De Lis Partners is a private investment partnership (controlled by Atlantic), also based in Tampa, Florida .\n\nContact:  Doug Licker, Atlantic Merchant Capital Investors 813-443-0745.  Doug.Licker@amci360.com.\n\nSOURCE Atlantic Merchant Capital Investors\n\nRELATED LINKS",
		MetaKeywords:    "Atlantic Merchant Capital Investors, florida, Financing Agreements, Banking & Financial Services, Social Media, Internet Technology, Computer & Electronics",
		CanonicalLink:   "http://www.prnewswire.com/news-releases/atlantic-merchant-capital-makes-lead-investment-in-social-quant-300074212.html",
		TopImage:        "http://content.prnewswire.com/designimages/logo-prn-01_PRN.gif",
	}
	//article.Links = []string{}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
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

func Test_SlideshareNet(t *testing.T) {
	article := Article{
		Domain:          "slideshare.net",
		Title:           "Lessons on Growth (Boston 2014)",
		MetaDescription: "This is the version of the talk \"Lessons in Growth\" I gave in Boston at the PayPal HQ on Feb 24, 2014.",
		CleanedText:     "Share\n\nPublished on Feb 24, 2014\n\nThis is the version of the talk \"Lessons in Growth\" I gave in Boston at the PayPal HQ on Feb 24, 2014.\n\n...\n\nBusiness\n\nTechnology\n\n0 Comments\n\ni\n\nStatistics\n\nNotes\n\n,\n\n,\n\nspan\n\nShow More\n\nViews\n\n0\n\nActions\n\n0\n\n0\n\n0\n\nEmbeds\n\n0\n\nwealthfront.com LESSONS IN GROWTH Adam\n\nNash @adamnash February 24, 2014\n\nLessons in Growth eBay Express\n\nBuild it and they will come?   LinkedIn  Viral growth just works… as long as you make it.  Wealthfront  An incredible service and business… if. wealthfront.com | 2\n\nWhy Growth Matters History is\n\nwritten by the victors Victory for consumer software is almost always deﬁned by your market reach Growth does not just happen, it has to be designed into your product and service The market is brutally competitive. If you don’t ﬁgure it out, your competitors will. wealthfront.com | 3\n\nThink About Non-Users Software companies\n\ntend to focus almost exclusively on their users Growth is about the experience you provide for nonusers (aka guests, visitors) When you are small and growing, almost all of your future users have not signed up yet Whats pages & features do you have for non-users? wealthfront.com | 4\n\nFive Sources of Web Traffic\n\nOrganic  People seek you out directly Email  You send out email, they click SEO  You expose pages that are indexed in Search Paid / Aﬃliate  Paid links to your content Social  Links to your content shared by users wealthfront.com | 5\n\nUnderstanding Virality One of the\n\nkey insights of our growth strategy from 2008. Extensible to literally all engagement features. Key measure used by applications on social platforms. This is an extremely useful frame. How does a new customer   today lead to a new customer   tomorrow? wealthfront.com | 6\n\nUnderstanding Virality At the heart\n\nof virality is an exponential based on branching factor and cycle time. Would your rather 10x every week, or 2x every day? Rabbits make lots of rabbits not because of big litters, but because they breed frequently. “n” matters more than “m”. n m number of cycles branching factor wealthfront.com | 7\n\nThree Steps to Virality Clearly\n\narticulate the ﬂows where content from a user can touch a non-user. Wireframes are ﬁne. Build a simple mathematical model for conversion rates and multipliers and cycle time. Instrument your ﬂows for these metrics. Develop, release, measure, iterate. Be prepared to execute 6-8 release cycles to get your factor high enough to contribute meaningfully. wealthfront.com | 8\n\nGrowth in a Mobile World\n\nNative apps work for a simple reason: they generate organic trafﬁc at scale. Transaction ﬂows are still heavily optimized for the web, and conversion rates are better. Non-users don’t have your app. It’s that simple. Engaged users, on your app, publishing content that reach non-users, likely converting on a web ﬂow. iOS, Android, Facebook & Twitter all working on bridging the native / transactional gap. wealthfront.com | 9\n\nValue First, Then Growth Engineering\n\ngrowth over a product that doesn’t provide real value will be short term, at best. Engagement and/or economic value must be at the heart of any sustainable growth curve. Don’t underestimate the power of word of mouth to drive your brand, your trafﬁc, and your conversion rates. wealthfront.com | 10\n\nFinal Thoughts We can be\n\nour own harshest critics. In the mirror we see every ﬂaw, every mistake, every imperfection. These are the very early years. Things that seem small now can and will be huge in 5 years. Each of you can and will have a profound impact on that future. Behavior matters. Values matter. wealthfront.com | 11",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.slideshare.net/adamnash/lessons-on-growth-boston-2014",
		TopImage:        "https://cdn.slidesharecdn.com/ss_thumbnails/lessonsongrowthv2-140224231617-phpapp01-thumbnail-4.jpg?cb=1393332344",
	}
	article.Links = []string{
		"http://www.linkedin.com/legal/copyright-policy",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_SoundCloudCom(t *testing.T) {
	article := Article{
		Domain:          "soundcloud2.com",
		Title:           "#18 Silence And Respect by Reply All",
		MetaDescription: "Stream #18 Silence And Respect by Reply All from desktop or your mobile device",
		CleanedText:     "In 2012, a woman named Lindsey Stone posted a picture she took as a joke to her Facebook page. A month later, she was under attack from all corners of the internet, out of a job, hounded by the press. The internet had targeted her for a public shaming. Jon Ronson, journalist and author of the new book \"So You've Been Publicly Shamed\", walks us through Lindsey's story and introduces us to the sometimes sketchy world of online reputation management.",
		MetaKeywords:    "record, sounds, share, sound, audio, tracks, music, soundcloud",
		CanonicalLink:   "https://soundcloud.com/replyall/18-silence-and-respect",
		TopImage:        "https://i1.sndcdn.com/artworks-000112044299-u970sx-t500x500.jpg",
	}
	//article.Links = []string{}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func Test_TechCrunchCom(t *testing.T) {
	article := Article{
		Domain:          "techcrunch.com",
		Title:           "Gmail Will Soon Warn Users When Emails Arrive Over Unencrypted Connections",
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
		Title:           "\"Disney is counting on 300 million tourists to flock to its new Shanghai theme park https",
		MetaDescription: "",
		CleanedText:     "Language:\n\nHave an account?\n\nBloomberg Business\n\n@\n\nBloomberg Business\n\n@\n\nThe first word in business news.\n\nDisney is counting on 300 million tourists to flock to its new Shanghai theme park span\n\nLoading seems to be taking a while.\n\nTwitter may be over capacity or experiencing a momentary hiccup. Try again or visit Twitter Status for more information.\n\nThis has already been marked as containing sensitive content.\n\nFlag this as containing potentially illegal content.\n\nList name\n\nDescription\n\nPublic · Anyone can follow this list\n\nPrivate · Only you can access this list\n\nThe URL of this tweet is below. Copy it to easily share with friends.\n\nAdd this Tweet to your website by copying the code below. Learn more\n\nAdd this video to your website by copying the code below. Learn more\n\nInclude parent Tweet\n\nInclude media\n\nForgot password?\n\nNot on Twitter? Sign up, tune into the things you care about, and get updates as they happen.\n\nSign up\n\nCountry\n\nCode\n\nFor customers of\n\nUnited States\n\n40404\n\n(any)\n\nCanada\n\n21212\n\n(any)\n\nUnited Kingdom\n\n86444\n\nVodafone, Orange, 3, O2\n\nBrazil\n\n40404\n\nNextel, TIM\n\nHaiti\n\n40404\n\nDigicel, Voila\n\nIreland\n\n51210\n\nVodafone, O2\n\nIndia\n\n53000\n\nBharti Airtel, Videocon, Reliance\n\nIndonesia\n\n89887\n\nAXIS, 3, Telkomsel, Indosat, XL Axiata\n\nItaly\n\n4880804\n\nWind\n\n3424486444\n\nVodafone\n\n» See SMS short codes for other countries\n\nHmm... Something went wrong. Please try again.",
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
		CleanedText:     "President Obama signed into law a bipartisan budget bill last week that, among other things, changes —\u00a0for better and worse — Social Security and Medicare laws. Here's a wrap-up:\n\n•\u00a0File and suspend.\u00a0Currently,\u00a0a married person — typically the higher wage earner in a couple — who's\u00a0at least full retirement age could file for his or her own Social Security benefits and then immediately suspend those benefits while the spouse could file\u00a0for spousal benefits. By doing this, the higher wage earner’s benefits would grow 8% per year. In the meantime, the couple still get\u00a0a Social Security check, and down the road the surviving spouse could get a higher benefit.\n\nThat option is ending for new filers starting May 1, 2016, so if you're\u00a0interested, now's the time to apply. People already using\u00a0this strategy will be grandfathered in until age 70.\n\nUSA TODAY\n\nFull retirement age is a magic number for Social Security benefits\n\n•\u00a0Restricted application.\u00a0\u00a0This is also being phased out.\u00a0Currently, individuals\u00a0eligible for both a spousal benefit based their spouse's work record and a retirement benefit based on his or her own work record could choose to elect only a spousal benefit at full retirement age, according to Social Security Timing. That would let them collect a higher benefit later on.\n\nUnder the new law, however, only those born Jan. 1, 1954, or earlier can use this option. Anyone younger will\u00a0just automatically get the larger of the two benefits,\u00a0according to Social Security Timing.\n\n•\u00a0Social Security Disability.\u00a0\u00a0The Social Security Disability trust was on pace\u00a0to run out money next year and, as a result, millions of Americans were going to receive an automatic 19% reduction in their disability benefits in the fourth quarter of 2016. The new law fixes that\u00a0by shifting payroll tax revenue from one Social Security trust fund —\u00a0the Old-Age and Survivors Insurance Trust fund —\u00a0to another,\u00a0the Disability Insurance Trust fund.\n\nUSA TODAY\n\nRetirement: When you should take Social Security\n\n•\u00a0Medicare Part B.\u00a0Some 30% of Medicare beneficiaries were expecting a 52% increase in their Medicare Part B medical insurance premiums and deductible\u00a0in 2016.\u00a0Under the new law, those beneficiaries —\u00a0an estimated 17 million Americans —\u00a0will pay about $119\u00a0per month, instead of $159.30, for Part B. (Some 70% of Medicare beneficiaries will continue to pay the same premium in 2016 as they did in 2015, $104.90.)\n\nBeneficiaries, however, will also have to pay an extra $3 per month to help pay down a loan the government gave to Medicare to offset lost revenue.\u00a0 Plus, all Part B beneficiaries will see their annual deductible increase by 15% to about $166\u00a0in 2016.\n\nRobert Powell is editor of Retirement Weekly, contributes regularly to MarketWatch, The Wall Street Journal, USA TODAY, and teaches at Boston University.",
		MetaKeywords:    "",
		CanonicalLink:   "http://www.usatoday.com/story/money/columnist/powell/2015/11/12/social-security-medicare-changes-budget-law-retirement/75164246/",
		TopImage:        "http://www.gannett-cdn.com/-mm-/eba3ab7ada1c4fcc1a671898ecfb68274260e9c9/c=0-48-508-335&r=x633&c=1200x630/local/-/media/2015/02/24/USATODAY/USATODAY/635603784536631512-177533853.jpg",
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
		CleanedText:     "For President Obama, it's legacy time.\n\nWith\u00a0less than a year before his successor is elected and he officially becomes a\u00a0lame-duck president, time is running short. Obama has moved the ball forward on a number of legacy items already this year. Some have solidified; others remain in limbo.\n\nHis 2010 health care reform law will already be mentioned at the top of the 44th president's Wikipedia page. But the Obama White House is moving quickly\u00a0on a number of issues that could be\u00a0listed\u00a0in the\u00a0first few paragraphs, too.\n\n\"You do get a sense they are aware of the legacy, and there is a kind of a presidential scorecard being filled out,\" says Gil Troy, a visiting fellow at the Brookings Institution.\n\nObama's ambitions are high. They start with one last shot at the seemingly impossible task of closing the prison in Guantanamo Bay, Cuba, and cover everything from signing an international climate change deal to finalizing\u00a0one of the world's largest free-trade agreements in a generation.\n\nIt's notable that most\u00a0of Obama's goals are abroad; that's because a Republican-controlled Congress has less authority to intervene.\u00a0But that doesn't mean crossing things off his final to-do list is going to be easy. The scope of what he wants to do means finishing it will take a\u00a0lot of late nights for Obama and his staff in their\u00a0final year, said Jacob Stokes, an associate fellow at the bipartisan Center for a New American Security.\n\nBut if the stars align for the president — as they seem to have done this summer — Stokes thinks Obama can get most of them done.\u00a0\"The president and the administration have a relatively large amount of agency to get these things done,\" he said. \"If they really focus on it.\"\n\nHere are seven things on Obama's final to-do list.\n\nShuttering Guantanamo is less of a legacy issue and more of a moral one for the president, Stokes said. Since the first days of his presidency, Obama has maintained the prison, where men can be held indefinitely, is a propaganda tool for terrorists. But congressional Republicans say closing it will create more risk than it's worth, and they — and the realities of what to do with existing prisoners there — have successfully blocked the president for six years from doing anything about it.\n\nThe clock's ticking for Obama to fulfill one of his oldest\u00a0campaign promises.\n\nHe's\u00a0planning a final standoff with Congress, by dropping a\u00a0plan as soon as this week to close Guantanamo without Congress's help. (On Tuesday, Congress passed its annual defense spending bill that, per usual, restricts the president from transferring detainees to the United States.)\n\nBut Obama must decide how badly he wants Guantanamo closed. Trying to transfer the remaining 112 prisoners by himself could start a much broader fight with Republicans over the president's constitutional power. Sen. John McCain (R-Ariz.) is threatening to sue the president if he acts over Congress's wishes.\n\nObama already scored a major legislative victory this summer when he persuaded\u00a0enough congressional Democrats — yes, he was working against much of his own party on this one — to give him authority to negotiate the Trans-Pacific Partnership without congressional say-so on\u00a0every little detail.\n\nHis job is only half done, though. After the United States and 11 other nations came to an agreement on the deal in October, Obama now needs to sway enough lawmakers on both sides to approve the whole package.\n\nLawmakers will soon review the sprawling deal and could vote on it this spring at the earliest. Getting it passed is going to be an uphill battle for Obama, reports The Washington Post's David Nakamura. Liberal Democrats are concerned about the trade deal's environmental impacts and potential drag on U.S. manufacturing jobs, while some Republicans worry the deal isn't strong enough.\n\nThe stakes are also higher for Obama than simply completing\u00a0the biggest free-trade deal in modern history. This trade agreement is a major economic cornerstone of Obama's pivot to Asia, Stokes said. Without TPP, he'll lose one of his most concrete examples\u00a0of a shift to Asia that has struggled to take shape.\n\nHere's one place Obama might\u00a0not need to do battle with Congress. Whatever comes out of a major United Nations summit on climate change held in Paris at the end of this month will likely not have to ratified by the Senate.\n\nThat's a good thing, because Obama didn't make any friends on Capitol Hill on Friday when he announced he won't approve an extension to the Keystone XL oil pipeline from Canada. In bucking a Republican priority and taking the environmentalists' side, Obama indicated he's ready to make some serious changes to U.S.\u00a0pollution levels, which rank among the top in the world.\n\nRejecting Keystone demonstrates to the rest of the world that Obama is \"willing to pull out all the stops on climate change,\" writes The Washington Post's Chris Mooney.\n\nPlus, a 2009 climate change meeting of world leaders in Copenhagen was kind of a bust, so Paris could be Obama's last chance to effect any meaningful change on the world stage.\n\n\"If they miss putting something together in Paris, it's going to be very tough to do anything beyond that,\" Stokes said.\n\nObama is leaning in on Syria in his final months in office. In addition to stepping up airstrikes on the Islamic State there, he announced in October he's putting 50 Special Ops forces\u00a0on the ground, appearing to go back on his past statements\u00a0he wouldn't commit ground troops to Syria.\n\nThis is happening as\u00a0Russia jumped into the\u00a0Syrian conflict, aiming to help Syrian President Bashar al-Assad keep control of his crumbling country. That's a major problem for the White House, which is facing an already no-win situation in an increasingly violent Middle East.\n\nBut Russia's sudden prioritization of the region could be just the crisis the world and Obama need to find a solution, Stokes said. Already, it\u00a0forced diplomats with stakes in the region to a hastily organized conference in Vienna last week to talk about what to do, he noted.\n\n\"I think there's a sense that a political agreement is not imminent by any stretch of the imagination,\" Stokes said, \"but that in the next year you may get parties into a region where they can start thinking more broadly.\"\n\nEither way, Obama would really like to leave office without a civil war in Syria — something which\u00a0is fueling the Islamic State's movement — still raging.\n\nYes, Obama announced a historic nuclear agreement with Iran in June, and yes, he managed to avoid a reluctant Congress from blocking it in September.\n\nBut the deal is still mostly on-paper, which means several GOP presidential candidates' campaign promises to \"rip it to shreds\" — as Sen. Ted Cruz (Tex.) likes to say — are legitimate threats.\n\nThat is, unless Obama can spend the next year or so setting\u00a0key elements of the deal in motion. That would make\u00a0it much tougher for another president to come along and undo one of his biggest foreign policy achievements, Stokes said.\n\nThese days, Obama and Republicans celebrate when they can agree on a budget just\u00a0to keep the government running. So there's little hope they'll come to an agreement on the president's other major domestic policy goals, like immigration reform.\n\nOne bright spot is criminal justice reform. A bipartisan bill to change federal sentencing mandates is moving quickly in the Senate and has the potential for bipartisan support in the House of Representatives, too.\n\nObama has made reforming sentencing laws a priority recently.\u00a0In 2014, then-Attorney General Eric Holder announced the department would stop charging nonviolent drug offenders with crimes that require judges to enact so-called mandatory minimum sentences. And the Justice Department recently released 6,000 federal prisoners, the largest one-time release ever, who were sentenced for non-violent drug crimes.\n\nNot just any court case, mind you. After 26 states challenged his executive actions on immigration, Obama is betting it all on the Supreme Court.\n\nA federal court upheld the states' challenge on Monday, and by Tuesday, the White House confirmed it would ask the Supreme Court to rule next year on whether he stayed within his constitutionally limited powers by deferring deportations for millions of young immigrants and some of their parents.\n\nIf the\u00a0Supreme Court takes up the case, it could rule by June, leaving just months for the administration to start enrolling immigrants and create a buffer for whoever comes into the White House next — and whatever vision of Obama's they might\u00a0try to undo.",
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
		Title:           "Strategy and Entrepreneurship",
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
		CleanedText:     "WASHINGTON—President Barack Obama ’s biggest campaign donors are mostly sitting on the sidelines of the 2016 Democratic presidential primary so far, not opening their wallets in support of Hillary Clinton or Bernie Sanders.\n\nAlmost four-fifths of the people who gave the 2012 maximum $5,000 to the president’s re-election committee hadn’t donated to a presidential candidate by Oct. 1, a Wall Street Journal analysis of federal campaign finance records found.\n\nIn interviews ahead of this Saturday’s Democratic debate in Iowa, donors said Mrs. Clinton, the party’s front-runner, hadn’t motivated them to give the way Mr. Obama and previous Democratic candidates had. Still others said they are put off by the larger role of super PACs and that their donations to candidates, which are limited in this election cycle to $5,400 for the eventual nominee, just don’t matter much anymore.\n\nRelated\n\nSome Candidates, Super PACs Draw Closer (Oct. 25)\n\nThe 2016 Money Race (Oct. 15)\n\n“I’m just not ready for Hillary yet,” said Robert Finnell, a Rome, Ga., lawyer who gave the maximum allowed contribution to Mr. Obama’s 2008 and 2012 campaigns and gave significant sums to 2008 hopeful John Edwards and 2004 Democratic nominee John Kerry. “It’s not that I don’t think she’s competent—she is competent, she’s just hard to like.”\n\nThe donors’ reluctance could be a troubling trend for Mrs. Clinton. They are some of the easiest prospective contributors to identify, given that their names are on Mr. Obama’s campaign disclosure reports, and that they’ve already made a habit of cutting checks to politicians.\n\nJulianna Smoot, finance director on President Obama’s 2008 campaign and deputy campaign manager of his re-election effort, said: “Most Democrats will be behind Hillary if she’s the nominee. Once that becomes clear, the rest of that money should be easy for her to get. I do think these folks will be there.”\n\nMrs. Clinton has outpaced Mr. Obama’s fundraising in the first two quarters of his initial presidential campaign. The former secretary of state has raised $77.5 million for those six months through October. By July 2007, when Mr. Obama had been in the race a comparable length of time, he had raised $58.9 million.\n\nIn 2012, roughly 4,000 individuals donated the maximum to Mr. Obama’s campaign committee, delivering $20 million to his account, according to disclosure reports filed with the Federal Election Commission. Of them, about 830 can be identified as having donated to a candidate in the 2016 presidential race. Mrs. Clinton is the largest recipient of their money at $1.8 million. The big Obama donors gave about $109,000 to Mr. Sanders, about $94,000 to former Maryland Gov. Martin O’Malley, and about $70,000 to Republican Jeb Bush.\n\nFor the analysis, The Wall Street Journal cross-referenced a list of individuals who had donated the maximum amount to Mr. Obama in 2012 with those who have given to candidates in the current presidential race. The maximum donation total is based on rules that allow a donor to give a candidate up to $2,700 each for the 2016 primary and general election.\n\nMichael Briggs, a spokesman for the campaign of Mr. Sanders, said he expected to be outspent and that “he’s taking on the establishment and does not expect the establishment to support that.” The Clinton campaign’s spokesman, Josh Schwerin, said: “Thanks to the support of hundreds of thousands of people, we have been able to raise a record amount for a nonincumbent during our first two quarters in the race.”\n\nSome people inclined to support Mrs. Clinton note that it is still early in the race and the Republican field remains unsettled. “I don’t think she needs the money right now,” said Jeff Choney, a retired high-school teacher in Wellesley, Mass., who gave $5,000 to Mr. Obama in 2012 and said he may contribute to Mrs. Clinton’s later in the cycle. “I like Bernie Sanders—he speaks the truth on a lot of things. But I don’t think he has a chance of beating her, so I’m not so worried about her campaign.”\n\nStill, the Clinton campaign is building a national campaign apparatus that will be expensive to maintain through the general election, should she win the party nomination. Mr. Obama built a similar operation incrementally during the extended 2007 Democratic primary contest. As of Oct. 1, the last records available to the public, Hillary for America had spent $44.5 million, compared with $14.3 million for Mr. Sanders and $20.1 million for retired neurosurgeon Ben Carson, the current fundraising leader in the Republican field.\n\nMrs. Clinton is also relying on support from super PACs, which can raise and spend unlimited sums as long as they don’t coordinate with her campaign. One of the largest of those, Priorities USA Action, raised $15.7 million as of July.\n\nThe limitations of super PACs have been on display in the GOP primary. Former Texas Gov. Rick Perry and Wisconsin Gov. Scott Walker dropped out of the race after struggling to pay campaign expenses that can’t be covered by an outside group.\n\nThough a super PAC backing former Florida Gov. Jeb Bush raised $103.2 million as of July, his campaign in October still had to cut staff salaries and trim the head count at its Miami headquarters.\n\nThe super PAC restrictions put pressure on Mrs. Clinton—and all candidates—to raise as much money as they can for their own campaign accounts. Meanwhile, some donors are demoralized witnessing the big checks pouring into super PACs.\n\n“Even though I gave the maximum [in 2012], it’s nothing compared with what these PACs do. I certainly don’t see my contribution as significant,” said Marilynn Duker, president of a Baltimore residential development and property management company. She gave to Mr. Obama’s 2008 campaign and Mr. Kerry in 2004, but has yet to donate to a White House hopeful in this cycle. “It has no real meaning relative to the gazillions of dollars that the PACs contribute to the races these days. I just don’t feel like the individual really makes a difference,” she said.\n\nDoug Curling, an Atlanta-area executive, gave significant sums to Mr. Obama’s campaigns and gave Mrs. Clinton the maximum donation in 2007. But he said he and his wife now plan to contribute to groups advocating for structural change in the political system. “Nobody needs our money,” Mr. Curling said. “I wouldn’t misinterpret it as we’re disenfranchised from our party, it’s more we’re disenfranchised from the system.”\n\nPeter Maroney, who was the Democratic National Committee’s national finance co-chairman in the 2004 presidential campaign, said many Democratic donors had been waiting on Vice President Joe Biden. “Now that the vice president has made his decision, this is an opportunity for candidates like Mrs. Clinton to proactively go after these donors and make them feel that they have a seat at her table,” he said.\n\nWrite to Daniel Nasaw at daniel.nasaw@wsj.com",
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
		CleanedText:     ".\n\nView photo\n\nLOS ANGELES (Reuters) - El Nino's warm currents have brought fish in an unexpected spectrum of shapes and colors from Mexican waters to the ocean off California's coast, thrilling scientists with the sight of bright tropical species and giving anglers the chance of a once-in-a-lifetime big catch.\n\nCreatures that have made a splash by venturing north in the past several weeks range from a whale shark, a gentle plankton-eating giant that ranks as the world's largest fish and was seen off Southern California, to two palm-sized pufferfish, a species with large and endearing eyes, that washed ashore on the state's central coast.\n\nScientists say El Nino, a periodic warming of ocean surface temperatures in the eastern and central Pacific, has sent warm waves to California's coastal waters that make them more hospitable to fish from the tropics.\n\nEl Nino is also expected to bring some relief to the state's devastating four-year drought by triggering heavy rains onshore.\n\nBut so far precipitation has been modest, and researchers say the northern migration of fish in the Pacific Ocean has been one of the most dynamic, albeit temporary, effects of the climate phenomenon.\n\nEven as marine biologists up and down the coast gleefully alert one another to each new, rare sighting, the arrival of large numbers of big fish such as wahoo and yellowtail has also invigorated California's saltwater sport fishing industry, which generates an estimated $1.8 billion a year.\n\n\"Every tropical fish seems to have punched their ticket for Southern California,\" said Milton Love, a marine science researcher at the University of California, Santa Barbara.\n\nSome fish made the journey north as larva, drifting on ocean currents, before they grew up, researchers said.\n\nThe first ever sighting off California's coast of a largemouth blenny fish was made over the summer near San Diego, said Phil Hastings, a curator of marine vertebrates at the Scripps Institution of Oceanography.\n\nThat species had previously only been seen further south, he said, off Mexico's Baja California.\n\nSmall, colorful cardinalfish were also spotted this year off San Diego, while spotfin burrfish, a rounded and spiny species, were sighted off the coast of Los Angeles, said Rick Feeney, a fish expert at the Natural History Museum of Los Angeles County.\n\nThose tropical species are hardly ever found in Californian waters, he said.\n\n'NEVER SEEN IT LIKE THIS'\n\nSome small tropical fish could remain in the state's waters over the coming months, researchers said, as El Nino is expected to last until early next year.\n\n\"As soon as the water gets cold, or as soon as they get eaten by something else, we'll never see them again,\" Love said.\n\nFor sports fishers, it was so-called pelagic zone fish like wahoo, that live neither close to the bottom nor near the shore, which made this year special.\n\nBefore the El Nino, California anglers only saw wahoo, a fish with a beak-like snout and a slim body that often measures more than 5 feet (1.5 meters) in length, when they made boat trips south to Mexican waters.\n\nThis year, there were 256 recorded catches of wahoo by sport fishing party boats from Southern California, with almost all of those being taken on the U.S. side of the border, said Chad Woods, founder of the tracking company Sportfishingreport.com.\n\nLast year, he said, the same boats made 42 wahoo catches.\n\nMichael Franklin, 56, a dock master for Marina Del Rey Sportfishing near Los Angeles in the Santa Monica Bay, said this was the best year he can remember, with plentiful catches of yellowtail and marlin.\n\n\"I've been fishing this bay all my life since I was old enough to fish, and I've never seen it like this,\" he said.\n\nMany hammerhead sharks also cruised into Californian waters because of El Nino, experts say.\n\nSport fisherman Rick DeVoe, 46, said he took a group of children out in his boat off the Southern California coast this September. A hammerhead followed them, chomping in half any tuna they tried to reel in.\n\n\"The kids were freaking out because the shark's going around our boat like 'Jaws',\" DeVoe said.\n\n(Reporting by Alex Dobuzinskis; Editing by Daniel Wallis and Andrew Hay)",
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

// https://jiradatasift.atlassian.net/browse/DEV-4510
func Test_Dev4510(t *testing.T) {
	article := Article{
		Domain:          "dev4510",
		Title:           "'I Was A Teenage Cyclist,' or How Anti-Bike-Lane Arguments Echo the Tea Party",
		MetaDescription: "Anti-bike-lane arguments often mirror the rhetorical tactics of the Tea Party: The appeal to an imagined golden age of yesteryear, reliance on dismissive shorthand and, most strikingly, warnings of a creeping, foreign-based anti-Americanism that’s plainly contrary to our core values.",
		CleanedText:     "If you’re itching to write an anti-bike-lane argument (and, if so, line up, because it’s a burgeoning literary genre), you could do no better than to follow the template laid out yesterday by The New Yorker’s John Cassidy in his blog post, “Battle of the Bike Lanes.”\n\nCassidy’s post — which has already been called “a seminal document of New York City’s bike lane backlash era” — helpfully includes all the requisite rhetorical tactics, thus providing an excellent blueprint. (You might even say “boilerplate.”) These include:\n\nPre-emptive self-exoneration: “I don’t have anything against bikes.”\n\nInvocation of humorlessness of cycling advocates, preferably with ironic comparison to homicidal political faction: “the bicycle lobby … pursues its agenda with about as much modesty and humor as the Jacobins pursued theirs.”\n\nReference to ominous encroachment of cycling-based anti-Americanism: “City Hall … sometimes seems intent on turning New York into Amsterdam, or perhaps Beijing.” (You know, Beijing: where the communists live!)\n\nInvocation of personal cycling bona fides: “As a student, I lived in the middle of Oxford, where cycling is the predominant mode of transport, and I cycled everywhere.”\n\nFond nostalgia for pre-lane New York City cycling perils, coupled with implied dismissal of today’s namby-pamby cyclists: “In those days … part of the thrill was avoiding cabs and other vehicles. … When I got back to my apartment on East 12th Street, I was sometimes shaking.”\n\nOddly self-contradictory declaration of support: “Generally speaking, I don’t have a problem with this movement; indeed, I support it.”\n\nInvocation of meddling government apparatchiks: “A classic case of regulatory capture by a small faddish minority.”\n\nInvocation of America’s long, sun-dappled love affair with cars: “Since 1989, when I nervously edged out of the Ford showroom on 11th Avenue and 57th Street, the proud leaser of a sporty Thunderbird coupe, I have owned and driven six cars in the city.”\n\nInvocation of obviously repellent stereotype: “I would put my knowledge of New York’s geography and topography up against most native residents’ — cycling members of the Park Slope food co-op included.” (To be fair, if you’ve ever been to the Park Slope food co-op, you know how its members are always prattling on about their topographical expertise.)\n\nBrief feint toward fact-based argument, unencumbered by actual facts: “From an economic perspective I also question whether the blanketing of the city with bike lanes … meets an objective cost-benefit criterion. … Beyond a certain point … the benefits of extra bike lanes must run into diminishing returns.” (Yes. They must. But when? At what point? Sorry — no time! Moving on!)\n\nFollowed by quick return to actual motivation: “Like many New Yorkers who don’t live in Manhattan, one of my favorite pastimes is to drive from Brooklyn … into the city for dinner to find a parking space once the 7 a.m. – 7 p.m. parking restrictions have lapsed. … These days, [this] is virtually impossible.” (A lack of parking spaces naturally serving as evidence of too many bike lanes, not too many parked cars.)\n\nInvocation of damnable scofflaw cyclists: “On those rare occasions when I do happen across a cyclist, or two, he or she invariably runs the red lights.” (On a related note, I personally witnessed three hit-and-run accidents outside my old apartment at Atlantic Ave. and 3rd Avenue in Brooklyn. I logically determined that drivers invariably get into accidents, and thus launched my campaign for the eradication of city streets.)\n\nOne last invocation of overreaching City Hall bureaucrats, for good measure: “[I]t is time to call a halt to Sadik-Kahn and her faceless road swipers.”\n\nSee? It’s easy. Or, if this all seems too strenuous or, you know, long-winded, you can simply reduce your argument to its four essential words: “I have been inconvenienced.”\n\nAs an occasional cycling commuter, I’m always struck (no pun intended) by the extent to which arguments like Cassidy’s mirror the rhetorical tactics of the Tea Party. (No small accusation, I understand.) For example: The appeal to an imagined golden age of yesteryear (gamely dodging cabs; Thunderbird coupes); the specter of bureaucracy run amok (the scourge of the faceless road swipers); reliance on dismissive shorthand (Park Slope co-op members); and, most strikingly, warnings of a creeping, foreign-based anti-Americanism that’s plainly contrary to our core values (They Came on Bikes From Beijing).\n\nThese parallel lines of reasoning were finally entangled last year in the gubernatorial campaign of the Colorado Republican Dan Maes, who warned that the pro-bike policies of his opponent,\u00a0 Mayor John Hickenlooper of Denver, were turning that city “into a United Nations community,” adding ominously, “This is bigger than it looks on the surface, and it could threaten our personal freedoms.” (Maes eventually lost the race for governor to Hickenlooper by a margin of 51 percent to 11 percent.)\n\nAll of which is to note: The discussion over cycling policy in New York has now taken on the tone (on both sides, sadly) of our culture wars: passion first, reason later (or, in most cases, never).\n\nSo in a spirit of understanding, I encourage you to read Cassidy’s article in full. You can also read these two (relatively) measured and enjoyable rebuttals, as well this well-balanced look at the bike-lane controversies in Brooklyn.\n\nAnd, if you’re interested in facts — yes! facts! — I would also point you toward this excellent long-form piece on cycling commuting by Tom Vanderbilt, author of the book “Traffic.” Here are two interesting statistics he mentions: 1) Portland, Ore., the American city with arguably the most progressive cycling policy, had exactly zero cycling traffic fatalities in 2010. (New York had 18.) And 2) closer to home, Vanderbilt points out that, since the implementation of New York’s Ninth Avenue dedicated bike lane, pedestrian injuries have gone down by 29 percent. That’s not accidents between bikes and people; that’s between cars and people.\n\nThese facts are interesting to contemplate. Or, failing that, there’s always: Road-swipers! Thunderbirds!! COMMUNISTS!!!\n\nAn earlier version of this posting misstated the given name of John Hickenlooper, the governor of Colorado.",
		MetaKeywords:    "Bicycles and Bicycling,New York City,Uncategorized",
		CanonicalLink:   "http://6thfloor.blogs.nytimes.com/2011/03/09/i-was-a-teenage-cyclist-or-how-anti-bike-lane-arguments-echo-the-tea-party/",
		TopImage:        "http://graphics8.nytimes.com/images/blogs_v5/../icons/t_logo_291_black.png",
	}
	article.Links = []string{
		"http://www.nypost.com/p/news/local/janette_big_transitway_road_to_ruin_V6obl2EErgSaSZtg04Lr1K",
		"http://www.nypost.com/p/news/opinion/editorials/we_janette_6ZhwHlxPxnIZzli8wjNrTM",
		"http://gothamist.com/2011/03/05/mayor_weiners_first_act_would_aboli.php",
		"http://www.newyorker.com/online/blogs/johncassidy/2011/03/battle-of-the-bike-lanes-im-with-mrs-schumer.html",
		"http://naparstek.com/2011/03/bike-lane-backlash-makes-no-sense/",
		"http://en.wikipedia.org/wiki/Reign_of_Terror",
		"http://www.politicsdaily.com/tag/Dan%20Maes/",
		"http://www.newyorker.com/online/blogs/johncassidy/2011/03/battle-of-the-bike-lanes-im-with-mrs-schumer.html",
		"http://blogs.reuters.com/felix-salmon/2011/03/09/john-cassidy-vs-bipeds/",
		"http://naparstek.com/2011/03/bike-lane-backlash-makes-no-sense/",
		"http://www.nytimes.com/2011/03/09/nyregion/09bike.html",
		"http://outsideonline.com/adventure/travel-ga-201103-new-york-bike-commuting-sidwcmdev_154507.html",
		"http://www.nytimes.com/2011/02/07/nyregion/07safety.html",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

// https://jiradatasift.atlassian.net/browse/DEV-4510
func Test_Dev4510b(t *testing.T) {
	article := Article{
		Domain:          "dev4510b",
		Title:           "In Suspects’ Brussels Neighborhood, a History of Petty Crimes and Missed Chances",
		MetaDescription: "Molenbeek, a known haven for extremists, was home to Ibrahim Abdeslam, who blew himself up in Paris. Although the authorities there had him in their sights, he slipped through their fingers.",
		CleanedText:     "BRUSSELS — Just eight days before Ibrahim Abdeslam blew himself up in Paris as part of an elaborate terrorist operation that killed 129 people on Friday, the authorities in the heavily immigrant Brussels district of Molenbeek already had the future terrorist in their sights.\n\nUnfortunately, they had identified him not as a potential killer, but as the proprietor of a bar that played host to drug dealers and drunks. Under an order signed by Molenbeek’s mayor, the bar was shut down on Nov. 5 “for compromising public security and tranquillity” through the spread of illegal substances.\n\nMolenbeek is well known as a haven for extremists, home to dozens of young men accused of leaving to wage jihad in Syria and, in some cases, plotting attacks against Europe. The area has now been linked to at least four terrorist attacks in two years.\n\nBut the inability to stop Mr. Abdeslam was just one example of the missed opportunities by the Belgian and French authorities and intelligence services, a list that also included allowing Mr. Abdeslam’s brother Salah, 26, another suspect, to slip through their fingers.\n\nSalah Abdeslam rented a car in Brussels that was apparently used to transport some of the gunmen who killed 89 people in a Paris concert hall. He had a criminal record, which outlines his suspected involvement in organized crime, but there was no arrest warrant linked to his file. Because of that record, his name popped up during a routine traffic stop by the French police on Saturday. But he was allowed to drive on because he had not yet been linked to the attacks. He remains at large.\n\nThe near misses raise troubling questions about the Belgian intelligence service and their French counterparts, not to mention concerns about Europe’s system of open internal borders, which has allowed terrorists to move freely between countries while outpacing the intelligence sharing needed to stop them.\n\n“Every time there is an attack, we discover that the perpetrators were known to the authorities,” said François Heisbourg, a counterterrorism expert and former defense official. “What this shows is that our intelligence is actually pretty good, but our ability to act on it is limited by the sheer numbers.”\n\nThe missed opportunities before and since the attacks in Paris, intelligence officials and experts say, pale when compared with the fact that the attackers were able, at least in part, to organize their plot under the noses of the authorities in Molenbeek.\n\nMr. Abdeslam, the suicide bomber in Paris, and Abdelhamid Abaaoud, the suspected architect of the attacks, had each lived barely 200 yards from Molenbeek’s main police station and had had brushes with the law.\n\nThe Paris attacks indicate that few real steps have been taken to keep the neighborhood under surveillance adequately and break up its small but lethal extremist underground.\n\nInvestigators believe that the massacres on Friday, the worst terrorist attacks in France, involved at least three people from Molenbeek, including the Abdeslams and Mr. Abaaoud, a foreign fighter in Syria for the Islamic State who investigators believe orchestrated the carnage.\n\nMr. Abaaoud has appeared regularly in gruesome recruiting videos issued by the Islamic State. But by his own account, he, too, managed to slip in and out of Belgium without being arrested, despite being stopped at one point by an officer who “let me go, as he did not see the resemblance” to photos of himself published in the Belgian news media.\n\nPosing for pictures holding an Islamic State flag and the Quran, Mr. Abaaoud boasted in an interview this year with the militant group’s magazine, Dabiq, of outsmarting security services. “We spent months trying to find a way into Europe, and by Allah’s strength, we succeeded in finally making our way to Belgium,” he said.\n\nHe added, “We were then able to obtain weapons and set up a safe house while we planned to carry out operations against the crusaders.”\n\nOne of their biggest allies, however, may have been a Belgian security system ill equipped to deal with a tight knit community like Molenbeek, where a mostly white police force has only tenuous links to a largely immigrant population resentful of being labeled potential terrorists.\n\nThe police have on occasion pounced, but mostly for petty crimes unrelated to Islamist extremism.\n\nIbrahim Abdeslam, the bar operator turned suicide bomber, was convicted of criminal activities as far back as 2010, and even stood trial for minor offenses with Mr. Abaaoud, according to a people briefed on information from the federal prosecutor’s office. The mayor of Molenbeek, Françoise Schepmans, said Mr. Abdeslam had been convicted of crimes involving drugs.\n\n“This is a small place; we all crossed paths with them,” the deputy mayor of Molenbeek, Ahmed el Khannouss, said Monday, referring to the three suspects from the district. But, he added, “we are all totally shocked that they could have been involved in something so terrible.”\n\nThe Expanding Web of Connections Among the Paris Attackers\n\nAs many as six of the assailants in the coordinated Islamic State terrorist assault in Paris were Europeans who had traveled to Syria.\n\nFamily members and friends voiced complete surprise, too, highlighting how difficult it is to penetrate a small group of determined jihadists.\n\nMohamed Abdeslam, the brother of one of the Paris suicide bombers, was picked up on Saturday by the police and released on Monday. An employee with the municipal government, he told reporters in Molenbeek on Monday that he first learned of his brother’s extremist affiliations from news media reports of the Paris attacks.\n\nHe said he previously knew “absolutely nothing, absolutely nothing.”\n\nThe father of Mr. Abaaoud was so shocked and appalled by his son’s embrace of violent jihadism that he sued him in Belgium in May, asserting that the son had “kidnapped” another sibling, who was just 13, and lured him to Syria.\n\n“If there is anyone who is certainly not aware of what is going on, it’s the father,” said the father’s lawyer, Nathalie Gallant.\n\nThe terrorists not only hid their radical views and intentions, but they also benefited from Belgium’s large pool of angry Muslim youths and its longstanding role as a center for the illegal arms trade. The tiny country has seen more than 400 of its citizens leave to fight in Syria, the highest per capita number in Europe.\n\n“Long before jihadism was on the scene, Belgium was a hub for illegal gunrunning,” Mr. Heisbourg said.\n\nTerrorists planning to carry out an attack find it much easier to coordinate and procure weapons from Belgium than anywhere else in Europe, said Jelle van Buuren, a lecturer in counterterrorism at Leiden University in the Netherlands. Belgium “is considered to be the weak link in Europe’s approach to tackling terrorism,” he said.\n\nBelgium is also divided into three different languages and cultures, Mr. Van Buuren said, which makes it hard for agencies to communicate and coordinate with one another, including the sharing of intelligence.\n\nCrucially, Belgium’s strict laws on surveillance, including the interception of telephone conversations, make it hard for the authorities to monitor potential terrorists.\n\nIn contrast, he said, France has wider surveillance on terrorist suspects. “You know that the chances of being discovered are not high in Belgium,” Mr. Van Buuren said. “Meanwhile, Brussels is just a few hundred kilometers away and shares the same French language, and you can plan very easily from there. So why wouldn’t you?”\n\nYet the French intelligence service has recently fared little better.\n\nA French official briefed on the investigation said Mr. Abaaoud had mentioned plans to attack “a concert hall” to a French citizen who had come back from Syria three months ago and was interrogated by security officials.\n\nMr. Abaaoud, the official said, was also in contact with Ismaël Omar Mostefaï, another of last week’s Paris attackers. Mr. Mostefaï traveled to Turkey in 2013, Turkish officials said, and is believed to have crossed into Syria.\n\nHis name was flagged to French officials twice, in December and again in June, Turkish officials said. But until after Friday’s attack, there was never any follow-up, they said.\n\nThe two brothers who fatally shot 12 people in the office of the satirical magazine Charlie Hebdo in Paris in January were well known to the authorities in the United States and France and had been under surveillance for long stretches. They struck just a few months after the authorities let the wiretap on their phones expire.\n\nIntelligence officials said there were now so many Europeans either in Syria or with links to Syria that to follow them all had become impossible. “We just haven’t got the resources,” one senior European official said.\n\nBut many of these extremists have clustered in a small area, a phenomenon that has long been noticed by terrorism experts. Few areas have been linked to quite so many bloody episodes as Molenbeek, which is relatively poor with a high unemployment rate but, far from a slum, is full of handsome homes, galleries and restaurants along with halal butchers and kebab houses.\n\nMany of the people linked to a suspected terrorist hide-out broken up in January in eastern Belgian came from Molenbeek.\n\nAmedy Coulibaly, who was involved in the Charlie Hebdo attack and an assault on shoppers in a Jewish supermarket in Paris, is believed to have bought weapons in Molenbeek. So did Mehdi Nemmouche, a Frenchman who targeted Jews at a Brussels museum in 2014, killing four. Ayoub el Khazzani, a Moroccan, who was thwarted in his attempt to attack passengers on a high-speed train traveling between Brussels in Paris in August, is also thought to have lived there for a while.\n\nBut many locals still believe their neighborhood has been unfairly maligned. “This is not a Molenbeek problem; it is a global problem,” said Mustafa Zoufri, who directs a local youth center.\n\nLocal officials on Monday denied turning a blind eye to extremism. The mayor, Ms. Schepmans, whose office looks out on the government-owned apartment building where Mr. Abdeslam’s family lived, insisted that while there might have been a period of “denial” under her predecessor, her own administration had “worked hard to fight radicalization.”\n\nStill, she acknowledged that while officials and police officers kept tabs on formally registered mosques, a plethora of small worship halls operated in the shadows with little supervision.\n\nOn Monday, scores of armed officers wearing black face masks sealed off a street that runs by the district’s biggest mosque as they hunted, in vain, for Salah Abdeslam, who rented the car used in the attacks.\n\nAs reporters from around the world swarmed into the borough’s cobblestoned square, a resident screamed at them, denouncing the news media and security officials for linking her neighborhood to terrorism.\n\n“You are all scum,” she shouted. “Do you have no shame?”",
		MetaKeywords:    "Paris Attacks (November 2015),Brussels (Belgium),Abdeslam  Ibrahim,Terrorism,Charlie Hebdo,Abdeslam  Salah,Abdeslam  Mohamed,Islamic State in Iraq and Syria (ISIS),Abaaoud  Abdelhamid",
		CanonicalLink:   "http://www.nytimes.com/2015/11/17/world/europe/in-suspects-brussels-neighborhood-a-history-of-petty-crimes-and-missed-chances.html",
		TopImage:        "http://static01.nyt.com/images/icons/t_logo_291_black.png",
	}
	article.Links = []string{
		"http://www.nytimes.com/interactive/2015/11/13/world/europe/paris-shooting-attacks.html",
		"http://www.nytimes.com/2015/01/25/world/europe/belgium-confronts-the-jihadist-danger-within.html",
		"http://www.nytimes.com/news-event/attacks-in-paris?inline=nyt-classifier",
		"http://www.nytimes.com/news-event/attacks-in-paris?inline=nyt-classifier",
		"http://www.theguardian.com/world/europe-news",
		"http://www.nytimes.com/2015/01/08/world/europe/charlie-hebdo-paris-shooting.html",
		"http://topics.nytimes.com/top/reference/timestopics/subjects/h/high_speed_rail_projects/index.html?inline=nyt-classifier",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

// https://jiradatasift.atlassian.net/browse/DEV-4510
func Test_Dev4510c(t *testing.T) {
	// https://www.rt.com/uk/322589-council-cuts-cameron-pmqs/
	article := Article{
		Domain:          "dev4510c",
		Title:           "‘Anti-austerity champion’ Cameron mocked in Parliament for council cuts hypocrisy (VIDEO)",
		MetaDescription: "Prime Minister David Cameron was mocked in Parliament for being “the new leader of the anti-austerity movement in Oxfordshire” after writing to his constituency council to complain about cuts to services.",
		CleanedText:     "Keep up with the news by installing RT’s extension for Chrome. Never miss a story with this clean and simple app that delivers the latest headlines to you.\n\nMoscow ready to help Hollande get Syria-Turkey border closed to stop fueling militants - Lavrov\n\n‘Anti-austerity champion’ Cameron mocked in Parliament for council cuts hypocrisy (VIDEO)\n\nGet short URL\n\nLabour MP Jonathan Reynolds delivered the blow during prime minister’s questions (PMQs) on Wednesday.\n\nReynolds was referring to a leaked letter in which Cameron protested about cuts to frontline public services by his own Tory-run local council.\n\nThe PM was ridiculed by Labour after the letter went public, with cabinet members welcoming him to the campaign against cuts.\n\n“As the new leader of the anti-austerity movement in Oxfordshire, can the prime minister tell us, how is his campaign going?” Reynolds asked Cameron in the House of Commons, much to MPs’ glee.\n\nI was as surprised as anyone to find out last week that the Prime Minister was an anti-austerity champion, following his...\n\nPosted by Jonathan Reynolds MP on Wednesday, 18 November 2015\n\nThe PM laughed in response, but insisted he wants local councils to make savings.\n\n“What I said to my local council is what I say to every council, which is you’ve got to get more for less, not less for more.\n\n“As I said on this side of the House, we want to make sure that every penny that is raised in council tax is well spent. And if his council would like to come in and get the same advice, I’d gladly oblige.”\n\nLabour MP Yvonne Fovargue also attacked the PM over council cuts.\n\n“Wigan Council has had over a 40 percent cut in its funding over the last five years and lost over a third of its staff,” she told the Commons.\n\n“Does the prime minister advise that I should write to the leader of the council regarding the consequent reduction in services? Or should I place the blame firmly where it belongs, in the hands of your government?”\n\nMP's on both sides laughing at Cameron/anti-austerity leader in Oxford joke. Have they mentioned the 600 austerity-related suicides? #PMQs\n\n— Steve Topple (@MrTopple) November 18, 2015\n\nCameron responded by shifting blame onto the previous Labour government.\n\n“I think if the Right Honorable Lady is looking for someone to blame, she might want to blame the Labour government that left this country with the biggest budget deficit anywhere in the Western world,” he said.\n\nCameron was accused of hypocrisy over his letter to the head of Oxfordshire Council, in which he complained about cuts to elderly day centers, libraries and museums.\n\nCouncil leader Ian Hudspeth replied, describing how the council had already made the cuts to office functions Cameron suggested and noting that new functions had been transferred to the authority, including public health and social care.\n\n#bbcdp Good question about Cameron leading anti-austerity movt in Oxford. PM looked embarrassed. No discussion on Daily Politics\n\n— sue owen (@sueowen3) November 18, 2015\n\n“Excluding schools, our total government grants have fallen from £194 million in 2009/10 to £122 million a year in 2015/16, and are projected to keep falling at a similar rate. I cannot accept your description of a drop in funding of £72 million or 37 percent as a ‘slight fall,’” Hudspeth said.\n\nLabour Shadow Chancellor John McDonnell issued a tongue-in-cheek response to the leaked letter, which was obtained by the Oxford Mail.\n\n“I’m backing David Cameron on this one. He is absolutely right that his chancellor’s cuts to local government are seriously damaging our communities and have to be opposed. I welcome the prime minister as another Tory MP joining our campaign against George Osborne’s cuts,” he quipped.",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "https://cdn.rt.com/files/2015.11/article/564c8e0fc46188150e8b4595.jpg",
	}
	article.Links = []string{
		"https://www.rt.com/news/323404-lavrov-syria-s24-turkey/",
		"http://on.rt.com/6wwt",
		"https://www.facebook.com/JonathanreynoldsMP/",
		"https://www.facebook.com/JonathanreynoldsMP/videos/443633579170356/",
		"https://twitter.com/hashtag/PMQs?src=hash",
		"https://twitter.com/MrTopple/status/666959584536469509",
		"https://twitter.com/hashtag/bbcdp?src=hash",
		"https://twitter.com/sueowen3/status/666970329638690817",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

//
func Test_NytEncodingIssues(t *testing.T) {
	// The following page has some encoding issues: when converting to UTF-8,
	// there are some invalid byte sequences that cause the page to be truncated
	// and the resulting main content is empty.
	// Test that the encoding decoder can skip those invalid byte sequences.
	//
	// http://www.nytimes.com/1977/05/26/movies/moviesspecial/26STAR.html?_r=0
	article := Article{
		Domain:          "nytencodingissues",
		Title:           "'Star Wars': A Trip to a Far Galaxy That's Fun and Funny. . .",
		MetaDescription: "font size=\"-1\" (May 26, 1977) \"Star Wars\" is the most elaborate, most expensive, most beautiful movie serial ever made.",
		CleanedText:     "tar Wars,\" George Lucas's first film since his terrifically successful \"American Graffiti,\" is the\n\nmovie that the teen-agers in \"American Graffiti\" would have broken their necks to see. It's also\n\nthe movie that's going to entertain a lot of contemporary folk who have a soft spot for the\n\nvirtually ritualized manners of comic-book adventure.\n\n\"Star Wars,\" which opened yesterday at the Astor Plaza, Orpheum and other theaters, is the most\n\nelaborate, most expensive, most beautiful movie serial ever made. It's both an apotheosis of\n\n\"Flash Gordon\" serials and a witty critique that makes associations with a variety of literature\n\nthat is nothing if not eclectic: \"Quo Vadis?\", \"Buck Rogers,\" \"Ivanhoe,\" \"Superman,\" \"The\n\nWizard of Oz,\" \"The Gospel According to St. Matthew,\" the legend of King Arthur and the\n\nknights of the Round Table.\n\nAll of these works, of course, had earlier left their marks on the kind of science-fiction comic\n\nstrips that Mr. Lucas, the writer as well as director of \"Star Wars,\" here remembers with affection\n\nof such cheerfulness that he avoids facetiousness. The way definitely not to approach \"Star\n\nWars,\" though, is to expect a film of cosmic implications or to footnote it with so many\n\nreferences that one anticipates it as if it were a literary duty. It's fun and funny.\n\nThe time, according to the opening credit card, is \"a long time ago\" and the setting \"a galaxy far\n\nfar away,\" which gives Mr. Lucas and his associates total freedom to come up with their own\n\nlandscapes, housing, vehicles, weapons, religion, politics--all of which are variations on the\n\nfamiliar.\n\nWhen the film opens, dark times have fallen upon the galactal empire once ruled, we are given to\n\nbelieve, from a kind of space-age Camelot. Against these evil tyrants there is, in progress, a\n\nrebellion led by a certain Princess Leia Organa, a pretty round-faced young woman of old-\n\nfashioned pluck who, before you can catch your breath, has been captured by the guardians of\n\nthe empire. Their object is to retrieve some secret plans that can be the empire's undoing.\n\nThat's about all the plot that anyone of voting age should be required to keep track of. The story\n\nof \"Star Wars\" could be written on the head of a pin and still leave room for the Bible. It is,\n\nrather, a breathless succession of escapes, pursuits, dangerous missions, unexpected encounters,\n\nwith each one ending in some kind of defeat until the final one.\n\nThese adventures involve, among others, an ever-optimistic young man named Luke Skywalker\n\n(Mark Hamill), who is innocent without being naive; Han Solo (Harrison Ford), a free-booting\n\nfreelance, space-ship captain who goes where he can make the most money, and an old mystic\n\nnamed Ben Kenobi (Alec Guinness), one of the last of the Old Guard, a fellow in possession of\n\nwhat's called \"the force,\" a mixture of what appears to be ESP and early Christian faith.\n\nAccompanying these three as they set out to liberate the princess and restore justice to the empire\n\nare a pair of Laurel-and-Hardyish robots. The thin one, who looks like a sort of brass woodman,\n\ntalks in the polished phrases of a valet (\"I'm adroit but I'm not very knowledgeable\"), while the\n\nsquat one, shaped like a portable washing machine, who is the one with the knowledge, simply\n\nsqueaks and blinks his lights. They are the year's best new comedy team.\n\nIn opposition to these good guys are the imperial forces led by someone called the Grand Moff\n\nTarkin (Peter Cushing) and his executive assistant, Lord Darth Vader (David Prowse), a former\n\nstudent of Ben Kenobi who elected to leave heaven sometime before to join the evil ones.\n\nThe true stars of \"Star Wars\" are John Barry, who was responsible for the production design, and\n\nthe people who were responsible for the incredible special effects--space ships, explosions of\n\nstars, space battles, hand-to-hand combat with what appear to be lethal neon swords. I have a\n\nparticular fondness for the look of the interior of a gigantic satellite called the Death Star, a place\n\nfull of the kind of waste space one finds today only in old Fifth Avenue mansions and public\n\nlibraries.\n\nThere's also a very funny sequence in a low-life bar on a remote planet, a frontierlike\n\nestablishment where they serve customers who look like turtles, apes, pythons and various\n\namalgams of same, but draw the line at robots. Says the bartender piously: \"We don't serve\n\ntheir kind here.\"\n\nIt's difficult to judge the performances in a film like this. I suspect that much of the time the\n\nactors had to perform with special effects that were later added in the laboratory. Yet everyone\n\ntreats his material with the proper combination of solemnity and good humor that avoids\n\ncondescension. One of Mr. Lucas's particular achievements is the manner in which he is able to\n\nrecall the tackiness of the old comic strips and serials he loves without making a movie that is,\n\nitself, tacky. \"Star Wars\" is good enough to convince the most skeptical 8-year-old sci-fi buff,\n\nwho is the toughest critic.\n\n\"Star Wars,\" which has been rated PG (\"Parental Guidance Suggested\"), contains a lot of\n\nexplosive action and not a bit of truly disturbing violence.",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "http://graphics8.nytimes.com/images/2002/05/10/movies/10STAR.1.jpg",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

func TestCharsetEucJp(t *testing.T) {
	article := Article{
		Domain:          "charset_euc_jp",
		Title:           "文字コード宣言は行いましょう(HTML)",
		MetaDescription: "",
		CleanedText:     "",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}

/*
func TestCharsetShiftJIS(t *testing.T) {
	article := Article{
		Domain:          "charset_shift_jis",
		Title:           "文字のエンコードを指定する：HTMLタグ辞典 - HTMLタグボード",
		MetaDescription: "",
		CleanedText:     "HTMLの記述形式（文字コード）を正しく設定することによって、ページが読み込まれたときの文字化けを防ぎます。",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "",
	}
	article.Links = []string{}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}
*/

func TestCharsetISO_8859_1(t *testing.T) {
	article := Article{
		Domain:          "charset_iso_8859_1",
		Title:           "httpd-2.2.31.tar.gz: .../manual/configuring.html.de",
		MetaDescription: "",
		CleanedText:     "Caution: In this restricted \"Fossies\" environment the current HTML page may not be correctly presentated and may have some non-functional links.\n\nAlternatively you can here view or download the uninterpreted raw source code.\n\nA member file download can also be achieved by clicking within a package contents listing on the according byte size field. See also the latest Fossies \"Diffs\" side-by-side code changes report for \"configuring.html.de\": 2.2.29_vs_2.2.31.\n\nApache > HTTP-Server > Dokumentation > Version 2.2 Konfigurationsdateien\n\nVerfügbare Sprachen:  de  |\n\nen  |\n\nfr  |\n\nja  |\n\nko  |\n\ntr\n\nDieses Dokument beschreibt die Dateien, die zur Konfiguration des Apache\n\nHTTP Servers verwendet werden.\n\nHauptkonfigurationsdateien\n\nSyntax der Konfigurationsdateien\n\nModule\n\nDer Gültigkeitsbereich von Direktiven\n\n.htaccess-Dateien\n\nKommentare\n\nHauptkonfigurationsdateien\n\nDer Apache wird konfiguriert, indem Direktiven in einfache Textdateien\n\neingetragen werden. Die Hauptkonfigurationsdatei heißt\n\nüblicherweise httpd.conf. Der Ablageort dieser Datei\n\nwird bei der Kompilierung festgelegt, kann jedoch mit der\n\nBefehlszeilenoption -f überschrieben werden. Durch\n\nVerwendung der Direktive Include\n\nkönnen außerdem weitere Konfigurationsdateien hinzugefügt\n\nwerden. Zum Einfügen von mehreren Konfigurationsdateien können\n\nPlatzhalter verwendet werden. Jede Direktive darf in jeder dieser\n\nKonfigurationsdateien angegeben werden. Änderungen in den\n\nHauptkonfigurationsdateien werden vom Apache nur beim Start oder Neustart\n\nerkannt.\n\nDer Server liest auch eine Datei mit MIME-Dokumenttypen ein. Der\n\nName dieser Datei wird durch die Direktive TypesConfig bestimmt. Die Voreinstellung\n\nist mime.types.\n\nSyntax der Konfigurationsdateien\n\nDie Konfigurationsdateien des Apache enthalten eine Direktive pro Zeile.\n\nDer Backslash \"\\\" läßt sich als letztes Zeichen in einer Zeile\n\ndazu verwenden, die Fortsetzung der Direktive in der nächsten Zeile\n\nanzuzeigen. Es darf kein weiteres Zeichen oder Whitespace zwischen dem\n\nBackslash und dem Zeilenende folgen.\n\nIn den Konfigurationsdateien wird bei den Direktiven nicht zwischen\n\nGroß- und Kleinschreibung unterschieden. Bei den Argumenten der\n\nDirektiven wird dagegen oftmals zwischen Groß- und Kleinschreibung\n\ndifferenziert. Zeilen, die mit dem Doppelkreuz \"#\" beginnen, werden als\n\nKommentare betrachtet und ignoriert. Kommentare dürfen\n\nnicht am Ende einer Zeile nach der Direktive\n\neingefügt werden. Leerzeilen und Whitespaces vor einer Direktive\n\nwerden ignoriert. Dadurch lassen sich Direktiven zur besseren Lesbarbeit\n\neinrücken.\n\nSie können die Syntax Ihrer Konfigurationsdateien auf Fehler\n\nprüfen, ohne den Server zu starten, indem Sie apachectl\n\nconfigtest oder die Befehlszeilenoption -t\n\nverwenden.\n\nModule\n\nDer Apache ist ein modularer Server. Das bedeutet, dass nur die abolute\n\nGrundfunktionalität im Kernserver enthalten ist. Weitergehende\n\nFähigkeiten sind mittels Modulen verfügbar,\n\ndie in den Apache geladen werden können. Standardmäßig\n\nwird bei der Kompilierung ein Satz von Basismodulen (Anm.d.Ü.: die so\n\ngenannten Base-Module) in den Server eingebunden. Wenn der\n\nServer für die Verwendung von dynamisch\n\nladbaren Modulen kompiliert wurde, dann können Module separat\n\nkompiliert und jederzeit mittels der Direktive LoadModule hinzugefügt werden.\n\nAndernfalls muss der Apache neu kompiliert werden, um Module\n\nhinzuzufügen oder zu entfernen. Konfigurationsanweisungen können\n\nabhängig vom Vorhandensein eines bestimmten Moduls eingesetzt werden,\n\nindem sie in einen <IfModule> -Block eingeschlossen werden.\n\nUm zu sehen, welche Module momentan in den Server einkompiliert sind,\n\nkann die Befehlszeilenoption -l verwendet werden.\n\nDer Gültigkeitsbereich von Direktiven\n\nDirektiven in den Hauptkonfigurationsdateien gelten für den\n\ngesamten Server. Wenn Sie die Konfiguration nur für einen Teil des\n\nServers verändern möchten, können Sie den\n\nGültigkeitsbereich der Direktiven beschränken, indem Sie diese\n\nin <Directory> -,\n\n<DirectoryMatch> -,\n\n<Files> -,\n\n<FilesMatch> -,\n\n<Location> - oder\n\n<LocationMatch> -Abschnitte eingefügen.\n\nDiese Abschnitte begrenzen die Anwendung der umschlossenen Direktiven\n\nauf bestimmte Pfade des Dateisystems oder auf\n\nbestimmte URLs. Sie können für eine fein abgestimmte\n\nKonfiguration auch ineinander verschachtelt werden.\n\nDer Apache besitzt die Fähigkeit, mehrere verschiedene Websites\n\ngleichzeitig zu bedienen. Dies wird virtuelles\n\nHosten genannt. Direktiven können auch in ihrem\n\nGültigkeitsgereich eingeschränkt werden, indem sie innerhalb\n\neines <VirtualHost> -Abschnittes angegeben werden.\n\nSie werden dann nur auf Anfragen für eine bestimmte Website\n\nangewendet.\n\nObwohl die meisten Direktiven in jedem dieser Abschnitte platziert\n\nwerden können, ergeben einige Direktiven in manchen Kontexten\n\nkeinen Sinn. Direktiven zur Prozesssteuerung beispielsweise\n\ndürfen nur im Kontext des Hauptservers angegeben werden. Prüfen\n\nSie den Kontext der\n\nDirektive, um herauszufinden, welche Direktiven in welche Abschnitte\n\neingefügt werden können. Weitere Informationen finden Sie unter\n\n\"Wie Directory-, Location- und Files-Abschnitte\n\narbeiten\".\n\n.htaccess-Dateien\n\nDer Apache ermöglicht die dezentrale Verwaltung der\n\nKonfiguration mittes spezieller Dateien innerhalb des\n\nWeb-Verzeichnisbaums. Diese speziellen Dateien heißen\n\ngewöhnlich .htaccess, mit der Direktive AccessFileName kann jedoch auch ein anderer\n\nName festgelegt werden. In .htaccess-Dateien angegebene\n\nDirektiven werden auf das Verzeichnis und dessen Unterverzeichnisse\n\nangewendet, in dem die Datei abgelegt ist. .htaccess-Dateien\n\nfolgen der gleichen Syntax wie die Hauptkonfigurationsdateien. Da\n\n.htaccess-Dateien bei jeder Anfrage eingelesen werden,\n\nwerden Änderungen in diesen Dateien sofort wirksam.\n\nPrüfen Sie den Kontext der Direktive, um\n\nherauszufinden, welche Direktiven in .htaccess-Dateien\n\nangegeben werden können. Darüber hinaus steuert der\n\nServeradministrator mit der Einstellung der Direktive AllowOverride in den\n\nHauptkonfigurationsdateien welche Direktiven in\n\n.htaccess-Dateien verwendet werden dürfen.\n\nWeitere Informationen über .htaccess-Dateien finden\n\nSie in der .htaccess-Einführung.\n\nVerfügbare Sprachen:  de  |\n\nen  |\n\nfr  |\n\nja  |\n\nko  |\n\ntr\n\nNotice:\n\nThis is not a Q&A section. Comments placed here should be pointed towards suggestions on improving the documentation or server, and may be removed again by our moderators if they are either implemented or considered invalid/off-topic. Questions on how to manage the Apache HTTP Server should be directed at either our IRC channel, #httpd, on Freenode, or sent to our mailing lists.",
		MetaKeywords:    "",
		CanonicalLink:   "",
		TopImage:        "/delta_answer_10.png",
	}
	article.Links = []string{
		"http://www.apache.org/",
		"http://httpd.apache.org/",
		"http://httpd.apache.org/docs/",
		"http://httpd.apache.org/lists.html",
	}

	removed := []string{"~~~REMOVED~~~"}
	err := ValidateArticle(article, &removed)
	if err != nil {
		t.Error(err)
	}
}
