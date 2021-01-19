<template>
  <a
    :key="video.title"
    :href="`https://www.youtube.com${video.url}`"
    target="_BLANK"
    rel="noopener noreferrer nofollow"
    class="grid grid-rows-1 mx-auto mt-5 md:grid-cols-3 hover:underline"
    style="max-width: 1000px"
  >
    <div class="w-full">
      <img
        class="object-contain w-full"
        :src="video.thumbnailUrl"
        alt="Thumbnail"
        loading="lazy"
      />
    </div>
    <div class="p-3 text-gray-700 bg-gray-200 md:col-span-2">
      <h2 class="text-2xl font-semibold break-all">
        {{ video.title }}
      </h2>
      <h3 class="mb-2">
        {{ video.channelTitle }} - {{ video.publishedAt }} -
        {{ video.views }}
      </h3>
    </div>
  </a>
</template>

<script lang="ts">
import { defineComponent } from "vue";

import Video from "../models/Video";

export default defineComponent({
  props: {
    video: {
      type: Video,
      required: true,
    },
  },
});
</script>