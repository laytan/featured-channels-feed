import Video from "./Video";

export default class Channel {
  subscriberCount: string;
  urlTitle: string;
  displayTitle: string;
  latestVideos: Array<Video>;

  constructor(
    subscriberCount: string,
    urlTitle: string,
    displayTitle: string,
    latestVideos: Array<Video>
  ) {
    this.subscriberCount = subscriberCount;
    this.urlTitle = urlTitle;
    this.displayTitle = displayTitle;
    this.latestVideos = latestVideos;
  }
}
