![](http://i.imgur.com/0ZWbc3T.png)

# Festivus

Festivus is a Slack app that will help you plan and celebrate the holidays easier.

This 48 hour project was made as a submission to the [gophergala 2016](http://gophergala.com/) [golang](http://golang.org) hackathon. 

Our last years submission to Gopher Gala [videq - High quality video encoding for modern web in golang](https://github.com/gophergala/videq) was fun, so we gave Gopher Gala another go this year. :)## TLDR; Wanna try it out?

Great! Just click the *Add to Slack* button and follow the lights. 
<a href="https://festivus.nivas.hr/add"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x"></a>
Once you have it in Slack:

* To see how long till Festivus:
/festivus

* To see upcoming holidays in United States of America:
/festivus us 

* In France
/festivus fr

* In UK
/festivus uk

* In ​India
/festivus in

* For more locales:
/festivus help

## InspirationOur team comes from Croatia, young and small European country in which (some say) it’s employed residents value non working holidays a bit more than elsewhere in the world (there are roumours about Italy and France closing in holiday race, feedback needed).

Over our short history we managed to invent many holidays, which on some years fall on week days - effectively making them – governmental non working days. > ***The sport comes from combining your vacation business off days with holiday non working days in order to maximize your vacation time potential (VTP).***The inspiration for the name comes from legendary Seinfeld episode in which Frank Costanza (Jerry Stiller) created Festivus as an alternative holiday in response to the commercialization of Christmas <https://en.wikipedia.org/wiki/Festivus>The inspiration in using Slack comes from article published in [Verge](http://www.theverge.com/2016/1/6/10718282/internet-bots-messaging-slack-facebook-m) some weeks back telling how year 2016. will be the year of Bots with killer AI.## The problem(s)Our teem started using Slack just recently, and we did not experiment with Slack integrations before whatsoever.We saw Gopher Gala as a great opportunity to experiment with golang and Slack. 
## Solution[Programming!](http://c00kiemon5ter.github.io/code/2011/04/16/Development-Methodologies.html)## How it works

### Install Festivus to your Slack team

![](http://i.imgur.com/etXd1dc.gif)
### Personal usage

![](http://i.imgur.com/jBM56k4.gif)

<!--
### Show help and available locales

![](http://i.imgur.com/8PB2OjE.gif)

### Show upcoming holidays for Croatia

![](http://i.imgur.com/vzx0550.gif)

### How many days until Festivus?

![](http://i.imgur.com/7ywryV0.gif)
-->

## How to build it for your self

You must create Slack app or & define the command and it's endpoint at <https://api.slack.com/applications>.
Then build and run from command line. 

We use following packages so you have to go get them:

```
go get github.com/nlopes/slack
go get golang.org/x/oauth2
```

After you have successfully go build, run it:

```
Usage: ./festivus --address ":8080" --client_id "YOUR_ID" --client_secret "YOUR_SECRET"
```
## Challenges we ran intoOriginally, our feature list for Festivus was huge, as you would expect from over the top full blown bot with AI and everything.  After we dug into Slack API, it became pretty clear that 48 hours is not enough for proper implementation of multi team support via websockets and Slack’s Real Time Messaging API.

We scoped it to be able to complete it in time and still have fun doing it.
## Accomplishments that we are proud ofWe had opportunity to test drive Slack apps, bots, RTM API, commands etc.

The list of holidays per country in the world is hard to find. Luckily our friends at Google have added some of holidays into calendar.google.com in form of custom calendars. @luzel found a technique how to fetch them, and @dvrkps built golang helper parser package around the data files.
## What we learned48 hours is not much, use it wisely.We should have automated git pull/go build/restart on staging server.## What's next for Festivus

Fork, code cleanup, refractoring should be done absolutely first; We are pretty anxious to transfer it to a form of a Slck Bot and continue exploring Slack platform further; Neural Networks and Collaborative Filtering.

## The TeamBig shout to fine Festivus team & supporters:

* [Neven Jacmenovic](https://twitter.com/guycalledseven)
* [Matej Baco](https://twitter.com/matejbaco)
* [Luka Uzel](https://twitter.com/LukaUzel)
* [Davor Kapsa](https://twitter.com/dvrkps)
* [Alen Cvitkovic](https://twitter.com/alencvitkovic)
