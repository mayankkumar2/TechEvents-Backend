package codechef
import (
"strings"
"github.com/gocolly/colly"
)
func GetContestDataFromCodeChef () []CodechefEvent {
	c := colly.NewCollector(
		colly.AllowedDomains("www.codechef.com", "codechef.com"),
	)
	events := make([]CodechefEvent,0,200)
	c.OnHTML("table tbody",func (e *colly.HTMLElement){
		e.ForEachWithBreak("tbody", func (i int,eh *colly.HTMLElement) bool{
			eh.ForEach("tr", func (j int, ej *colly.HTMLElement) {
				var event CodechefEvent
				ej.ForEach("td", func (k int, ek *colly.HTMLElement) {
					if k == 0 {
						event.Code = ek.Text
					}
					if k == 1 {
						ek.ForEach("a", func ( _ int, el *colly.HTMLElement){
							event.HrefAddress = "https://www.codechef.com"+ el.Attr("href")
							event.Name = el.Text
						});
					}
					if k == 2 {
						event.StartDate = strings.TrimSpace(ek.Text)[:len(strings.TrimSpace(ek.Text))-9]
						event.StartTime = strings.TrimSpace(ek.Text)[len(strings.TrimSpace(ek.Text))-9:]
					}
					if k == 3 {
						event.EndDate = strings.TrimSpace(ek.Text)[:len(strings.TrimSpace(ek.Text))-9]
						event.EndTime = strings.TrimSpace(ek.Text)[len(strings.TrimSpace(ek.Text))-9:]
					}
				})
				events = append(events,event)
			})
			if i >= 1 {
				return false
			}
			return true
		})

	})
	c.Visit("https://www.codechef.com/contests")
	return events
}
type CodechefEvent struct {
	Code string `json:"code"`
	Name string `json:"name"`
	HrefAddress string `json:"href"`
	StartTime string `json:"startTime"`
	StartDate string `json:"startDate"`
	EndTime string `json:"endTime"`
	EndDate string `json:"endDate"`
}