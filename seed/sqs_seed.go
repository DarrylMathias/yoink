package seed

import (
	"fmt"
	"math/rand/v2"
	mysqs "yoink/utils/myaws/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func SeedSQS(){
	var SeedURLs = []string{
		"https://en.wikipedia.org",
		"https://www.nature.com",
		"https://techcrunch.com",
		"https://github.com",
		"https://www.mit.edu",
		"https://www.reuters.com",
		"https://www.bloomberg.com",
		"https://www.data.gov",
		"https://www.apache.org",
		"https://www.imdb.com",

		"https://www.britannica.com",
		"https://www.science.org",
		"https://www.theverge.com",
		"https://stackoverflow.com",
		"https://www.stanford.edu",
		"https://apnews.com",
		"https://www.investopedia.com",
		"https://www.census.gov",
		"https://www.linuxfoundation.org",
		"https://www.goodreads.com",

		"https://www.howstuffworks.com",
		"https://arxiv.org",
		"https://arstechnica.com",
		"https://go.dev",
		"https://www.harvard.edu",
		"https://www.bbc.com",
		"https://www.forbes.com",
		"https://www.usa.gov",
		"https://www.mozilla.org",
		"https://www.metmuseum.org",

		"https://www.worldhistory.org",
		"https://www.cern.ch",
		"https://www.wired.com",
		"https://developer.mozilla.org",
		"https://www.cam.ac.uk",
		"https://www.npr.org",
		"https://www.marketwatch.com",
		"https://www.gov.uk",
		"https://www.python.org",
		"https://www.moma.org",

		"https://www.nationalgeographic.com",
		"https://www.nasa.gov",
		"https://www.tomshardware.com",
		"https://docs.python.org",
		"https://www.ox.ac.uk",
		"https://www.theguardian.com",
		"https://www.cnbc.com",
		"https://europa.eu",
		"https://wordpress.org",
		"https://www.loc.gov",

		"https://www.smithsonianmag.com",
		"https://www.noaa.gov",
		"https://www.zdnet.com",
		"https://kubernetes.io",
		"https://www.coursera.org",
		"https://www.aljazeera.com",
		"https://www.sec.gov",
		"https://www.un.org",
		"https://www.fsf.org",
		"https://www.gutenberg.org",

		"https://www.livescience.com",
		"https://www.jpl.nasa.gov",
		"https://www.engadget.com",
		"https://www.postgresql.org",
		"https://www.edx.org",
		"https://www.economist.com",
		"https://www.imf.org",
		"https://www.who.int",
		"https://www.eclipse.org",
		"https://www.ted.com",

		"https://www.scientificamerican.com",
		"https://www.esa.int",
		"https://www.anandtech.com",
		"https://redis.io",
		"https://ocw.mit.edu",
		"https://www.cnn.com",
		"https://www.worldbank.org",
		"https://www.unesco.org",
		"https://www.gnome.org",
		"https://www.bbc.co.uk",

		"https://www.discovery.com",
		"https://www.aps.org",
		"https://www.bleepingcomputer.com",
		"https://www.docker.com",
		"https://www.khanacademy.org",
		"https://www.nytimes.com",
		"https://www.federalreserve.gov",
		"https://www.oecd.org",
		"https://www.kernel.org",
		"https://www.nationalgeographic.com",

		"https://www.history.com",
		"https://www.aaas.org",
		"https://www.howtogeek.com",
		"https://www.gnu.org",
		"https://www.udacity.com",
		"https://www.wsj.com",
		"https://www.mckinsey.com",
		"https://www.nist.gov",
		"https://www.openstreetmap.org",
		"https://www.projecteuclid.org",
	}
	mysqs.GetQueueURL()

	rand.Shuffle(len(SeedURLs), func(i, j int) {
		SeedURLs[i], SeedURLs[j] = SeedURLs[j], SeedURLs[i]
	})
	for i, url := range SeedURLs{
		output, err := mysqs.SendMessage(url)
		if err != nil{
			panic(fmt.Errorf("error in sqs send message --- %s", err.Error()))
		}
		fmt.Printf("success - url %d: %v\n", i+1, aws.ToString(output.MD5OfMessageBody))
	}
}