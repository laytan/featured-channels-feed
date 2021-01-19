# Featured Channel Feed

Featured Channel Feed allows you to see all the latest videos of a channels featured channels.

Try it out on [s.laytanlaats.com/channels](https://featured-channels-feed.laytanlaats.com)

## Screenshots

![Screenshot]("https://github.com/laytan/featured-channels-feed/raw/main/assets/screen1.png")

![Screenshot 2]("https://github.com/laytan/featured-channels-feed/raw/main/assets/screen2.png")

![Screenshot 3]("https://github.com/laytan/featured-channels-feed/raw/main/assets/screen3.png")

## Goals

I used this project to get better at Golang, VueJS and TailwindCSS.

It was the first time i scraped using Golang and i was surprised by the simplicity of the geziyor library which handles everything from starting up chrome and the javascript engine to being able to navigate the DOM very easily with go query.

Using Docker to deploy and develop the back-end service was also a good way to learn Docker.

## The project

I came up with this project because i noticed there was no way of showing this in the original Youtube UI.

The first thing i checked was if the Youtube API could help me but i unfortunately found a bug in it where the featured channels would not be added to the response. So i filed a bug report and started going the route of scraping Youtube.


