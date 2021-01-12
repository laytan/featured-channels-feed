export default class Video {
  id: string;
  channelId: string;
  channelTitle: string;
  description: string;
  publishedAt: Date;
  thumbnail: string;
  title: string;

  constructor(
    id: string,
    channelId: string,
    channelTitle: string,
    description: string,
    publishedAt: string,
    thumbnail: string,
    title: string
  ) {
    this.id = id;
    this.channelId = channelId;
    this.channelTitle = channelTitle;
    this.description = description;
    this.publishedAt = new Date(publishedAt);
    this.thumbnail = thumbnail;
    this.title = title;
  }
}
