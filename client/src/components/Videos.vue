<template>
  <div class="mt-2">
    <p class="ml-4 text-sm">
      Showing {{ videos.length }} video{{
        videos.length > 1 ? "s" : ""
      }}
      published by featured channels of {{ currChan }} (the 8 latest videos per
      channel).
    </p>
    <Video v-for="video in videos" :key="video.title" :video="video"></Video>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";

import Video from "./Video.vue";
import VideoModel from "../models/Video";

export default defineComponent({
  components: {
    Video,
  },
  props: {
    videos: {
      type: Array as PropType<Array<VideoModel>>,
      required: true,
    },
  },
});
</script>