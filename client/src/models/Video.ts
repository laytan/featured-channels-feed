export default class Video {
  url: string;
  publishedAt: string;
  thumbnailUrl: string;
  title: string;
  views: string;
  channelTitle: string;

  constructor(
    url: string,
    publishedAt: string,
    thumbnailUrl: string,
    title: string,
    views: string,
    channelTitle: string
  ) {
    this.url = url;
    this.publishedAt = publishedAt;
    this.thumbnailUrl = thumbnailUrl;
    this.title = title;
    this.views = views;
    this.channelTitle = channelTitle;
  }

  // Turns 3 months ago, 1 year ago etc. into the amount of seconds ago
  getSeconds(): number {
    const split = this.publishedAt.split(" ");
    if (split[0] === "Streamed") {
      split.shift();
    }

    switch (split[1]) {
      case "second":
      case "seconds":
        return Number(split[0]);
      case "minute":
      case "minutes":
        return Number(split[0]) * 60;
      case "hour":
      case "hours":
        return Number(split[0]) * 60 * 60;
      case "day":
      case "days":
        return Number(split[0]) * 60 * 60 * 24;
      case "week":
      case "weeks":
        return Number(split[0]) * 60 * 60 * 24 * 7;
      case "month":
      case "months":
        return Number(split[0]) * 60 * 60 * 24 * 7 * 30;
      case "year":
      case "years":
        return Number(split[0]) * 60 * 60 * 24 * 7 * 30 * 365;
      default:
        console.error(
          `${this.title} getSeconds did not get valid publishedAt: ${split[1]}`
        );
        return Number.MAX_SAFE_INTEGER;
    }
  }
}
